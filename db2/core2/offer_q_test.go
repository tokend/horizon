package core2

import (
	sq "github.com/Masterminds/squirrel"
	"testing"
)

// TODO remove this file in master

func TestCountCount(t *testing.T) {
	q := sq.Select().From("offer offers").Columns("COUNT(*)")
	sql, _, _ := q.ToSql()
	if sql != "SELECT COUNT(*) FROM offer offers" {
		t.Error("Query for select count don't match")
	}
}

func TestSelectForDefaultColumns(t *testing.T) {
	defaultQ := sq.Select("offers.offer_id",
		"offers.owner_id",
		"offers.order_book_id",
		"offers.base_asset_code",
		"offers.quote_asset_code",
		"offers.base_balance_id",
		"offers.quote_balance_id",
		"offers.fee",
		"offers.is_buy",
		"offers.created_at",
		"offers.base_amount",
		"offers.quote_amount",
		"offers.price").From("offer offers")
	defQSql, _, _ := defaultQ.ToSql()

	modifiedQ := sq.Select().From("offer offers").Columns("offers.offer_id",
		"offers.owner_id",
		"offers.order_book_id",
		"offers.base_asset_code",
		"offers.quote_asset_code",
		"offers.base_balance_id",
		"offers.quote_balance_id",
		"offers.fee",
		"offers.is_buy",
		"offers.created_at",
		"offers.base_amount",
		"offers.quote_amount",
		"offers.price")
	modQSql, _, _ := modifiedQ.ToSql()

	if defQSql != modQSql {
		t.Error("Columns inserted with sq.Select() don't match to columns selected with sq.Columns()")
	}
}

func TestApplyFilters(t *testing.T) {

}
