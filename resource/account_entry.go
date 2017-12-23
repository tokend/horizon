package resource

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type AccountEntry struct {
	AccountID     string            `json:"account_id"`
	AccountTypeI  int32             `json:"account_type_i"`
	AccountType   string            `json:"account_type"`
	BlockReasonsI uint32            `json:"block_reasons_i"`
	BlockReasons  []base.Flag       `json:"block_reasons"`
	Limits        *Limits           `json:"limits"`
	Policies      AccountPolicies   `json:"policies"`
	Signers       []Signer          `json:"signers"`
	Thresholds    AccountThresholds `json:"thresholds"`
}

type LedgerKeyAccount struct {
	AccountID string `json:"account_id"`
}

func (r *AccountEntry) Populate(entry xdr.AccountEntry) {
	r.AccountID = entry.AccountId.Address()
	r.AccountTypeI = int32(entry.AccountType)
	r.AccountType = entry.AccountType.String()
	r.BlockReasonsI = uint32(entry.BlockReasons)
	r.BlockReasons = base.FlagFromXdrBlockReasons(int32(entry.BlockReasons), xdr.BlockReasonsAll)

	r.Policies.Populate(int32(entry.Policies))
	r.Thresholds.Populate(entry.Thresholds)

	if entry.Limits != nil {
		r.Limits.FromXDR(*entry.Limits)
	}
	r.Signers = make([]Signer, 0)
	for _, xSigner := range entry.Signers {
		sgn := Signer{}
		sgn.FromXDR(xSigner)
		r.Signers = append(r.Signers, sgn)
	}

}

func (r *LedgerKeyAccount) Populate(xdrAcc xdr.LedgerKeyAccount) {
	r.AccountID = xdrAcc.AccountId.Address()
}
