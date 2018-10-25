package jsonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"gitlab.com/tokend/horizon/context/requestid"
	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/render/problem"
	"net/http"
	"strconv"
)

type ErrorObject struct {
	// ID is a unique identifier for this particular occurrence of a problem.
	ID string `json:"id,omitempty"`

	// Title is a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence
	// of the problem, except for purposes of localization.
	Title string `json:"title,omitempty"`

	// Detail is a human-readable explanation specific to this occurrence of the problem. Like title, this field’s
	// value can be localized.
	Detail string `json:"detail,omitempty"`

	// Status is the HTTP status code applicable to this problem, expressed as a string value.
	Status string `json:"status,omitempty"`

	// Code is an application-specific error code, expressed as a string value.
	Code string `json:"code,omitempty"`

	// Meta is an object containing non-standard meta-information about the error.
	Meta map[string]interface{} `json:"meta,omitempty"`
}

var (
	errToJsonApiMap = map[error]ErrorObject{}
)

// RegisterError records an error -> ErrorObject mapping, allowing the app to register
// specific errors that may occur in other packages to be rendered as a specific
// ErrorObject instance.
func RegisterError(err error, errObj ErrorObject) {
	errToJsonApiMap[err] = errObj
}

// Inflate expands a problem with contextual information.
// At present it adds the request's id as the problem's ID, if available.
func Inflate(ctx context.Context, errObj *ErrorObject) {
	errObj.ID = requestid.FromContext(ctx)
}

func RenderErr(ctx context.Context, w http.ResponseWriter, err interface{}) {
	switch err := err.(type) {
	case problem.P:
		render(ctx, w, FromProblem(err))
	case *problem.P:
		render(ctx, w, FromProblem(*err))
	case problem.HasProblem:
		render(ctx, w, FromProblem(err.Problem()))
	case error:
		render(ctx, w, FromError(err))
	default:
		panic(fmt.Sprintf("Invalid error: %v+", err))
	}
}

func render(ctx context.Context, w http.ResponseWriter, errObj ErrorObject) {
	Inflate(ctx, &errObj)

	w.Header().Set("Content-Type", "application/vnd.api+json")

	js, err := json.MarshalIndent(errObj, "", "  ")

	if err != nil {
		err := errors.Wrap(err, 1)
		log.Ctx(ctx).WithStack(err).Error(err)
		http.Error(w, "error rendering problem", http.StatusInternalServerError)
		return
	}

	status, err := strconv.Atoi(errObj.Status)
	if err != nil {
		// TODO
		panic(fmt.Sprintf("Invalid status: %d+", errObj.Status))
	}

	w.WriteHeader(status)
	w.Write(js)
}
