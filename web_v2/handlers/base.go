package handlers

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type Base struct {}

type Allowable interface {
	IsAllowed() (bool, error)
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

func (b *Base) RenderResource(w http.ResponseWriter, r *http.Request, resource Resource) error {
	err := resource.Prepare(r)
	if err != nil {
		return problems.NotAllowed()
	}

	err = b.CheckAllowed(resource)
	if err != nil {
		return problems.NotAllowed()
	}

	err = resource.Fetch()
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

func (b *Base) RenderCollection(w http.ResponseWriter, r *http.Request, collection Collection) error {
	err := collection.Prepare(r)
	if err != nil {
		return problems.NotAllowed()
	}

	err = b.CheckAllowed(collection)
	if err != nil {
		return problems.NotAllowed()
	}

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
