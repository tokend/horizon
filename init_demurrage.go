package horizon

import (
	"time"

	"math/rand"

	"encoding/hex"

	"bullioncoin.githost.io/development/go/keypair"
	"bullioncoin.githost.io/development/go/network"
	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/codes"
	"gitlab.com/distributed_lab/tokend/horizon/corer"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/errors"
	"gitlab.com/distributed_lab/tokend/horizon/log"
	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/txsub"
	"golang.org/x/net/context"
)

type DemurrageManager struct {
	coreDB            core.QInterface
	demurrageOperator string
	randomGenerator   *rand.Rand
	coreInfo          *corer.Info
	submitter         *txsub.System

	masterAccountID xdr.AccountId

	Log *log.Entry
	ctx context.Context
}

func NewDemurrageManager(ctx context.Context, demurrageOperator string, coreInfo *corer.Info, coreDB core.QInterface, submitter *txsub.System) *DemurrageManager {
	manager := DemurrageManager{
		coreInfo:          coreInfo,
		demurrageOperator: demurrageOperator,
		coreDB:            coreDB,
		submitter:         submitter,
		randomGenerator:   rand.New(rand.NewSource(time.Now().UnixNano())),
		ctx:               ctx,
		Log:               log.WithField("service", "demurrage_manager"),
	}

	err := manager.masterAccountID.SetAddress(coreInfo.MasterAccountID)
	if err != nil {
		manager.Log.WithError(err).Panic("Failed to set master address")
	}

	return &manager
}

func (m *DemurrageManager) Run() {
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		// have to stop
		case <-m.ctx.Done():
			return
		case <-ticker:
			// do runOnce
		}

		err := m.runOnce()
		if err != nil {
			m.Log.WithError(err).Error("Failed to run demurrage")
		}
	}
}

func (m *DemurrageManager) haveToStop() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}

func (m *DemurrageManager) runOnce() error {
	if m.haveToStop() {
		return nil
	}

	defer func() {
		if rec := recover(); rec != nil {
			err := errors.FromPanic(rec)
			m.Log.WithError(err).Error("Demurrage panicked")
		}
	}()

	assets, err := m.coreDB.Assets()
	if err != nil {
		return logan.Wrap(err, "Failed to get assets")
	}
	for _, asset := range assets {
		if m.haveToStop() {
			return nil
		}

		err = m.demurrageForAsset(asset)
		if err != nil {
			return logan.Wrap(err, "Failed to run demurrage for asset").WithField("asset_code", asset.Code)
		}
	}

	return nil
}

func (m *DemurrageManager) submit(ctx context.Context, operation xdr.Operation) (*txsub.Result, error) {

	currentTime := time.Now().Unix()
	salt := m.randomGenerator.Int63()
	transaction := xdr.Transaction{
		Salt:       xdr.Salt(salt),
		Operations: []xdr.Operation{operation},
		TimeBounds: xdr.TimeBounds{
			MaxTime: xdr.Uint64(currentTime + m.coreInfo.TxExpirationPeriod),
		},
	}
	transaction.SourceAccount.SetAddress(m.coreInfo.MasterAccountID)

	hash, err := network.HashTransaction(&transaction, m.coreInfo.NetworkPassphrase)
	if err != nil {
		return nil, logan.Wrap(err, "Failed to get hash for tx")
	}

	kp, err := keypair.Parse(m.demurrageOperator)
	if err != nil {
		return nil, logan.Wrap(err, "Failed to parse demurrage operator kp")
	}

	signature, err := kp.SignDecorated(hash[:])
	if err != nil {
		return nil, logan.Wrap(err, "Failed to sign tx")
	}

	envelope := xdr.TransactionEnvelope{
		Tx:         transaction,
		Signatures: []xdr.DecoratedSignature{signature},
	}
	env, err := xdr.MarshalBase64(envelope)
	if err != nil {
		return nil, logan.Wrap(err, "Failed to marshal envelope")
	}

	txResult := m.submitter.Submit(ctx, &txsub.EnvelopeInfo{
		ContentHash:   hex.EncodeToString(hash[:]),
		SourceAddress: m.coreInfo.MasterAccountID,
		RawBlob:       env,
	})
	if txResult.HasInternalError() {
		return nil, logan.Wrap(txResult.Err, "Failed to submit tx")
	}

	return &txResult, nil
}

func (m *DemurrageManager) demurrageForAsset(asset core.Asset) error {
	demurrageOp := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeDemurrage,
			DemurrageOp: &xdr.DemurrageOp{
				Asset: xdr.AssetCode(asset.Code),
			},
		},
		SourceAccount: &m.masterAccountID,
	}
	ctx, _ := context.WithTimeout(m.ctx, 30*time.Second)
	result, err := m.submit(ctx, demurrageOp)
	if err != nil {
		return logan.Wrap(err, "Failed to submit demurrage op")
	}

	if result.Err == nil {
		return nil
	}

	txError, ok := result.Err.(txsub.Error)
	if !ok {
		return logan.Wrap(result.Err, "Expected error to be txsub.Error type")
	}

	switch txError.Type() {
	case txsub.Timeout:
		return logan.Wrap(txError, "Failed to submit demurrage op")
	case txsub.RejectedTx:
		var txResult xdr.TransactionResult
		err = xdr.SafeUnmarshalBase64(txError.ResultXDR(), &txResult)
		if err != nil {
			return logan.Wrap(err, "Failed to get parser tx result")
		}

		txResultCode, opResultCodes, err := codes.ForTxResult(txResult)
		if err != nil {
			return logan.Wrap(err, "Failed to convert tx result to codes")
		}

		if len(opResultCodes) != 1 {
			return logan.NewError("Expected opResultCodes to have 1 elem").WithField("txResultCode", txResultCode).WithField("opResultCodes", opResultCodes)
		}

		if opResultCodes[0] == codes.OpDemurrageNotRequired {
			return nil
		}

		return logan.NewError("Failed to submit demurrage op").WithField("txResultCode", txResultCode).WithField("opResultCodes", opResultCodes)

	default:
		return logan.Wrap(txError, "Unexpected error type").WithField("type", txError.Type())
	}
}

func initDemurrage(app *App) {
	if app.config.Core.DemurrageOperator == "" {
		log.Warn("Demurrage Operator is not set, so demurrage is not active")
		return
	}

	ctx, _ := context.WithCancel(app.ctx)
	manager := NewDemurrageManager(ctx, app.config.Core.DemurrageOperator, app.CoreInfo, app.CoreQ(), app.submitter)
	go manager.Run()
}

func init() {
	appInit.Add("demurrage", initDemurrage, "core-db", "txsub", "stellarCoreInfo")
}
