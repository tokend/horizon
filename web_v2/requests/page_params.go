package requests

import (
	"fmt"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math"
	"net/url"
	"strconv"
	"strings"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamCursor = "page[cursor]"
	pageParamOrder  = "page[order]"
)
const defaultLimit uint64 = 15

type offsetBasedPageParams struct {
	limit      uint64
	pageNumber uint64
}

func newOffsetBasedPageParams(limit, pageNumber uint64) *offsetBasedPageParams {
	return &offsetBasedPageParams{
		limit,
		pageNumber,
	}
}

// Limit - returns us the limit we can should use for sql query
func (p *offsetBasedPageParams) Limit() uint64 {
	if p.limit == 0 {
		return defaultLimit
	}

	return p.limit
}

// Offset - calculates the actual offset we should use for sql query
func (p *offsetBasedPageParams) Offset() uint64 {
	return p.Limit() * p.pageNumber
}

// Links - returns pagination links we should render to the client
func (p *offsetBasedPageParams) Links(url *url.URL) *jsonapi.Links {
	var query strings.Builder

	for key, values := range url.Query() {
		switch key {
		case pageParamNumber, pageParamLimit:
			continue
		default:
			query.WriteString(fmt.Sprintf("%s=%s&", key, strings.Join(values, ",")))
		}
	}

	format := url.Path + "?" + query.String() + "&page[number]=%d&page[limit]=%d"

	return &jsonapi.Links{
		"self": fmt.Sprintf(format, p.pageNumber, p.Limit()),
		"next": fmt.Sprintf(format, p.pageNumber+1, p.Limit()),
	}
}

type pageOrder string

const (
	pageOrderDesc = "desc"
	pageOrderAsc  = "asc"
)

type cursorBasedPageParams struct {
	// with string cursor we can properly iterate when `order=desc&cursor=`
	cursor string
	order  string
	limit  uint64
}

func newCursorBasedPageParams(limit uint64, cursor, order string) *cursorBasedPageParams {
	return &cursorBasedPageParams{
		cursor,
		order,
		limit,
	}
}

// Limit - returns the limit we can should use for sql query
func (p *cursorBasedPageParams) Limit() uint64 {
	if p.limit == 0 {
		return defaultLimit
	}

	return p.limit
}

// Order - returns the wished order of response records
func (p *cursorBasedPageParams) Order() string {
	if p.order == "" {
		return pageOrderAsc
	}

	return p.order
}

// CursorStr - returns cursor as string
func (p *cursorBasedPageParams) CursorStr() string {
	return p.cursor
}

// CursorUInt64 - returns cursor as uint64
func (p *cursorBasedPageParams) CursorUInt64() (uint64, error) {
	if p.cursor == "" {
		switch p.Order() {
		case pageOrderAsc:
			return 0, nil
		case pageOrderDesc:
			return math.MaxInt64, nil
		default:
			return 0, errors.New("Invalid order")
		}
	}

	i, err := strconv.ParseUint(p.cursor, 10, 64)

	if err != nil {
		return 0, errors.New("Invalid cursor")
	}

	if i < 0 {
		return 0, errors.New("Invalid cursor")
	}

	return i, nil
}

type pageable interface {
	PagingToken() string
}

// Links - returns pagination links we should render to the client
func (p *cursorBasedPageParams) Links(linkBase string, records []pageable) *jsonapi.Links {
	format := linkBase + "&page[cursor]=%d&page[limit]=%d&page[order]"
	return &jsonapi.Links{
		"self": fmt.Sprintf(format, p.cursor, p.Limit(), p.Order()),
		"next": fmt.Sprintf(format, records[len(records)-1].PagingToken(), p.Limit(), p.Order()),
	}
}
