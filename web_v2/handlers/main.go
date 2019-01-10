package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"net/http"
	"strconv"
)

func Includes(string) bool {
	// TODO
	return true
}

func GetUint64(r *http.Request, name string) (uint64, error) {
	str := chi.URLParam(r, name)
	if len(str) == 0 {
		return 0, nil
	}

	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to get "+name+"query param")
	}

	return res, nil
}

func GetPageQuery(r *http.Request) (*db2.PageQueryV2, error) {
	limit, err := GetUint64(r, "limit")
	if err != nil {
		return nil, err
	}

	page, err := GetUint64(r, "page")
	if err != nil {
		return nil, err
	}

	pageQuery, err := db2.NewPageQueryV2(page, limit)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init page query")
	}

	return &pageQuery, nil
}


func Render(w http.ResponseWriter, resource interface{}) error {
	js, err := json.MarshalIndent(resource, "", "	")
	if err != nil {
		return errors.Wrap(err, "Failed to marshal the response")
	}

	_, err = w.Write(js)
	if err != nil {
		return errors.Wrap(err, "Failed to write the response")
	}

	return nil
}

func RenderErr(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(type) {
	case *jsonapi.ErrorObject:
		ape.RenderErr(w, e)
		break
	default:
		ctx.Log(r).Error(err)
		ape.RenderErr(w, problems.InternalError())
	}
}
