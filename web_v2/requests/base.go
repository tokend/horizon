package requests

import (
	"gitlab.com/distributed_lab/figure"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//base - base struct describing params provided by user for the request handler
type base struct {
	include map[string]struct{}
	filter  map[string]string

	queryValues *url.Values
	request     *http.Request
}

type baseOpts struct {
	supportedIncludes map[string]struct{}
	supportedFilters  map[string]struct{}
}

// newRequest - creates new instance of request
func newBase(r *http.Request, opts baseOpts) (*base, error) {
	request := base{
		request: r,
	}

	err := request.unmarshalQuery(opts)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *base) unmarshalQuery(opts baseOpts) error {
	r.queryValues = new(url.Values)
	*r.queryValues = r.request.URL.Query()

	var err error
	r.filter, err = r.getFilters(opts.supportedFilters)
	if err != nil {
		return err
	}

	r.include, err = r.getIncludes(opts.supportedIncludes)
	if err != nil {
		return err
	}

	return nil
}

func (r *base) URL() *url.URL {
	return r.request.URL
}

func (r *base) marshalQuery() string {
	var builder strings.Builder

	for key, values := range *r.queryValues {
		if !strings.Contains(key, "page[") {
			builder.WriteString(key + "=" + strings.Join(values, ",") + "&")
		}
	}

	return strings.TrimSuffix(builder.String(), "&")
}

func (r *base) populateFilters(target interface{}) error {
	filter := make(map[string]interface{})
	for k, v := range r.filter {
		filter[k] = v
	}

	err := figure.Out(target).From(filter).Please()

	if err != nil {
		return err
	}

	return nil
}

func (r *base) ShouldFilter(name string) bool {
	_, ok := r.filter[name]
	return ok
}

//ShouldInclude - returns true if user requested to include resource
func (r *base) ShouldInclude(name string) bool {
	_, ok := r.include[name]
	return ok
}

//ShouldIncludeAny - returns true if user requested to include one of the provided resources
func (r *base) ShouldIncludeAny(names ...string) bool {
	for _, name := range names {
		_, ok := r.include[name]
		if ok {
			return ok
		}
	}

	return false
}

// getString - tries to get string from URL param, if empty gets from query values
func (r *base) getString(name string) string {
	result := chi.URLParam(r.request, name)
	if result != "" {
		return strings.TrimSpace(result)
	}

	return strings.TrimSpace(r.queryValues.Get(name))
}

func (r *base) getUint64(name string) (uint64, error) {
	strVal := r.getString(name)
	if strVal == "" {
		return 0, nil
	}

	return strconv.ParseUint(strVal, 0, 64)
}

func (r *base) getFilters(supportedFilters map[string]struct{}) (map[string]string, error) {
	filters := make(map[string]string)
	for queryParam, values := range *r.queryValues {
		if strings.Contains(queryParam, "filter") {
			filterKey := strings.TrimPrefix(queryParam, "filter[")
			filterKey = strings.TrimSuffix(filterKey, "]")
			if _, supported := supportedFilters[filterKey]; !supported {
				return nil, validation.Errors{
					queryParam: errors.New(
						fmt.Sprintf("filter is not supported; supported values: %v",
							getSliceOfSupportedIncludes(supportedFilters)),
					),
				}
			}

			filters[filterKey] = values[0]
		}
	}

	return filters, nil
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

func (r *base) GetOffsetBasedPageParams() (*offsetBasedPageParams, error) {
	limit, err := r.getUint64(pageParamLimit)
	if err != nil {
		return nil, err
	}

	number, err := r.getUint64(pageParamNumber)
	if err != nil {
		return nil, err
	}

	return newOffsetBasedPageParams(limit, number), nil
}

func (r *base) GetCursorBasedPageParams() (*cursorBasedPageParams, error) {
	limit, err := r.getUint64(pageParamLimit)
	if err != nil {
		return nil, validation.Errors{
			pageParamLimit: errors.New("Must be a valid uint64 value"),
		}
	}

	order := r.getString(pageParamOrder)
	if order != pageOrderAsc && order != pageOrderDesc {
		return nil, validation.Errors{
			pageParamOrder: errors.New("Must be a valid uint64 value"),
		}
	}

	cursor := r.getString(pageParamCursor)

	return newCursorBasedPageParams(limit, cursor, order), nil
}

func getSliceOfSupportedIncludes(includes map[string]struct{}) []string {
	result := make([]string, 0, len(includes))
	for include := range includes {
		result = append(result, include)
	}

	return result
}
