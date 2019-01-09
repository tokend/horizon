package handlers

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/middleware"
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
		return errors.New("Resource is not allowed")
	}

	if err != nil {
		return errors.New("Failed to check if resource is allowed to expose")
	}

	return nil
}

func BuildResource(r *http.Request, resource Resource) (*resource.Response, error) {
	err := resource.Prepare(r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to prepare the resource")
		return nil, problems.BadRequest(err)[0]
	}

	err = CheckAllowed(resource)
	if err != nil {
		Log(r).WithError(err).Error("Failed to check if resource is allowed")
		return nil, problems.NotAllowed()
	}

	err = resource.Fetch()
	if err != nil {
		Log(r).WithError(err).Error("Failed to fetch the resource data")
		return nil, problems.InternalError()
	}

	err = resource.Populate()
	if err != nil {
		Log(r).WithError(err).Error("Failed to populate the resource")
		return nil, problems.InternalError()
	}

	response, err := resource.Response()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get the resource")
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
		return err
	}

	_, err = w.Write(js)
	if err != nil {
		return err
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

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(middleware.LogCtxKey).(*logan.Entry)
}

func RenderErr(r *http.Request, w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *jsonapi.ErrorObject:
		ape.RenderErr(w, e)
		break
	default:
		Log(r).Error(err)
		ape.RenderErr(w, problems.InternalError())
	}
}
