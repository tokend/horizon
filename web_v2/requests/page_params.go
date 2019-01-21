package requests

import (
	"fmt"
	"math"
	"strconv"

	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources/v2"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamCursor = "page[cursor]"
	pageParamOrder  = "page[order]"
)
const defaultLimit uint64 = 15

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

func (p *OffsetBasedPageParams) Limit() uint64 {
	if p.limit == 0 {
		return defaultLimit
	}

	return p.limit
}

func (p *OffsetBasedPageParams) Offset() uint64 {
	return p.Limit() * p.pageNumber
}

// TODO: accept net.URL instead of string
func (p *OffsetBasedPageParams) GetLinks(linkBase string) *regources.Links {
	format := linkBase + "&page[number]=%d&page[limit]=%d"
	return &regources.Links{
		Self: fmt.Sprintf(format, p.pageNumber, p.Limit()),
		Next: fmt.Sprintf(format, p.pageNumber+1, p.Limit()),
	}
}

type pageOrder string

const (
	PageOrderDesc = "desc"
	PageOrderAsc  = "asc"
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

func (p *cursorBasedPageParams) Limit() uint64 {
	if p.limit == 0 {
		return defaultLimit
	}

	return p.limit
}

func (p *cursorBasedPageParams) Order() string {
	if p.order == "" {
		return PageOrderAsc
	}

	return p.order
}

func (p *cursorBasedPageParams) CursorStr() string {
	return p.cursor
}

func (p *cursorBasedPageParams) CursorUInt64() (uint64, error) {
	if p.cursor == "" {
		switch p.Order() {
		case PageOrderAsc:
			return 0, nil
		case PageOrderDesc:
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

type Pageable interface {
	PagingToken() string
}

func (p *cursorBasedPageParams) GetLinks(linkBase string, records []Pageable) *jsonapi.Links {
	format := linkBase + "&page[cursor]=%d&page[limit]=%d&page[order]"
	return &jsonapi.Links{
		"self": fmt.Sprintf(format, p.cursor, p.Limit(), p.Order()),
		"next": fmt.Sprintf(format, records[len(records)-1].PagingToken(), p.Limit(), p.Order()),
	}
}
