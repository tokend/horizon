package requests

import (
	"net/http"

	amount2 "gitlab.com/tokend/go/amount"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/regources/rgenerated"
)

//GetCalculatedFees - represents params to be specified for Get Fees handler
type GetCalculatedFees struct {
	*base
	Address string
	Asset   string
	Subtype int64
	FeeType int32
	Amount  rgenerated.Amount
}

// NewGetCalculatedFees returns the new instance of GetCalculatedFees request
func NewGetCalculatedFees(r *http.Request) (*GetCalculatedFees, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}
	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	asset := b.getString("asset")
	if asset == "" {
		return nil, errors.New("asset code is required")
	}

	subtype, err := b.getInt64("subtype")
	if err != nil {
		return nil, errors.Wrap(err, "fee subtype is required")
	}

	feeType, err := b.getInt32("fee_type")
	if err != nil {
		return nil, errors.Wrap(err, "fee type is required")
	}

	amountRaw := b.getString("amount")
	if amountRaw == "" {
		return nil, errors.New("amount is required")
	}

	amount, err := amount2.Parse(amountRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse amount")
	}

	request := GetCalculatedFees{
		base:    b,
		Address: address,
		Asset:   asset,
		Subtype: subtype,
		FeeType: feeType,
		Amount:  rgenerated.Amount(amount),
	}

	return &request, nil
}
