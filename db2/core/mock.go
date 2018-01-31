package core

import (
	"github.com/jmoiron/sqlx"
	sq "github.com/lann/squirrel"
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/horizon/db2"
)

type CoreQMock struct {
	mock.Mock
}

func (q *CoreQMock) AccountTypeLimits() AccountTypeLimitsQI {
	args := q.Called()
	return args.Get(0).(AccountTypeLimitsQI)
}

func (q *CoreQMock) BalanceByID(dest interface{}, bid string) error {
	return nil
}

func (q *CoreQMock) Balances() BalancesQI {
	args := q.Called()
	return args.Get(0).(BalancesQI)
}

func (q *CoreQMock) BalancesByAddress(dest interface{}, addy string) error {
	return nil
}

func (q *CoreQMock) FeesByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account) ([]FeeEntry, error) {
	return nil, nil
}

func (q *CoreQMock) LimitsByAccountType(accountType int32) (*AccountTypeLimits, error) {
	return nil, nil
}

func (q *CoreQMock) ExternalSystemAccountID() ExternalSystemAccountIDQI {
	args := q.Called()
	return args.Get(0).(ExternalSystemAccountIDQI)
}

func (q *CoreQMock) LimitsForAccount(accountID string, accountType int32) (Limits, error) {
	args := q.Called()
	limits := args.Get(0).(Limits)
	return limits, args.Error(1)
}

func (q *CoreQMock) Offers() *OfferQ {
	return nil
}

func (q *CoreQMock) OrderBook() *OrderBookQ {
	return nil
}
func (q *CoreQMock) Trusts() *TrustQ {
	return nil
}

func (q *CoreQMock) LimitsByAddress(addy string) (*AccountLimits, error) {
	return nil, nil
}

func (q *CoreQMock) FeeByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account, amount int64) (*FeeEntry, error) {
	return nil, nil
}

func (q *CoreQMock) AssetPairs() AssetPairsQ {
	args := q.Called()
	return args.Get(0).(AssetPairsQ)
}

func (q *CoreQMock) Accounts() AccountQI {
	args := q.Called()
	return args.Get(0).(AccountQI)
}

func (q *CoreQMock) Assets() AssetQI {
	args := q.Called()
	return args.Get(0).(AssetQI)
}

func (q *CoreQMock) GetRepo() *db2.Repo {
	args := q.Called()
	return args.Get(0).(*db2.Repo)
}

func (q *CoreQMock) AccountByAddress(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) ExchangeName(addy string) (*string, error) {
	args := q.Called(addy)
	return args.Get(0).(*string), args.Error(1)
}

func (q *CoreQMock) LedgerHeaderBySequence(dest interface{}, seq int32) error {
	args := q.Called(dest, seq)
	return args.Error(0)
}
func (q *CoreQMock) ElderLedger(dest *int32) error {
	args := q.Called(dest)
	return args.Error(0)
}
func (q *CoreQMock) LatestLedger(dest interface{}) error {
	args := q.Called(dest)
	return args.Error(0)
}

func (q *CoreQMock) SignersByAddress(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) PoliciesByExchangeID(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) TransactionByHash(dest interface{}, hash string) error {
	args := q.Called(dest, hash)
	return args.Error(0)
}
func (q *CoreQMock) TransactionsByLedger(dest interface{}, seq int32) error {
	args := q.Called(dest, seq)
	return args.Error(0)
}

func (q *CoreQMock) TransactionFeesByLedger(dest interface{}, seq int32) error {
	args := q.Called(dest, seq)
	return args.Error(0)
}

func (q *CoreQMock) FeeEntries() FeeEntryQI {
	args := q.Called()
	return args.Get(0).(FeeEntryQI)
}

func (q *CoreQMock) Query(query sq.Sqlizer) (*sqlx.Rows, error) {
	args := q.Called(query)
	return args.Get(0).(*sqlx.Rows), args.Error(1)
}
func (q *CoreQMock) NoRows(err error) bool {
	return false
}
func (q *CoreQMock) FeeByOperationType(dest interface{}, operationType int) error {
	args := q.Called()
	return args.Error(0)
}
