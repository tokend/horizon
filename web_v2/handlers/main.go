package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/horizon/web_v2/ctx"
)

// ensureAllowed - returns false if user is not allowed to access requested data or failed to check - renders all
// corresponding error; returns true - if allowed
func isAllowed(r *http.Request, w http.ResponseWriter, dataOwners ...string) bool {
	constraints := make([]doorman.SignerConstraint, 0, len(dataOwners))
	for _, dataOwner := range dataOwners {
		constraints = append(constraints, doorman.SignerOf(dataOwner))
	}
	constraints = append(constraints, doorman.SignerOf(ctx.CoreInfo(r).AdminAccountID))

	err := ctx.Doorman(r, constraints...)
	switch err {
	case nil:
		return true
	case signcontrol.ErrNotAllowed, signcontrol.ErrNotSigned, signcontrol.ErrValidUntil, signcontrol.ErrExpired,
		signcontrol.ErrSignerKey, signcontrol.ErrSignature:

		notAllowed := problems.NotAllowed()
		notAllowed.Meta = &map[string]interface{}{
			"cause": err.Error(),
		}
		ape.RenderErr(w, notAllowed)
		return false
	default:
		ctx.Log(r).WithError(err).Error("failed to perform signature check")
		ape.RenderErr(w, problems.InternalError())
		return false
	}
}
