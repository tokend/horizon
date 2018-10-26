package horizon

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/zenazn/goji/web"
	"gitlab.com/tokend/horizon/actions"
	"gitlab.com/tokend/horizon/cache"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/httpx"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/log"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/toid"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/resources"
	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/go/xdr"
)

// Action is the "base type" for all actions in horizon.  It provides
// structs that embed it with access to the App struct.
//
// Additionally, this type is a trigger for go-codegen and causes
// the file at Action.tmpl to be instantiated for each struct that
// embeds Action.
type Action struct {
	actions.Base
	App *App
	Log *log.Entry

	hq history.QInterface
	cq core.QInterface

	cachedQ cache.QInterface
}

func (action *Action) GetAccountIdByBalance(balanceID string) (*string, error) {
	var balance core.Balance
	err := action.CoreQ().BalanceByID(&balance, balanceID)
	if err != nil {
		return nil, err
	}
	return &balance.AccountID, nil
}

func (action *Action) IsAllowed(ownersOfData ...string) {
	if action.Err != nil {
		return
	}

	if len(ownersOfData) == 0 {
		action.Err = errors.New("ownersOfData must not be empty")
		action.Log.WithError(action.Err)
		return
	}

	for _, ownerOfData := range ownersOfData {
		if action.Err != nil && action.Err.Error() != problem.NotAllowed.Error() {
			return
		}
		action.Err = nil
		action.isAllowed(ownerOfData)
		if action.Err == nil {
			return
		}
	}
}

func (action *Action) isAllowed(ownerOfData string) {
	//return if develop mode without signatures is used
	if action.App.config.SkipCheck {
		action.IsAdmin = true
		return
	}

	isSigner := action.IsAccountSigner(action.App.CoreInfo.MasterAccountID, action.Signer)
	if action.Err != nil {
		return
	}

	if isSigner != nil && *isSigner {
		action.IsAdmin = true
		return
	}

	// only master or master signers can access this data
	if ownerOfData == "" || ownerOfData == action.App.CoreInfo.MasterAccountID {
		action.Err = &problem.NotAllowed
		return
	}

	isSigner = action.IsAccountSigner(ownerOfData, action.Signer)
	if action.Err != nil {
		return
	}

	if ownerOfData == action.Signer && isSigner == nil {
		return
	}

	if isSigner != nil && *isSigner {
		return
	}

	action.Err = &problem.NotAllowed
}

// CoreQ provides access to queries that access the stellar core database.
func (action *Action) CoreQ() core.QInterface {
	if action.cq == nil {
		action.cq = &core.Q{Repo: action.App.CoreRepo(action.Ctx)}
	}
	return action.cq
}

// HistoryQ provides access to queries that access the history portion of
// horizon's database.
func (action *Action) HistoryQ() history.QInterface {
	if action.hq == nil {
		action.hq = &history.Q{Repo: action.App.HistoryRepo(action.Ctx)}
	}

	return action.hq
}

func (action *Action) CachedQ() cache.QInterface {
	if action.cachedQ == nil {
		action.cachedQ = cache.NewQ(action.CoreQ(), action.HistoryQ(), action.App.cacheProvider)
	}

	return action.cachedQ
}

// GetPagingParams modifies the base GetPagingParams method to replace
// cursors that are "now" with the last seen ledger's cursor.
func (action *Action) GetPagingParams() (cursor string, order string, limit uint64) {
	if action.Err != nil {
		return
	}

	cursor, order, limit = action.Base.GetPagingParams()

	if cursor == "now" {
		tid := toid.ID{
			LedgerSequence:   ledger.CurrentState().HistoryLatest,
			TransactionOrder: toid.TransactionMask,
			OperationOrder:   toid.OperationMask,
		}
		cursor = tid.String()
	}

	return
}

func (action *Action) GetPagingParamsV2 () (page uint64, limit uint64) {
	if action.Err != nil {
		return
	}

	page, limit = action.Base.GetPagingParamsV2()

	return
}

