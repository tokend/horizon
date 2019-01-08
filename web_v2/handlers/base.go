package handlers

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"net/http"
)

type Base struct {
	CoreQ    core.QInterface
	HistoryQ history.QInterface
}

func (b *Base) CheckAllowed(request *http.Request, resource Resource) error {
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

func (b *Base) RenderResource(w http.ResponseWriter, r *http.Request, resource Resource) {
	err := resource.Fetch()
	if err != nil {
		b.RenderErr()
		return
	}

	err = resource.PopulateAttributes()
	if err != nil {
		b.RenderErr()
		return
	}

	response, err := resource.Response()
	if err != nil {
		b.RenderErr()
		return
	}

	js, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		b.RenderErr()
		return
	}

	_, err = w.Write(js)
	if err != nil {
		b.RenderErr()
		return
	}
}

func (b *Base) RenderCollection(w http.ResponseWriter, r *http.Request, collection Resource) {
	b.RenderResource(w, r, collection)
}

func (b *Base) RenderErr() {
	// TODO
}
