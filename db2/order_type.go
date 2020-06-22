package db2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//OrderType - represents sorting order of the query
type OrderType string

//Invert - inverts order by
func (o OrderType) Invert() OrderType {
	switch o {
	case OrderDesc:
		return "asc"
	case OrderAsc:
		return "desc"
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": o,
		}))
	}
}

const (
	// OrderAsc - ascending order
	OrderAsc OrderType = "asc"
	// OrderDesc - descending order
	OrderDesc OrderType = "desc"
)
