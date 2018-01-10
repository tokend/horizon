package horizon

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/charts"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/log"
)

func initCharts(app *App) {
	listener := NewMetaListener(
		log.WithField("service", "chart-listener"),
		app.HistoryQ().Transactions,
	)

	app.charts = NewCharts()

	// sales initial value
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeCreated {
			return
		}
		if change.Created.Data.Type != xdr.LedgerEntryTypeSale {
			return
		}
		data := change.Created.Data.Sale
		app.charts.Set(string(data.BaseAsset), ts, int64(data.CurrentCap))
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
		app.charts.Set(string(data.BaseAsset), ts, int64(data.CurrentCap))
	})

	// sun issued
	listener.Subscribe(func(ts time.Time, change xdr.LedgerEntryChange) {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			return
		}
		if change.Updated.Data.Type != xdr.LedgerEntryTypeAsset {
			return
		}
		data := change.Updated.Data.Asset
		if string(data.Code) != "SUN" {
			return
		}
		app.charts.Set("SUN", ts, int64(data.Issued))
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
