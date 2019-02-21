package xdrbuild

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
)

type CreateAccountRole struct {
	Details json.Marshaler
	Rules   []uint64
}

func (op *CreateAccountRole) XDR() (*xdr.Operation, error) {
	details, err := op.Details.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	rules := make([]xdr.Uint64, 0, len(op.Rules))
	for _, rule := range op.Rules {
		rules = append(rules, xdr.Uint64(rule))
	}
	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeManageAccountRole,
			ManageAccountRoleOp: &xdr.ManageAccountRoleOp{
				Data: xdr.ManageAccountRoleOpData{
					Action: xdr.ManageAccountRoleActionCreate,
					CreateData: &xdr.CreateAccountRoleData{
						Details:        xdr.Longstring(details),
						AccountRuleIDs: rules,
					},
				},
			},
		},
	}, nil
}
