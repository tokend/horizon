package resource

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/regources"
)

type AccountEntry struct {
	AccountID     string            `json:"account_id"`
	AccountTypeI  int32             `json:"account_type_i"`
	AccountType   string            `json:"account_type"`
	BlockReasonsI uint32            `json:"block_reasons_i"`
	BlockReasons  []regources.Flag  `json:"block_reasons"`
	LimitsV2      []LimitsV2        `json:"limits"`
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
