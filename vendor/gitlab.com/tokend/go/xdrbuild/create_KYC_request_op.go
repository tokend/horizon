package xdrbuild

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
)

type CreateChangeRoleRequestOp struct {
	RequestID          uint64
	DestiantionAccount string
	AccountRoleToSet   uint64
	KYCData            string
	AllTasks           *uint32
}

func (op CreateChangeRoleRequestOp) Validate() error {
	return validation.ValidateStruct(&op,
		validation.Field(&op.AccountRoleToSet, validation.Required),
		validation.Field(&op.KYCData, validation.Required),
		validation.Field(&op.DestiantionAccount, validation.Required),
	)
}

func (op CreateChangeRoleRequestOp) XDR() (*xdr.Operation, error) {
	var accountToUpdateKYC xdr.AccountId
	if err := accountToUpdateKYC.SetAddress(op.DestiantionAccount); err != nil {
		return nil, errors.Wrap(err, "failed to set updated account")
	}

	var allTasksXDR xdr.Uint32
	var allTasksXDRPointer *xdr.Uint32

	if op.AllTasks != nil {
		allTasksXDR = xdr.Uint32(*op.AllTasks)
		allTasksXDRPointer = &allTasksXDR
	} else {
		allTasksXDRPointer = nil
	}

	xdrop := &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeCreateChangeRoleRequest,
			CreateChangeRoleRequestOp: &xdr.CreateChangeRoleRequestOp{
				RequestId:          xdr.Uint64(op.RequestID),
				DestinationAccount: accountToUpdateKYC,
				AccountRoleToSet:   xdr.Uint64(op.AccountRoleToSet),
				KycData:            xdr.Longstring(op.KYCData),
				AllTasks:           allTasksXDRPointer,
			},
		},
	}
	return xdrop, nil
}
