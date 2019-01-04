package horizon

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
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
	coreRepo := app.CoreRepo(ctx)
	log := log.WithField("service", "ingest_data_producer")
	coreRepo.Log = log
	txProvider := struct {
		*core2.LedgerHeaderQ
		*core2.TransactionQ
	}{
		LedgerHeaderQ: core2.NewLedgerHeaderQ(coreRepo),
		TransactionQ:  core2.NewTransactionQ(coreRepo),
	}
	// TODO part with current state of the ledger must be refactored
	err := app.updateLedgerState()
	if err != nil {
		panic(errors.Wrap(err, "ingest failed to update ledger state"))
	}

	ledgersChan := ingest2.NewProducer(txProvider, log).Start(ctx, 100, ledger.CurrentState())

	hRepo := app.HistoryRepo(ctx)
	accountStorage := storage.NewAccount(hRepo, coreRepo)
	balanceStorage := storage.NewBalance(hRepo, coreRepo, accountStorage)

	ledgerChangesHandler := changes.NewHandler(accountStorage, balanceStorage,
		storage.NewReviewableRequest(hRepo), storage.NewSale(hRepo))

	idProvider := struct {
		*storage.Account
		*storage.Balance
	}{
		Account: accountStorage,
		Balance: balanceStorage,
	}
	opHandler := operations.NewOperationsHandler(storage.NewOperationDetails(hRepo), storage.NewOpParticipants(hRepo), &idProvider, balanceStorage)

	consumer := ingest2.NewConsumer(log.WithField("service", "ingest_data_consumer"), hRepo, []ingest2.Handler{
		ingest2.NewLedgerHandler(storage.NewLedger(hRepo)),
		ingest2.NewTxSaver(storage.NewTx(hRepo)),
		ingest2.NewLedgerChangesHandler(storage.NewLedgerChange(hRepo)),
		ledgerChangesHandler,
		opHandler,
	}, ledgersChan)

	consumer.Start(ctx)
}

func init() {
	appInit.Add("ingester2", initIngester2, "core_connector", "app-context", "log", "horizon-db", "core-db", "stellarCoreInfo")
}
