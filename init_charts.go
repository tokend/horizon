package horizon

import (
	"time"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/charts"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/exchange"
	"gitlab.com/swarmfund/horizon/log"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

func initCharts(app *App) {
	listener := NewMetaListener(
		log.WithField("service", "chart-listener"),
		app.HistoryQ().Transactions,
	)

	app.charts = NewCharts()

	// asset price initial value
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeCreated {
			return
		}
		if change.Created.Data.Type != xdr.LedgerEntryTypeAssetPair {
			return
		}
		data := change.Created.Data.AssetPair
		app.charts.Set(fmt.Sprintf("%s-%s", data.Base, data.Quote), ts, int64(data.CurrentPrice))
	})

	// asset prices
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			return
		}
		if change.Updated.Data.Type != xdr.LedgerEntryTypeAssetPair {
			return
		}
		data := change.Updated.Data.AssetPair
		app.charts.Set(fmt.Sprintf("%s-%s", data.Base, data.Quote), ts, int64(data.CurrentPrice))
	})

	// sales initial value
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeCreated {
			return
		}
		if change.Created.Data.Type != xdr.LedgerEntryTypeSale {
			return
		}
		data := change.Created.Data.Sale
		app.charts.Set(string(data.BaseAsset), ts, int64(0))
	})

	// sales current cap charts
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			return
		}
		if change.Updated.Data.Type != xdr.LedgerEntryTypeSale {
			return
		}
		data := change.Updated.Data.Sale
		// TODO: fix me
		logger := log.WithField("listener", "sales_current_cap")
		converter, err := exchange.NewConverter(app.CoreQ())
		if err != nil {
			logger.WithError(err).Error("Failed to init converter")
			return
		}

		cupInQuote, err := getCurrentCapInDefaultQuoteForEntry(*data, converter)
		if err != nil {
			logger.WithError(err).Error("Failed to convert")
			return
		}

		app.charts.Set(string(data.BaseAsset), ts, cupInQuote)
	})

	// sun issued
	prevIssued := map[string]*int64{
		"ETH": nil,
		"BTC": nil,
	}
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			return
		}
		if change.Updated.Data.Type != xdr.LedgerEntryTypeAsset {
			return
		}
		data := change.Updated.Data.Asset
		code := string(data.Code)
		if _, ok := prevIssued[code]; !ok {
			return
		}

		issued := int64(data.Issued)
		prevIssued[code] = &issued
		logger := log.WithField("listener", "sales_current_cap")
		converter, err := exchange.NewConverter(app.CoreQ())
		if err != nil {
			logger.WithError(err).Error("Failed to init converter")
			return
		}

		totalIssued, err := convertMap(prevIssued, "SUN", converter)
		if err != nil {
			logger.WithError(err).Error("Failed to convert map of ETH BTC to SUN")
			return
		}

		app.charts.Set("SUN", ts, totalIssued)
	})

	if err := listener.Init(); err != nil {
		panic(errors.Wrap(err, "failed to init chart listener"))
	}

	go listener.Run()
}

type Charts struct {
	histograms map[string]map[string]*charts.Histogram
}

func NewCharts() *Charts {
	return &Charts{
		histograms: make(map[string]map[string]*charts.Histogram),
	}
}

func (c *Charts) Get(name string) map[string]*charts.Histogram {
	return c.histograms[name]
}

func (c *Charts) Set(name string, ts time.Time, value int64) {
	histograms, ok := c.histograms[name]
	if !ok {
		histograms = make(map[string]*charts.Histogram)
		histograms["hour"] = charts.NewHistogram(1*time.Hour, 360)
		histograms["day"] = charts.NewHistogram(24*time.Hour, 360)
		histograms["month"] = charts.NewHistogram(30*24*time.Hour, 360)
		histograms["year"] = charts.NewHistogram(365*24*time.Hour, 360)
		c.histograms[name] = histograms
	}

	for _, histogram := range histograms {
		histogram.Run(value, ts)
	}
}

