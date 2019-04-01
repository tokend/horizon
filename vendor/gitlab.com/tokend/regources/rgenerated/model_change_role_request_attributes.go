package rgenerated

type ChangeRoleRequestAttributes struct {
	AccountRoleToSet uint64  `json:"account_role_to_set"`
	CreatorDetails   Details `json:"creator_details"`
	// Sequence number
	SequenceNumber uint32 `json:"sequence_number"`
}
