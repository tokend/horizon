package rgenerated

type ManageLimitsOpAttributes struct {
	// * 0: \"create\", * 1: \"remove\"
	Action xdr.ManageLimitsAction  `json:"action"`
	Create *ManageLimitsCreationOp `json:"create,omitempty"`
	Remove *ManageLimitsRemovalOp  `json:"remove,omitempty"`
}
