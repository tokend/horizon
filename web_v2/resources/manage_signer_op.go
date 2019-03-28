package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

func newManageSigner(id int64, details history2.ManageSignerDetails,
) *rgenerated.ManageSignerOp {
	switch details.Action {
	case xdr.ManageSignerActionCreate:
		return &rgenerated.ManageSignerOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_CREATE_SIGNER),
			Attributes: &rgenerated.ManageSignerOpAttributes{
				Details:  details.CreateDetails.Details,
				Weight:   details.CreateDetails.Weight,
				Identity: details.CreateDetails.Identity,
			},
			Relationships: &rgenerated.ManageSignerOpRelationships{
				Role:   NewSignerRoleKey(details.CreateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionUpdate:
		return &rgenerated.ManageSignerOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_UPDATE_SIGNER),
			Attributes: &rgenerated.ManageSignerOpAttributes{
				Details:  details.UpdateDetails.Details,
				Weight:   details.UpdateDetails.Weight,
				Identity: details.UpdateDetails.Identity,
			},
			Relationships: &rgenerated.ManageSignerOpRelationships{
				Role:   NewSignerRoleKey(details.UpdateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionRemove:
		return &rgenerated.ManageSignerOp{
			Key: rgenerated.NewKeyInt64(id, rgenerated.OPERATIONS_REMOVE_SIGNER),
			Relationships: &rgenerated.ManageSignerOpRelationships{
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