// GetPageQuery is a helper that returns a new db.PageQuery struct initialized
// using the results from a call to GetPagingParams()
func (action *Action) GetPageQuery() db2.PageQuery {
	if action.Err != nil {
		return db2.PageQuery{}
	}

	r, err := db2.NewPageQuery(action.GetPagingParams())

	if err != nil {
		action.Err = err
	}

	return r
}

func (action *Action) GetPageQueryV2() db2.PageQueryV2 {
	if action.Err != nil {
		return db2.PageQueryV2{}
	}

	r, err := db2.NewPageQueryV2(action.GetPagingParamsV2())

	if err != nil {
		action.Err = err
	}

	return r
}

// Prepare sets the action's App field based upon the goji context
func (action *Action) Prepare(c web.C, w http.ResponseWriter, r *http.Request) {
	base := &action.Base
	base.Prepare(c, w, r)
	action.App = action.GojiCtx.Env["app"].(*App)

	base.SkipCheck = action.App.config.SkipCheck //pass config variable to base (since base can't read one)

	base.Signer, _ = signcontrol.CheckSignature(r)

	if action.Ctx != nil {
		action.Log = log.Ctx(action.Ctx)
	} else {
		action.Log = log.DefaultLogger
	}
}

// ValidateCursorAsDefault ensures that the cursor parameter is valid in the way
// it is normally used, i.e. it is either the string "now" or a string of
// numerals that can be parsed as an int64.
func (action *Action) ValidateCursorAsDefault() {
	if action.Err != nil {
		return
	}

	if action.GetString(actions.ParamCursor) == "now" {
		return
	}

	action.GetInt64(actions.ParamCursor)
}

// ValidateCursorWithinHistory compares the requested page of data against the
// ledger state of the history database.  In the event that the cursor is
// guaranteed to return no results, we return a 410 GONE http response.
func (action *Action) ValidateCursorWithinHistory() {
	if action.Err != nil {
		return
	}

	pq := action.GetPageQuery()
	if action.Err != nil {
		return
	}

	// an ascending query should never return a gone response:  An ascending query
	// prior to known history should return results at the beginning of history,
	// and an ascending query beyond the end of history should not error out but
	// rather return an empty page (allowing code that tracks the procession of
	// some resource more easily).
	if pq.Order != "desc" {
		return
	}

	var cursor int64
	var err error

	// HACK: checking for the presence of "-" to see whether we should use
	// CursorInt64 or CursorInt64Pair is gross.
	if strings.Contains(pq.Cursor, "-") {
		cursor, _, err = pq.CursorInt64Pair("-")
	} else {
		cursor, err = pq.CursorInt64()
	}

	if err != nil {
		action.Err = err
		return
	}

	elder := toid.New(ledger.CurrentState().HistoryElder, 0, 0)

	if cursor <= elder.ToInt64() {
		action.Err = &problem.BeforeHistory
	}
}

// EnsureHistoryFreshness halts processing and raises
func (action *Action) EnsureHistoryFreshness() {
	if action.Err != nil {
		return
	}

	if action.App.IsHistoryStale() {
		ls := ledger.CurrentState()
		err := problem.StaleHistory
		err.Extras = map[string]interface{}{
			"history_latest_ledger": ls.HistoryLatest,
			"core_latest_ledger":    ls.CoreLatest,
		}
		action.Err = &err
	}
}

// BaseURL returns the base url for this requestion, defined as a url containing
// the Host and Scheme portions of the request uri.
func (action *Action) BaseURL() *url.URL {
	return httpx.BaseURL(action.Ctx)
}

// IsAccountSigner load core account by `accountId` and checks to see if any of the signers is`signer`
func (action *Action) IsAccountSigner(accountId, signer string) *bool {
	isSigner := new(bool)
	account, err := action.CoreQ().Accounts().ByAddress(accountId)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load account")
		action.Err = &problem.ServerError
		*isSigner = false
		return isSigner
	}

	if account == nil {
		return nil
	}

	signers, err := action.GetSigners(account)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load signers")
		action.Err = &problem.ServerError
		*isSigner = false
		return isSigner
	}

	for i := range signers {
		if signer == signers[i].Publickey {
			*isSigner = true
			return isSigner
		}
	}
	*isSigner = false
	return isSigner
}

