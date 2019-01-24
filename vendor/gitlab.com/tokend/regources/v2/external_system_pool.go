package regources

import "gitlab.com/tokend/go/xdr"

//ManageExternalSystemPool - details of corresponding op
type ManageExternalSystemPool struct {
	Key
	Attributes ManageExternalSystemPoolAttrs `json:"attributes"`
}

//ManageExternalSystemPoolAttrs - details of corresponding op
type ManageExternalSystemPoolAttrs struct {
	Action xdr.ManageExternalSystemAccountIdPoolEntryAction `json:"action"`
	Create *CreateExternalSystemPool                        `json:"create"`
	Remove *RemoveExternalSystemPool                        `json:"remove"`
}

//CreateExternalSystemPool - details of corresponding op
type CreateExternalSystemPool struct {
	PoolID             uint64 `json:"pool_id"`
	Data               string `json:"data"`
	Parent             uint64 `json:"parent"`
	ExternalSystemType int32  `json:"external_system_type"`
}

//RemoveExternalSystemPool - details of corresponding op
type RemoveExternalSystemPool struct {
	PoolID uint64 `json:"pool_id"`
}

//BindExternalSystemAccountAttrs - details of corresponding op
type BindExternalSystemAccount struct {
	Key
	Attributes BindExternalSystemAccountAttrs `json:"attributes"`
}

//BindExternalSystemAccountAttrs - details of corresponding op
type BindExternalSystemAccountAttrs struct {
	ExternalSystemType int32 `json:"external_system_type"`
}
