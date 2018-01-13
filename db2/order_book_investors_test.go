package db2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderBookInvestors_String(t *testing.T) {
	testObj := OrderBookInvestors{
		OrderBookId: 1,
		Quantity:    42,
	}
	valid := "(1,42)"
	assert.Equal(t, valid, testObj.String())
}

func TestOrderBooksInvestors_String(t *testing.T) {
	testObj := OrderBooksInvestors{
		{OrderBookId: 1, Quantity: 42},
		{OrderBookId: 2, Quantity: 43},
		{OrderBookId: 3, Quantity: 44},
	}
	valid := "(1,42),(2,43),(3,44)"
	assert.Equal(t, valid, testObj.String())
}
