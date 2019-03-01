package requests

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/distributed_lab/figure"

	"fmt"

	"math"

	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/regources/v2"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamCursor = "page[cursor]"
	pageParamOrder  = "page[order]"
)

const defaultLimit uint64 = 15
const maxLimit uint64 = 100

//base - base struct describing params provided by user for the request handler
type base struct {
	include           map[string]struct{}
	supportedIncludes map[string]struct{}
	filter            map[string]string

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
		request:           r,
		supportedIncludes: opts.supportedIncludes,
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

var amountHook = figure.Hooks{
	"regources.Amount": func(value interface{}) (reflect.Value, error) {
		strVal, ok := value.(string)
		if !ok {
			return reflect.Value{}, errors.New("Failed to parse value as string")
		}

		intVal, err := amount.Parse(strVal)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to parse value as int64")
		}

		result := regources.Amount(intVal)

		return reflect.ValueOf(result), nil
	},
}

func (r *base) populateFilters(target interface{}) error {
	filter := make(map[string]interface{})
	for k, v := range r.filter {
		filter[k] = v
	}

	err := figure.Out(target).With(figure.BaseHooks, amountHook).From(filter).Please()

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
// panics if tries to check not supported include
func (r *base) ShouldInclude(name string) bool {
	if _, ok := r.supportedIncludes[name]; !ok {
		panic(errors.From(errors.New("unexpected include check of the request"), logan.F{
			"supported_includes": getSliceOfSupportedIncludes(r.supportedIncludes),
			"checking_include":   name,
		}))
	}
	_, ok := r.include[name]
	return ok
}

//ShouldIncludeAny - returns true if user requested to include one of the provided resources
func (r *base) ShouldIncludeAny(names ...string) bool {
	for _, name := range names {
		if r.ShouldInclude(name) {
			return true
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

func (r *base) getUint64ID() (uint64, error) {
	result, err := r.getUint64("id")
	if err != nil {
		return 0, validation.Errors{
			"id": err,
		}
	}

	return result, nil
}

func (r *base) getInt64(name string) (int64, error) {
	strVal := r.getString(name)
	if strVal == "" {
		return 0, nil
	}

	return strconv.ParseInt(strVal, 0, 64)
}

func (r *base) getInt32(name string) (int32, error) {
	strVal := r.getString(name)
	if strVal == "" {
		return 0, nil
	}

	raw, err := strconv.ParseInt(strVal, 0, 32)
	if err != nil {
		return 0, errors.Wrap(err, "overflow during int32 parsing")
	}

	return int32(raw), nil
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

			if len(values) == 0 {
				continue
			}

			if len(values) > 1 {
				return nil, validation.Errors{
					queryParam: errors.New("multiple values per one filter are not supported"),
				}
			}

			if values[0] == "" {
				continue
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

		// Note: Because compound documents require full linkage (except when relationship linkage is excluded by sparse fieldsets),
		// intermediate resources in a multi-part path must be returned along with the leaf nodes. For example,
		// a response to a request for comments.author should include comments as well as the author of each of those comments.
		subIncludes := strings.Split(include, ".")
		parentInclude := ""
		for i := range subIncludes {
			if parentInclude == "" {
				parentInclude = subIncludes[i]
			} else {
				parentInclude += "." + subIncludes[i]
			}
			requestIncludes[parentInclude] = struct{}{}
		}
	}

	return requestIncludes, nil
}

func (r *base) getOffsetBasedPageParams() (*db2.OffsetPageParams, error) {
	limit, err := r.getLimit(defaultLimit, maxLimit)
	if err != nil {
		return nil, err
	}

	pageNumber, err := r.getPageNumber()
	if err != nil {
		return nil, err
	}

	order, err := r.getOrder()
	if err != nil {
		return nil, err
	}

	return &db2.OffsetPageParams{
		Order:      order,
		Limit:      limit,
		PageNumber: pageNumber,
	}, nil
}

func (r *base) getCursorBasedPageParams() (*db2.CursorPageParams, error) {
	limit, err := r.getLimit(defaultLimit, maxLimit)
	if err != nil {
		return nil, err
	}

	order, err := r.getOrder()
	if err != nil {
		return nil, err
	}

	cursor, err := r.getCursor()
	if err != nil {
		return nil, err
	}

	if order == db2.OrderDescending && cursor == 0 {
		cursor = math.MaxInt64
	}

	return &db2.CursorPageParams{
		Limit:  limit,
		Order:  order,
		Cursor: cursor,
	}, nil
}

//GetCursorLinks - returns links for cursor based page params
func (r *base) GetCursorLinks(p db2.CursorPageParams, last string) *regources.Links {
	result := regources.Links{
		Self: r.getCursorLink(p.Cursor, p.Limit, p.Order),
		Prev: r.getCursorLink(p.Cursor, p.Limit, p.Order.Invert()),
	}

	if last != "" {
		lastI, err := strconv.ParseUint(last, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse cursor", logan.F{
				"last": last,
			}))
		}
		result.Next = r.getCursorLink(lastI, p.Limit, p.Order)
	}

	return &result
}

//GetOffsetLinks - returns links for offset based page params
func (r *base) GetOffsetLinks(p db2.OffsetPageParams) *regources.Links {
	result := regources.Links{
		Next: r.getOffsetLink(p.PageNumber+1, p.Limit, p.Order),
		Self: r.getOffsetLink(p.PageNumber, p.Limit, p.Order),
	}

	return &result
}

func (r *base) getCursorLink(cursor, limit uint64, order db2.OrderType) string {
	u := r.URL()
	query := u.Query()
	query.Set(pageParamCursor, strconv.FormatUint(cursor, 10))
	query.Set(pageParamLimit, strconv.FormatUint(limit, 10))
	query.Set(pageParamOrder, string(order))
	u.RawQuery = query.Encode()
	return u.String()
}

func (r *base) getOffsetLink(pageNumber, limit uint64, order db2.OrderType) string {
	u := r.URL()
	query := u.Query()
	query.Set(pageParamNumber, strconv.FormatUint(pageNumber, 10))
	query.Set(pageParamLimit, strconv.FormatUint(limit, 10))
	query.Set(pageParamOrder, string(order))
	u.RawQuery = query.Encode()
	return u.String()
}

func getSliceOfSupportedIncludes(includes map[string]struct{}) []string {
	result := make([]string, 0, len(includes))
	for include := range includes {
		result = append(result, include)
	}

	return result
}

func (r *base) getLimit(defaultLimit, maxLimit uint64) (uint64, error) {
	result, err := r.getUint64(pageParamLimit)
	if err != nil {
		return 0, validation.Errors{
			pageParamLimit: errors.New("Must be a valid uint64 value"),
		}
	}

	if result == 0 {
		return defaultLimit, nil
	}

	if result > maxLimit {
		return 0, validation.Errors{
			pageParamLimit: fmt.Errorf("limit must not exceed %d", maxLimit),
		}
	}

	return result, nil
}

func (r *base) getOrder() (db2.OrderType, error) {
	order := r.getString(pageParamOrder)
	switch order {
	case db2.OrderAscending, db2.OrderDescending:
		return db2.OrderType(order), nil
	case "":
		return db2.OrderAscending, nil
	default:
		return db2.OrderDescending, validation.Errors{
			pageParamOrder: fmt.Errorf("allowed order types: %s, %s", db2.OrderAscending, db2.OrderDescending),
		}
	}
}

func (r *base) getPageNumber() (uint64, error) {
	result, err := r.getUint64(pageParamNumber)
	if err != nil {
		return 0, validation.Errors{
			pageParamNumber: err,
		}
	}

	return result, nil
}

func (r *base) getCursor() (uint64, error) {
	result, err := r.getUint64(pageParamCursor)
	if err != nil {
		return 0, validation.Errors{
			pageParamCursor: err,
		}
	}

	if result > math.MaxInt64 {
		return 0, validation.Errors{
			pageParamCursor: fmt.Errorf("cursor %d exceed max allowed %d", result, math.MaxInt64),
		}
	}

	return result, nil
}
