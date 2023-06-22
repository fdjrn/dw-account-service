package validator

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

	case *entity.MerchantTrxRequest:
		if p.PartnerID == "" {
			msg = append(msg, "partnerId cannot be empty.")
		}

		if p.MerchantID == "" {
			msg = append(msg, "merchantId cannot be empty.")
		}

		if p.Amount == 0 {
			msg = append(msg, "topupAmount must be greater than 0.")
		}

		if p.PartnerRefNumber == "" {
			msg = append(msg, "partnerRefNumber cannot be empty.")
		}

		if p.PartnerTransDate == "" {
			msg = append(msg, "partnerTransDate cannot be empty.")
		}

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

	case *entity.AccountBalance:

		if p.Type == 0 {
			msg = append(msg, "type cannot be empty. eg: 1 (regular) or 2 (admin)")
			return msg, errors.New("request validation status failed")
		}

		if p.Type > 2 {
			msg = append(msg, "unsupported account type. only: 1 (regular) or 2 (merchant)")
			return msg, errors.New("request validation status failed")
		}

		// if non-merchant account type
		if p.Type == 1 && p.TerminalID == "" {
			msg = append(msg, "terminalId cannot be empty.")
		}

		if p.PartnerID == "" {
			msg = append(msg, "partnerId cannot be empty.")
		}

		if p.MerchantID == "" {
			msg = append(msg, "merchantId cannot be empty.")
		}

	//case *entity.DefaultMerchantRequest:
	//	if p.PartnerID == "" {
	//		msg = append(msg, "partnerId cannot be empty.")
	//	}
	//
	//	if p.MerchantID == "" {
	//		msg = append(msg, "merchantId cannot be empty.")
	//	}
	default:
	}

	if len(msg) > 0 {
		return msg, errors.New("request validation status failed")
	}
	return msg, nil

}

func ValidateAccountDetailRequest(p *entity.AccountBalance) (interface{}, error) {
	var errMsq []string

	// if non-merchant account type
	if p.TerminalID == "" {
		errMsq = append(errMsq, "terminalId cannot be empty.")
	}

	if p.PartnerID == "" {
		errMsq = append(errMsq, "partnerId cannot be empty.")
	}

	if p.MerchantID == "" {
		errMsq = append(errMsq, "merchantId cannot be empty.")
	}

	if len(errMsq) > 0 {
		return errMsq, errors.New("request validation status failed")
	}

	return errMsq, nil
}
