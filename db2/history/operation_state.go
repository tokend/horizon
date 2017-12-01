package history

type OperationState int32

const (
	OperationStatePending OperationState = 1 + iota
	OperationStateSuccess
	OperationStateRejected
	OperationStateCanceled
	OperationStateFailed
)

var operationStateStr = map[OperationState]string{
	OperationStatePending: "pending",
	OperationStateSuccess: "success",
	OperationStateRejected: "rejected",
	OperationStateCanceled: "canceled",
	OperationStateFailed: "failed",
}

func (s OperationState) String() string {
	return operationStateStr[s]
}
