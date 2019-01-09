package handlers

import (
	"encoding/json"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type Allowable interface {
	IsAllowed() (bool, error)
}

type ResourceBase interface {
	Allowable
	Prepare(r *http.Request) error
	Fetch() error
	Populate() error
}

type Resource interface {
	ResourceBase
	Response() (resource.Response, error)
}

type Collection interface {
	ResourceBase
	Response() ([]resource.Response, error)
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

func BuildResourceBase (r *http.Request, resource ResourceBase) error {
	err := resource.Prepare(r)
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to prepare the resource")
		return problems.BadRequest(err)[0]
	}

	err = CheckAllowed(resource)
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to check if resource is allowed")
		return problems.NotAllowed()
	}

	err = resource.Fetch()
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to fetch the resource data")
		return problems.InternalError()
	}

	err = resource.Populate()
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to populate the resource")
		return problems.InternalError()
	}

	return nil
}

func BuildResource(r *http.Request, resource Resource) (*resource.Response, error) {
	err := BuildResourceBase(r, resource)
	if err != nil {
		return nil, err
	}

	response, err := resource.Response()
	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to get the resource")
		return nil, problems.InternalError()
	}

	return &response, nil
}

func BuildCollection(r *http.Request, collection Collection) (*[]resource.Response, error) {
	err := BuildResourceBase(r, collection)
	if err != nil {
		return nil, err
	}

	response, err := collection.Response()

	if err != nil {
		ctx.Log(r).WithError(err).Error("Failed to get the resource")
		return nil, problems.InternalError()
	}

	return &response, nil
}

func RenderResource(w http.ResponseWriter, data *resource.Response) error {
	var response struct {
		Data *resource.Response `json:"data"`
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

func RenderCollection(w http.ResponseWriter, data *[]resource.Response) error {
	var response struct {
		Data *[]resource.Response `json:"data"`
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

func RenderErr(r *http.Request, w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *jsonapi.ErrorObject:
		ape.RenderErr(w, e)
		break
	default:
		ctx.Log(r).Error(err)
		ape.RenderErr(w, problems.InternalError())
	}
}
