package resources

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageSigner(id int64, details history2.ManageSignerDetails,
) *regources.ManageSigner {
	switch details.Action {
	case xdr.ManageSignerActionCreate:
		return &regources.ManageSigner{
			Key: regources.NewKeyInt64(id, regources.TypeCreateSigner),
			Attributes: &regources.ManageSignerAttrs{
				Details:  details.CreateDetails.Details,
				Weight:   details.CreateDetails.Weight,
				Identity: details.CreateDetails.Identity,
			},
			Relationships: &regources.ManageSignerRelation{
				Role:   NewSignerRoleKey(details.CreateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionUpdate:
		return &regources.ManageSigner{
			Key: regources.NewKeyInt64(id, regources.TypeUpdateSigner),
			Attributes: &regources.ManageSignerAttrs{
				Details:  details.UpdateDetails.Details,
				Weight:   details.UpdateDetails.Weight,
				Identity: details.UpdateDetails.Identity,
			},
			Relationships: &regources.ManageSignerRelation{
				Role:   NewSignerRoleKey(details.UpdateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionRemove:
		return &regources.ManageSigner{
			Key: regources.NewKeyInt64(id, regources.TypeRemoveSigner),
			Relationships: &regources.ManageSignerRelation{
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