func getSystemAccountTypes() []xdr.AccountType {
	return []xdr.AccountType{xdr.AccountTypeOperational, xdr.AccountTypeCommission, xdr.AccountTypeMaster}
}

func isSystemAccount(accountType int32) bool {
	sysAccountTypes := getSystemAccountTypes()
	for _, sysAccountType := range sysAccountTypes {
		if accountType == int32(sysAccountType) {
			return true
		}
	}

	return false
}

func (action *Action) Doorman() doorman.Doorman {
	return doorman.New(false, action)
}

// Signers used by doorman, basically just a connector to existing signers check logic
func (action *Action) Signers(address string) ([]resources.Signer, error) {
	// just to ensure backwards compatibility with checkAllowed
	if address == "" {
		address = action.App.CoreInfo.MasterAccountID
	}
	// get core account
	account, err := action.CoreQ().Accounts().ByAddress(address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}
	// pass it to legacy routine
	signers, err := action.GetSigners(account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}
	// convert structs
	result := make([]resources.Signer, 0, len(signers))
	for _, signer := range signers {
		result = append(result, resources.Signer{
			AccountID:  signer.Publickey,
			Weight:     int(signer.Weight),
			SignerType: int(signer.SignerType),
			Identity:   int(signer.Identity),
			Name:       signer.Name,
		})
	}
	return result, nil
}

func (action *Action) GetSigners(account *core.Account) ([]core.Signer, error) {
	// all system accounts are managed by master account signers
	if isSystemAccount(account.AccountType) && account.AccountType != int32(xdr.AccountTypeMaster) {
		masterAccount, err := action.CoreQ().Accounts().ByAddress(action.App.CoreInfo.MasterAccountID)
		if err != nil || masterAccount == nil {
			if err == nil {
				err = errors.New("Not found")
			}

			action.Log.WithError(err).Error("Failed to get master account from db")
			return nil, err
		}

		return action.GetSigners(masterAccount)
	}

	var signers []core.Signer
	err := action.CoreQ().SignersByAddress(&signers, account.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get signers")
		return nil, err
	}

	if !isSystemAccount(account.AccountType) {
		// add recovery signer
		signers = append(signers, core.Signer{
			Accountid:  account.RecoveryID,
			Publickey:  account.RecoveryID,
			Weight:     255,
			SignerType: action.getMasterSignerType(),
			Identity:   0,
		})
	}

	// is master key allowed
	if account.Thresholds[0] <= 0 {
		return signers, nil
	}

	signers = append(signers, core.Signer{
		Accountid:  account.AccountID,
		Publickey:  account.AccountID,
		Weight:     int32(account.Thresholds[0]),
		SignerType: action.getMasterSignerType(),
		Identity:   0,
	})

	return signers, nil
}

func (action *Action) getMasterSignerType() int32 {
	result := int32(0)
	for i := range xdr.SignerTypeAll {
		result |= int32(xdr.SignerTypeAll[i])
	}
	return result
}

func (action *Action) LoadParticipants(ownerAccountID string, participants map[int64]*history.OperationParticipants) {
	if action.IsAdmin {
		ownerAccountID = ""
	}

	// loading what we can from history database
	err := action.HistoryQ().Operations().Participants(participants)
	if err != nil {
		action.Log.WithError(err).Error("failed to get participant details from history")
		action.Err = &problem.ServerError
		return
	}

	// needs filter
	if !action.IsAdmin && !action.SkipCheck {
		for _, opParticipants := range participants {
			if opParticipants.OpType != xdr.OperationTypeManageOffer {
				continue
			}
			filteredParticipants := []*history.Participant{}
			for _, participant := range opParticipants.Participants {
				if participant.AccountID == ownerAccountID {
					filteredParticipants = append(filteredParticipants, participant)
				}
			}
			opParticipants.Participants = filteredParticipants
		}
	}
}
