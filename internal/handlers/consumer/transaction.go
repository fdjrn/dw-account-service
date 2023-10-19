package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/db/repository"
	"github.com/dw-account-service/internal/utilities"
	"github.com/dw-account-service/internal/utilities/crypt"
	"github.com/dw-account-service/internal/utilities/str"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

type TransactionHandler struct {
	transactionRepository repository.TransactionRepository
	accountRepository     repository.AccountRepository
}

func NewTransactionHandler() TransactionHandler {
	return TransactionHandler{
		transactionRepository: repository.NewTransactionRepository(),
		accountRepository:     repository.NewAccountRepository(),
	}
}

func modifyBalance(current, amount int64, keys string, transType int) (int64, string) {
	var last int64

	if transType == utilities.TransTypeTopUp {
		last = current + amount
	} else {
		last = current - amount
	}

	encrypted, err := crypt.Encrypt(
		[]byte(keys),
		fmt.Sprintf("%016s", strconv.FormatInt(last, 10)),
	)

	if err != nil {
		return last, "-"
	}

	return last, encrypted
}

func (t *TransactionHandler) doValidation(data *entity.BalanceTransaction) (*entity.BalanceTransaction, error) {

	t.accountRepository.Entity.PartnerID = data.PartnerID
	t.accountRepository.Entity.MerchantID = data.MerchantID
	t.accountRepository.Entity.TerminalID = data.TerminalID

	if data.TerminalID == "" {
		t.accountRepository.Entity.Type = utilities.AccountTypeMerchant
	}

	account, err := t.accountRepository.FindOne()
	if err != nil {
		// invalid account infos
		data.Status = utilities.TrxStatusInvalidAccount
		if errors.Is(err, mongo.ErrNoDocuments) {
			return data, errors.New("unable to find account detail with supplied parameters")
		}
		return data, err
	}

	// validate account is in active status
	if !account.Active {
		utilities.Log.Println("| account deactivated, balance update cannot be processed ")
		data.Status = utilities.TrxStatusInvalidAccount
		return data, err
	}

	t.accountRepository.Entity = account

	data.LastBalance = account.LastBalanceNumeric

	if data.TransType == utilities.TransTypeDistribution {
		memberCount, err2 := t.accountRepository.CountMembers()
		if err2 != nil {
			data.Status = utilities.TrxStatusFailed
			return data, errors.New("failed to get total active members")
		}

		data.TotalAmount = memberCount * data.Items[0].Amount
		data.Items[0].Qty = int(memberCount)
	}

	if data.TransType != utilities.TransTypeTopUp && data.LastBalance < data.TotalAmount {
		data.Status = utilities.TrxStatusInsufficientFund
		return data, errors.New("insufficient account balance")
	}

	return data, nil
}

func (t *TransactionHandler) DoHandleTransactionRequest(message *sarama.ConsumerMessage) (*entity.BalanceTransaction, error) {

	var err error

	data := new(entity.BalanceTransaction)
	err = json.Unmarshal(message.Value, &data)
	if err != nil {
		return nil, err
	}

	// validate account partner, merchant and terminal
	data, err = t.doValidation(data)
	if err != nil {
		return data, err
	}

	// modify last balance with amount of transaction, based on transType value
	lb, encLb := modifyBalance(
		t.accountRepository.Entity.LastBalanceNumeric,
		data.TotalAmount,
		t.accountRepository.Entity.SecretKey,
		data.TransType,
	)

	data.LastBalance = lb
	data.LastBalanceEncrypted = encLb

	// update account last balance
	t.transactionRepository.Entity = data
	updatedAccount, err := t.transactionRepository.UpdateBalance()
	if err != nil {
		utilities.Log.Println("| failed to update balance, with err: ", err.Error())
		data.Status = utilities.TrxStatusFailed
		return data, err
	}

	// return entity.BalanceTransaction data with status Success ("00")
	trxDate := time.Now()
	data.TransDateNumeric = trxDate.UnixMilli()
	data.TransDate = trxDate.Format("20060102150405")

	data.ReceiptNumber = str.GenerateReceiptNumber(data.TransType, "")
	data.LastBalance = updatedAccount.LastBalanceNumeric
	data.Status = utilities.TrxStatusSuccess

	return data, nil
}
