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
		// invalid account address will make doorman return 401 w/o considering other constraints
		if dataOwner == "" {
			continue
		}
		constraints = append(constraints, doorman.SignerOf(dataOwner))
	}
	constraints = append(constraints, doorman.SignerOf(ctx.CoreInfo(r).AdminAccountID))

	switch err := ctx.Doorman(r, constraints...); err.(type) {
	case nil:
		return true
	case *signcontrol.Error, *doorman.Error:
		ape.RenderErr(w, problems.NotAllowed(err))
		return false
	default:
		// while problems.NotAllowed will handle that as well,
		// there is no easy way to get that log in case of error
		ctx.Log(r).WithError(err).Error("failed to perform signature check")
		ape.RenderErr(w, problems.InternalError())
		return false
	}
}
