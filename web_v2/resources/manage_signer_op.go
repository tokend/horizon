package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

func newManageSigner(id int64, details history2.ManageSignerDetails,
) *regources.ManageSignerOp {
	switch details.Action {
	case xdr.ManageSignerActionCreate:
		return &regources.ManageSignerOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_CREATE_SIGNER),
			Attributes: &regources.ManageSignerOpAttributes{
				Details:  details.CreateDetails.Details,
				Weight:   details.CreateDetails.Weight,
				Identity: details.CreateDetails.Identity,
			},
			Relationships: &regources.ManageSignerOpRelationships{
				Role:   NewSignerRoleKey(details.CreateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionUpdate:
		return &regources.ManageSignerOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_UPDATE_SIGNER),
			Attributes: &regources.ManageSignerOpAttributes{
				Details:  details.UpdateDetails.Details,
				Weight:   details.UpdateDetails.Weight,
				Identity: details.UpdateDetails.Identity,
			},
			Relationships: &regources.ManageSignerOpRelationships{
				Role:   NewSignerRoleKey(details.UpdateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionRemove:
		return &regources.ManageSignerOp{
			Key: regources.NewKeyInt64(id, regources.OPERATIONS_REMOVE_SIGNER),
			Relationships: &regources.ManageSignerOpRelationships{
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