type Subscriber func(time.Time, xdr.LedgerEntryChange)

type MetaListener struct {
	log         *log.Entry
	cursor      db2.PageQuery
	txq         func() history.TransactionsQI
	subscribers []Subscriber
}

func NewMetaListener(log *log.Entry, txq func() history.TransactionsQI) *MetaListener {
	return &MetaListener{
		log: log,
		txq: txq,
		cursor: db2.PageQuery{
			Cursor: "",
			Order:  "asc",
			Limit:  200,
		},
	}
}

func (l *MetaListener) Subscribe(subscriber Subscriber) {
	l.subscribers = append(l.subscribers, subscriber)
}

func (l *MetaListener) Init() error {
	for {
		var records []history.Transaction
		err := l.txq().Page(l.cursor).Select(&records)
		if err != nil {
			return errors.Wrap(err, "failed to get records")
		}
		for _, record := range records {
			if err := l.processRecord(record); err != nil {
				return errors.Wrap(err, "failed to process tx")
			}
			l.cursor.Cursor = record.PagingToken()
		}

		if len(records) == 0 {
			return nil
		}
	}
}

func (l *MetaListener) processRecord(tx history.Transaction) error {
	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshalBase64(tx.TxMeta, &meta); err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}

	for _, changes := range meta.MustOperations() {
		for _, change := range changes.Changes {
			for _, subscriber := range l.subscribers {
				subscriber(tx.LedgerCloseTime, change)
			}
		}
	}

	return nil
}

func (l *MetaListener) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		if rvr := recover(); rvr != nil {
			l.log.WithError(errors.FromPanic(rvr)).Error("panicked")
		}
		ticker.Stop()
	}()

	for ; ; <-ticker.C {
		var records []history.Transaction
		err := l.txq().Page(l.cursor).Select(&records)
		if err != nil {
			l.log.WithError(err).Error("failed to get records")
			continue
		}
		for _, record := range records {
			if err := l.processRecord(record); err != nil {
				panic(errors.Wrap(err, "failed to process tx"))
			}
			l.cursor.Cursor = record.PagingToken()
		}
	}
}

func init() {
	appInit.Add("swarm_horizon", initCharts, "horizon-db")
}

func convertMap(data map[string]*int64, destAsset string, converter *exchange.Converter) (int64, error) {
	var result int64
	for key, value := range data {
		if value == nil {
			continue
		}

		converted, err := converter.TryToConvertWithOneHop(*value, key, destAsset)
		if err != nil || converted == nil {
			if err == nil {
				err = errors.New("failed to find path to convert asset")
			}
			return 0, errors.Wrap(err, "failed to convert asset")
		}

		var ok bool
		result, ok = amount.SafePositiveSum(result, *converted)
		if !ok {
			return 0, errors.New("overflow on conversion")
		}
	}

	return result, nil
}

func getCurrentCapInDefaultQuoteForEntry(sale xdr.SaleEntry, converter *exchange.Converter) (int64, error) {
	totalCapInDefaultQuoteAsset := int64(0)
	for _, quoteAsset := range sale.QuoteAssets {
		currentCapInDefaultQuoteAsset, err := converter.TryToConvertWithOneHop(int64(quoteAsset.CurrentCap),
			string(quoteAsset.QuoteAsset), string(sale.DefaultQuoteAsset))
		if err != nil {
			return 0, errors.Wrap(err, "failed to convert current cap to default quote asset")
		}

		if currentCapInDefaultQuoteAsset == nil {
			return 0, errors.New("failed to convert to current cap: no path found")
		}

		var isOk bool
		totalCapInDefaultQuoteAsset, isOk = amount.SafePositiveSum(totalCapInDefaultQuoteAsset, *currentCapInDefaultQuoteAsset)
		if !isOk {
			return 0, errors.New("failed to find total cap in default quote asset: overflow")
		}
	}

	return totalCapInDefaultQuoteAsset, nil
}
