package validator

import (
	"errors"
	"github.com/dw-account-service/internal/db/entity"
)

func ValidateRequest(payload interface{}) (interface{}, error) {
	var msg []string

	switch p := payload.(type) {

	case *entity.AccountBalance:

		if p.Type == 0 {
			msg = append(msg, "type cannot be empty. eg: 1 (regular) or 2 (merchant)")
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

	default:
	}

	if len(msg) > 0 {
		return msg, errors.New("request validation status failed")
	}
	return msg, nil

}

//
//func ValidateAccountDetailRequest(p *entity.AccountBalance) (interface{}, error) {
//	var errMsq []string
//
//	// if non-merchant account type
//	if p.TerminalID == "" {
//		errMsq = append(errMsq, "terminalId cannot be empty.")
//	}
//
//	if p.PartnerID == "" {
//		errMsq = append(errMsq, "partnerId cannot be empty.")
//	}
//
//	if p.MerchantID == "" {
//		errMsq = append(errMsq, "merchantId cannot be empty.")
//	}
//
//	if len(errMsq) > 0 {
//		return errMsq, errors.New("request validation status failed")
//	}
//
//	return errMsq, nil
//}
