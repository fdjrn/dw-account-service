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
	accountRepo repository.AccountRepository
	trxRepo     repository.TransactionRepository
}

func NewDistributionTrx() DistributionTrx {
	return DistributionTrx{
		accountRepo: repository.NewAccountRepository(),
		trxRepo:     repository.NewTransactionRepository(),
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
	chanJobIndex := generateDataIndexes(members)

	// pipeline 2: update balance
	workerCount := 10
	tWorker := len(members) / 4

	if tWorker < 10 {
		workerCount = 2
	}

	chanUpdateResult := d.doBatchUpdateBalance(chanJobIndex, workerCount, data)

	for result := range chanUpdateResult {
		if result.Err != nil {
			utilities.Log.Println("| error on update balance on account id: ", result.Data.ID)
			continue
		}

		payload, _ := json.Marshal(result.Data)
		err = ProduceMsg(topic.DistributionResult, payload)
		if err != nil {
			utilities.Log.Println("| cannot produce message for topic: ", topic.DistributionResult, ", with err: ", err.Error())
			continue
		}

		utilities.Log.Printf("| %s with RefNo: %s, has been successfully processed with receipt number: %s\n",
			"balance distribution",
			result.Data.ReferenceNo,
			result.Data.ReceiptNumber,
		)

	}

	return nil
}

func (d *DistributionTrx) doBatchUpdateBalance(chanIn <-chan entity.AccountBalanceDistInfo, workerCount int, data *entity.BalanceTransaction) <-chan entity.BalanceDistributionInfo {
	chanOut := make(chan entity.BalanceDistributionInfo)

	wgUpdateBalance := new(sync.WaitGroup)
	wgUpdateBalance.Add(workerCount)

	go func() {
		for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
			go func(workerIdx int) {
				for distInfo := range chanIn {

					// update balance
					distInfo.Account.LastBalanceNumeric += data.Items[0].Amount
					encrypted, err := crypt.Encrypt(
						[]byte(distInfo.Account.SecretKey),
						fmt.Sprintf("%016s", strconv.FormatInt(distInfo.Account.LastBalanceNumeric, 10)),
					)
					if err != nil {
						encrypted = "-"
					}

					distInfo.Account.LastBalance = encrypted

					//t := NewTransactionHandler()
					d.trxRepo.Entity = new(entity.BalanceTransaction) // distInfo.Account
					d.trxRepo.Entity.MerchantID = distInfo.Account.MerchantID
					d.trxRepo.Entity.PartnerID = distInfo.Account.PartnerID
					d.trxRepo.Entity.TerminalID = distInfo.Account.TerminalID
					d.trxRepo.Entity.LastBalance = distInfo.Account.LastBalanceNumeric
					d.trxRepo.Entity.LastBalanceEncrypted = encrypted

					updAccount, err := d.trxRepo.UpdateBalance()
					trxDate := time.Now()

					// populate chanOut Data
					var items []entity.TransactionItem
					items = append(items, entity.TransactionItem{
						Name:   "Receiving Balance From: " + updAccount.PartnerID + "-" + updAccount.MerchantID,
						Amount: data.Items[0].Amount,
						Qty:    1,
					})

					//utilities.Log.Println("| worker", workerIdx, "working on update Account Balance ID: ", updAccount.ID)

					chanOut <- entity.BalanceDistributionInfo{
						Data: entity.BalanceTransaction{
							TransDate:        trxDate.Format("20060102150405"),
							TransDateNumeric: trxDate.UnixMilli(),
							ReferenceNo:      data.ReferenceNo,
							ReceiptNumber:    str.GenerateReceiptNumber(data.TransType, ""),
							LastBalance:      updAccount.LastBalanceNumeric,
							Status:           data.Status,
							TransType:        data.TransType,
							PartnerTransDate: data.PartnerTransDate,
							PartnerRefNumber: data.PartnerRefNumber,
							PartnerID:        updAccount.PartnerID,
							MerchantID:       updAccount.MerchantID,
							TerminalID:       updAccount.TerminalID,
							TerminalName:     updAccount.TerminalName,
							TotalAmount:      data.Items[0].Amount,
							Items:            items,
							CreatedAt:        trxDate.UnixMilli(),
							UpdatedAt:        trxDate.UnixMilli(),
						},
						WorkerIndex: workerIdx,
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

func generateDataIndexes(members []entity.AccountBalance) <-chan entity.AccountBalanceDistInfo {
	chanOut := make(chan entity.AccountBalanceDistInfo)

	go func() {
		for i := 0; i < len(members); i++ {
			chanOut <- entity.AccountBalanceDistInfo{
				Index:   i,
				Account: &members[i],
			}
		}

		close(chanOut)
	}()

	return chanOut
}
