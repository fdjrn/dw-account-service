package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/kafka/topic"
	"github.com/dw-account-service/internal/utilities"
	"github.com/dw-account-service/internal/utilities/crypt"
	"github.com/dw-account-service/internal/utilities/str"
	"strconv"
	"sync"
	"time"
)

type DistributionTrx struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewDistributionTrx() DistributionTrx {
	return DistributionTrx{
		accountRepo:     repository.NewAccountRepository(),
		transactionRepo: repository.NewTransactionRepository(),
	}
}

func (d *DistributionTrx) DoBalanceDistribution(data *entity.BalanceTransaction) error {

	params := new(entity.PaginatedAccountRequest)
	params.PartnerID = data.PartnerID
	params.MerchantID = data.MerchantID
	params.Type = utilities.AccountTypeRegular
	params.Status = utilities.AccountStatusActive

	members, err := d.accountRepo.FindMembers(params)
	if err != nil {
		return err
	}

	// pipeline 1: job distribution
	chanJobIndex := generateWorkerData(members)

	// pipeline 2: update balance
	workerCount := 10
	tWorker := len(members) / 4

	if tWorker < 10 {
		workerCount = 5
	}

	chanUpdateResult := d.doBatchUpdateBalance(chanJobIndex, workerCount, data)

	totalJob := 0
	successJob := 0
	successProduce := 0
	for result := range chanUpdateResult {
		if result.Err != nil {
			utilities.Log.Println("| error on update balance on account id: ", result.Data.ID)
		} else {
			payload, _ := json.Marshal(result.Data)
			err = ProduceMsg(topic.DistributionResultMembers, payload)
			if err != nil {
				utilities.Log.Printf("| cannot produce message (%s) for topic: %s, with err: %s",
					result.Data.ReceiptNumber,
					topic.DistributionResultMembers,
					err.Error(),
				)
			} else {
				//utilities.Log.Printf("| account balance with id: %s, has been successfully processed with receipt number: %s\n",
				//	result.Data.ID,
				//	result.Data.ReceiptNumber,
				//)
				successProduce++
			}
			successJob++
		}
		totalJob++
	}

	utilities.Log.Printf("| %d/%d of member balances has been successfully updated", successJob, totalJob)
	utilities.Log.Printf("| %d/%d of success messages has been successfully produced", successProduce, successJob)

	return nil
}

func (d *DistributionTrx) doBatchUpdateBalance(chanIn <-chan entity.AccountBalance, workerCount int, data *entity.BalanceTransaction) <-chan entity.BalanceDistributionInfo {
	chanOut := make(chan entity.BalanceDistributionInfo)

	wgUpdateBalance := new(sync.WaitGroup)
	wgUpdateBalance.Add(workerCount)

	go func() {
		for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
			go func(idx int) {
				for accountBalance := range chanIn {

					// update balance
					accountBalance.LastBalanceNumeric += data.Items[0].Amount
					encrypted, err := crypt.Encrypt(
						[]byte(accountBalance.SecretKey),
						fmt.Sprintf("%016s", strconv.FormatInt(accountBalance.LastBalanceNumeric, 10)),
					)
					if err != nil {
						encrypted = "-"
					}

					accountBalance.LastBalance = encrypted

					d.transactionRepo.Entity = new(entity.BalanceTransaction)
					d.transactionRepo.Entity.MerchantID = accountBalance.MerchantID
					d.transactionRepo.Entity.PartnerID = accountBalance.PartnerID
					d.transactionRepo.Entity.TerminalID = accountBalance.TerminalID
					d.transactionRepo.Entity.LastBalance = accountBalance.LastBalanceNumeric
					d.transactionRepo.Entity.LastBalanceEncrypted = encrypted

					account, err := d.transactionRepo.UpdateBalance()

					// populate chanOut Data
					var items []entity.TransactionItem
					items = append(items, entity.TransactionItem{
						Name:   "Receiving Balance From: " + account.PartnerID + "-" + account.MerchantID,
						Amount: data.Items[0].Amount,
						Qty:    1,
					})

					trxDate := time.Now()
					//utilities.Log.Println("| using worker", idx, "to update account balance with id: ", account.ID)

					chanOut <- entity.BalanceDistributionInfo{
						Data: entity.BalanceTransaction{
							TransDate:        trxDate.Format("20060102150405"),
							TransDateNumeric: trxDate.UnixMilli(),
							ReferenceNo:      data.ReferenceNo,
							ReceiptNumber:    str.GenerateReceiptNumber(data.TransType, ""),
							LastBalance:      account.LastBalanceNumeric,
							Status:           data.Status,
							TransType:        data.TransType,
							PartnerTransDate: data.PartnerTransDate,
							PartnerRefNumber: data.PartnerRefNumber,
							PartnerID:        account.PartnerID,
							MerchantID:       account.MerchantID,
							TerminalID:       account.TerminalID,
							TerminalName:     account.TerminalName,
							TotalAmount:      data.Items[0].Amount,
							Items:            items,
							CreatedAt:        trxDate.UnixMilli(),
							UpdatedAt:        trxDate.UnixMilli(),
						},
						WorkerIndex: idx,
						Err:         err,
					}
				}
				wgUpdateBalance.Done()
			}(workerIdx)
		}
	}()

	go func() {
		wgUpdateBalance.Wait()
		close(chanOut)
	}()

	return chanOut
}

func generateWorkerData(members []entity.AccountBalance) <-chan entity.AccountBalance {
	chanOut := make(chan entity.AccountBalance)

	go func() {
		for i := 0; i < len(members); i++ {
			chanOut <- members[i]
		}

		close(chanOut)
	}()

	return chanOut
}
