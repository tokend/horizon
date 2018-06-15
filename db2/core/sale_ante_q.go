package core

import sq "github.com/lann/squirrel"

var _ SaleAnteQI = &SaleAnteQ{}

type SaleAnteQI interface {
	// filters by sale id
	ForSale(saleID string) SaleAnteQI
	// filters by balance id
	ForBalance(balanceID string) SaleAnteQI
	// Select loads the results of the query specified by `q`
	Select() ([]SaleAnte, error)
}

type SaleAnteQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) SaleAntes() *SaleAnteQ {
	return &SaleAnteQ{
		parent: q,
		sql:    selectSaleAnte,
	}
}

func (q *SaleAnteQ) ForSale(saleID string) SaleAnteQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("sa.sale_id = ?", saleID)
	return q
}

func (q *SaleAnteQ) ForBalance(balanceID string) SaleAnteQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("sa.participant_balance_id = ?", balanceID)
	return q
}

func (q *SaleAnteQ) Select() ([]SaleAnte, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	result := make([]SaleAnte, 0)
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

var selectSaleAnte = sq.Select(
	"sa.sale_id",
	"sa.participant_balance_id",
	"sa.amount",
).From("sale_ante sa")
