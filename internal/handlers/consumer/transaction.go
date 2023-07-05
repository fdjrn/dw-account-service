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

func (t *TransactionHandler) DoHandleTopupTransaction(message *sarama.ConsumerMessage) (*entity.BalanceTransaction, error) {

	data := new(entity.BalanceTransaction)
	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		return nil, err
	}

	// validate account partner, merchant and terminal
	t.accountRepository.Entity.PartnerID = data.PartnerID
	t.accountRepository.Entity.MerchantID = data.MerchantID
	t.accountRepository.Entity.TerminalID = data.TerminalID

	if data.TerminalID == "" {
		t.accountRepository.Entity.Type = repository.AccountTypeMerchant
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

	// add last balance with amount of topup
	data.LastBalance = account.LastBalanceNumeric + data.Items[0].Amount
	data.LastBalanceEncrypted, err = crypt.Encrypt(
		[]byte(account.SecretKey),
		fmt.Sprintf("%016s", strconv.FormatInt(data.LastBalance, 10)),
	)

	// update account last balance
	t.transactionRepository.Entity = data
	account, err = t.transactionRepository.UpdateBalance()
	if err != nil {
		data.Status = utilities.TrxStatusFailed
	}

	// return entity.BalanceTransaction data with status Success ("00")
	data.TransDate = time.Now().Format("20060102150405")
	data.ReceiptNumber = str.GenerateReceiptNumber(utilities.TransTopUp, "")
	data.LastBalance = account.LastBalanceNumeric
	data.Status = utilities.TrxStatusSuccess

	return data, nil
}

func (t *TransactionHandler) DoHandleDeductTransaction(message *sarama.ConsumerMessage) (*entity.BalanceTransaction, error) {

	data := new(entity.BalanceTransaction)
	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		return nil, err
	}

	// validate account partner, merchant and terminal
	t.accountRepository.Entity.PartnerID = data.PartnerID
	t.accountRepository.Entity.MerchantID = data.MerchantID
	t.accountRepository.Entity.TerminalID = data.TerminalID

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

	// subtract last balance with amount of deduct
	data.LastBalance = account.LastBalanceNumeric - data.Items[0].Amount
	data.LastBalanceEncrypted, err = crypt.Encrypt(
		[]byte(account.SecretKey),
		fmt.Sprintf("%016s", strconv.FormatInt(data.LastBalance, 10)),
	)

	if data.LastBalance < 0 {
		// avoid update balance.
		// return entity.BalanceTransaction data with status Insufficient Funds ("06")
		data.Status = utilities.TrxStatusInsufficientFund
		return data, errors.New("insufficient account balance")
	}

	// update account last balance
	t.transactionRepository.Entity = data
	account, err = t.transactionRepository.UpdateBalance()
	if err != nil {
		data.Status = utilities.TrxStatusFailed
	}

	// return entity.BalanceTransaction data with status Success ("00")
	data.TransDate = time.Now().Format("20060102150405")
	data.ReceiptNumber = str.GenerateReceiptNumber(utilities.TransTopUp, "")
	data.LastBalance = account.LastBalanceNumeric
	data.Status = utilities.TrxStatusSuccess

	return data, nil
}
