package ape

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
	"github.com/pkg/errors"
)

func RenderErr(w http.ResponseWriter, errs ...*jsonapi.ErrorObject) {
	if len(errs) == 0 {
		panic("expected non-empty errors slice")
	}

	// getting status of first occurred error
	status, err := strconv.ParseInt(errs[0].Status, 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse status"))
	}
	w.Header().Set("content-type", jsonapi.MediaType)
	w.WriteHeader(int(status))
	jsonapi.MarshalErrors(w, errs)
}

func Render(w http.ResponseWriter, res interface{}) {
	w.Header().Set("content-type", jsonapi.MediaType)
	err := jsonapi.MarshalPayload(w, res)
	if err != nil {
		panic(errors.Wrap(err, "failed to render response"))
	}
}

func RenderPage(w http.ResponseWriter, res interface{}, links *jsonapi.Links) {
	w.Header().Set("content-type", jsonapi.MediaType)

	payload, err := jsonapi.Marshal(res)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal response"))
	}

	manyPayload, ok := payload.(*jsonapi.ManyPayload)
	if !ok {
		panic(errors.Wrap(err, "failed to populate links"))
	}
	manyPayload.Links = links

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		panic(errors.Wrap(err, "failed to render response"))
	}
}
