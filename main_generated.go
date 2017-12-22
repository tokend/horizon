package horizon

import (
	"github.com/zenazn/goji/web"
	"net/http"
)

// ServeHTTPC is a method for web.Handler
func (action AccountBalancesAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AccountBalancesAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AccountIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AccountIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AccountShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AccountShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AccountSummaryAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AccountSummaryAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AccountTypeLimitsAllAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AccountTypeLimitsAllAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AccountTypeLimitsShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AccountTypeLimitsShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AssetPairsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AssetPairsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AssetPairsConverterAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AssetPairsConverterAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AssetsIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AssetsIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action AssetsShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "AssetsShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action BalanceAccountAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "BalanceAccountAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action BalanceAssetAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "BalanceAssetAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action BalanceIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "BalanceIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action BalanceTrustsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "BalanceTrustsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action FeesAllAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "FeesAllAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action FeesShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "FeesShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action HistoryOperationIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "HistoryOperationIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action HistoryOperationShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "HistoryOperationShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action KdfParamsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "KdfParamsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action LedgerIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "LedgerIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action LedgerShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "LedgerShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action MetricsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "MetricsAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action NotFoundAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "NotFoundAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action NotImplementedAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "NotImplementedAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action OffersAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "OffersAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action OperationIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "OperationIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action OperationShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "OperationShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action OrderBookAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "OrderBookAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PaymentRequestIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PaymentRequestIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PaymentRequestShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PaymentRequestShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action PricesHistoryAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "PricesHistoryAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action RateLimitExceededAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "RateLimitExceededAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action ReviewableRequestIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "ReviewableRequestIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action ReviewableRequestShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "ReviewableRequestShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action RootAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "RootAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action SaleIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "SaleIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action SaleShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "SaleShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action SignerShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "SignerShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action SignersIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "SignersIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action TradesAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "TradesAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action TransactionCreateAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "TransactionCreateAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action TransactionIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "TransactionIndexAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action TransactionShowAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "TransactionShowAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action WithdrawalIndexAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "WithdrawalIndexAction")
	ap.Execute(&action)
}
