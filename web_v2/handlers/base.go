package handlers

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type Base struct {
	CoreQ    core.QInterface
	HistoryQ history.QInterface
}

func (b *Base) CheckAllowed(request *http.Request, resource resource.Resource) error {
	err := resource.FindOwner()
	if err != nil {
		return errors.Wrap(err, "Failed to define the owner of data") // TODO: 401
	}

	isAllowed, err := resource.IsAllowed()
	if !isAllowed {
		return errors.New("Resource is not allowed") // TODO: 401
	}

	if err != nil {
		return errors.New("Failed to check if resource is allowed to expose")
	}

	return nil
}
