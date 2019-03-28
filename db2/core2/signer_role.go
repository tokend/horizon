package core2

import "gitlab.com/tokend/regources/rgenerated"

//SignerRole - represents role of the signer
type SignerRole struct {
	ID      uint64             `db:"id"`
	RuleIDs []uint64           `db:"rule_ids"`
	OwnerID string             `db:"owner_id"`
	Details rgenerated.Details `db:"details"`
}
