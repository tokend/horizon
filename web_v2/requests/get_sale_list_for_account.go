package requests

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"gitlab.com/tokend/horizon/db2"
	"math"
	"net/http"
)

// GetSaleList - represents params to be specified by user for getSaleList handler
type GetSaleListForAccount struct {
	SalesBase
	Address    string
	PageParams *pgdb.CursorPageParams
}



// NewGetSaleListForAccount returns new instance of GetSaleListForAccount request
func NewGetSaleListForAccount(r *http.Request) (*GetSaleListForAccount, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleListAll,
		supportedFilters:  filterTypeSaleListAll,
	})
	if err != nil {
		return nil, err
	}

	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err}
	request:=GetSaleListForAccount{
		Address: address,
		SalesBase: SalesBase{
			base:b,
		},
		PageParams: &pgdb.CursorPageParams{},
	}

	err=urlval.Decode(r.URL.Query(),&request)
	if err != nil {
		return nil, err
	}

	if request.PageParams.Cursor > math.MaxInt64 {
		request.PageParams.Cursor,err= 0, validation.Errors{
			pageParamCursor: fmt.Errorf("cursor %d exceed max allowed %d", request.PageParams.Cursor, math.MaxInt64),
		}
	}
	if err != nil {
		return nil, err
	}
	if request.PageParams.Order == db2.OrderDescending && request.PageParams.Cursor == 0 {
		request.PageParams.Cursor = math.MaxInt64
	}

	if request.PageParams.Limit == 0 {
		request.PageParams.Limit,err= defaultLimit, nil
	}

	if request.PageParams.Limit > maxLimit {
		request.PageParams.Limit,err= 0, validation.Errors{
			pageParamLimit: fmt.Errorf("limit must not exceed %d", maxLimit),
		}
	}

	order := request.PageParams.Order
	switch order {
	case pgdb.OrderTypeAsc, pgdb.OrderTypeDesc:
		request.PageParams.Order,err= order, nil
	case "":
		request.PageParams.Order,err= pgdb.OrderTypeAsc, nil
	default:
		request.PageParams.Order,err= pgdb.OrderTypeDesc, validation.Errors{
			pageParamOrder: fmt.Errorf("allowed order types: %s, %s", pgdb.OrderTypeAsc, pgdb.OrderTypeDesc),
		}
	}
	if err != nil {
		return nil, err
	}

	return &request,nil
}
