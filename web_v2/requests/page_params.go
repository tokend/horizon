package requests

import (
	"fmt"
	"net/url"
	"strings"

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

type OffsetBasedPageParams struct {
	limit      uint64
	pageNumber uint64
}

func newOffsetBasedPageParams(limit, pageNumber uint64) *OffsetBasedPageParams {
	return &OffsetBasedPageParams{
		limit,
		pageNumber,
	}
}

// Limit - returns us the limit we can should use for sql query
func (p *OffsetBasedPageParams) Limit() uint64 {
	if p.limit == 0 {
		return defaultLimit
	}

	return p.limit
}

// Offset - calculates the actual offset we should use for sql query
func (p *OffsetBasedPageParams) Offset() uint64 {
	return p.Limit() * p.pageNumber
}

// Links - returns pagination links we should render to the client
func (p *OffsetBasedPageParams) Links(url *url.URL) *regources.Links {
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

	return &regources.Links{
		Self: fmt.Sprintf(format, p.pageNumber, p.Limit()),
		Next: fmt.Sprintf(format, p.pageNumber+1, p.Limit()),
	}
}
