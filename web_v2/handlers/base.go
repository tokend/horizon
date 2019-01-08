package handlers

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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

type Allowable interface {
	IsAllowed() (bool, error)
}

func (b *Base) PrepareResource(request *http.Request, resource Resource) error {
	err := resource.Prepare(request)
	if err != nil {
		return problems.NotAllowed()
	}

	err = b.CheckAllowed(resource)
	if err != nil {
		return problems.NotAllowed()
	}

	return nil
}

func (b *Base) PrepareCollection (request *http.Request, collection Collection) error {
	err := collection.Prepare(request)
	if err != nil {
		return problems.NotAllowed()
	}

	err = b.CheckAllowed(collection)
	if err != nil {
		return problems.NotAllowed()
	}

	return nil
}

func (b *Base) CheckAllowed(resource Allowable) error {
	isAllowed, err := resource.IsAllowed()
	if !isAllowed {
		return errors.New("Resource is not allowed") // TODO: 401
	}

	if err != nil {
		return errors.New("Failed to check if resource is allowed to expose")
	}

	return nil
}

func (b *Base) RenderResource(w http.ResponseWriter, r *http.Request, id string, resource Resource) error {
	err := resource.Fetch(id)
	if err != nil {
		return problems.InternalError()
	}

	err = resource.Populate()
	if err != nil {
		return problems.InternalError()
	}

	response, err := resource.Response()
	if err != nil {
		return problems.InternalError()
	}

	js, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		return problems.InternalError()
	}

	_, err = w.Write(js)
	if err != nil {
		return problems.InternalError()
	}

	return nil
}

func (b *Base) RenderCollection(w http.ResponseWriter, r *http.Request, pp resource.PagingParams, collection Collection) error {
	return nil
}

func (b *Base) RenderErr(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *jsonapi.ErrorObject:
		ape.RenderErr(w, e)
		break
	default:
		ape.RenderErr(w, problems.InternalError())
	}
}
