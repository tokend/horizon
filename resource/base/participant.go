package base

import (
	"time"

	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/db2/history"
	"github.com/go-errors/errors"
)

type BaseEffects interface {
}

type ManageOfferEffects struct {
	BaseAssetCode  string `json:"base_asset_code"`
	QuoteAssetCode string `json:"quote_asset_code"`
	IsBuy          bool   `json:"is_buy"`
	BaseAmount     string `json:"base_amount"`
	QuoteAmount    string `json:"quote_amount"`
	Price          string `json:"price"`
}

type DemurrageEffects struct {
	Asset      string    `json:"asset"`
	Amount     string    `json:"amount"`
	PeriodFrom time.Time `json:"period_from"`
	PeriodTo   time.Time `json:"period_to"`
}

type Participant struct {
	AccountID string      `json:"account_id,omitempty"`
	BalanceID string      `json:"balance_id,omitempty"`
	Email     string      `json:"email,omitempty"`
	FullName  string      `json:"full_name,omitempty"`
	Nickname  string      `json:"nickname,omitempty"`
	Effects   BaseEffects `json:"effects,omitempty"`
}

func (f *Participant) Populate(p *history.Participant, opType xdr.OperationType, public bool) error {
	if !public {
		f.AccountID = p.AccountID
		f.BalanceID = p.BalanceID
		f.Nickname = p.Nickname
		f.Email = p.Email
	}
	if p.Effects != nil {
		err := f.PopulateEffects(p, opType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Participant) PopulateEffects(p *history.Participant, opType xdr.OperationType) error {
	var err error
	switch opType {
	case xdr.OperationTypeDemurrage:
		f.Effects = DemurrageEffects{}
		err = p.UnmarshalEffects(&f.Effects)
	case xdr.OperationTypeManageOffer:
		f.Effects = MatchEffects{}
		err = p.UnmarshalEffects(&f.Effects)
	default:
		err = errors.New("Unexpected effects type")
	}

	return err
}
