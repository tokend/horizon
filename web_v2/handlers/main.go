package handlers

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type Resource interface {
	Prepare(r *http.Request) error
	IsAllowed() (bool, error)
	Fetch() error
	Populate() error
	Response() (resource.Response, error)
}

type Collection interface {
	Prepare(r *http.Request) error
	IsAllowed() (bool, error)
	Fetch(resource.PagingParams) error
	Populate() error
	Response() ([]interface{}, error)
}

type Allowable interface {
	IsAllowed() (bool, error)
}

func CheckAllowed(resource Allowable) error {
	isAllowed, err := resource.IsAllowed()
	if !isAllowed {
		return errors.New("Resource is not allowed") // TODO: 401
	}

	if err != nil {
		return errors.New("Failed to check if resource is allowed to expose")
	}

	return nil
}

func BuildResource(r *http.Request, resource Resource) (*resource.Response, error) {
	err := resource.Prepare(r)
	if err != nil {
		return nil, problems.NotAllowed()
	}

	err = CheckAllowed(resource)
	if err != nil {
		return nil, problems.NotAllowed()
	}

	err = resource.Fetch()
	if err != nil {
		return nil, problems.InternalError()
	}

	err = resource.Populate()
	if err != nil {
		return nil, problems.InternalError()
	}

	response, err := resource.Response()
	if err != nil {
		return nil, problems.InternalError()
	}

	return &response, nil
}

func RenderResource(w http.ResponseWriter, data resource.Response) error {
	var response struct {
		Data resource.Response `json:"data"`
	}
	response.Data = data

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

func RenderCollection(w http.ResponseWriter, r *http.Request, collection Collection) error {
	err := collection.Prepare(r)
	if err != nil {
		return problems.NotAllowed()
	}

	err = CheckAllowed(collection)
	if err != nil {
		return problems.NotAllowed()
	}

	return nil
}

func RenderErr(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *jsonapi.ErrorObject:
		ape.RenderErr(w, e)
		break
	default:
		ape.RenderErr(w, problems.InternalError())
	}
}
