package rgenerated

type ManageExternalSystemAccountIdPoolEntryOpAttributes struct {
	// * 0: \"create\" * 1: \"remove\"
	Action xdr.ManageExternalSystemAccountIdPoolEntryAction `json:"action"`
	Create *CreateExternalSystemPoolOp                      `json:"create,omitempty"`
	Remove *RemoveExternalSystemPoolOp                      `json:"remove,omitempty"`
}
