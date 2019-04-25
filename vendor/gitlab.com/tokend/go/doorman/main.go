package doorman

import (
	"net/http"

	"gitlab.com/tokend/go/doorman/data"
)

type SignerConstraint func(*http.Request, Doorman) error

func New(passAllChecks bool, accountQ data.AccountQ) Doorman {
	return &doorman{
		SkipChecker: data.NewChecker(passAllChecks),
		AccountQ:    accountQ,
	}
}

func NewWithChecker(skipChecker data.SkipChecker, accountQ data.AccountQ) Doorman {
	return &doorman{
		SkipChecker: skipChecker,
		AccountQ:    accountQ,
	}
}
