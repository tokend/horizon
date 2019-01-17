package requests

import (
	"net/http"
	"net/url"

	"strings"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//base - base struct describing params provided by user for the request handler
type base struct {
	include map[string]struct{}

	queryValues *url.Values
	request     *http.Request
}

// newRequest - creates new instance of request
func newBase(r *http.Request, supportedIncludes map[string]struct{}) (*base, error) {
	request := base{
		request: r,
	}
	var err error
	request.include, err = request.getIncludes(supportedIncludes)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

//shouldInclude - returns true if user requested to include resource
func (r *base) shouldInclude(name string) bool {
	_, ok := r.include[name]
	return ok
}

// getString - tries to get string from URL param, if empty gets from query values
func (r *base) getString(name string) string {
	result := chi.URLParam(r.request, name)
	if result != "" {
		return strings.TrimSpace(result)
	}

	if r.queryValues == nil {
		r.queryValues = new(url.Values)
		*r.queryValues = r.request.URL.Query()
	}

	return strings.TrimSpace(r.queryValues.Get(name))
}

func (r *base) getIncludes(supportedIncludes map[string]struct{}) (map[string]struct{}, error) {
	const fieldName = "include"
	rawIncludes := r.getString(fieldName)
	if rawIncludes == "" {
		return nil, nil
	}
	includes := strings.Split(rawIncludes, ",")
	requestIncludes := map[string]struct{}{}
	for _, include := range includes {
		if _, supported := supportedIncludes[include]; !supported {
			return nil, validation.Errors{
				fieldName: errors.New(fmt.Sprintf("`%s` is not supported; supported values`: %v", include,
					getSliceOfSupportedIncludes(supportedIncludes))),
			}
		}

		requestIncludes[include] = struct{}{}
	}

	return requestIncludes, nil

}

func getSliceOfSupportedIncludes(includes map[string]struct{}) []string {
	result := make([]string, 0, len(includes))
	for include := range includes {
		result = append(result, include)
	}

	return result
}
