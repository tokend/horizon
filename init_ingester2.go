package horizon

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2"
	"gitlab.com/tokend/horizon/ingest2/changes"
	"gitlab.com/tokend/horizon/ingest2/operations"
	"gitlab.com/tokend/horizon/ingest2/storage"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/log"
)

func initIngester2(app *App) {
	if !app.config.Ingest {
		return
	}

	ctx := app.ctx
	logger := log.DefaultLogger.Entry.WithField("service", "ingest_v2")
	coreRepo := app.CoreRepoLogged(logger)

	txProvider := struct {
		*core2.LedgerHeaderQ
		*core2.TransactionQ
	}{
		LedgerHeaderQ: core2.NewLedgerHeaderQ(coreRepo),
		TransactionQ:  core2.NewTransactionQ(coreRepo),
	}

	hRepo := app.HistoryRepoLogged(logger)
	ledgersChan := ingest2.NewProducer(txProvider, history2.NewLedgerQ(hRepo), logger).Start(ctx, 1000, ledger.CurrentState())

	accountStorage := storage.NewAccount(hRepo, coreRepo)
	balanceStorage := storage.NewBalance(hRepo, coreRepo, accountStorage)

	ledgerChangesHandler := changes.NewHandler(
		accountStorage,
		balanceStorage,
		storage.NewReviewableRequest(hRepo),
		storage.NewSale(hRepo),
		storage.NewAssertPair(hRepo),
		storage.NewPoll(hRepo),
		storage.NewVote(hRepo),
		storage.NewAccountSpecificRules(hRepo),
		storage.NewSaleParticipation(hRepo),
	)

	idProvider := struct {
		*storage.Account
		*storage.Balance
	}{
		Account: accountStorage,
		Balance: balanceStorage,
	}
	opHandler := operations.NewOperationsHandler(storage.NewOperationDetails(hRepo), storage.NewOpParticipants(hRepo), &idProvider, balanceStorage)

	consumer := ingest2.NewConsumer(logger.WithField("service", "ingest_data_consumer"), hRepo, app.CoreConnector, []ingest2.Handler{
		ingest2.NewLedgerHandler(storage.NewLedger(hRepo)),
		ingest2.NewTxSaver(storage.NewTx(hRepo)),
		ingest2.NewLedgerChangesHandler(storage.NewLedgerChange(hRepo)),
		ledgerChangesHandler,
		opHandler,
	}, ledgersChan)

	consumer.Start(ctx)
}

func init() {
	appInit.Add("ingester2", initIngester2, "app-context", "log", "horizon-db", "core-db", "core_connector", "core-info", "ledger-state")
}
