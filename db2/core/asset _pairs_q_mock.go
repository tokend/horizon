package core

import "github.com/stretchr/testify/mock"

type AssetPairsQMock struct {
	mock.Mock
}

func (q *AssetPairsQMock) ByCode(base, quote string) (*AssetPair, error) {
	args := q.Called(base, quote)
	return args.Get(0).(*AssetPair), args.Error(1)
}

func (q *AssetPairsQMock) ForAssets(baseAssets, quoteAssets []string) AssetPairsQ {
	args := q.Called(baseAssets, quoteAssets)
	return args.Get(0).(AssetPairsQ)
}

func (q *AssetPairsQMock) Select() ([]AssetPair, error) {
	args := q.Called()
	return args.Get(0).([]AssetPair), args.Error(1)
}
