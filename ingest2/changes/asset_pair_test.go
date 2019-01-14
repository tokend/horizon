package changes

import (
	"testing"
	"time"

	"gitlab.com/tokend/horizon/db2/history2"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"gitlab.com/tokend/go/xdr"
)

func TestAssetPairHandler(t *testing.T) {
	entry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: 12,
			AssetPair: &xdr.AssetPairEntry{
				Base:          "TEST1",
				Quote:         "TEST2",
				CurrentPrice:  2,
				PhysicalPrice: 2,
			},
		},
	}
	ts := time.Now()
	storage := mockAssetPairStorage{}
	handler := NewHandler(nil, nil, nil, nil, &storage)

	t.Run("created", func(t *testing.T) {
		lc := ledgerChange{
			LedgerCloseTime: ts,
			LedgerChange: xdr.LedgerEntryChange{
				Type:    0,
				Created: &entry},
		}
		defer storage.AssertExpectations(t)
		storage.On("InsertAssetPair", mock.AnythingOfType("history2.AssetPair")).Run(func(arguments mock.Arguments) {
			args := arguments.Get(0).(history2.AssetPair)
			assert.Equal(t, "TEST1", args.Base)
			assert.Equal(t, "TEST2", args.Quote)
			assert.Equal(t, int64(2), args.CurrentPrice)
			assert.Equal(t, ts, args.LedgerCloseTime)
		}).Return(nil).Once()

		err := handler.handle(lc)
		assert.NoError(t, err)
	})

	t.Run("updated", func(t *testing.T) {
		lc := ledgerChange{
			LedgerCloseTime: ts,
			LedgerChange: xdr.LedgerEntryChange{
				Type:    1,
				Updated: &entry},
		}
		defer storage.AssertExpectations(t)
		storage.On("InsertAssetPair", mock.AnythingOfType("history2.AssetPair")).Run(func(arguments mock.Arguments) {
			args := arguments.Get(0).(history2.AssetPair)
			assert.Equal(t, "TEST1", args.Base)
			assert.Equal(t, "TEST2", args.Quote)
			assert.Equal(t, int64(2), args.CurrentPrice)
			assert.Equal(t, ts, args.LedgerCloseTime)
		}).Return(nil).Once()

		err := handler.handle(lc)
		assert.NoError(t, err)
	})
}
