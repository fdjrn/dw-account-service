package tools

import (
	"errors"
	"github.com/dw-account-service/internal/db/entity"
)

func ValidateRequest(payload interface{}) (interface{}, error) {
	var msg []string

	switch p := payload.(type) {
	case *entity.BalanceTopUp:
		if p.UniqueID == "" {
			msg = append(msg, "uniqueId cannot be empty.")
		}

		if p.Amount == 0 {
			msg = append(msg, "topup amount must be greater than 0.")
		}

		if p.InRefNumber == "" {
			msg = append(msg, "inRefNumber cannot be empty.")
		}

		if p.ExRefNumber == "" {
			msg = append(msg, "exRefNumber cannot be empty.")
		}

		if p.TransDate == 0 {
			msg = append(msg, "exRefNumber must be greater than 0.")
		}

		if len(msg) > 0 {
			return msg, errors.New("request validation status failed")
		}

		return msg, nil

	case *entity.BalanceDeduction:
		if p.UniqueID == "" {
			msg = append(msg, "uniqueId cannot be empty.")
		}

		if p.Amount == 0 {
			msg = append(msg, "topup amount must be greater than 0.")
		}

		if p.TransType == 0 {
			msg = append(msg, "transType cannot be empty.")
		}

		if p.Description == "" {
			msg = append(msg, "transaction description cannot be empty.")
		}

		if p.InvoiceNumber == "" {
			msg = append(msg, "invoiceNumber cannot be empty.")
		}

		if len(msg) > 0 {
			return msg, errors.New("request validation status failed")
		}

		return msg, nil

	case *entity.AccountBalance:
		if p.UniqueID == "" {
			msg = append(msg, "uniqueId cannot be empty.")
		}

		if p.PartnerID == "" {
			msg = append(msg, "partnerId cannot be empty.")
		}

		if p.MerchantID == "" {
			msg = append(msg, "merchantId cannot be empty.")
		}

		if p.Type == 0 {
			msg = append(msg, "type cannot be empty. eg: 1 >> regular account | 2 >> sub-account")
		}

		if p.Type == 2 && p.MainAccountID == "" {
			msg = append(msg, "mainAccountId cannot be empty.")
		}

		if len(msg) > 0 {
			return msg, errors.New("request validation status failed")
		}

		return msg, nil
	default:
		return msg, nil
	}

}
