// Package xdr is generated from:
//
//  xdr/raw/Stellar-ledger-entries-offer.x
//  xdr/raw/Stellar-operation-direct-debit.x
//  xdr/raw/Stellar-operation-manage-coins-emission-request.x
//  xdr/raw/Stellar-operation-manage-offer.x
//  xdr/raw/Stellar-ledger-entries-fee.x
//  xdr/raw/Stellar-ledger-entries-reference.x
//  xdr/raw/Stellar-operation-upload-preemissions.x
//  xdr/raw/Stellar-ledger-entries-account-type-limits.x
//  xdr/raw/Stellar-ledger-entries-invoice.x
//  xdr/raw/Stellar-SCP.x
//  xdr/raw/Stellar-ledger.x
//  xdr/raw/Stellar-ledger-entries-balance.x
//  xdr/raw/Stellar-overlay.x
//  xdr/raw/Stellar-operation-create-account.x
//  xdr/raw/Stellar-operation-manage-account.x
//  xdr/raw/Stellar-transaction.x
//  xdr/raw/Stellar-ledger-entries-payment-request.x
//  xdr/raw/Stellar-operation-manage-asset.x
//  xdr/raw/Stellar-operation-manage-asset-pair.x
//  xdr/raw/Stellar-ledger-entries-asset.x
//  xdr/raw/Stellar-operation-set-fees.x
//  xdr/raw/Stellar-operation-manage-balance.x
//  xdr/raw/Stellar-operation-review-payment-request.x
//  xdr/raw/Stellar-types.x
//  xdr/raw/Stellar-operation-manage-invoice.x
//  xdr/raw/Stellar-ledger-entries-coins-emission-request.x
//  xdr/raw/Stellar-ledger-entries-account-limits.x
//  xdr/raw/Stellar-operation-payment.x
//  xdr/raw/Stellar-operation-recover.x
//  xdr/raw/Stellar-operation-manage-forfeit-request.x
//  xdr/raw/Stellar-ledger-entries-account.x
//  xdr/raw/Stellar-operation-set-options.x
//  xdr/raw/Stellar-operation-review-coins-emission-request.x
//  xdr/raw/Stellar-ledger-entries-statistics.x
//  xdr/raw/Stellar-operation-set-limits.x
//  xdr/raw/Stellar-ledger-entries-asset-pair.x
//  xdr/raw/Stellar-ledger-entries.x
//
// DO NOT EDIT or your changes may be overwritten
package xdr

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nullstyle/go-xdr/xdr3"
)

// Unmarshal reads an xdr element from `r` into `v`.
func Unmarshal(r io.Reader, v interface{}) (int, error) {
	// delegate to xdr package's Unmarshal
	return xdr.Unmarshal(r, v)
}

// Marshal writes an xdr element `v` into `w`.
func Marshal(w io.Writer, v interface{}) (int, error) {
	// delegate to xdr package's Marshal
	return xdr.Marshal(w, v)
}

// OfferEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type OfferEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OfferEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OfferEntryExt
func (u OfferEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewOfferEntryExt creates a new  OfferEntryExt.
func NewOfferEntryExt(v LedgerVersion, value interface{}) (result OfferEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// OfferEntry is an XDR Struct defines as:
//
//   struct OfferEntry
//    {
//        uint64 offerID;
//    	AccountID ownerID;
//    	bool isBuy;
//        AssetCode base; // A
//        AssetCode quote;  // B
//    	BalanceID baseBalance;
//    	BalanceID quoteBalance;
//        int64 baseAmount;
//    	int64 quoteAmount;
//    	uint64 createdAt;
//    	int64 fee;
//
//        int64 percentFee;
//
//    	// price of A in terms of B
//        int64 price;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type OfferEntry struct {
	OfferId      Uint64        `json:"offerID,omitempty"`
	OwnerId      AccountId     `json:"ownerID,omitempty"`
	IsBuy        bool          `json:"isBuy,omitempty"`
	Base         AssetCode     `json:"base,omitempty"`
	Quote        AssetCode     `json:"quote,omitempty"`
	BaseBalance  BalanceId     `json:"baseBalance,omitempty"`
	QuoteBalance BalanceId     `json:"quoteBalance,omitempty"`
	BaseAmount   Int64         `json:"baseAmount,omitempty"`
	QuoteAmount  Int64         `json:"quoteAmount,omitempty"`
	CreatedAt    Uint64        `json:"createdAt,omitempty"`
	Fee          Int64         `json:"fee,omitempty"`
	PercentFee   Int64         `json:"percentFee,omitempty"`
	Price        Int64         `json:"price,omitempty"`
	Ext          OfferEntryExt `json:"ext,omitempty"`
}

// DirectDebitOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type DirectDebitOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DirectDebitOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DirectDebitOpExt
func (u DirectDebitOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewDirectDebitOpExt creates a new  DirectDebitOpExt.
func NewDirectDebitOpExt(v LedgerVersion, value interface{}) (result DirectDebitOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// DirectDebitOp is an XDR Struct defines as:
//
//   struct DirectDebitOp
//    {
//        AccountID from;
//        PaymentOp paymentOp;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type DirectDebitOp struct {
	From      AccountId        `json:"from,omitempty"`
	PaymentOp PaymentOp        `json:"paymentOp,omitempty"`
	Ext       DirectDebitOpExt `json:"ext,omitempty"`
}

// DirectDebitResultCode is an XDR Enum defines as:
//
//   enum DirectDebitResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0, // payment successfuly completed
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,       // bad input
//        UNDERFUNDED = -2,     // not enough funds in source account
//        LINE_FULL = -3,       // destination would go above their limit
//    	FEE_MISMATCHED = -4,   // fee is not equal to expected fee
//        BALANCE_NOT_FOUND = -5, // destination balance not found
//        BALANCE_ACCOUNT_MISMATCHED = -6,
//        BALANCE_ASSETS_MISMATCHED = -7,
//    	SRC_BALANCE_NOT_FOUND = -8, // source balance not found
//        REFERENCE_DUPLICATION = -9,
//        STATS_OVERFLOW = -10,
//        LIMITS_EXCEEDED = -11,
//        NOT_ALLOWED_BY_ASSET_POLICY = -12,
//        NO_TRUST = -13
//    };
//
type DirectDebitResultCode int32

const (
	DirectDebitResultCodeSuccess                  DirectDebitResultCode = 0
	DirectDebitResultCodeMalformed                DirectDebitResultCode = -1
	DirectDebitResultCodeUnderfunded              DirectDebitResultCode = -2
	DirectDebitResultCodeLineFull                 DirectDebitResultCode = -3
	DirectDebitResultCodeFeeMismatched            DirectDebitResultCode = -4
	DirectDebitResultCodeBalanceNotFound          DirectDebitResultCode = -5
	DirectDebitResultCodeBalanceAccountMismatched DirectDebitResultCode = -6
	DirectDebitResultCodeBalanceAssetsMismatched  DirectDebitResultCode = -7
	DirectDebitResultCodeSrcBalanceNotFound       DirectDebitResultCode = -8
	DirectDebitResultCodeReferenceDuplication     DirectDebitResultCode = -9
	DirectDebitResultCodeStatsOverflow            DirectDebitResultCode = -10
	DirectDebitResultCodeLimitsExceeded           DirectDebitResultCode = -11
	DirectDebitResultCodeNotAllowedByAssetPolicy  DirectDebitResultCode = -12
	DirectDebitResultCodeNoTrust                  DirectDebitResultCode = -13
)

var DirectDebitResultCodeAll = []DirectDebitResultCode{
	DirectDebitResultCodeSuccess,
	DirectDebitResultCodeMalformed,
	DirectDebitResultCodeUnderfunded,
	DirectDebitResultCodeLineFull,
	DirectDebitResultCodeFeeMismatched,
	DirectDebitResultCodeBalanceNotFound,
	DirectDebitResultCodeBalanceAccountMismatched,
	DirectDebitResultCodeBalanceAssetsMismatched,
	DirectDebitResultCodeSrcBalanceNotFound,
	DirectDebitResultCodeReferenceDuplication,
	DirectDebitResultCodeStatsOverflow,
	DirectDebitResultCodeLimitsExceeded,
	DirectDebitResultCodeNotAllowedByAssetPolicy,
	DirectDebitResultCodeNoTrust,
}

var directDebitResultCodeMap = map[int32]string{
	0:   "DirectDebitResultCodeSuccess",
	-1:  "DirectDebitResultCodeMalformed",
	-2:  "DirectDebitResultCodeUnderfunded",
	-3:  "DirectDebitResultCodeLineFull",
	-4:  "DirectDebitResultCodeFeeMismatched",
	-5:  "DirectDebitResultCodeBalanceNotFound",
	-6:  "DirectDebitResultCodeBalanceAccountMismatched",
	-7:  "DirectDebitResultCodeBalanceAssetsMismatched",
	-8:  "DirectDebitResultCodeSrcBalanceNotFound",
	-9:  "DirectDebitResultCodeReferenceDuplication",
	-10: "DirectDebitResultCodeStatsOverflow",
	-11: "DirectDebitResultCodeLimitsExceeded",
	-12: "DirectDebitResultCodeNotAllowedByAssetPolicy",
	-13: "DirectDebitResultCodeNoTrust",
}

var directDebitResultCodeShortMap = map[int32]string{
	0:   "success",
	-1:  "malformed",
	-2:  "underfunded",
	-3:  "line_full",
	-4:  "fee_mismatched",
	-5:  "balance_not_found",
	-6:  "balance_account_mismatched",
	-7:  "balance_assets_mismatched",
	-8:  "src_balance_not_found",
	-9:  "reference_duplication",
	-10: "stats_overflow",
	-11: "limits_exceeded",
	-12: "not_allowed_by_asset_policy",
	-13: "no_trust",
}

var directDebitResultCodeRevMap = map[string]int32{
	"DirectDebitResultCodeSuccess":                  0,
	"DirectDebitResultCodeMalformed":                -1,
	"DirectDebitResultCodeUnderfunded":              -2,
	"DirectDebitResultCodeLineFull":                 -3,
	"DirectDebitResultCodeFeeMismatched":            -4,
	"DirectDebitResultCodeBalanceNotFound":          -5,
	"DirectDebitResultCodeBalanceAccountMismatched": -6,
	"DirectDebitResultCodeBalanceAssetsMismatched":  -7,
	"DirectDebitResultCodeSrcBalanceNotFound":       -8,
	"DirectDebitResultCodeReferenceDuplication":     -9,
	"DirectDebitResultCodeStatsOverflow":            -10,
	"DirectDebitResultCodeLimitsExceeded":           -11,
	"DirectDebitResultCodeNotAllowedByAssetPolicy":  -12,
	"DirectDebitResultCodeNoTrust":                  -13,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for DirectDebitResultCode
func (e DirectDebitResultCode) ValidEnum(v int32) bool {
	_, ok := directDebitResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e DirectDebitResultCode) String() string {
	name, _ := directDebitResultCodeMap[int32(e)]
	return name
}

func (e DirectDebitResultCode) ShortString() string {
	name, _ := directDebitResultCodeShortMap[int32(e)]
	return name
}

func (e DirectDebitResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *DirectDebitResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := directDebitResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = DirectDebitResultCode(value)
	return nil
}

// DirectDebitSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type DirectDebitSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DirectDebitSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DirectDebitSuccessExt
func (u DirectDebitSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewDirectDebitSuccessExt creates a new  DirectDebitSuccessExt.
func NewDirectDebitSuccessExt(v LedgerVersion, value interface{}) (result DirectDebitSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// DirectDebitSuccess is an XDR Struct defines as:
//
//   struct DirectDebitSuccess {
//    	PaymentResponse paymentResponse;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type DirectDebitSuccess struct {
	PaymentResponse PaymentResponse       `json:"paymentResponse,omitempty"`
	Ext             DirectDebitSuccessExt `json:"ext,omitempty"`
}

// DirectDebitResult is an XDR Union defines as:
//
//   union DirectDebitResult switch (DirectDebitResultCode code)
//    {
//    case SUCCESS:
//        DirectDebitSuccess success;
//    default:
//        void;
//    };
//
type DirectDebitResult struct {
	Code    DirectDebitResultCode `json:"code,omitempty"`
	Success *DirectDebitSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DirectDebitResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DirectDebitResult
func (u DirectDebitResult) ArmForSwitch(sw int32) (string, bool) {
	switch DirectDebitResultCode(sw) {
	case DirectDebitResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewDirectDebitResult creates a new  DirectDebitResult.
func NewDirectDebitResult(code DirectDebitResultCode, value interface{}) (result DirectDebitResult, err error) {
	result.Code = code
	switch DirectDebitResultCode(code) {
	case DirectDebitResultCodeSuccess:
		tv, ok := value.(DirectDebitSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be DirectDebitSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u DirectDebitResult) MustSuccess() DirectDebitSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u DirectDebitResult) GetSuccess() (result DirectDebitSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageCoinsEmissionRequestAction is an XDR Enum defines as:
//
//   enum ManageCoinsEmissionRequestAction
//    {
//        CREATE = 0,
//        DELETE = 1
//    };
//
type ManageCoinsEmissionRequestAction int32

const (
	ManageCoinsEmissionRequestActionCreate ManageCoinsEmissionRequestAction = 0
	ManageCoinsEmissionRequestActionDelete ManageCoinsEmissionRequestAction = 1
)

var ManageCoinsEmissionRequestActionAll = []ManageCoinsEmissionRequestAction{
	ManageCoinsEmissionRequestActionCreate,
	ManageCoinsEmissionRequestActionDelete,
}

var manageCoinsEmissionRequestActionMap = map[int32]string{
	0: "ManageCoinsEmissionRequestActionCreate",
	1: "ManageCoinsEmissionRequestActionDelete",
}

var manageCoinsEmissionRequestActionShortMap = map[int32]string{
	0: "create",
	1: "delete",
}

var manageCoinsEmissionRequestActionRevMap = map[string]int32{
	"ManageCoinsEmissionRequestActionCreate": 0,
	"ManageCoinsEmissionRequestActionDelete": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageCoinsEmissionRequestAction
func (e ManageCoinsEmissionRequestAction) ValidEnum(v int32) bool {
	_, ok := manageCoinsEmissionRequestActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageCoinsEmissionRequestAction) String() string {
	name, _ := manageCoinsEmissionRequestActionMap[int32(e)]
	return name
}

func (e ManageCoinsEmissionRequestAction) ShortString() string {
	name, _ := manageCoinsEmissionRequestActionShortMap[int32(e)]
	return name
}

func (e ManageCoinsEmissionRequestAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageCoinsEmissionRequestAction) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageCoinsEmissionRequestActionRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageCoinsEmissionRequestAction(value)
	return nil
}

// ManageCoinsEmissionRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageCoinsEmissionRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageCoinsEmissionRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageCoinsEmissionRequestOpExt
func (u ManageCoinsEmissionRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageCoinsEmissionRequestOpExt creates a new  ManageCoinsEmissionRequestOpExt.
func NewManageCoinsEmissionRequestOpExt(v LedgerVersion, value interface{}) (result ManageCoinsEmissionRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageCoinsEmissionRequestOp is an XDR Struct defines as:
//
//   struct ManageCoinsEmissionRequestOp
//    {
//    	// 0=create a new request, otherwise edit an existing offer
//        ManageCoinsEmissionRequestAction action;
//    	uint64 requestID;
//        int64 amount;        // amount being issued. if set to 0, delete the offer
//        BalanceID receiver;
//        AssetCode asset;
//        string64 reference;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageCoinsEmissionRequestOp struct {
	Action    ManageCoinsEmissionRequestAction `json:"action,omitempty"`
	RequestId Uint64                           `json:"requestID,omitempty"`
	Amount    Int64                            `json:"amount,omitempty"`
	Receiver  BalanceId                        `json:"receiver,omitempty"`
	Asset     AssetCode                        `json:"asset,omitempty"`
	Reference String64                         `json:"reference,omitempty"`
	Ext       ManageCoinsEmissionRequestOpExt  `json:"ext,omitempty"`
}

// ManageCoinsEmissionRequestResultCode is an XDR Enum defines as:
//
//   enum ManageCoinsEmissionRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        INVALID_AMOUNT = -1,      // amount is negative
//    	INVALID_REQUEST_ID = -2, // not 0 for delete etc
//    	NOT_FOUND = -3,           // failed to find emission request with such ID
//    	ALREADY_REVIEWED = -4,    // emission request have been already reviewed - can't edit
//        ASSET_NOT_FOUND = -5,
//        BALANCE_NOT_FOUND = -6,
//        ASSET_MISMATCH = -7,
//        INVALID_ASSET = -8,
//        REFERENCE_DUPLICATION = -9,
//        LINE_FULL = -10,
//        INVALID_REFERENCE = -11
//    };
//
type ManageCoinsEmissionRequestResultCode int32

const (
	ManageCoinsEmissionRequestResultCodeSuccess              ManageCoinsEmissionRequestResultCode = 0
	ManageCoinsEmissionRequestResultCodeInvalidAmount        ManageCoinsEmissionRequestResultCode = -1
	ManageCoinsEmissionRequestResultCodeInvalidRequestId     ManageCoinsEmissionRequestResultCode = -2
	ManageCoinsEmissionRequestResultCodeNotFound             ManageCoinsEmissionRequestResultCode = -3
	ManageCoinsEmissionRequestResultCodeAlreadyReviewed      ManageCoinsEmissionRequestResultCode = -4
	ManageCoinsEmissionRequestResultCodeAssetNotFound        ManageCoinsEmissionRequestResultCode = -5
	ManageCoinsEmissionRequestResultCodeBalanceNotFound      ManageCoinsEmissionRequestResultCode = -6
	ManageCoinsEmissionRequestResultCodeAssetMismatch        ManageCoinsEmissionRequestResultCode = -7
	ManageCoinsEmissionRequestResultCodeInvalidAsset         ManageCoinsEmissionRequestResultCode = -8
	ManageCoinsEmissionRequestResultCodeReferenceDuplication ManageCoinsEmissionRequestResultCode = -9
	ManageCoinsEmissionRequestResultCodeLineFull             ManageCoinsEmissionRequestResultCode = -10
	ManageCoinsEmissionRequestResultCodeInvalidReference     ManageCoinsEmissionRequestResultCode = -11
)

var ManageCoinsEmissionRequestResultCodeAll = []ManageCoinsEmissionRequestResultCode{
	ManageCoinsEmissionRequestResultCodeSuccess,
	ManageCoinsEmissionRequestResultCodeInvalidAmount,
	ManageCoinsEmissionRequestResultCodeInvalidRequestId,
	ManageCoinsEmissionRequestResultCodeNotFound,
	ManageCoinsEmissionRequestResultCodeAlreadyReviewed,
	ManageCoinsEmissionRequestResultCodeAssetNotFound,
	ManageCoinsEmissionRequestResultCodeBalanceNotFound,
	ManageCoinsEmissionRequestResultCodeAssetMismatch,
	ManageCoinsEmissionRequestResultCodeInvalidAsset,
	ManageCoinsEmissionRequestResultCodeReferenceDuplication,
	ManageCoinsEmissionRequestResultCodeLineFull,
	ManageCoinsEmissionRequestResultCodeInvalidReference,
}

var manageCoinsEmissionRequestResultCodeMap = map[int32]string{
	0:   "ManageCoinsEmissionRequestResultCodeSuccess",
	-1:  "ManageCoinsEmissionRequestResultCodeInvalidAmount",
	-2:  "ManageCoinsEmissionRequestResultCodeInvalidRequestId",
	-3:  "ManageCoinsEmissionRequestResultCodeNotFound",
	-4:  "ManageCoinsEmissionRequestResultCodeAlreadyReviewed",
	-5:  "ManageCoinsEmissionRequestResultCodeAssetNotFound",
	-6:  "ManageCoinsEmissionRequestResultCodeBalanceNotFound",
	-7:  "ManageCoinsEmissionRequestResultCodeAssetMismatch",
	-8:  "ManageCoinsEmissionRequestResultCodeInvalidAsset",
	-9:  "ManageCoinsEmissionRequestResultCodeReferenceDuplication",
	-10: "ManageCoinsEmissionRequestResultCodeLineFull",
	-11: "ManageCoinsEmissionRequestResultCodeInvalidReference",
}

var manageCoinsEmissionRequestResultCodeShortMap = map[int32]string{
	0:   "success",
	-1:  "invalid_amount",
	-2:  "invalid_request_id",
	-3:  "not_found",
	-4:  "already_reviewed",
	-5:  "asset_not_found",
	-6:  "balance_not_found",
	-7:  "asset_mismatch",
	-8:  "invalid_asset",
	-9:  "reference_duplication",
	-10: "line_full",
	-11: "invalid_reference",
}

var manageCoinsEmissionRequestResultCodeRevMap = map[string]int32{
	"ManageCoinsEmissionRequestResultCodeSuccess":              0,
	"ManageCoinsEmissionRequestResultCodeInvalidAmount":        -1,
	"ManageCoinsEmissionRequestResultCodeInvalidRequestId":     -2,
	"ManageCoinsEmissionRequestResultCodeNotFound":             -3,
	"ManageCoinsEmissionRequestResultCodeAlreadyReviewed":      -4,
	"ManageCoinsEmissionRequestResultCodeAssetNotFound":        -5,
	"ManageCoinsEmissionRequestResultCodeBalanceNotFound":      -6,
	"ManageCoinsEmissionRequestResultCodeAssetMismatch":        -7,
	"ManageCoinsEmissionRequestResultCodeInvalidAsset":         -8,
	"ManageCoinsEmissionRequestResultCodeReferenceDuplication": -9,
	"ManageCoinsEmissionRequestResultCodeLineFull":             -10,
	"ManageCoinsEmissionRequestResultCodeInvalidReference":     -11,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageCoinsEmissionRequestResultCode
func (e ManageCoinsEmissionRequestResultCode) ValidEnum(v int32) bool {
	_, ok := manageCoinsEmissionRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageCoinsEmissionRequestResultCode) String() string {
	name, _ := manageCoinsEmissionRequestResultCodeMap[int32(e)]
	return name
}

func (e ManageCoinsEmissionRequestResultCode) ShortString() string {
	name, _ := manageCoinsEmissionRequestResultCodeShortMap[int32(e)]
	return name
}

func (e ManageCoinsEmissionRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageCoinsEmissionRequestResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageCoinsEmissionRequestResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageCoinsEmissionRequestResultCode(value)
	return nil
}

// ManageCoinsEmissionRequestResultManageRequestInfoExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ManageCoinsEmissionRequestResultManageRequestInfoExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageCoinsEmissionRequestResultManageRequestInfoExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageCoinsEmissionRequestResultManageRequestInfoExt
func (u ManageCoinsEmissionRequestResultManageRequestInfoExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageCoinsEmissionRequestResultManageRequestInfoExt creates a new  ManageCoinsEmissionRequestResultManageRequestInfoExt.
func NewManageCoinsEmissionRequestResultManageRequestInfoExt(v LedgerVersion, value interface{}) (result ManageCoinsEmissionRequestResultManageRequestInfoExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageCoinsEmissionRequestResultManageRequestInfo is an XDR NestedStruct defines as:
//
//   struct {
//            uint64 requestID;
//            bool fulfilled;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type ManageCoinsEmissionRequestResultManageRequestInfo struct {
	RequestId Uint64                                               `json:"requestID,omitempty"`
	Fulfilled bool                                                 `json:"fulfilled,omitempty"`
	Ext       ManageCoinsEmissionRequestResultManageRequestInfoExt `json:"ext,omitempty"`
}

// ManageCoinsEmissionRequestResult is an XDR Union defines as:
//
//   union ManageCoinsEmissionRequestResult switch (ManageCoinsEmissionRequestResultCode code)
//    {
//    case SUCCESS:
//        struct {
//            uint64 requestID;
//            bool fulfilled;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } manageRequestInfo;
//    default:
//        void;
//    };
//
type ManageCoinsEmissionRequestResult struct {
	Code              ManageCoinsEmissionRequestResultCode               `json:"code,omitempty"`
	ManageRequestInfo *ManageCoinsEmissionRequestResultManageRequestInfo `json:"manageRequestInfo,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageCoinsEmissionRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageCoinsEmissionRequestResult
func (u ManageCoinsEmissionRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageCoinsEmissionRequestResultCode(sw) {
	case ManageCoinsEmissionRequestResultCodeSuccess:
		return "ManageRequestInfo", true
	default:
		return "", true
	}
}

// NewManageCoinsEmissionRequestResult creates a new  ManageCoinsEmissionRequestResult.
func NewManageCoinsEmissionRequestResult(code ManageCoinsEmissionRequestResultCode, value interface{}) (result ManageCoinsEmissionRequestResult, err error) {
	result.Code = code
	switch ManageCoinsEmissionRequestResultCode(code) {
	case ManageCoinsEmissionRequestResultCodeSuccess:
		tv, ok := value.(ManageCoinsEmissionRequestResultManageRequestInfo)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageCoinsEmissionRequestResultManageRequestInfo")
			return
		}
		result.ManageRequestInfo = &tv
	default:
		// void
	}
	return
}

// MustManageRequestInfo retrieves the ManageRequestInfo value from the union,
// panicing if the value is not set.
func (u ManageCoinsEmissionRequestResult) MustManageRequestInfo() ManageCoinsEmissionRequestResultManageRequestInfo {
	val, ok := u.GetManageRequestInfo()

	if !ok {
		panic("arm ManageRequestInfo is not set")
	}

	return val
}

// GetManageRequestInfo retrieves the ManageRequestInfo value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageCoinsEmissionRequestResult) GetManageRequestInfo() (result ManageCoinsEmissionRequestResultManageRequestInfo, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "ManageRequestInfo" {
		result = *u.ManageRequestInfo
		ok = true
	}

	return
}

// ManageOfferOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageOfferOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferOpExt
func (u ManageOfferOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferOpExt creates a new  ManageOfferOpExt.
func NewManageOfferOpExt(v LedgerVersion, value interface{}) (result ManageOfferOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferOp is an XDR Struct defines as:
//
//   struct ManageOfferOp
//    {
//        BalanceID baseBalance; // balance for base asset
//    	BalanceID quoteBalance; // balance for quote asset
//    	bool isBuy;
//        int64 amount; // if set to 0, delete the offer
//        int64 price;  // price of base asset in terms of quote
//
//        int64 fee;
//
//        // 0=create a new offer, otherwise edit an existing offer
//        uint64 offerID;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageOfferOp struct {
	BaseBalance  BalanceId        `json:"baseBalance,omitempty"`
	QuoteBalance BalanceId        `json:"quoteBalance,omitempty"`
	IsBuy        bool             `json:"isBuy,omitempty"`
	Amount       Int64            `json:"amount,omitempty"`
	Price        Int64            `json:"price,omitempty"`
	Fee          Int64            `json:"fee,omitempty"`
	OfferId      Uint64           `json:"offerID,omitempty"`
	Ext          ManageOfferOpExt `json:"ext,omitempty"`
}

// ManageOfferResultCode is an XDR Enum defines as:
//
//   enum ManageOfferResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,     // generated offer would be invalid
//        PAIR_NOT_TRADED = -2, // it's not allowed to trage with this pair
//        BALANCE_NOT_FOUND = -3,  // does not own balance for buying or selling
//        UNDERFUNDED = -4,    // doesn't hold what it's trying to sell
//        CROSS_SELF = -5,     // would cross an offer from the same user
//    	OFFER_OVERFLOW = -6,
//    	ASSET_PAIR_NOT_TRADABLE = -7,
//    	PHYSICAL_PRICE_RESTRICTION = -8, // offer price violates physical price restriction
//    	CURRENT_PRICE_RESTRICTION = -9,
//        NOT_FOUND = -10, // offerID does not match an existing offer
//        INVALID_PERCENT_FEE = -11,
//    	INSUFFICIENT_PRICE = -12
//    };
//
type ManageOfferResultCode int32

const (
	ManageOfferResultCodeSuccess                  ManageOfferResultCode = 0
	ManageOfferResultCodeMalformed                ManageOfferResultCode = -1
	ManageOfferResultCodePairNotTraded            ManageOfferResultCode = -2
	ManageOfferResultCodeBalanceNotFound          ManageOfferResultCode = -3
	ManageOfferResultCodeUnderfunded              ManageOfferResultCode = -4
	ManageOfferResultCodeCrossSelf                ManageOfferResultCode = -5
	ManageOfferResultCodeOfferOverflow            ManageOfferResultCode = -6
	ManageOfferResultCodeAssetPairNotTradable     ManageOfferResultCode = -7
	ManageOfferResultCodePhysicalPriceRestriction ManageOfferResultCode = -8
	ManageOfferResultCodeCurrentPriceRestriction  ManageOfferResultCode = -9
	ManageOfferResultCodeNotFound                 ManageOfferResultCode = -10
	ManageOfferResultCodeInvalidPercentFee        ManageOfferResultCode = -11
	ManageOfferResultCodeInsufficientPrice        ManageOfferResultCode = -12
)

var ManageOfferResultCodeAll = []ManageOfferResultCode{
	ManageOfferResultCodeSuccess,
	ManageOfferResultCodeMalformed,
	ManageOfferResultCodePairNotTraded,
	ManageOfferResultCodeBalanceNotFound,
	ManageOfferResultCodeUnderfunded,
	ManageOfferResultCodeCrossSelf,
	ManageOfferResultCodeOfferOverflow,
	ManageOfferResultCodeAssetPairNotTradable,
	ManageOfferResultCodePhysicalPriceRestriction,
	ManageOfferResultCodeCurrentPriceRestriction,
	ManageOfferResultCodeNotFound,
	ManageOfferResultCodeInvalidPercentFee,
	ManageOfferResultCodeInsufficientPrice,
}

var manageOfferResultCodeMap = map[int32]string{
	0:   "ManageOfferResultCodeSuccess",
	-1:  "ManageOfferResultCodeMalformed",
	-2:  "ManageOfferResultCodePairNotTraded",
	-3:  "ManageOfferResultCodeBalanceNotFound",
	-4:  "ManageOfferResultCodeUnderfunded",
	-5:  "ManageOfferResultCodeCrossSelf",
	-6:  "ManageOfferResultCodeOfferOverflow",
	-7:  "ManageOfferResultCodeAssetPairNotTradable",
	-8:  "ManageOfferResultCodePhysicalPriceRestriction",
	-9:  "ManageOfferResultCodeCurrentPriceRestriction",
	-10: "ManageOfferResultCodeNotFound",
	-11: "ManageOfferResultCodeInvalidPercentFee",
	-12: "ManageOfferResultCodeInsufficientPrice",
}

var manageOfferResultCodeShortMap = map[int32]string{
	0:   "success",
	-1:  "malformed",
	-2:  "pair_not_traded",
	-3:  "balance_not_found",
	-4:  "underfunded",
	-5:  "cross_self",
	-6:  "offer_overflow",
	-7:  "asset_pair_not_tradable",
	-8:  "physical_price_restriction",
	-9:  "current_price_restriction",
	-10: "not_found",
	-11: "invalid_percent_fee",
	-12: "insufficient_price",
}

var manageOfferResultCodeRevMap = map[string]int32{
	"ManageOfferResultCodeSuccess":                  0,
	"ManageOfferResultCodeMalformed":                -1,
	"ManageOfferResultCodePairNotTraded":            -2,
	"ManageOfferResultCodeBalanceNotFound":          -3,
	"ManageOfferResultCodeUnderfunded":              -4,
	"ManageOfferResultCodeCrossSelf":                -5,
	"ManageOfferResultCodeOfferOverflow":            -6,
	"ManageOfferResultCodeAssetPairNotTradable":     -7,
	"ManageOfferResultCodePhysicalPriceRestriction": -8,
	"ManageOfferResultCodeCurrentPriceRestriction":  -9,
	"ManageOfferResultCodeNotFound":                 -10,
	"ManageOfferResultCodeInvalidPercentFee":        -11,
	"ManageOfferResultCodeInsufficientPrice":        -12,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageOfferResultCode
func (e ManageOfferResultCode) ValidEnum(v int32) bool {
	_, ok := manageOfferResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageOfferResultCode) String() string {
	name, _ := manageOfferResultCodeMap[int32(e)]
	return name
}

func (e ManageOfferResultCode) ShortString() string {
	name, _ := manageOfferResultCodeShortMap[int32(e)]
	return name
}

func (e ManageOfferResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageOfferResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageOfferResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageOfferResultCode(value)
	return nil
}

// ManageOfferEffect is an XDR Enum defines as:
//
//   enum ManageOfferEffect
//    {
//        CREATED = 0,
//        UPDATED = 1,
//        DELETED = 2
//    };
//
type ManageOfferEffect int32

const (
	ManageOfferEffectCreated ManageOfferEffect = 0
	ManageOfferEffectUpdated ManageOfferEffect = 1
	ManageOfferEffectDeleted ManageOfferEffect = 2
)

var ManageOfferEffectAll = []ManageOfferEffect{
	ManageOfferEffectCreated,
	ManageOfferEffectUpdated,
	ManageOfferEffectDeleted,
}

var manageOfferEffectMap = map[int32]string{
	0: "ManageOfferEffectCreated",
	1: "ManageOfferEffectUpdated",
	2: "ManageOfferEffectDeleted",
}

var manageOfferEffectShortMap = map[int32]string{
	0: "created",
	1: "updated",
	2: "deleted",
}

var manageOfferEffectRevMap = map[string]int32{
	"ManageOfferEffectCreated": 0,
	"ManageOfferEffectUpdated": 1,
	"ManageOfferEffectDeleted": 2,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageOfferEffect
func (e ManageOfferEffect) ValidEnum(v int32) bool {
	_, ok := manageOfferEffectMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageOfferEffect) String() string {
	name, _ := manageOfferEffectMap[int32(e)]
	return name
}

func (e ManageOfferEffect) ShortString() string {
	name, _ := manageOfferEffectShortMap[int32(e)]
	return name
}

func (e ManageOfferEffect) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageOfferEffect) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageOfferEffectRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageOfferEffect(value)
	return nil
}

// ClaimOfferAtomExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ClaimOfferAtomExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimOfferAtomExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimOfferAtomExt
func (u ClaimOfferAtomExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewClaimOfferAtomExt creates a new  ClaimOfferAtomExt.
func NewClaimOfferAtomExt(v LedgerVersion, value interface{}) (result ClaimOfferAtomExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ClaimOfferAtom is an XDR Struct defines as:
//
//   struct ClaimOfferAtom
//    {
//        // emitted to identify the offer
//        AccountID bAccountID; // Account that owns the offer
//        uint64 offerID;
//    	int64 baseAmount;
//    	int64 quoteAmount;
//    	int64 bFeePaid;
//    	int64 aFeePaid;
//    	BalanceID baseBalance;
//    	BalanceID quoteBalance;
//
//    	int64 currentPrice;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ClaimOfferAtom struct {
	BAccountId   AccountId         `json:"bAccountID,omitempty"`
	OfferId      Uint64            `json:"offerID,omitempty"`
	BaseAmount   Int64             `json:"baseAmount,omitempty"`
	QuoteAmount  Int64             `json:"quoteAmount,omitempty"`
	BFeePaid     Int64             `json:"bFeePaid,omitempty"`
	AFeePaid     Int64             `json:"aFeePaid,omitempty"`
	BaseBalance  BalanceId         `json:"baseBalance,omitempty"`
	QuoteBalance BalanceId         `json:"quoteBalance,omitempty"`
	CurrentPrice Int64             `json:"currentPrice,omitempty"`
	Ext          ClaimOfferAtomExt `json:"ext,omitempty"`
}

// ManageOfferSuccessResultOffer is an XDR NestedUnion defines as:
//
//   union switch (ManageOfferEffect effect)
//        {
//        case CREATED:
//        case UPDATED:
//            OfferEntry offer;
//        default:
//            void;
//        }
//
type ManageOfferSuccessResultOffer struct {
	Effect ManageOfferEffect `json:"effect,omitempty"`
	Offer  *OfferEntry       `json:"offer,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferSuccessResultOffer) SwitchFieldName() string {
	return "Effect"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferSuccessResultOffer
func (u ManageOfferSuccessResultOffer) ArmForSwitch(sw int32) (string, bool) {
	switch ManageOfferEffect(sw) {
	case ManageOfferEffectCreated:
		return "Offer", true
	case ManageOfferEffectUpdated:
		return "Offer", true
	default:
		return "", true
	}
}

// NewManageOfferSuccessResultOffer creates a new  ManageOfferSuccessResultOffer.
func NewManageOfferSuccessResultOffer(effect ManageOfferEffect, value interface{}) (result ManageOfferSuccessResultOffer, err error) {
	result.Effect = effect
	switch ManageOfferEffect(effect) {
	case ManageOfferEffectCreated:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	case ManageOfferEffectUpdated:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	default:
		// void
	}
	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u ManageOfferSuccessResultOffer) MustOffer() OfferEntry {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferSuccessResultOffer) GetOffer() (result OfferEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Effect))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// ManageOfferSuccessResultExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageOfferSuccessResultExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferSuccessResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferSuccessResultExt
func (u ManageOfferSuccessResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferSuccessResultExt creates a new  ManageOfferSuccessResultExt.
func NewManageOfferSuccessResultExt(v LedgerVersion, value interface{}) (result ManageOfferSuccessResultExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferSuccessResult is an XDR Struct defines as:
//
//   struct ManageOfferSuccessResult
//    {
//
//        // offers that got claimed while creating this offer
//        ClaimOfferAtom offersClaimed<>;
//    	AssetCode baseAsset;
//    	AssetCode quoteAsset;
//
//        union switch (ManageOfferEffect effect)
//        {
//        case CREATED:
//        case UPDATED:
//            OfferEntry offer;
//        default:
//            void;
//        }
//        offer;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageOfferSuccessResult struct {
	OffersClaimed []ClaimOfferAtom              `json:"offersClaimed,omitempty"`
	BaseAsset     AssetCode                     `json:"baseAsset,omitempty"`
	QuoteAsset    AssetCode                     `json:"quoteAsset,omitempty"`
	Offer         ManageOfferSuccessResultOffer `json:"offer,omitempty"`
	Ext           ManageOfferSuccessResultExt   `json:"ext,omitempty"`
}

// ManageOfferResultPhysicalPriceRestrictionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ManageOfferResultPhysicalPriceRestrictionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferResultPhysicalPriceRestrictionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferResultPhysicalPriceRestrictionExt
func (u ManageOfferResultPhysicalPriceRestrictionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferResultPhysicalPriceRestrictionExt creates a new  ManageOfferResultPhysicalPriceRestrictionExt.
func NewManageOfferResultPhysicalPriceRestrictionExt(v LedgerVersion, value interface{}) (result ManageOfferResultPhysicalPriceRestrictionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferResultPhysicalPriceRestriction is an XDR NestedStruct defines as:
//
//   struct {
//    		int64 physicalPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type ManageOfferResultPhysicalPriceRestriction struct {
	PhysicalPrice Int64                                        `json:"physicalPrice,omitempty"`
	Ext           ManageOfferResultPhysicalPriceRestrictionExt `json:"ext,omitempty"`
}

// ManageOfferResultCurrentPriceRestrictionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ManageOfferResultCurrentPriceRestrictionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferResultCurrentPriceRestrictionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferResultCurrentPriceRestrictionExt
func (u ManageOfferResultCurrentPriceRestrictionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferResultCurrentPriceRestrictionExt creates a new  ManageOfferResultCurrentPriceRestrictionExt.
func NewManageOfferResultCurrentPriceRestrictionExt(v LedgerVersion, value interface{}) (result ManageOfferResultCurrentPriceRestrictionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferResultCurrentPriceRestriction is an XDR NestedStruct defines as:
//
//   struct {
//    		int64 currentPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type ManageOfferResultCurrentPriceRestriction struct {
	CurrentPrice Int64                                       `json:"currentPrice,omitempty"`
	Ext          ManageOfferResultCurrentPriceRestrictionExt `json:"ext,omitempty"`
}

// ManageOfferResult is an XDR Union defines as:
//
//   union ManageOfferResult switch (ManageOfferResultCode code)
//    {
//    case SUCCESS:
//        ManageOfferSuccessResult success;
//    case PHYSICAL_PRICE_RESTRICTION:
//    	struct {
//    		int64 physicalPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} physicalPriceRestriction;
//    case CURRENT_PRICE_RESTRICTION:
//    	struct {
//    		int64 currentPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} currentPriceRestriction;
//
//    default:
//        void;
//    };
//
type ManageOfferResult struct {
	Code                     ManageOfferResultCode                      `json:"code,omitempty"`
	Success                  *ManageOfferSuccessResult                  `json:"success,omitempty"`
	PhysicalPriceRestriction *ManageOfferResultPhysicalPriceRestriction `json:"physicalPriceRestriction,omitempty"`
	CurrentPriceRestriction  *ManageOfferResultCurrentPriceRestriction  `json:"currentPriceRestriction,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferResult
func (u ManageOfferResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageOfferResultCode(sw) {
	case ManageOfferResultCodeSuccess:
		return "Success", true
	case ManageOfferResultCodePhysicalPriceRestriction:
		return "PhysicalPriceRestriction", true
	case ManageOfferResultCodeCurrentPriceRestriction:
		return "CurrentPriceRestriction", true
	default:
		return "", true
	}
}

// NewManageOfferResult creates a new  ManageOfferResult.
func NewManageOfferResult(code ManageOfferResultCode, value interface{}) (result ManageOfferResult, err error) {
	result.Code = code
	switch ManageOfferResultCode(code) {
	case ManageOfferResultCodeSuccess:
		tv, ok := value.(ManageOfferSuccessResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferSuccessResult")
			return
		}
		result.Success = &tv
	case ManageOfferResultCodePhysicalPriceRestriction:
		tv, ok := value.(ManageOfferResultPhysicalPriceRestriction)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferResultPhysicalPriceRestriction")
			return
		}
		result.PhysicalPriceRestriction = &tv
	case ManageOfferResultCodeCurrentPriceRestriction:
		tv, ok := value.(ManageOfferResultCurrentPriceRestriction)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferResultCurrentPriceRestriction")
			return
		}
		result.CurrentPriceRestriction = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageOfferResult) MustSuccess() ManageOfferSuccessResult {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferResult) GetSuccess() (result ManageOfferSuccessResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// MustPhysicalPriceRestriction retrieves the PhysicalPriceRestriction value from the union,
// panicing if the value is not set.
func (u ManageOfferResult) MustPhysicalPriceRestriction() ManageOfferResultPhysicalPriceRestriction {
	val, ok := u.GetPhysicalPriceRestriction()

	if !ok {
		panic("arm PhysicalPriceRestriction is not set")
	}

	return val
}

// GetPhysicalPriceRestriction retrieves the PhysicalPriceRestriction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferResult) GetPhysicalPriceRestriction() (result ManageOfferResultPhysicalPriceRestriction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "PhysicalPriceRestriction" {
		result = *u.PhysicalPriceRestriction
		ok = true
	}

	return
}

// MustCurrentPriceRestriction retrieves the CurrentPriceRestriction value from the union,
// panicing if the value is not set.
func (u ManageOfferResult) MustCurrentPriceRestriction() ManageOfferResultCurrentPriceRestriction {
	val, ok := u.GetCurrentPriceRestriction()

	if !ok {
		panic("arm CurrentPriceRestriction is not set")
	}

	return val
}

// GetCurrentPriceRestriction retrieves the CurrentPriceRestriction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferResult) GetCurrentPriceRestriction() (result ManageOfferResultCurrentPriceRestriction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "CurrentPriceRestriction" {
		result = *u.CurrentPriceRestriction
		ok = true
	}

	return
}

// FeeType is an XDR Enum defines as:
//
//   enum FeeType
//    {
//        PAYMENT_FEE = 0,
//        REFERRAL_FEE = 1,
//    	OFFER_FEE = 2,
//        FORFEIT_FEE = 3,
//        EMISSION_FEE = 4
//    };
//
type FeeType int32

const (
	FeeTypePaymentFee  FeeType = 0
	FeeTypeReferralFee FeeType = 1
	FeeTypeOfferFee    FeeType = 2
	FeeTypeForfeitFee  FeeType = 3
	FeeTypeEmissionFee FeeType = 4
)

var FeeTypeAll = []FeeType{
	FeeTypePaymentFee,
	FeeTypeReferralFee,
	FeeTypeOfferFee,
	FeeTypeForfeitFee,
	FeeTypeEmissionFee,
}

var feeTypeMap = map[int32]string{
	0: "FeeTypePaymentFee",
	1: "FeeTypeReferralFee",
	2: "FeeTypeOfferFee",
	3: "FeeTypeForfeitFee",
	4: "FeeTypeEmissionFee",
}

var feeTypeShortMap = map[int32]string{
	0: "payment_fee",
	1: "referral_fee",
	2: "offer_fee",
	3: "forfeit_fee",
	4: "emission_fee",
}

var feeTypeRevMap = map[string]int32{
	"FeeTypePaymentFee":  0,
	"FeeTypeReferralFee": 1,
	"FeeTypeOfferFee":    2,
	"FeeTypeForfeitFee":  3,
	"FeeTypeEmissionFee": 4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for FeeType
func (e FeeType) ValidEnum(v int32) bool {
	_, ok := feeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e FeeType) String() string {
	name, _ := feeTypeMap[int32(e)]
	return name
}

func (e FeeType) ShortString() string {
	name, _ := feeTypeShortMap[int32(e)]
	return name
}

func (e FeeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *FeeType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := feeTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = FeeType(value)
	return nil
}

// EmissionFeeType is an XDR Enum defines as:
//
//   enum EmissionFeeType
//    {
//    	PRIMARY_MARKET = 1,
//    	SECONDARY_MARKET = 2
//    };
//
type EmissionFeeType int32

const (
	EmissionFeeTypePrimaryMarket   EmissionFeeType = 1
	EmissionFeeTypeSecondaryMarket EmissionFeeType = 2
)

var EmissionFeeTypeAll = []EmissionFeeType{
	EmissionFeeTypePrimaryMarket,
	EmissionFeeTypeSecondaryMarket,
}

var emissionFeeTypeMap = map[int32]string{
	1: "EmissionFeeTypePrimaryMarket",
	2: "EmissionFeeTypeSecondaryMarket",
}

var emissionFeeTypeShortMap = map[int32]string{
	1: "primary_market",
	2: "secondary_market",
}

var emissionFeeTypeRevMap = map[string]int32{
	"EmissionFeeTypePrimaryMarket":   1,
	"EmissionFeeTypeSecondaryMarket": 2,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EmissionFeeType
func (e EmissionFeeType) ValidEnum(v int32) bool {
	_, ok := emissionFeeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e EmissionFeeType) String() string {
	name, _ := emissionFeeTypeMap[int32(e)]
	return name
}

func (e EmissionFeeType) ShortString() string {
	name, _ := emissionFeeTypeShortMap[int32(e)]
	return name
}

func (e EmissionFeeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *EmissionFeeType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := emissionFeeTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = EmissionFeeType(value)
	return nil
}

// FeeEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type FeeEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeEntryExt
func (u FeeEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewFeeEntryExt creates a new  FeeEntryExt.
func NewFeeEntryExt(v LedgerVersion, value interface{}) (result FeeEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// FeeEntry is an XDR Struct defines as:
//
//   struct FeeEntry
//    {
//        FeeType feeType;
//        AssetCode asset;
//        int64 fixedFee; // fee paid for operation
//    	int64 percentFee; // percent of transfer amount to be charged
//
//        AccountID* accountID;
//        AccountType* accountType;
//        int64 subtype; // for example, different withdrawals — bars or coins
//
//        int64 lowerBound;
//        int64 upperBound;
//
//        Hash hash;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//
//    };
//
type FeeEntry struct {
	FeeType     FeeType      `json:"feeType,omitempty"`
	Asset       AssetCode    `json:"asset,omitempty"`
	FixedFee    Int64        `json:"fixedFee,omitempty"`
	PercentFee  Int64        `json:"percentFee,omitempty"`
	AccountId   *AccountId   `json:"accountID,omitempty"`
	AccountType *AccountType `json:"accountType,omitempty"`
	Subtype     Int64        `json:"subtype,omitempty"`
	LowerBound  Int64        `json:"lowerBound,omitempty"`
	UpperBound  Int64        `json:"upperBound,omitempty"`
	Hash        Hash         `json:"hash,omitempty"`
	Ext         FeeEntryExt  `json:"ext,omitempty"`
}

// ReferenceEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ReferenceEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReferenceEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReferenceEntryExt
func (u ReferenceEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReferenceEntryExt creates a new  ReferenceEntryExt.
func NewReferenceEntryExt(v LedgerVersion, value interface{}) (result ReferenceEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReferenceEntry is an XDR Struct defines as:
//
//   struct ReferenceEntry
//    {
//    	AccountID sender;
//        string64 reference;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ReferenceEntry struct {
	Sender    AccountId         `json:"sender,omitempty"`
	Reference String64          `json:"reference,omitempty"`
	Ext       ReferenceEntryExt `json:"ext,omitempty"`
}

// PreEmissionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type PreEmissionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PreEmissionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PreEmissionExt
func (u PreEmissionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPreEmissionExt creates a new  PreEmissionExt.
func NewPreEmissionExt(v LedgerVersion, value interface{}) (result PreEmissionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PreEmission is an XDR Struct defines as:
//
//   struct PreEmission
//    {
//        string64 serialNumber;
//        AssetCode asset;
//        int64 amount;
//        DecoratedSignature signatures<20>;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type PreEmission struct {
	SerialNumber String64             `json:"serialNumber,omitempty"`
	Asset        AssetCode            `json:"asset,omitempty"`
	Amount       Int64                `json:"amount,omitempty"`
	Signatures   []DecoratedSignature `json:"signatures,omitempty" xdrmaxsize:"20"`
	Ext          PreEmissionExt       `json:"ext,omitempty"`
}

// UploadPreemissionsOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type UploadPreemissionsOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u UploadPreemissionsOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of UploadPreemissionsOpExt
func (u UploadPreemissionsOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewUploadPreemissionsOpExt creates a new  UploadPreemissionsOpExt.
func NewUploadPreemissionsOpExt(v LedgerVersion, value interface{}) (result UploadPreemissionsOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// UploadPreemissionsOp is an XDR Struct defines as:
//
//   struct UploadPreemissionsOp
//    {
//        PreEmission preEmissions<>;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type UploadPreemissionsOp struct {
	PreEmissions []PreEmission           `json:"preEmissions,omitempty"`
	Ext          UploadPreemissionsOpExt `json:"ext,omitempty"`
}

// UploadPreemissionsResultCode is an XDR Enum defines as:
//
//   enum UploadPreemissionsResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,
//        SERIAL_DUPLICATION = -2,    // serial is already used
//        MALFORMED_PREEMISSIONS = -3, // if pre-emissions has empty signatures or zero amount etc
//        ASSET_NOT_FOUND = -4,
//        LINE_FULL = -5
//    };
//
type UploadPreemissionsResultCode int32

const (
	UploadPreemissionsResultCodeSuccess               UploadPreemissionsResultCode = 0
	UploadPreemissionsResultCodeMalformed             UploadPreemissionsResultCode = -1
	UploadPreemissionsResultCodeSerialDuplication     UploadPreemissionsResultCode = -2
	UploadPreemissionsResultCodeMalformedPreemissions UploadPreemissionsResultCode = -3
	UploadPreemissionsResultCodeAssetNotFound         UploadPreemissionsResultCode = -4
	UploadPreemissionsResultCodeLineFull              UploadPreemissionsResultCode = -5
)

var UploadPreemissionsResultCodeAll = []UploadPreemissionsResultCode{
	UploadPreemissionsResultCodeSuccess,
	UploadPreemissionsResultCodeMalformed,
	UploadPreemissionsResultCodeSerialDuplication,
	UploadPreemissionsResultCodeMalformedPreemissions,
	UploadPreemissionsResultCodeAssetNotFound,
	UploadPreemissionsResultCodeLineFull,
}

var uploadPreemissionsResultCodeMap = map[int32]string{
	0:  "UploadPreemissionsResultCodeSuccess",
	-1: "UploadPreemissionsResultCodeMalformed",
	-2: "UploadPreemissionsResultCodeSerialDuplication",
	-3: "UploadPreemissionsResultCodeMalformedPreemissions",
	-4: "UploadPreemissionsResultCodeAssetNotFound",
	-5: "UploadPreemissionsResultCodeLineFull",
}

var uploadPreemissionsResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "malformed",
	-2: "serial_duplication",
	-3: "malformed_preemissions",
	-4: "asset_not_found",
	-5: "line_full",
}

var uploadPreemissionsResultCodeRevMap = map[string]int32{
	"UploadPreemissionsResultCodeSuccess":               0,
	"UploadPreemissionsResultCodeMalformed":             -1,
	"UploadPreemissionsResultCodeSerialDuplication":     -2,
	"UploadPreemissionsResultCodeMalformedPreemissions": -3,
	"UploadPreemissionsResultCodeAssetNotFound":         -4,
	"UploadPreemissionsResultCodeLineFull":              -5,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for UploadPreemissionsResultCode
func (e UploadPreemissionsResultCode) ValidEnum(v int32) bool {
	_, ok := uploadPreemissionsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e UploadPreemissionsResultCode) String() string {
	name, _ := uploadPreemissionsResultCodeMap[int32(e)]
	return name
}

func (e UploadPreemissionsResultCode) ShortString() string {
	name, _ := uploadPreemissionsResultCodeShortMap[int32(e)]
	return name
}

func (e UploadPreemissionsResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *UploadPreemissionsResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := uploadPreemissionsResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = UploadPreemissionsResultCode(value)
	return nil
}

// UploadPreemissionsResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type UploadPreemissionsResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u UploadPreemissionsResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of UploadPreemissionsResultSuccessExt
func (u UploadPreemissionsResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewUploadPreemissionsResultSuccessExt creates a new  UploadPreemissionsResultSuccessExt.
func NewUploadPreemissionsResultSuccessExt(v LedgerVersion, value interface{}) (result UploadPreemissionsResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// UploadPreemissionsResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type UploadPreemissionsResultSuccess struct {
	Ext UploadPreemissionsResultSuccessExt `json:"ext,omitempty"`
}

// UploadPreemissionsResult is an XDR Union defines as:
//
//   union UploadPreemissionsResult switch (UploadPreemissionsResultCode code)
//    {
//    case SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type UploadPreemissionsResult struct {
	Code    UploadPreemissionsResultCode     `json:"code,omitempty"`
	Success *UploadPreemissionsResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u UploadPreemissionsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of UploadPreemissionsResult
func (u UploadPreemissionsResult) ArmForSwitch(sw int32) (string, bool) {
	switch UploadPreemissionsResultCode(sw) {
	case UploadPreemissionsResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewUploadPreemissionsResult creates a new  UploadPreemissionsResult.
func NewUploadPreemissionsResult(code UploadPreemissionsResultCode, value interface{}) (result UploadPreemissionsResult, err error) {
	result.Code = code
	switch UploadPreemissionsResultCode(code) {
	case UploadPreemissionsResultCodeSuccess:
		tv, ok := value.(UploadPreemissionsResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be UploadPreemissionsResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u UploadPreemissionsResult) MustSuccess() UploadPreemissionsResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u UploadPreemissionsResult) GetSuccess() (result UploadPreemissionsResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// AccountTypeLimitsEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AccountTypeLimitsEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountTypeLimitsEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountTypeLimitsEntryExt
func (u AccountTypeLimitsEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAccountTypeLimitsEntryExt creates a new  AccountTypeLimitsEntryExt.
func NewAccountTypeLimitsEntryExt(v LedgerVersion, value interface{}) (result AccountTypeLimitsEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AccountTypeLimitsEntry is an XDR Struct defines as:
//
//   struct AccountTypeLimitsEntry
//    {
//    	AccountType accountType;
//        Limits limits;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AccountTypeLimitsEntry struct {
	AccountType AccountType               `json:"accountType,omitempty"`
	Limits      Limits                    `json:"limits,omitempty"`
	Ext         AccountTypeLimitsEntryExt `json:"ext,omitempty"`
}

// InvoiceState is an XDR Enum defines as:
//
//   enum InvoiceState
//    {
//        INVOICE_NEEDS_PAYMENT = 0,
//        INVOICE_NEEDS_PAYMENT_REVIEW = 1
//    };
//
type InvoiceState int32

const (
	InvoiceStateInvoiceNeedsPayment       InvoiceState = 0
	InvoiceStateInvoiceNeedsPaymentReview InvoiceState = 1
)

var InvoiceStateAll = []InvoiceState{
	InvoiceStateInvoiceNeedsPayment,
	InvoiceStateInvoiceNeedsPaymentReview,
}

var invoiceStateMap = map[int32]string{
	0: "InvoiceStateInvoiceNeedsPayment",
	1: "InvoiceStateInvoiceNeedsPaymentReview",
}

var invoiceStateShortMap = map[int32]string{
	0: "invoice_needs_payment",
	1: "invoice_needs_payment_review",
}

var invoiceStateRevMap = map[string]int32{
	"InvoiceStateInvoiceNeedsPayment":       0,
	"InvoiceStateInvoiceNeedsPaymentReview": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for InvoiceState
func (e InvoiceState) ValidEnum(v int32) bool {
	_, ok := invoiceStateMap[v]
	return ok
}

// String returns the name of `e`
func (e InvoiceState) String() string {
	name, _ := invoiceStateMap[int32(e)]
	return name
}

func (e InvoiceState) ShortString() string {
	name, _ := invoiceStateShortMap[int32(e)]
	return name
}

func (e InvoiceState) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *InvoiceState) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := invoiceStateRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = InvoiceState(value)
	return nil
}

// InvoiceEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type InvoiceEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InvoiceEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InvoiceEntryExt
func (u InvoiceEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewInvoiceEntryExt creates a new  InvoiceEntryExt.
func NewInvoiceEntryExt(v LedgerVersion, value interface{}) (result InvoiceEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// InvoiceEntry is an XDR Struct defines as:
//
//   struct InvoiceEntry
//    {
//        uint64 invoiceID;
//        AccountID receiverAccount;
//        BalanceID receiverBalance;
//    	AccountID sender;
//        int64 amount;
//
//        InvoiceState state;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type InvoiceEntry struct {
	InvoiceId       Uint64          `json:"invoiceID,omitempty"`
	ReceiverAccount AccountId       `json:"receiverAccount,omitempty"`
	ReceiverBalance BalanceId       `json:"receiverBalance,omitempty"`
	Sender          AccountId       `json:"sender,omitempty"`
	Amount          Int64           `json:"amount,omitempty"`
	State           InvoiceState    `json:"state,omitempty"`
	Ext             InvoiceEntryExt `json:"ext,omitempty"`
}

// Value is an XDR Typedef defines as:
//
//   typedef opaque Value<>;
//
type Value []byte

// ScpBallot is an XDR Struct defines as:
//
//   struct SCPBallot
//    {
//        uint32 counter; // n
//        Value value;    // x
//    };
//
type ScpBallot struct {
	Counter Uint32 `json:"counter,omitempty"`
	Value   Value  `json:"value,omitempty"`
}

// ScpStatementType is an XDR Enum defines as:
//
//   enum SCPStatementType
//    {
//        PREPARE = 0,
//        CONFIRM = 1,
//        EXTERNALIZE = 2,
//        NOMINATE = 3
//    };
//
type ScpStatementType int32

const (
	ScpStatementTypePrepare     ScpStatementType = 0
	ScpStatementTypeConfirm     ScpStatementType = 1
	ScpStatementTypeExternalize ScpStatementType = 2
	ScpStatementTypeNominate    ScpStatementType = 3
)

var ScpStatementTypeAll = []ScpStatementType{
	ScpStatementTypePrepare,
	ScpStatementTypeConfirm,
	ScpStatementTypeExternalize,
	ScpStatementTypeNominate,
}

var scpStatementTypeMap = map[int32]string{
	0: "ScpStatementTypePrepare",
	1: "ScpStatementTypeConfirm",
	2: "ScpStatementTypeExternalize",
	3: "ScpStatementTypeNominate",
}

var scpStatementTypeShortMap = map[int32]string{
	0: "prepare",
	1: "confirm",
	2: "externalize",
	3: "nominate",
}

var scpStatementTypeRevMap = map[string]int32{
	"ScpStatementTypePrepare":     0,
	"ScpStatementTypeConfirm":     1,
	"ScpStatementTypeExternalize": 2,
	"ScpStatementTypeNominate":    3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ScpStatementType
func (e ScpStatementType) ValidEnum(v int32) bool {
	_, ok := scpStatementTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ScpStatementType) String() string {
	name, _ := scpStatementTypeMap[int32(e)]
	return name
}

func (e ScpStatementType) ShortString() string {
	name, _ := scpStatementTypeShortMap[int32(e)]
	return name
}

func (e ScpStatementType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ScpStatementType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := scpStatementTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ScpStatementType(value)
	return nil
}

// ScpNomination is an XDR Struct defines as:
//
//   struct SCPNomination
//    {
//        Hash quorumSetHash; // D
//        Value votes<>;      // X
//        Value accepted<>;   // Y
//    };
//
type ScpNomination struct {
	QuorumSetHash Hash    `json:"quorumSetHash,omitempty"`
	Votes         []Value `json:"votes,omitempty"`
	Accepted      []Value `json:"accepted,omitempty"`
}

// ScpStatementPrepare is an XDR NestedStruct defines as:
//
//   struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            }
//
type ScpStatementPrepare struct {
	QuorumSetHash Hash       `json:"quorumSetHash,omitempty"`
	Ballot        ScpBallot  `json:"ballot,omitempty"`
	Prepared      *ScpBallot `json:"prepared,omitempty"`
	PreparedPrime *ScpBallot `json:"preparedPrime,omitempty"`
	NC            Uint32     `json:"nC,omitempty"`
	NH            Uint32     `json:"nH,omitempty"`
}

// ScpStatementConfirm is an XDR NestedStruct defines as:
//
//   struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            }
//
type ScpStatementConfirm struct {
	Ballot        ScpBallot `json:"ballot,omitempty"`
	NPrepared     Uint32    `json:"nPrepared,omitempty"`
	NCommit       Uint32    `json:"nCommit,omitempty"`
	NH            Uint32    `json:"nH,omitempty"`
	QuorumSetHash Hash      `json:"quorumSetHash,omitempty"`
}

// ScpStatementExternalize is an XDR NestedStruct defines as:
//
//   struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            }
//
type ScpStatementExternalize struct {
	Commit              ScpBallot `json:"commit,omitempty"`
	NH                  Uint32    `json:"nH,omitempty"`
	CommitQuorumSetHash Hash      `json:"commitQuorumSetHash,omitempty"`
}

// ScpStatementPledges is an XDR NestedUnion defines as:
//
//   union switch (SCPStatementType type)
//        {
//        case PREPARE:
//            struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            } prepare;
//        case CONFIRM:
//            struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            } confirm;
//        case EXTERNALIZE:
//            struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            } externalize;
//        case NOMINATE:
//            SCPNomination nominate;
//        }
//
type ScpStatementPledges struct {
	Type        ScpStatementType         `json:"type,omitempty"`
	Prepare     *ScpStatementPrepare     `json:"prepare,omitempty"`
	Confirm     *ScpStatementConfirm     `json:"confirm,omitempty"`
	Externalize *ScpStatementExternalize `json:"externalize,omitempty"`
	Nominate    *ScpNomination           `json:"nominate,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ScpStatementPledges) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ScpStatementPledges
func (u ScpStatementPledges) ArmForSwitch(sw int32) (string, bool) {
	switch ScpStatementType(sw) {
	case ScpStatementTypePrepare:
		return "Prepare", true
	case ScpStatementTypeConfirm:
		return "Confirm", true
	case ScpStatementTypeExternalize:
		return "Externalize", true
	case ScpStatementTypeNominate:
		return "Nominate", true
	}
	return "-", false
}

// NewScpStatementPledges creates a new  ScpStatementPledges.
func NewScpStatementPledges(aType ScpStatementType, value interface{}) (result ScpStatementPledges, err error) {
	result.Type = aType
	switch ScpStatementType(aType) {
	case ScpStatementTypePrepare:
		tv, ok := value.(ScpStatementPrepare)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementPrepare")
			return
		}
		result.Prepare = &tv
	case ScpStatementTypeConfirm:
		tv, ok := value.(ScpStatementConfirm)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementConfirm")
			return
		}
		result.Confirm = &tv
	case ScpStatementTypeExternalize:
		tv, ok := value.(ScpStatementExternalize)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementExternalize")
			return
		}
		result.Externalize = &tv
	case ScpStatementTypeNominate:
		tv, ok := value.(ScpNomination)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpNomination")
			return
		}
		result.Nominate = &tv
	}
	return
}

// MustPrepare retrieves the Prepare value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustPrepare() ScpStatementPrepare {
	val, ok := u.GetPrepare()

	if !ok {
		panic("arm Prepare is not set")
	}

	return val
}

// GetPrepare retrieves the Prepare value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetPrepare() (result ScpStatementPrepare, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Prepare" {
		result = *u.Prepare
		ok = true
	}

	return
}

// MustConfirm retrieves the Confirm value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustConfirm() ScpStatementConfirm {
	val, ok := u.GetConfirm()

	if !ok {
		panic("arm Confirm is not set")
	}

	return val
}

// GetConfirm retrieves the Confirm value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetConfirm() (result ScpStatementConfirm, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Confirm" {
		result = *u.Confirm
		ok = true
	}

	return
}

// MustExternalize retrieves the Externalize value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustExternalize() ScpStatementExternalize {
	val, ok := u.GetExternalize()

	if !ok {
		panic("arm Externalize is not set")
	}

	return val
}

// GetExternalize retrieves the Externalize value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetExternalize() (result ScpStatementExternalize, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Externalize" {
		result = *u.Externalize
		ok = true
	}

	return
}

// MustNominate retrieves the Nominate value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustNominate() ScpNomination {
	val, ok := u.GetNominate()

	if !ok {
		panic("arm Nominate is not set")
	}

	return val
}

// GetNominate retrieves the Nominate value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetNominate() (result ScpNomination, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Nominate" {
		result = *u.Nominate
		ok = true
	}

	return
}

// ScpStatement is an XDR Struct defines as:
//
//   struct SCPStatement
//    {
//        NodeID nodeID;    // v
//        uint64 slotIndex; // i
//
//        union switch (SCPStatementType type)
//        {
//        case PREPARE:
//            struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            } prepare;
//        case CONFIRM:
//            struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            } confirm;
//        case EXTERNALIZE:
//            struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            } externalize;
//        case NOMINATE:
//            SCPNomination nominate;
//        }
//        pledges;
//    };
//
type ScpStatement struct {
	NodeId    NodeId              `json:"nodeID,omitempty"`
	SlotIndex Uint64              `json:"slotIndex,omitempty"`
	Pledges   ScpStatementPledges `json:"pledges,omitempty"`
}

// ScpEnvelope is an XDR Struct defines as:
//
//   struct SCPEnvelope
//    {
//        SCPStatement statement;
//        Signature signature;
//    };
//
type ScpEnvelope struct {
	Statement ScpStatement `json:"statement,omitempty"`
	Signature Signature    `json:"signature,omitempty"`
}

// ScpQuorumSet is an XDR Struct defines as:
//
//   struct SCPQuorumSet
//    {
//        uint32 threshold;
//        PublicKey validators<>;
//        SCPQuorumSet innerSets<>;
//    };
//
type ScpQuorumSet struct {
	Threshold  Uint32         `json:"threshold,omitempty"`
	Validators []PublicKey    `json:"validators,omitempty"`
	InnerSets  []ScpQuorumSet `json:"innerSets,omitempty"`
}

// UpgradeType is an XDR Typedef defines as:
//
//   typedef opaque UpgradeType<128>;
//
type UpgradeType []byte

// StellarValueExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type StellarValueExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u StellarValueExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of StellarValueExt
func (u StellarValueExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewStellarValueExt creates a new  StellarValueExt.
func NewStellarValueExt(v LedgerVersion, value interface{}) (result StellarValueExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// StellarValue is an XDR Struct defines as:
//
//   struct StellarValue
//    {
//        Hash txSetHash;   // transaction set to apply to previous ledger
//        uint64 closeTime; // network close time
//
//        // upgrades to apply to the previous ledger (usually empty)
//        // this is a vector of encoded 'LedgerUpgrade' so that nodes can drop
//        // unknown steps during consensus if needed.
//        // see notes below on 'LedgerUpgrade' for more detail
//        // max size is dictated by number of upgrade types (+ room for future)
//        UpgradeType upgrades<6>;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type StellarValue struct {
	TxSetHash Hash            `json:"txSetHash,omitempty"`
	CloseTime Uint64          `json:"closeTime,omitempty"`
	Upgrades  []UpgradeType   `json:"upgrades,omitempty" xdrmaxsize:"6"`
	Ext       StellarValueExt `json:"ext,omitempty"`
}

// LedgerHeaderExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LedgerHeaderExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderExt
func (u LedgerHeaderExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerHeaderExt creates a new  LedgerHeaderExt.
func NewLedgerHeaderExt(v LedgerVersion, value interface{}) (result LedgerHeaderExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerHeader is an XDR Struct defines as:
//
//   struct LedgerHeader
//    {
//        uint32 ledgerVersion;    // the protocol version of the ledger
//        Hash previousLedgerHash; // hash of the previous ledger header
//        StellarValue scpValue;   // what consensus agreed to
//        Hash txSetResultHash;    // the TransactionResultSet that led to this ledger
//        Hash bucketListHash;     // hash of the ledger state
//
//        uint32 ledgerSeq; // sequence number of this ledger
//
//        uint64 idPool; // last used global ID, used for generating objects
//
//        uint32 baseFee;     // base fee per operation in stroops
//        uint32 baseReserve; // account base reserve in stroops
//
//        uint32 maxTxSetSize; // maximum size a transaction set can be
//
//        PublicKey issuanceKeys<>;
//        int64 txExpirationPeriod;
//
//        Hash skipList[4]; // hashes of ledgers in the past. allows you to jump back
//                          // in time without walking the chain back ledger by ledger
//                          // each slot contains the oldest ledger that is mod of
//                          // either 50  5000  50000 or 500000 depending on index
//                          // skipList[0] mod(50), skipList[1] mod(5000), etc
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type LedgerHeader struct {
	LedgerVersion      Uint32          `json:"ledgerVersion,omitempty"`
	PreviousLedgerHash Hash            `json:"previousLedgerHash,omitempty"`
	ScpValue           StellarValue    `json:"scpValue,omitempty"`
	TxSetResultHash    Hash            `json:"txSetResultHash,omitempty"`
	BucketListHash     Hash            `json:"bucketListHash,omitempty"`
	LedgerSeq          Uint32          `json:"ledgerSeq,omitempty"`
	IdPool             Uint64          `json:"idPool,omitempty"`
	BaseFee            Uint32          `json:"baseFee,omitempty"`
	BaseReserve        Uint32          `json:"baseReserve,omitempty"`
	MaxTxSetSize       Uint32          `json:"maxTxSetSize,omitempty"`
	IssuanceKeys       []PublicKey     `json:"issuanceKeys,omitempty"`
	TxExpirationPeriod Int64           `json:"txExpirationPeriod,omitempty"`
	SkipList           [4]Hash         `json:"skipList,omitempty"`
	Ext                LedgerHeaderExt `json:"ext,omitempty"`
}

// LedgerUpgradeType is an XDR Enum defines as:
//
//   enum LedgerUpgradeType
//    {
//        VERSION = 1,
//        MAX_TX_SET_SIZE = 2,
//        ISSUANCE_KEYS = 3,
//        TX_EXPIRATION_PERIOD = 4
//    };
//
type LedgerUpgradeType int32

const (
	LedgerUpgradeTypeVersion            LedgerUpgradeType = 1
	LedgerUpgradeTypeMaxTxSetSize       LedgerUpgradeType = 2
	LedgerUpgradeTypeIssuanceKeys       LedgerUpgradeType = 3
	LedgerUpgradeTypeTxExpirationPeriod LedgerUpgradeType = 4
)

var LedgerUpgradeTypeAll = []LedgerUpgradeType{
	LedgerUpgradeTypeVersion,
	LedgerUpgradeTypeMaxTxSetSize,
	LedgerUpgradeTypeIssuanceKeys,
	LedgerUpgradeTypeTxExpirationPeriod,
}

var ledgerUpgradeTypeMap = map[int32]string{
	1: "LedgerUpgradeTypeVersion",
	2: "LedgerUpgradeTypeMaxTxSetSize",
	3: "LedgerUpgradeTypeIssuanceKeys",
	4: "LedgerUpgradeTypeTxExpirationPeriod",
}

var ledgerUpgradeTypeShortMap = map[int32]string{
	1: "version",
	2: "max_tx_set_size",
	3: "issuance_keys",
	4: "tx_expiration_period",
}

var ledgerUpgradeTypeRevMap = map[string]int32{
	"LedgerUpgradeTypeVersion":            1,
	"LedgerUpgradeTypeMaxTxSetSize":       2,
	"LedgerUpgradeTypeIssuanceKeys":       3,
	"LedgerUpgradeTypeTxExpirationPeriod": 4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerUpgradeType
func (e LedgerUpgradeType) ValidEnum(v int32) bool {
	_, ok := ledgerUpgradeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerUpgradeType) String() string {
	name, _ := ledgerUpgradeTypeMap[int32(e)]
	return name
}

func (e LedgerUpgradeType) ShortString() string {
	name, _ := ledgerUpgradeTypeShortMap[int32(e)]
	return name
}

func (e LedgerUpgradeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *LedgerUpgradeType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := ledgerUpgradeTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = LedgerUpgradeType(value)
	return nil
}

// LedgerUpgrade is an XDR Union defines as:
//
//   union LedgerUpgrade switch (LedgerUpgradeType type)
//    {
//    case VERSION:
//        uint32 newLedgerVersion; // update ledgerVersion
//    case MAX_TX_SET_SIZE:
//        uint32 newMaxTxSetSize; // update maxTxSetSize
//    case ISSUANCE_KEYS:
//        PublicKey newIssuanceKeys<>;
//    case TX_EXPIRATION_PERIOD:
//        int64 newTxExpirationPeriod;
//    };
//
type LedgerUpgrade struct {
	Type                  LedgerUpgradeType `json:"type,omitempty"`
	NewLedgerVersion      *Uint32           `json:"newLedgerVersion,omitempty"`
	NewMaxTxSetSize       *Uint32           `json:"newMaxTxSetSize,omitempty"`
	NewIssuanceKeys       *[]PublicKey      `json:"newIssuanceKeys,omitempty"`
	NewTxExpirationPeriod *Int64            `json:"newTxExpirationPeriod,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerUpgrade) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerUpgrade
func (u LedgerUpgrade) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerUpgradeType(sw) {
	case LedgerUpgradeTypeVersion:
		return "NewLedgerVersion", true
	case LedgerUpgradeTypeMaxTxSetSize:
		return "NewMaxTxSetSize", true
	case LedgerUpgradeTypeIssuanceKeys:
		return "NewIssuanceKeys", true
	case LedgerUpgradeTypeTxExpirationPeriod:
		return "NewTxExpirationPeriod", true
	}
	return "-", false
}

// NewLedgerUpgrade creates a new  LedgerUpgrade.
func NewLedgerUpgrade(aType LedgerUpgradeType, value interface{}) (result LedgerUpgrade, err error) {
	result.Type = aType
	switch LedgerUpgradeType(aType) {
	case LedgerUpgradeTypeVersion:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewLedgerVersion = &tv
	case LedgerUpgradeTypeMaxTxSetSize:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewMaxTxSetSize = &tv
	case LedgerUpgradeTypeIssuanceKeys:
		tv, ok := value.([]PublicKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be []PublicKey")
			return
		}
		result.NewIssuanceKeys = &tv
	case LedgerUpgradeTypeTxExpirationPeriod:
		tv, ok := value.(Int64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Int64")
			return
		}
		result.NewTxExpirationPeriod = &tv
	}
	return
}

// MustNewLedgerVersion retrieves the NewLedgerVersion value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewLedgerVersion() Uint32 {
	val, ok := u.GetNewLedgerVersion()

	if !ok {
		panic("arm NewLedgerVersion is not set")
	}

	return val
}

// GetNewLedgerVersion retrieves the NewLedgerVersion value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewLedgerVersion() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewLedgerVersion" {
		result = *u.NewLedgerVersion
		ok = true
	}

	return
}

// MustNewMaxTxSetSize retrieves the NewMaxTxSetSize value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewMaxTxSetSize() Uint32 {
	val, ok := u.GetNewMaxTxSetSize()

	if !ok {
		panic("arm NewMaxTxSetSize is not set")
	}

	return val
}

// GetNewMaxTxSetSize retrieves the NewMaxTxSetSize value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewMaxTxSetSize() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewMaxTxSetSize" {
		result = *u.NewMaxTxSetSize
		ok = true
	}

	return
}

// MustNewIssuanceKeys retrieves the NewIssuanceKeys value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewIssuanceKeys() []PublicKey {
	val, ok := u.GetNewIssuanceKeys()

	if !ok {
		panic("arm NewIssuanceKeys is not set")
	}

	return val
}

// GetNewIssuanceKeys retrieves the NewIssuanceKeys value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewIssuanceKeys() (result []PublicKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewIssuanceKeys" {
		result = *u.NewIssuanceKeys
		ok = true
	}

	return
}

// MustNewTxExpirationPeriod retrieves the NewTxExpirationPeriod value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewTxExpirationPeriod() Int64 {
	val, ok := u.GetNewTxExpirationPeriod()

	if !ok {
		panic("arm NewTxExpirationPeriod is not set")
	}

	return val
}

// GetNewTxExpirationPeriod retrieves the NewTxExpirationPeriod value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewTxExpirationPeriod() (result Int64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewTxExpirationPeriod" {
		result = *u.NewTxExpirationPeriod
		ok = true
	}

	return
}

// LedgerKeyAccountExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAccountExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAccountExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAccountExt
func (u LedgerKeyAccountExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAccountExt creates a new  LedgerKeyAccountExt.
func NewLedgerKeyAccountExt(v LedgerVersion, value interface{}) (result LedgerKeyAccountExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAccount is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAccount struct {
	AccountId AccountId           `json:"accountID,omitempty"`
	Ext       LedgerKeyAccountExt `json:"ext,omitempty"`
}

// LedgerKeyCoinsEmissionRequestExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyCoinsEmissionRequestExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyCoinsEmissionRequestExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyCoinsEmissionRequestExt
func (u LedgerKeyCoinsEmissionRequestExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyCoinsEmissionRequestExt creates a new  LedgerKeyCoinsEmissionRequestExt.
func NewLedgerKeyCoinsEmissionRequestExt(v LedgerVersion, value interface{}) (result LedgerKeyCoinsEmissionRequestExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyCoinsEmissionRequest is an XDR NestedStruct defines as:
//
//   struct
//    	{
//    		AccountID issuer;
//    		uint64 requestID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type LedgerKeyCoinsEmissionRequest struct {
	Issuer    AccountId                        `json:"issuer,omitempty"`
	RequestId Uint64                           `json:"requestID,omitempty"`
	Ext       LedgerKeyCoinsEmissionRequestExt `json:"ext,omitempty"`
}

// LedgerKeyCoinsEmissionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyCoinsEmissionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyCoinsEmissionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyCoinsEmissionExt
func (u LedgerKeyCoinsEmissionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyCoinsEmissionExt creates a new  LedgerKeyCoinsEmissionExt.
func NewLedgerKeyCoinsEmissionExt(v LedgerVersion, value interface{}) (result LedgerKeyCoinsEmissionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyCoinsEmission is an XDR NestedStruct defines as:
//
//   struct
//    	{
//    		string64 serialNumber;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type LedgerKeyCoinsEmission struct {
	SerialNumber String64                  `json:"serialNumber,omitempty"`
	Ext          LedgerKeyCoinsEmissionExt `json:"ext,omitempty"`
}

// LedgerKeyFeeStateExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyFeeStateExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyFeeStateExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyFeeStateExt
func (u LedgerKeyFeeStateExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyFeeStateExt creates a new  LedgerKeyFeeStateExt.
func NewLedgerKeyFeeStateExt(v LedgerVersion, value interface{}) (result LedgerKeyFeeStateExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyFeeState is an XDR NestedStruct defines as:
//
//   struct {
//            Hash hash;
//    		int64 lowerBound;
//    		int64 upperBound;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyFeeState struct {
	Hash       Hash                 `json:"hash,omitempty"`
	LowerBound Int64                `json:"lowerBound,omitempty"`
	UpperBound Int64                `json:"upperBound,omitempty"`
	Ext        LedgerKeyFeeStateExt `json:"ext,omitempty"`
}

// LedgerKeyBalanceExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyBalanceExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyBalanceExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyBalanceExt
func (u LedgerKeyBalanceExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyBalanceExt creates a new  LedgerKeyBalanceExt.
func NewLedgerKeyBalanceExt(v LedgerVersion, value interface{}) (result LedgerKeyBalanceExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyBalance is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		BalanceID balanceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyBalance struct {
	BalanceId BalanceId           `json:"balanceID,omitempty"`
	Ext       LedgerKeyBalanceExt `json:"ext,omitempty"`
}

// LedgerKeyPaymentRequestExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyPaymentRequestExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyPaymentRequestExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyPaymentRequestExt
func (u LedgerKeyPaymentRequestExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyPaymentRequestExt creates a new  LedgerKeyPaymentRequestExt.
func NewLedgerKeyPaymentRequestExt(v LedgerVersion, value interface{}) (result LedgerKeyPaymentRequestExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyPaymentRequest is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		uint64 paymentID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyPaymentRequest struct {
	PaymentId Uint64                     `json:"paymentID,omitempty"`
	Ext       LedgerKeyPaymentRequestExt `json:"ext,omitempty"`
}

// LedgerKeyAssetExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAssetExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAssetExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAssetExt
func (u LedgerKeyAssetExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAssetExt creates a new  LedgerKeyAssetExt.
func NewLedgerKeyAssetExt(v LedgerVersion, value interface{}) (result LedgerKeyAssetExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAsset is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		AssetCode code;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAsset struct {
	Code AssetCode         `json:"code,omitempty"`
	Ext  LedgerKeyAssetExt `json:"ext,omitempty"`
}

// LedgerKeyPaymentExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyPaymentExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyPaymentExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyPaymentExt
func (u LedgerKeyPaymentExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyPaymentExt creates a new  LedgerKeyPaymentExt.
func NewLedgerKeyPaymentExt(v LedgerVersion, value interface{}) (result LedgerKeyPaymentExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyPayment is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		AccountID sender;
//    		string64 reference;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyPayment struct {
	Sender    AccountId           `json:"sender,omitempty"`
	Reference String64            `json:"reference,omitempty"`
	Ext       LedgerKeyPaymentExt `json:"ext,omitempty"`
}

// LedgerKeyAccountTypeLimitsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAccountTypeLimitsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAccountTypeLimitsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAccountTypeLimitsExt
func (u LedgerKeyAccountTypeLimitsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAccountTypeLimitsExt creates a new  LedgerKeyAccountTypeLimitsExt.
func NewLedgerKeyAccountTypeLimitsExt(v LedgerVersion, value interface{}) (result LedgerKeyAccountTypeLimitsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAccountTypeLimits is an XDR NestedStruct defines as:
//
//   struct {
//            AccountType accountType;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAccountTypeLimits struct {
	AccountType AccountType                   `json:"accountType,omitempty"`
	Ext         LedgerKeyAccountTypeLimitsExt `json:"ext,omitempty"`
}

// LedgerKeyStatsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyStatsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyStatsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyStatsExt
func (u LedgerKeyStatsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyStatsExt creates a new  LedgerKeyStatsExt.
func NewLedgerKeyStatsExt(v LedgerVersion, value interface{}) (result LedgerKeyStatsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyStats is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyStats struct {
	AccountId AccountId         `json:"accountID,omitempty"`
	Ext       LedgerKeyStatsExt `json:"ext,omitempty"`
}

// LedgerKeyTrustExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyTrustExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyTrustExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyTrustExt
func (u LedgerKeyTrustExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyTrustExt creates a new  LedgerKeyTrustExt.
func NewLedgerKeyTrustExt(v LedgerVersion, value interface{}) (result LedgerKeyTrustExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyTrust is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID allowedAccount;
//            BalanceID balanceToUse;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyTrust struct {
	AllowedAccount AccountId         `json:"allowedAccount,omitempty"`
	BalanceToUse   BalanceId         `json:"balanceToUse,omitempty"`
	Ext            LedgerKeyTrustExt `json:"ext,omitempty"`
}

// LedgerKeyAccountLimitsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAccountLimitsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAccountLimitsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAccountLimitsExt
func (u LedgerKeyAccountLimitsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAccountLimitsExt creates a new  LedgerKeyAccountLimitsExt.
func NewLedgerKeyAccountLimitsExt(v LedgerVersion, value interface{}) (result LedgerKeyAccountLimitsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAccountLimits is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAccountLimits struct {
	AccountId AccountId                 `json:"accountID,omitempty"`
	Ext       LedgerKeyAccountLimitsExt `json:"ext,omitempty"`
}

// LedgerKeyAssetPairExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAssetPairExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAssetPairExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAssetPairExt
func (u LedgerKeyAssetPairExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAssetPairExt creates a new  LedgerKeyAssetPairExt.
func NewLedgerKeyAssetPairExt(v LedgerVersion, value interface{}) (result LedgerKeyAssetPairExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAssetPair is an XDR NestedStruct defines as:
//
//   struct {
//             AssetCode base;
//    		 AssetCode quote;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAssetPair struct {
	Base  AssetCode             `json:"base,omitempty"`
	Quote AssetCode             `json:"quote,omitempty"`
	Ext   LedgerKeyAssetPairExt `json:"ext,omitempty"`
}

// LedgerKeyOffer is an XDR NestedStruct defines as:
//
//   struct {
//    		uint64 offerID;
//    		AccountID ownerID;
//    	}
//
type LedgerKeyOffer struct {
	OfferId Uint64    `json:"offerID,omitempty"`
	OwnerId AccountId `json:"ownerID,omitempty"`
}

// LedgerKeyInvoiceExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyInvoiceExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyInvoiceExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyInvoiceExt
func (u LedgerKeyInvoiceExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyInvoiceExt creates a new  LedgerKeyInvoiceExt.
func NewLedgerKeyInvoiceExt(v LedgerVersion, value interface{}) (result LedgerKeyInvoiceExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyInvoice is an XDR NestedStruct defines as:
//
//   struct {
//            uint64 invoiceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyInvoice struct {
	InvoiceId Uint64              `json:"invoiceID,omitempty"`
	Ext       LedgerKeyInvoiceExt `json:"ext,omitempty"`
}

// LedgerKey is an XDR Union defines as:
//
//   union LedgerKey switch (LedgerEntryType type)
//    {
//    case ACCOUNT:
//        struct
//        {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } account;
//    case COINS_EMISSION_REQUEST:
//    	struct
//    	{
//    		AccountID issuer;
//    		uint64 requestID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} coinsEmissionRequest;
//    case COINS_EMISSION:
//    	struct
//    	{
//    		string64 serialNumber;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} coinsEmission;
//    case FEE:
//        struct {
//            Hash hash;
//    		int64 lowerBound;
//    		int64 upperBound;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } feeState;
//    case BALANCE:
//        struct
//        {
//    		BalanceID balanceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } balance;
//    case PAYMENT_REQUEST:
//        struct
//        {
//    		uint64 paymentID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } paymentRequest;
//    case ASSET:
//        struct
//        {
//    		AssetCode code;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } asset;
//    case REFERENCE_ENTRY:
//        struct
//        {
//    		AccountID sender;
//    		string64 reference;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } payment;
//    case ACCOUNT_TYPE_LIMITS:
//        struct {
//            AccountType accountType;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } accountTypeLimits;
//    case STATISTICS:
//        struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } stats;
//    case TRUST:
//        struct {
//            AccountID allowedAccount;
//            BalanceID balanceToUse;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } trust;
//    case ACCOUNT_LIMITS:
//        struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } accountLimits;
//    case ASSET_PAIR:
//    	struct {
//             AssetCode base;
//    		 AssetCode quote;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } assetPair;
//    case OFFER_ENTRY:
//    	struct {
//    		uint64 offerID;
//    		AccountID ownerID;
//    	} offer;
//    case INVOICE:
//        struct {
//            uint64 invoiceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } invoice;
//    };
//
type LedgerKey struct {
	Type                 LedgerEntryType                `json:"type,omitempty"`
	Account              *LedgerKeyAccount              `json:"account,omitempty"`
	CoinsEmissionRequest *LedgerKeyCoinsEmissionRequest `json:"coinsEmissionRequest,omitempty"`
	CoinsEmission        *LedgerKeyCoinsEmission        `json:"coinsEmission,omitempty"`
	FeeState             *LedgerKeyFeeState             `json:"feeState,omitempty"`
	Balance              *LedgerKeyBalance              `json:"balance,omitempty"`
	PaymentRequest       *LedgerKeyPaymentRequest       `json:"paymentRequest,omitempty"`
	Asset                *LedgerKeyAsset                `json:"asset,omitempty"`
	Payment              *LedgerKeyPayment              `json:"payment,omitempty"`
	AccountTypeLimits    *LedgerKeyAccountTypeLimits    `json:"accountTypeLimits,omitempty"`
	Stats                *LedgerKeyStats                `json:"stats,omitempty"`
	Trust                *LedgerKeyTrust                `json:"trust,omitempty"`
	AccountLimits        *LedgerKeyAccountLimits        `json:"accountLimits,omitempty"`
	AssetPair            *LedgerKeyAssetPair            `json:"assetPair,omitempty"`
	Offer                *LedgerKeyOffer                `json:"offer,omitempty"`
	Invoice              *LedgerKeyInvoice              `json:"invoice,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKey
func (u LedgerKey) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeCoinsEmissionRequest:
		return "CoinsEmissionRequest", true
	case LedgerEntryTypeCoinsEmission:
		return "CoinsEmission", true
	case LedgerEntryTypeFee:
		return "FeeState", true
	case LedgerEntryTypeBalance:
		return "Balance", true
	case LedgerEntryTypePaymentRequest:
		return "PaymentRequest", true
	case LedgerEntryTypeAsset:
		return "Asset", true
	case LedgerEntryTypeReferenceEntry:
		return "Payment", true
	case LedgerEntryTypeAccountTypeLimits:
		return "AccountTypeLimits", true
	case LedgerEntryTypeStatistics:
		return "Stats", true
	case LedgerEntryTypeTrust:
		return "Trust", true
	case LedgerEntryTypeAccountLimits:
		return "AccountLimits", true
	case LedgerEntryTypeAssetPair:
		return "AssetPair", true
	case LedgerEntryTypeOfferEntry:
		return "Offer", true
	case LedgerEntryTypeInvoice:
		return "Invoice", true
	}
	return "-", false
}

// NewLedgerKey creates a new  LedgerKey.
func NewLedgerKey(aType LedgerEntryType, value interface{}) (result LedgerKey, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(LedgerKeyAccount)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccount")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeCoinsEmissionRequest:
		tv, ok := value.(LedgerKeyCoinsEmissionRequest)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyCoinsEmissionRequest")
			return
		}
		result.CoinsEmissionRequest = &tv
	case LedgerEntryTypeCoinsEmission:
		tv, ok := value.(LedgerKeyCoinsEmission)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyCoinsEmission")
			return
		}
		result.CoinsEmission = &tv
	case LedgerEntryTypeFee:
		tv, ok := value.(LedgerKeyFeeState)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyFeeState")
			return
		}
		result.FeeState = &tv
	case LedgerEntryTypeBalance:
		tv, ok := value.(LedgerKeyBalance)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyBalance")
			return
		}
		result.Balance = &tv
	case LedgerEntryTypePaymentRequest:
		tv, ok := value.(LedgerKeyPaymentRequest)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyPaymentRequest")
			return
		}
		result.PaymentRequest = &tv
	case LedgerEntryTypeAsset:
		tv, ok := value.(LedgerKeyAsset)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAsset")
			return
		}
		result.Asset = &tv
	case LedgerEntryTypeReferenceEntry:
		tv, ok := value.(LedgerKeyPayment)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyPayment")
			return
		}
		result.Payment = &tv
	case LedgerEntryTypeAccountTypeLimits:
		tv, ok := value.(LedgerKeyAccountTypeLimits)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccountTypeLimits")
			return
		}
		result.AccountTypeLimits = &tv
	case LedgerEntryTypeStatistics:
		tv, ok := value.(LedgerKeyStats)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyStats")
			return
		}
		result.Stats = &tv
	case LedgerEntryTypeTrust:
		tv, ok := value.(LedgerKeyTrust)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyTrust")
			return
		}
		result.Trust = &tv
	case LedgerEntryTypeAccountLimits:
		tv, ok := value.(LedgerKeyAccountLimits)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccountLimits")
			return
		}
		result.AccountLimits = &tv
	case LedgerEntryTypeAssetPair:
		tv, ok := value.(LedgerKeyAssetPair)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAssetPair")
			return
		}
		result.AssetPair = &tv
	case LedgerEntryTypeOfferEntry:
		tv, ok := value.(LedgerKeyOffer)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyOffer")
			return
		}
		result.Offer = &tv
	case LedgerEntryTypeInvoice:
		tv, ok := value.(LedgerKeyInvoice)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyInvoice")
			return
		}
		result.Invoice = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccount() LedgerKeyAccount {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccount() (result LedgerKeyAccount, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustCoinsEmissionRequest() LedgerKeyCoinsEmissionRequest {
	val, ok := u.GetCoinsEmissionRequest()

	if !ok {
		panic("arm CoinsEmissionRequest is not set")
	}

	return val
}

// GetCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetCoinsEmissionRequest() (result LedgerKeyCoinsEmissionRequest, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmissionRequest" {
		result = *u.CoinsEmissionRequest
		ok = true
	}

	return
}

// MustCoinsEmission retrieves the CoinsEmission value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustCoinsEmission() LedgerKeyCoinsEmission {
	val, ok := u.GetCoinsEmission()

	if !ok {
		panic("arm CoinsEmission is not set")
	}

	return val
}

// GetCoinsEmission retrieves the CoinsEmission value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetCoinsEmission() (result LedgerKeyCoinsEmission, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmission" {
		result = *u.CoinsEmission
		ok = true
	}

	return
}

// MustFeeState retrieves the FeeState value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustFeeState() LedgerKeyFeeState {
	val, ok := u.GetFeeState()

	if !ok {
		panic("arm FeeState is not set")
	}

	return val
}

// GetFeeState retrieves the FeeState value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetFeeState() (result LedgerKeyFeeState, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeState" {
		result = *u.FeeState
		ok = true
	}

	return
}

// MustBalance retrieves the Balance value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustBalance() LedgerKeyBalance {
	val, ok := u.GetBalance()

	if !ok {
		panic("arm Balance is not set")
	}

	return val
}

// GetBalance retrieves the Balance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetBalance() (result LedgerKeyBalance, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Balance" {
		result = *u.Balance
		ok = true
	}

	return
}

// MustPaymentRequest retrieves the PaymentRequest value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustPaymentRequest() LedgerKeyPaymentRequest {
	val, ok := u.GetPaymentRequest()

	if !ok {
		panic("arm PaymentRequest is not set")
	}

	return val
}

// GetPaymentRequest retrieves the PaymentRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetPaymentRequest() (result LedgerKeyPaymentRequest, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentRequest" {
		result = *u.PaymentRequest
		ok = true
	}

	return
}

// MustAsset retrieves the Asset value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAsset() LedgerKeyAsset {
	val, ok := u.GetAsset()

	if !ok {
		panic("arm Asset is not set")
	}

	return val
}

// GetAsset retrieves the Asset value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAsset() (result LedgerKeyAsset, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Asset" {
		result = *u.Asset
		ok = true
	}

	return
}

// MustPayment retrieves the Payment value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustPayment() LedgerKeyPayment {
	val, ok := u.GetPayment()

	if !ok {
		panic("arm Payment is not set")
	}

	return val
}

// GetPayment retrieves the Payment value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetPayment() (result LedgerKeyPayment, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Payment" {
		result = *u.Payment
		ok = true
	}

	return
}

// MustAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccountTypeLimits() LedgerKeyAccountTypeLimits {
	val, ok := u.GetAccountTypeLimits()

	if !ok {
		panic("arm AccountTypeLimits is not set")
	}

	return val
}

// GetAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccountTypeLimits() (result LedgerKeyAccountTypeLimits, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountTypeLimits" {
		result = *u.AccountTypeLimits
		ok = true
	}

	return
}

// MustStats retrieves the Stats value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustStats() LedgerKeyStats {
	val, ok := u.GetStats()

	if !ok {
		panic("arm Stats is not set")
	}

	return val
}

// GetStats retrieves the Stats value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetStats() (result LedgerKeyStats, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Stats" {
		result = *u.Stats
		ok = true
	}

	return
}

// MustTrust retrieves the Trust value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustTrust() LedgerKeyTrust {
	val, ok := u.GetTrust()

	if !ok {
		panic("arm Trust is not set")
	}

	return val
}

// GetTrust retrieves the Trust value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetTrust() (result LedgerKeyTrust, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Trust" {
		result = *u.Trust
		ok = true
	}

	return
}

// MustAccountLimits retrieves the AccountLimits value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccountLimits() LedgerKeyAccountLimits {
	val, ok := u.GetAccountLimits()

	if !ok {
		panic("arm AccountLimits is not set")
	}

	return val
}

// GetAccountLimits retrieves the AccountLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccountLimits() (result LedgerKeyAccountLimits, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountLimits" {
		result = *u.AccountLimits
		ok = true
	}

	return
}

// MustAssetPair retrieves the AssetPair value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAssetPair() LedgerKeyAssetPair {
	val, ok := u.GetAssetPair()

	if !ok {
		panic("arm AssetPair is not set")
	}

	return val
}

// GetAssetPair retrieves the AssetPair value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAssetPair() (result LedgerKeyAssetPair, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetPair" {
		result = *u.AssetPair
		ok = true
	}

	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustOffer() LedgerKeyOffer {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetOffer() (result LedgerKeyOffer, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// MustInvoice retrieves the Invoice value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustInvoice() LedgerKeyInvoice {
	val, ok := u.GetInvoice()

	if !ok {
		panic("arm Invoice is not set")
	}

	return val
}

// GetInvoice retrieves the Invoice value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetInvoice() (result LedgerKeyInvoice, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Invoice" {
		result = *u.Invoice
		ok = true
	}

	return
}

// BucketEntryType is an XDR Enum defines as:
//
//   enum BucketEntryType
//    {
//        LIVEENTRY = 0,
//        DEADENTRY = 1
//    };
//
type BucketEntryType int32

const (
	BucketEntryTypeLiveentry BucketEntryType = 0
	BucketEntryTypeDeadentry BucketEntryType = 1
)

var BucketEntryTypeAll = []BucketEntryType{
	BucketEntryTypeLiveentry,
	BucketEntryTypeDeadentry,
}

var bucketEntryTypeMap = map[int32]string{
	0: "BucketEntryTypeLiveentry",
	1: "BucketEntryTypeDeadentry",
}

var bucketEntryTypeShortMap = map[int32]string{
	0: "liveentry",
	1: "deadentry",
}

var bucketEntryTypeRevMap = map[string]int32{
	"BucketEntryTypeLiveentry": 0,
	"BucketEntryTypeDeadentry": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BucketEntryType
func (e BucketEntryType) ValidEnum(v int32) bool {
	_, ok := bucketEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e BucketEntryType) String() string {
	name, _ := bucketEntryTypeMap[int32(e)]
	return name
}

func (e BucketEntryType) ShortString() string {
	name, _ := bucketEntryTypeShortMap[int32(e)]
	return name
}

func (e BucketEntryType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *BucketEntryType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := bucketEntryTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = BucketEntryType(value)
	return nil
}

// BucketEntry is an XDR Union defines as:
//
//   union BucketEntry switch (BucketEntryType type)
//    {
//    case LIVEENTRY:
//        LedgerEntry liveEntry;
//
//    case DEADENTRY:
//        LedgerKey deadEntry;
//    };
//
type BucketEntry struct {
	Type      BucketEntryType `json:"type,omitempty"`
	LiveEntry *LedgerEntry    `json:"liveEntry,omitempty"`
	DeadEntry *LedgerKey      `json:"deadEntry,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BucketEntry) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BucketEntry
func (u BucketEntry) ArmForSwitch(sw int32) (string, bool) {
	switch BucketEntryType(sw) {
	case BucketEntryTypeLiveentry:
		return "LiveEntry", true
	case BucketEntryTypeDeadentry:
		return "DeadEntry", true
	}
	return "-", false
}

// NewBucketEntry creates a new  BucketEntry.
func NewBucketEntry(aType BucketEntryType, value interface{}) (result BucketEntry, err error) {
	result.Type = aType
	switch BucketEntryType(aType) {
	case BucketEntryTypeLiveentry:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.LiveEntry = &tv
	case BucketEntryTypeDeadentry:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.DeadEntry = &tv
	}
	return
}

// MustLiveEntry retrieves the LiveEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustLiveEntry() LedgerEntry {
	val, ok := u.GetLiveEntry()

	if !ok {
		panic("arm LiveEntry is not set")
	}

	return val
}

// GetLiveEntry retrieves the LiveEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetLiveEntry() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiveEntry" {
		result = *u.LiveEntry
		ok = true
	}

	return
}

// MustDeadEntry retrieves the DeadEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustDeadEntry() LedgerKey {
	val, ok := u.GetDeadEntry()

	if !ok {
		panic("arm DeadEntry is not set")
	}

	return val
}

// GetDeadEntry retrieves the DeadEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetDeadEntry() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DeadEntry" {
		result = *u.DeadEntry
		ok = true
	}

	return
}

// TransactionSet is an XDR Struct defines as:
//
//   struct TransactionSet
//    {
//        Hash previousLedgerHash;
//        TransactionEnvelope txs<>;
//    };
//
type TransactionSet struct {
	PreviousLedgerHash Hash                  `json:"previousLedgerHash,omitempty"`
	Txs                []TransactionEnvelope `json:"txs,omitempty"`
}

// TransactionResultPair is an XDR Struct defines as:
//
//   struct TransactionResultPair
//    {
//        Hash transactionHash;
//        TransactionResult result; // result for the transaction
//    };
//
type TransactionResultPair struct {
	TransactionHash Hash              `json:"transactionHash,omitempty"`
	Result          TransactionResult `json:"result,omitempty"`
}

// TransactionResultSet is an XDR Struct defines as:
//
//   struct TransactionResultSet
//    {
//        TransactionResultPair results<>;
//    };
//
type TransactionResultSet struct {
	Results []TransactionResultPair `json:"results,omitempty"`
}

// TransactionHistoryEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionHistoryEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionHistoryEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionHistoryEntryExt
func (u TransactionHistoryEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionHistoryEntryExt creates a new  TransactionHistoryEntryExt.
func NewTransactionHistoryEntryExt(v LedgerVersion, value interface{}) (result TransactionHistoryEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TransactionHistoryEntry is an XDR Struct defines as:
//
//   struct TransactionHistoryEntry
//    {
//        uint32 ledgerSeq;
//        TransactionSet txSet;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TransactionHistoryEntry struct {
	LedgerSeq Uint32                     `json:"ledgerSeq,omitempty"`
	TxSet     TransactionSet             `json:"txSet,omitempty"`
	Ext       TransactionHistoryEntryExt `json:"ext,omitempty"`
}

// TransactionHistoryResultEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionHistoryResultEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionHistoryResultEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionHistoryResultEntryExt
func (u TransactionHistoryResultEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionHistoryResultEntryExt creates a new  TransactionHistoryResultEntryExt.
func NewTransactionHistoryResultEntryExt(v LedgerVersion, value interface{}) (result TransactionHistoryResultEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TransactionHistoryResultEntry is an XDR Struct defines as:
//
//   struct TransactionHistoryResultEntry
//    {
//        uint32 ledgerSeq;
//        TransactionResultSet txResultSet;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TransactionHistoryResultEntry struct {
	LedgerSeq   Uint32                           `json:"ledgerSeq,omitempty"`
	TxResultSet TransactionResultSet             `json:"txResultSet,omitempty"`
	Ext         TransactionHistoryResultEntryExt `json:"ext,omitempty"`
}

// LedgerHeaderHistoryEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LedgerHeaderHistoryEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderHistoryEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderHistoryEntryExt
func (u LedgerHeaderHistoryEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerHeaderHistoryEntryExt creates a new  LedgerHeaderHistoryEntryExt.
func NewLedgerHeaderHistoryEntryExt(v LedgerVersion, value interface{}) (result LedgerHeaderHistoryEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerHeaderHistoryEntry is an XDR Struct defines as:
//
//   struct LedgerHeaderHistoryEntry
//    {
//        Hash hash;
//        LedgerHeader header;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type LedgerHeaderHistoryEntry struct {
	Hash   Hash                        `json:"hash,omitempty"`
	Header LedgerHeader                `json:"header,omitempty"`
	Ext    LedgerHeaderHistoryEntryExt `json:"ext,omitempty"`
}

// LedgerScpMessages is an XDR Struct defines as:
//
//   struct LedgerSCPMessages
//    {
//        uint32 ledgerSeq;
//        SCPEnvelope messages<>;
//    };
//
type LedgerScpMessages struct {
	LedgerSeq Uint32        `json:"ledgerSeq,omitempty"`
	Messages  []ScpEnvelope `json:"messages,omitempty"`
}

// ScpHistoryEntryV0 is an XDR Struct defines as:
//
//   struct SCPHistoryEntryV0
//    {
//        SCPQuorumSet quorumSets<>; // additional quorum sets used by ledgerMessages
//        LedgerSCPMessages ledgerMessages;
//    };
//
type ScpHistoryEntryV0 struct {
	QuorumSets     []ScpQuorumSet    `json:"quorumSets,omitempty"`
	LedgerMessages LedgerScpMessages `json:"ledgerMessages,omitempty"`
}

// ScpHistoryEntry is an XDR Union defines as:
//
//   union SCPHistoryEntry switch (LedgerVersion v)
//    {
//    case EMPTY_VERSION:
//        SCPHistoryEntryV0 v0;
//    };
//
type ScpHistoryEntry struct {
	V  LedgerVersion      `json:"v,omitempty"`
	V0 *ScpHistoryEntryV0 `json:"v0,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ScpHistoryEntry) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ScpHistoryEntry
func (u ScpHistoryEntry) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "V0", true
	}
	return "-", false
}

// NewScpHistoryEntry creates a new  ScpHistoryEntry.
func NewScpHistoryEntry(v LedgerVersion, value interface{}) (result ScpHistoryEntry, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		tv, ok := value.(ScpHistoryEntryV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpHistoryEntryV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u ScpHistoryEntry) MustV0() ScpHistoryEntryV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpHistoryEntry) GetV0() (result ScpHistoryEntryV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// LedgerEntryChangeType is an XDR Enum defines as:
//
//   enum LedgerEntryChangeType
//    {
//        CREATED = 0, // entry was added to the ledger
//        UPDATED = 1, // entry was modified in the ledger
//        REMOVED = 2, // entry was removed from the ledger
//        STATE = 3    // value of the entry
//    };
//
type LedgerEntryChangeType int32

const (
	LedgerEntryChangeTypeCreated LedgerEntryChangeType = 0
	LedgerEntryChangeTypeUpdated LedgerEntryChangeType = 1
	LedgerEntryChangeTypeRemoved LedgerEntryChangeType = 2
	LedgerEntryChangeTypeState   LedgerEntryChangeType = 3
)

var LedgerEntryChangeTypeAll = []LedgerEntryChangeType{
	LedgerEntryChangeTypeCreated,
	LedgerEntryChangeTypeUpdated,
	LedgerEntryChangeTypeRemoved,
	LedgerEntryChangeTypeState,
}

var ledgerEntryChangeTypeMap = map[int32]string{
	0: "LedgerEntryChangeTypeCreated",
	1: "LedgerEntryChangeTypeUpdated",
	2: "LedgerEntryChangeTypeRemoved",
	3: "LedgerEntryChangeTypeState",
}

var ledgerEntryChangeTypeShortMap = map[int32]string{
	0: "created",
	1: "updated",
	2: "removed",
	3: "state",
}

var ledgerEntryChangeTypeRevMap = map[string]int32{
	"LedgerEntryChangeTypeCreated": 0,
	"LedgerEntryChangeTypeUpdated": 1,
	"LedgerEntryChangeTypeRemoved": 2,
	"LedgerEntryChangeTypeState":   3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryChangeType
func (e LedgerEntryChangeType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryChangeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryChangeType) String() string {
	name, _ := ledgerEntryChangeTypeMap[int32(e)]
	return name
}

func (e LedgerEntryChangeType) ShortString() string {
	name, _ := ledgerEntryChangeTypeShortMap[int32(e)]
	return name
}

func (e LedgerEntryChangeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *LedgerEntryChangeType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := ledgerEntryChangeTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = LedgerEntryChangeType(value)
	return nil
}

// LedgerEntryChange is an XDR Union defines as:
//
//   union LedgerEntryChange switch (LedgerEntryChangeType type)
//    {
//    case CREATED:
//        LedgerEntry created;
//    case UPDATED:
//        LedgerEntry updated;
//    case REMOVED:
//        LedgerKey removed;
//    case STATE:
//        LedgerEntry state;
//    };
//
type LedgerEntryChange struct {
	Type    LedgerEntryChangeType `json:"type,omitempty"`
	Created *LedgerEntry          `json:"created,omitempty"`
	Updated *LedgerEntry          `json:"updated,omitempty"`
	Removed *LedgerKey            `json:"removed,omitempty"`
	State   *LedgerEntry          `json:"state,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryChange) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryChange
func (u LedgerEntryChange) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryChangeType(sw) {
	case LedgerEntryChangeTypeCreated:
		return "Created", true
	case LedgerEntryChangeTypeUpdated:
		return "Updated", true
	case LedgerEntryChangeTypeRemoved:
		return "Removed", true
	case LedgerEntryChangeTypeState:
		return "State", true
	}
	return "-", false
}

// NewLedgerEntryChange creates a new  LedgerEntryChange.
func NewLedgerEntryChange(aType LedgerEntryChangeType, value interface{}) (result LedgerEntryChange, err error) {
	result.Type = aType
	switch LedgerEntryChangeType(aType) {
	case LedgerEntryChangeTypeCreated:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.Created = &tv
	case LedgerEntryChangeTypeUpdated:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.Updated = &tv
	case LedgerEntryChangeTypeRemoved:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.Removed = &tv
	case LedgerEntryChangeTypeState:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.State = &tv
	}
	return
}

// MustCreated retrieves the Created value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustCreated() LedgerEntry {
	val, ok := u.GetCreated()

	if !ok {
		panic("arm Created is not set")
	}

	return val
}

// GetCreated retrieves the Created value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetCreated() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Created" {
		result = *u.Created
		ok = true
	}

	return
}

// MustUpdated retrieves the Updated value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustUpdated() LedgerEntry {
	val, ok := u.GetUpdated()

	if !ok {
		panic("arm Updated is not set")
	}

	return val
}

// GetUpdated retrieves the Updated value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetUpdated() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Updated" {
		result = *u.Updated
		ok = true
	}

	return
}

// MustRemoved retrieves the Removed value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustRemoved() LedgerKey {
	val, ok := u.GetRemoved()

	if !ok {
		panic("arm Removed is not set")
	}

	return val
}

// GetRemoved retrieves the Removed value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetRemoved() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Removed" {
		result = *u.Removed
		ok = true
	}

	return
}

// MustState retrieves the State value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustState() LedgerEntry {
	val, ok := u.GetState()

	if !ok {
		panic("arm State is not set")
	}

	return val
}

// GetState retrieves the State value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetState() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "State" {
		result = *u.State
		ok = true
	}

	return
}

// LedgerEntryChanges is an XDR Typedef defines as:
//
//   typedef LedgerEntryChange LedgerEntryChanges<>;
//
type LedgerEntryChanges []LedgerEntryChange

// OperationMeta is an XDR Struct defines as:
//
//   struct OperationMeta
//    {
//        LedgerEntryChanges changes;
//    };
//
type OperationMeta struct {
	Changes LedgerEntryChanges `json:"changes,omitempty"`
}

// TransactionMeta is an XDR Union defines as:
//
//   union TransactionMeta switch (LedgerVersion v)
//    {
//    case EMPTY_VERSION:
//        OperationMeta operations<>;
//    };
//
type TransactionMeta struct {
	V          LedgerVersion    `json:"v,omitempty"`
	Operations *[]OperationMeta `json:"operations,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionMeta) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionMeta
func (u TransactionMeta) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "Operations", true
	}
	return "-", false
}

// NewTransactionMeta creates a new  TransactionMeta.
func NewTransactionMeta(v LedgerVersion, value interface{}) (result TransactionMeta, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		tv, ok := value.([]OperationMeta)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationMeta")
			return
		}
		result.Operations = &tv
	}
	return
}

// MustOperations retrieves the Operations value from the union,
// panicing if the value is not set.
func (u TransactionMeta) MustOperations() []OperationMeta {
	val, ok := u.GetOperations()

	if !ok {
		panic("arm Operations is not set")
	}

	return val
}

// GetOperations retrieves the Operations value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionMeta) GetOperations() (result []OperationMeta, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Operations" {
		result = *u.Operations
		ok = true
	}

	return
}

// BalanceEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type BalanceEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BalanceEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BalanceEntryExt
func (u BalanceEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewBalanceEntryExt creates a new  BalanceEntryExt.
func NewBalanceEntryExt(v LedgerVersion, value interface{}) (result BalanceEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// BalanceEntry is an XDR Struct defines as:
//
//   struct BalanceEntry
//    {
//        BalanceID balanceID;
//        AssetCode asset;
//        AccountID accountID;
//        int64 amount;
//        int64 locked;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type BalanceEntry struct {
	BalanceId BalanceId       `json:"balanceID,omitempty"`
	Asset     AssetCode       `json:"asset,omitempty"`
	AccountId AccountId       `json:"accountID,omitempty"`
	Amount    Int64           `json:"amount,omitempty"`
	Locked    Int64           `json:"locked,omitempty"`
	Ext       BalanceEntryExt `json:"ext,omitempty"`
}

// ErrorCode is an XDR Enum defines as:
//
//   enum ErrorCode
//    {
//        MISC = 0, // Unspecific error
//        DATA = 1, // Malformed data
//        CONF = 2, // Misconfiguration error
//        AUTH = 3, // Authentication failure
//        LOAD = 4  // System overloaded
//    };
//
type ErrorCode int32

const (
	ErrorCodeMisc ErrorCode = 0
	ErrorCodeData ErrorCode = 1
	ErrorCodeConf ErrorCode = 2
	ErrorCodeAuth ErrorCode = 3
	ErrorCodeLoad ErrorCode = 4
)

var ErrorCodeAll = []ErrorCode{
	ErrorCodeMisc,
	ErrorCodeData,
	ErrorCodeConf,
	ErrorCodeAuth,
	ErrorCodeLoad,
}

var errorCodeMap = map[int32]string{
	0: "ErrorCodeMisc",
	1: "ErrorCodeData",
	2: "ErrorCodeConf",
	3: "ErrorCodeAuth",
	4: "ErrorCodeLoad",
}

var errorCodeShortMap = map[int32]string{
	0: "misc",
	1: "data",
	2: "conf",
	3: "auth",
	4: "load",
}

var errorCodeRevMap = map[string]int32{
	"ErrorCodeMisc": 0,
	"ErrorCodeData": 1,
	"ErrorCodeConf": 2,
	"ErrorCodeAuth": 3,
	"ErrorCodeLoad": 4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ErrorCode
func (e ErrorCode) ValidEnum(v int32) bool {
	_, ok := errorCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ErrorCode) String() string {
	name, _ := errorCodeMap[int32(e)]
	return name
}

func (e ErrorCode) ShortString() string {
	name, _ := errorCodeShortMap[int32(e)]
	return name
}

func (e ErrorCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ErrorCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := errorCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ErrorCode(value)
	return nil
}

// Error is an XDR Struct defines as:
//
//   struct Error
//    {
//        ErrorCode code;
//        string msg<100>;
//    };
//
type Error struct {
	Code ErrorCode `json:"code,omitempty"`
	Msg  string    `json:"msg,omitempty" xdrmaxsize:"100"`
}

// AuthCert is an XDR Struct defines as:
//
//   struct AuthCert
//    {
//        Curve25519Public pubkey;
//        uint64 expiration;
//        Signature sig;
//    };
//
type AuthCert struct {
	Pubkey     Curve25519Public `json:"pubkey,omitempty"`
	Expiration Uint64           `json:"expiration,omitempty"`
	Sig        Signature        `json:"sig,omitempty"`
}

// Hello is an XDR Struct defines as:
//
//   struct Hello
//    {
//        uint32 ledgerVersion;
//        uint32 overlayVersion;
//        uint32 overlayMinVersion;
//        Hash networkID;
//        string versionStr<100>;
//        int listeningPort;
//        NodeID peerID;
//        AuthCert cert;
//        uint256 nonce;
//    };
//
type Hello struct {
	LedgerVersion     Uint32   `json:"ledgerVersion,omitempty"`
	OverlayVersion    Uint32   `json:"overlayVersion,omitempty"`
	OverlayMinVersion Uint32   `json:"overlayMinVersion,omitempty"`
	NetworkId         Hash     `json:"networkID,omitempty"`
	VersionStr        string   `json:"versionStr,omitempty" xdrmaxsize:"100"`
	ListeningPort     int32    `json:"listeningPort,omitempty"`
	PeerId            NodeId   `json:"peerID,omitempty"`
	Cert              AuthCert `json:"cert,omitempty"`
	Nonce             Uint256  `json:"nonce,omitempty"`
}

// Auth is an XDR Struct defines as:
//
//   struct Auth
//    {
//        // Empty message, just to confirm
//        // establishment of MAC keys.
//        int unused;
//    };
//
type Auth struct {
	Unused int32 `json:"unused,omitempty"`
}

// IpAddrType is an XDR Enum defines as:
//
//   enum IPAddrType
//    {
//        IPv4 = 0,
//        IPv6 = 1
//    };
//
type IpAddrType int32

const (
	IpAddrTypeIPv4 IpAddrType = 0
	IpAddrTypeIPv6 IpAddrType = 1
)

var IpAddrTypeAll = []IpAddrType{
	IpAddrTypeIPv4,
	IpAddrTypeIPv6,
}

var ipAddrTypeMap = map[int32]string{
	0: "IpAddrTypeIPv4",
	1: "IpAddrTypeIPv6",
}

var ipAddrTypeShortMap = map[int32]string{
	0: "i_pv4",
	1: "i_pv6",
}

var ipAddrTypeRevMap = map[string]int32{
	"IpAddrTypeIPv4": 0,
	"IpAddrTypeIPv6": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for IpAddrType
func (e IpAddrType) ValidEnum(v int32) bool {
	_, ok := ipAddrTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e IpAddrType) String() string {
	name, _ := ipAddrTypeMap[int32(e)]
	return name
}

func (e IpAddrType) ShortString() string {
	name, _ := ipAddrTypeShortMap[int32(e)]
	return name
}

func (e IpAddrType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *IpAddrType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := ipAddrTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = IpAddrType(value)
	return nil
}

// PeerAddressIp is an XDR NestedUnion defines as:
//
//   union switch (IPAddrType type)
//        {
//        case IPv4:
//            opaque ipv4[4];
//        case IPv6:
//            opaque ipv6[16];
//        }
//
type PeerAddressIp struct {
	Type IpAddrType `json:"type,omitempty"`
	Ipv4 *[4]byte   `json:"ipv4,omitempty"`
	Ipv6 *[16]byte  `json:"ipv6,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PeerAddressIp) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PeerAddressIp
func (u PeerAddressIp) ArmForSwitch(sw int32) (string, bool) {
	switch IpAddrType(sw) {
	case IpAddrTypeIPv4:
		return "Ipv4", true
	case IpAddrTypeIPv6:
		return "Ipv6", true
	}
	return "-", false
}

// NewPeerAddressIp creates a new  PeerAddressIp.
func NewPeerAddressIp(aType IpAddrType, value interface{}) (result PeerAddressIp, err error) {
	result.Type = aType
	switch IpAddrType(aType) {
	case IpAddrTypeIPv4:
		tv, ok := value.([4]byte)
		if !ok {
			err = fmt.Errorf("invalid value, must be [4]byte")
			return
		}
		result.Ipv4 = &tv
	case IpAddrTypeIPv6:
		tv, ok := value.([16]byte)
		if !ok {
			err = fmt.Errorf("invalid value, must be [16]byte")
			return
		}
		result.Ipv6 = &tv
	}
	return
}

// MustIpv4 retrieves the Ipv4 value from the union,
// panicing if the value is not set.
func (u PeerAddressIp) MustIpv4() [4]byte {
	val, ok := u.GetIpv4()

	if !ok {
		panic("arm Ipv4 is not set")
	}

	return val
}

// GetIpv4 retrieves the Ipv4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PeerAddressIp) GetIpv4() (result [4]byte, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ipv4" {
		result = *u.Ipv4
		ok = true
	}

	return
}

// MustIpv6 retrieves the Ipv6 value from the union,
// panicing if the value is not set.
func (u PeerAddressIp) MustIpv6() [16]byte {
	val, ok := u.GetIpv6()

	if !ok {
		panic("arm Ipv6 is not set")
	}

	return val
}

// GetIpv6 retrieves the Ipv6 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PeerAddressIp) GetIpv6() (result [16]byte, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ipv6" {
		result = *u.Ipv6
		ok = true
	}

	return
}

// PeerAddress is an XDR Struct defines as:
//
//   struct PeerAddress
//    {
//        union switch (IPAddrType type)
//        {
//        case IPv4:
//            opaque ipv4[4];
//        case IPv6:
//            opaque ipv6[16];
//        }
//        ip;
//        uint32 port;
//        uint32 numFailures;
//    };
//
type PeerAddress struct {
	Ip          PeerAddressIp `json:"ip,omitempty"`
	Port        Uint32        `json:"port,omitempty"`
	NumFailures Uint32        `json:"numFailures,omitempty"`
}

// MessageType is an XDR Enum defines as:
//
//   enum MessageType
//    {
//        ERROR_MSG = 0,
//        AUTH = 2,
//        DONT_HAVE = 3,
//
//        GET_PEERS = 4, // gets a list of peers this guy knows about
//        PEERS = 5,
//
//        GET_TX_SET = 6, // gets a particular txset by hash
//        TX_SET = 7,
//
//        TRANSACTION = 8, // pass on a tx you have heard about
//
//        // SCP
//        GET_SCP_QUORUMSET = 9,
//        SCP_QUORUMSET = 10,
//        SCP_MESSAGE = 11,
//        GET_SCP_STATE = 12,
//
//        // new messages
//        HELLO = 13
//    };
//
type MessageType int32

const (
	MessageTypeErrorMsg        MessageType = 0
	MessageTypeAuth            MessageType = 2
	MessageTypeDontHave        MessageType = 3
	MessageTypeGetPeers        MessageType = 4
	MessageTypePeers           MessageType = 5
	MessageTypeGetTxSet        MessageType = 6
	MessageTypeTxSet           MessageType = 7
	MessageTypeTransaction     MessageType = 8
	MessageTypeGetScpQuorumset MessageType = 9
	MessageTypeScpQuorumset    MessageType = 10
	MessageTypeScpMessage      MessageType = 11
	MessageTypeGetScpState     MessageType = 12
	MessageTypeHello           MessageType = 13
)

var MessageTypeAll = []MessageType{
	MessageTypeErrorMsg,
	MessageTypeAuth,
	MessageTypeDontHave,
	MessageTypeGetPeers,
	MessageTypePeers,
	MessageTypeGetTxSet,
	MessageTypeTxSet,
	MessageTypeTransaction,
	MessageTypeGetScpQuorumset,
	MessageTypeScpQuorumset,
	MessageTypeScpMessage,
	MessageTypeGetScpState,
	MessageTypeHello,
}

var messageTypeMap = map[int32]string{
	0:  "MessageTypeErrorMsg",
	2:  "MessageTypeAuth",
	3:  "MessageTypeDontHave",
	4:  "MessageTypeGetPeers",
	5:  "MessageTypePeers",
	6:  "MessageTypeGetTxSet",
	7:  "MessageTypeTxSet",
	8:  "MessageTypeTransaction",
	9:  "MessageTypeGetScpQuorumset",
	10: "MessageTypeScpQuorumset",
	11: "MessageTypeScpMessage",
	12: "MessageTypeGetScpState",
	13: "MessageTypeHello",
}

var messageTypeShortMap = map[int32]string{
	0:  "error_msg",
	2:  "auth",
	3:  "dont_have",
	4:  "get_peers",
	5:  "peers",
	6:  "get_tx_set",
	7:  "tx_set",
	8:  "transaction",
	9:  "get_scp_quorumset",
	10: "scp_quorumset",
	11: "scp_message",
	12: "get_scp_state",
	13: "hello",
}

var messageTypeRevMap = map[string]int32{
	"MessageTypeErrorMsg":        0,
	"MessageTypeAuth":            2,
	"MessageTypeDontHave":        3,
	"MessageTypeGetPeers":        4,
	"MessageTypePeers":           5,
	"MessageTypeGetTxSet":        6,
	"MessageTypeTxSet":           7,
	"MessageTypeTransaction":     8,
	"MessageTypeGetScpQuorumset": 9,
	"MessageTypeScpQuorumset":    10,
	"MessageTypeScpMessage":      11,
	"MessageTypeGetScpState":     12,
	"MessageTypeHello":           13,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MessageType
func (e MessageType) ValidEnum(v int32) bool {
	_, ok := messageTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MessageType) String() string {
	name, _ := messageTypeMap[int32(e)]
	return name
}

func (e MessageType) ShortString() string {
	name, _ := messageTypeShortMap[int32(e)]
	return name
}

func (e MessageType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *MessageType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := messageTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = MessageType(value)
	return nil
}

// DontHave is an XDR Struct defines as:
//
//   struct DontHave
//    {
//        MessageType type;
//        uint256 reqHash;
//    };
//
type DontHave struct {
	Type    MessageType `json:"type,omitempty"`
	ReqHash Uint256     `json:"reqHash,omitempty"`
}

// StellarMessage is an XDR Union defines as:
//
//   union StellarMessage switch (MessageType type)
//    {
//    case ERROR_MSG:
//        Error error;
//    case HELLO:
//        Hello hello;
//    case AUTH:
//        Auth auth;
//    case DONT_HAVE:
//        DontHave dontHave;
//    case GET_PEERS:
//        void;
//    case PEERS:
//        PeerAddress peers<>;
//
//    case GET_TX_SET:
//        uint256 txSetHash;
//    case TX_SET:
//        TransactionSet txSet;
//
//    case TRANSACTION:
//        TransactionEnvelope transaction;
//
//    // SCP
//    case GET_SCP_QUORUMSET:
//        uint256 qSetHash;
//    case SCP_QUORUMSET:
//        SCPQuorumSet qSet;
//    case SCP_MESSAGE:
//        SCPEnvelope envelope;
//    case GET_SCP_STATE:
//        uint32 getSCPLedgerSeq; // ledger seq requested ; if 0, requests the latest
//    };
//
type StellarMessage struct {
	Type            MessageType          `json:"type,omitempty"`
	Error           *Error               `json:"error,omitempty"`
	Hello           *Hello               `json:"hello,omitempty"`
	Auth            *Auth                `json:"auth,omitempty"`
	DontHave        *DontHave            `json:"dontHave,omitempty"`
	Peers           *[]PeerAddress       `json:"peers,omitempty"`
	TxSetHash       *Uint256             `json:"txSetHash,omitempty"`
	TxSet           *TransactionSet      `json:"txSet,omitempty"`
	Transaction     *TransactionEnvelope `json:"transaction,omitempty"`
	QSetHash        *Uint256             `json:"qSetHash,omitempty"`
	QSet            *ScpQuorumSet        `json:"qSet,omitempty"`
	Envelope        *ScpEnvelope         `json:"envelope,omitempty"`
	GetScpLedgerSeq *Uint32              `json:"getSCPLedgerSeq,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u StellarMessage) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of StellarMessage
func (u StellarMessage) ArmForSwitch(sw int32) (string, bool) {
	switch MessageType(sw) {
	case MessageTypeErrorMsg:
		return "Error", true
	case MessageTypeHello:
		return "Hello", true
	case MessageTypeAuth:
		return "Auth", true
	case MessageTypeDontHave:
		return "DontHave", true
	case MessageTypeGetPeers:
		return "", true
	case MessageTypePeers:
		return "Peers", true
	case MessageTypeGetTxSet:
		return "TxSetHash", true
	case MessageTypeTxSet:
		return "TxSet", true
	case MessageTypeTransaction:
		return "Transaction", true
	case MessageTypeGetScpQuorumset:
		return "QSetHash", true
	case MessageTypeScpQuorumset:
		return "QSet", true
	case MessageTypeScpMessage:
		return "Envelope", true
	case MessageTypeGetScpState:
		return "GetScpLedgerSeq", true
	}
	return "-", false
}

// NewStellarMessage creates a new  StellarMessage.
func NewStellarMessage(aType MessageType, value interface{}) (result StellarMessage, err error) {
	result.Type = aType
	switch MessageType(aType) {
	case MessageTypeErrorMsg:
		tv, ok := value.(Error)
		if !ok {
			err = fmt.Errorf("invalid value, must be Error")
			return
		}
		result.Error = &tv
	case MessageTypeHello:
		tv, ok := value.(Hello)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hello")
			return
		}
		result.Hello = &tv
	case MessageTypeAuth:
		tv, ok := value.(Auth)
		if !ok {
			err = fmt.Errorf("invalid value, must be Auth")
			return
		}
		result.Auth = &tv
	case MessageTypeDontHave:
		tv, ok := value.(DontHave)
		if !ok {
			err = fmt.Errorf("invalid value, must be DontHave")
			return
		}
		result.DontHave = &tv
	case MessageTypeGetPeers:
		// void
	case MessageTypePeers:
		tv, ok := value.([]PeerAddress)
		if !ok {
			err = fmt.Errorf("invalid value, must be []PeerAddress")
			return
		}
		result.Peers = &tv
	case MessageTypeGetTxSet:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.TxSetHash = &tv
	case MessageTypeTxSet:
		tv, ok := value.(TransactionSet)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionSet")
			return
		}
		result.TxSet = &tv
	case MessageTypeTransaction:
		tv, ok := value.(TransactionEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionEnvelope")
			return
		}
		result.Transaction = &tv
	case MessageTypeGetScpQuorumset:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.QSetHash = &tv
	case MessageTypeScpQuorumset:
		tv, ok := value.(ScpQuorumSet)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpQuorumSet")
			return
		}
		result.QSet = &tv
	case MessageTypeScpMessage:
		tv, ok := value.(ScpEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpEnvelope")
			return
		}
		result.Envelope = &tv
	case MessageTypeGetScpState:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.GetScpLedgerSeq = &tv
	}
	return
}

// MustError retrieves the Error value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustError() Error {
	val, ok := u.GetError()

	if !ok {
		panic("arm Error is not set")
	}

	return val
}

// GetError retrieves the Error value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetError() (result Error, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Error" {
		result = *u.Error
		ok = true
	}

	return
}

// MustHello retrieves the Hello value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustHello() Hello {
	val, ok := u.GetHello()

	if !ok {
		panic("arm Hello is not set")
	}

	return val
}

// GetHello retrieves the Hello value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetHello() (result Hello, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hello" {
		result = *u.Hello
		ok = true
	}

	return
}

// MustAuth retrieves the Auth value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustAuth() Auth {
	val, ok := u.GetAuth()

	if !ok {
		panic("arm Auth is not set")
	}

	return val
}

// GetAuth retrieves the Auth value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetAuth() (result Auth, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Auth" {
		result = *u.Auth
		ok = true
	}

	return
}

// MustDontHave retrieves the DontHave value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustDontHave() DontHave {
	val, ok := u.GetDontHave()

	if !ok {
		panic("arm DontHave is not set")
	}

	return val
}

// GetDontHave retrieves the DontHave value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetDontHave() (result DontHave, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DontHave" {
		result = *u.DontHave
		ok = true
	}

	return
}

// MustPeers retrieves the Peers value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustPeers() []PeerAddress {
	val, ok := u.GetPeers()

	if !ok {
		panic("arm Peers is not set")
	}

	return val
}

// GetPeers retrieves the Peers value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetPeers() (result []PeerAddress, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Peers" {
		result = *u.Peers
		ok = true
	}

	return
}

// MustTxSetHash retrieves the TxSetHash value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustTxSetHash() Uint256 {
	val, ok := u.GetTxSetHash()

	if !ok {
		panic("arm TxSetHash is not set")
	}

	return val
}

// GetTxSetHash retrieves the TxSetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetTxSetHash() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TxSetHash" {
		result = *u.TxSetHash
		ok = true
	}

	return
}

// MustTxSet retrieves the TxSet value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustTxSet() TransactionSet {
	val, ok := u.GetTxSet()

	if !ok {
		panic("arm TxSet is not set")
	}

	return val
}

// GetTxSet retrieves the TxSet value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetTxSet() (result TransactionSet, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TxSet" {
		result = *u.TxSet
		ok = true
	}

	return
}

// MustTransaction retrieves the Transaction value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustTransaction() TransactionEnvelope {
	val, ok := u.GetTransaction()

	if !ok {
		panic("arm Transaction is not set")
	}

	return val
}

// GetTransaction retrieves the Transaction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetTransaction() (result TransactionEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Transaction" {
		result = *u.Transaction
		ok = true
	}

	return
}

// MustQSetHash retrieves the QSetHash value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustQSetHash() Uint256 {
	val, ok := u.GetQSetHash()

	if !ok {
		panic("arm QSetHash is not set")
	}

	return val
}

// GetQSetHash retrieves the QSetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetQSetHash() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "QSetHash" {
		result = *u.QSetHash
		ok = true
	}

	return
}

// MustQSet retrieves the QSet value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustQSet() ScpQuorumSet {
	val, ok := u.GetQSet()

	if !ok {
		panic("arm QSet is not set")
	}

	return val
}

// GetQSet retrieves the QSet value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetQSet() (result ScpQuorumSet, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "QSet" {
		result = *u.QSet
		ok = true
	}

	return
}

// MustEnvelope retrieves the Envelope value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustEnvelope() ScpEnvelope {
	val, ok := u.GetEnvelope()

	if !ok {
		panic("arm Envelope is not set")
	}

	return val
}

// GetEnvelope retrieves the Envelope value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetEnvelope() (result ScpEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Envelope" {
		result = *u.Envelope
		ok = true
	}

	return
}

// MustGetScpLedgerSeq retrieves the GetScpLedgerSeq value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustGetScpLedgerSeq() Uint32 {
	val, ok := u.GetGetScpLedgerSeq()

	if !ok {
		panic("arm GetScpLedgerSeq is not set")
	}

	return val
}

// GetGetScpLedgerSeq retrieves the GetScpLedgerSeq value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetGetScpLedgerSeq() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "GetScpLedgerSeq" {
		result = *u.GetScpLedgerSeq
		ok = true
	}

	return
}

// AuthenticatedMessageV0 is an XDR NestedStruct defines as:
//
//   struct
//        {
//            uint64 sequence;
//            StellarMessage message;
//            HmacSha256Mac mac;
//        }
//
type AuthenticatedMessageV0 struct {
	Sequence Uint64         `json:"sequence,omitempty"`
	Message  StellarMessage `json:"message,omitempty"`
	Mac      HmacSha256Mac  `json:"mac,omitempty"`
}

// AuthenticatedMessage is an XDR Union defines as:
//
//   union AuthenticatedMessage switch (LedgerVersion v)
//    {
//    case EMPTY_VERSION:
//        struct
//        {
//            uint64 sequence;
//            StellarMessage message;
//            HmacSha256Mac mac;
//        } v0;
//    };
//
type AuthenticatedMessage struct {
	V  LedgerVersion           `json:"v,omitempty"`
	V0 *AuthenticatedMessageV0 `json:"v0,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AuthenticatedMessage) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AuthenticatedMessage
func (u AuthenticatedMessage) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "V0", true
	}
	return "-", false
}

// NewAuthenticatedMessage creates a new  AuthenticatedMessage.
func NewAuthenticatedMessage(v LedgerVersion, value interface{}) (result AuthenticatedMessage, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		tv, ok := value.(AuthenticatedMessageV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be AuthenticatedMessageV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u AuthenticatedMessage) MustV0() AuthenticatedMessageV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AuthenticatedMessage) GetV0() (result AuthenticatedMessageV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// CreateAccountOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CreateAccountOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountOpExt
func (u CreateAccountOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCreateAccountOpExt creates a new  CreateAccountOpExt.
func NewCreateAccountOpExt(v LedgerVersion, value interface{}) (result CreateAccountOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CreateAccountOp is an XDR Struct defines as:
//
//   struct CreateAccountOp
//    {
//        AccountID destination; // account to create
//        AccountID* referrer;     // parent account
//    	AccountType accountType;
//
//    	uint32 policies; //account policies for the account
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CreateAccountOp struct {
	Destination AccountId          `json:"destination,omitempty"`
	Referrer    *AccountId         `json:"referrer,omitempty"`
	AccountType AccountType        `json:"accountType,omitempty"`
	Policies    Uint32             `json:"policies,omitempty"`
	Ext         CreateAccountOpExt `json:"ext,omitempty"`
}

// CreateAccountResultCode is an XDR Enum defines as:
//
//   enum CreateAccountResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0, // account was created
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,       // invalid destination
//    	ACCOUNT_TYPE_MISMATCHED = -2, // account already exist and change of account type is not allowed
//    	TYPE_NOT_ALLOWED = -3, // master or commission account types are not allowed
//        NAME_DUPLICATION = -4,
//        REFERRER_NOT_FOUND = -5,
//    	INVALID_ACCOUNT_VERSION = -6 // if account version is higher than ledger version
//    };
//
type CreateAccountResultCode int32

const (
	CreateAccountResultCodeSuccess               CreateAccountResultCode = 0
	CreateAccountResultCodeMalformed             CreateAccountResultCode = -1
	CreateAccountResultCodeAccountTypeMismatched CreateAccountResultCode = -2
	CreateAccountResultCodeTypeNotAllowed        CreateAccountResultCode = -3
	CreateAccountResultCodeNameDuplication       CreateAccountResultCode = -4
	CreateAccountResultCodeReferrerNotFound      CreateAccountResultCode = -5
	CreateAccountResultCodeInvalidAccountVersion CreateAccountResultCode = -6
)

var CreateAccountResultCodeAll = []CreateAccountResultCode{
	CreateAccountResultCodeSuccess,
	CreateAccountResultCodeMalformed,
	CreateAccountResultCodeAccountTypeMismatched,
	CreateAccountResultCodeTypeNotAllowed,
	CreateAccountResultCodeNameDuplication,
	CreateAccountResultCodeReferrerNotFound,
	CreateAccountResultCodeInvalidAccountVersion,
}

var createAccountResultCodeMap = map[int32]string{
	0:  "CreateAccountResultCodeSuccess",
	-1: "CreateAccountResultCodeMalformed",
	-2: "CreateAccountResultCodeAccountTypeMismatched",
	-3: "CreateAccountResultCodeTypeNotAllowed",
	-4: "CreateAccountResultCodeNameDuplication",
	-5: "CreateAccountResultCodeReferrerNotFound",
	-6: "CreateAccountResultCodeInvalidAccountVersion",
}

var createAccountResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "malformed",
	-2: "account_type_mismatched",
	-3: "type_not_allowed",
	-4: "name_duplication",
	-5: "referrer_not_found",
	-6: "invalid_account_version",
}

var createAccountResultCodeRevMap = map[string]int32{
	"CreateAccountResultCodeSuccess":               0,
	"CreateAccountResultCodeMalformed":             -1,
	"CreateAccountResultCodeAccountTypeMismatched": -2,
	"CreateAccountResultCodeTypeNotAllowed":        -3,
	"CreateAccountResultCodeNameDuplication":       -4,
	"CreateAccountResultCodeReferrerNotFound":      -5,
	"CreateAccountResultCodeInvalidAccountVersion": -6,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CreateAccountResultCode
func (e CreateAccountResultCode) ValidEnum(v int32) bool {
	_, ok := createAccountResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e CreateAccountResultCode) String() string {
	name, _ := createAccountResultCodeMap[int32(e)]
	return name
}

func (e CreateAccountResultCode) ShortString() string {
	name, _ := createAccountResultCodeShortMap[int32(e)]
	return name
}

func (e CreateAccountResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *CreateAccountResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := createAccountResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = CreateAccountResultCode(value)
	return nil
}

// CreateAccountSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CreateAccountSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountSuccessExt
func (u CreateAccountSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCreateAccountSuccessExt creates a new  CreateAccountSuccessExt.
func NewCreateAccountSuccessExt(v LedgerVersion, value interface{}) (result CreateAccountSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CreateAccountSuccess is an XDR Struct defines as:
//
//   struct CreateAccountSuccess
//    {
//    	int64 referrerFee;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CreateAccountSuccess struct {
	ReferrerFee Int64                   `json:"referrerFee,omitempty"`
	Ext         CreateAccountSuccessExt `json:"ext,omitempty"`
}

// CreateAccountResult is an XDR Union defines as:
//
//   union CreateAccountResult switch (CreateAccountResultCode code)
//    {
//    case SUCCESS:
//        CreateAccountSuccess success;
//    default:
//        void;
//    };
//
type CreateAccountResult struct {
	Code    CreateAccountResultCode `json:"code,omitempty"`
	Success *CreateAccountSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountResult
func (u CreateAccountResult) ArmForSwitch(sw int32) (string, bool) {
	switch CreateAccountResultCode(sw) {
	case CreateAccountResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewCreateAccountResult creates a new  CreateAccountResult.
func NewCreateAccountResult(code CreateAccountResultCode, value interface{}) (result CreateAccountResult, err error) {
	result.Code = code
	switch CreateAccountResultCode(code) {
	case CreateAccountResultCodeSuccess:
		tv, ok := value.(CreateAccountSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u CreateAccountResult) MustSuccess() CreateAccountSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u CreateAccountResult) GetSuccess() (result CreateAccountSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageAccountOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAccountOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAccountOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAccountOpExt
func (u ManageAccountOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAccountOpExt creates a new  ManageAccountOpExt.
func NewManageAccountOpExt(v LedgerVersion, value interface{}) (result ManageAccountOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAccountOp is an XDR Struct defines as:
//
//   struct ManageAccountOp
//    {
//        AccountID account; // account to manage
//        AccountType accountType;
//        uint32 blockReasonsToAdd;
//        uint32 blockReasonsToRemove;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAccountOp struct {
	Account              AccountId          `json:"account,omitempty"`
	AccountType          AccountType        `json:"accountType,omitempty"`
	BlockReasonsToAdd    Uint32             `json:"blockReasonsToAdd,omitempty"`
	BlockReasonsToRemove Uint32             `json:"blockReasonsToRemove,omitempty"`
	Ext                  ManageAccountOpExt `json:"ext,omitempty"`
}

// ManageAccountResultCode is an XDR Enum defines as:
//
//   enum ManageAccountResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0, // account was created
//
//        // codes considered as "failure" for the operation
//        NOT_FOUND = -1,         // account does not exists
//        MALFORMED = -2,
//    	NOT_ALLOWED = -3,         // manage account operation is not allowed on this account
//        TYPE_MISMATCH = -4
//    };
//
type ManageAccountResultCode int32

const (
	ManageAccountResultCodeSuccess      ManageAccountResultCode = 0
	ManageAccountResultCodeNotFound     ManageAccountResultCode = -1
	ManageAccountResultCodeMalformed    ManageAccountResultCode = -2
	ManageAccountResultCodeNotAllowed   ManageAccountResultCode = -3
	ManageAccountResultCodeTypeMismatch ManageAccountResultCode = -4
)

var ManageAccountResultCodeAll = []ManageAccountResultCode{
	ManageAccountResultCodeSuccess,
	ManageAccountResultCodeNotFound,
	ManageAccountResultCodeMalformed,
	ManageAccountResultCodeNotAllowed,
	ManageAccountResultCodeTypeMismatch,
}

var manageAccountResultCodeMap = map[int32]string{
	0:  "ManageAccountResultCodeSuccess",
	-1: "ManageAccountResultCodeNotFound",
	-2: "ManageAccountResultCodeMalformed",
	-3: "ManageAccountResultCodeNotAllowed",
	-4: "ManageAccountResultCodeTypeMismatch",
}

var manageAccountResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "not_found",
	-2: "malformed",
	-3: "not_allowed",
	-4: "type_mismatch",
}

var manageAccountResultCodeRevMap = map[string]int32{
	"ManageAccountResultCodeSuccess":      0,
	"ManageAccountResultCodeNotFound":     -1,
	"ManageAccountResultCodeMalformed":    -2,
	"ManageAccountResultCodeNotAllowed":   -3,
	"ManageAccountResultCodeTypeMismatch": -4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAccountResultCode
func (e ManageAccountResultCode) ValidEnum(v int32) bool {
	_, ok := manageAccountResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAccountResultCode) String() string {
	name, _ := manageAccountResultCodeMap[int32(e)]
	return name
}

func (e ManageAccountResultCode) ShortString() string {
	name, _ := manageAccountResultCodeShortMap[int32(e)]
	return name
}

func (e ManageAccountResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageAccountResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageAccountResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageAccountResultCode(value)
	return nil
}

// ManageAccountSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAccountSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAccountSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAccountSuccessExt
func (u ManageAccountSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAccountSuccessExt creates a new  ManageAccountSuccessExt.
func NewManageAccountSuccessExt(v LedgerVersion, value interface{}) (result ManageAccountSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAccountSuccess is an XDR Struct defines as:
//
//   struct ManageAccountSuccess {
//    	uint32 blockReasons;
//     // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAccountSuccess struct {
	BlockReasons Uint32                  `json:"blockReasons,omitempty"`
	Ext          ManageAccountSuccessExt `json:"ext,omitempty"`
}

// ManageAccountResult is an XDR Union defines as:
//
//   union ManageAccountResult switch (ManageAccountResultCode code)
//    {
//    case SUCCESS:
//        ManageAccountSuccess success;
//    default:
//        void;
//    };
//
type ManageAccountResult struct {
	Code    ManageAccountResultCode `json:"code,omitempty"`
	Success *ManageAccountSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAccountResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAccountResult
func (u ManageAccountResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageAccountResultCode(sw) {
	case ManageAccountResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageAccountResult creates a new  ManageAccountResult.
func NewManageAccountResult(code ManageAccountResultCode, value interface{}) (result ManageAccountResult, err error) {
	result.Code = code
	switch ManageAccountResultCode(code) {
	case ManageAccountResultCodeSuccess:
		tv, ok := value.(ManageAccountSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAccountSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageAccountResult) MustSuccess() ManageAccountSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageAccountResult) GetSuccess() (result ManageAccountSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// OperationBody is an XDR NestedUnion defines as:
//
//   union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountOp createAccountOp;
//        case PAYMENT:
//            PaymentOp paymentOp;
//        case SET_OPTIONS:
//            SetOptionsOp setOptionsOp;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestOp manageCoinsEmissionRequestOp;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestOp reviewCoinsEmissionRequestOp;
//        case SET_FEES:
//            SetFeesOp setFeesOp;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountOp manageAccountOp;
//    	case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestOp manageForfeitRequestOp;
//    	case RECOVER:
//    		RecoverOp recoverOp;
//    	case MANAGE_BALANCE:
//    		ManageBalanceOp manageBalanceOp;
//    	case REVIEW_PAYMENT_REQUEST:
//    		ReviewPaymentRequestOp reviewPaymentRequestOp;
//        case MANAGE_ASSET:
//            ManageAssetOp manageAssetOp;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsOp uploadPreemissionsOp;
//        case SET_LIMITS:
//            SetLimitsOp setLimitsOp;
//        case DIRECT_DEBIT:
//            DirectDebitOp directDebitOp;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairOp manageAssetPairOp;
//    	case MANAGE_OFFER:
//    		ManageOfferOp manageOfferOp;
//        case MANAGE_INVOICE:
//            ManageInvoiceOp manageInvoiceOp;
//        }
//
type OperationBody struct {
	Type                         OperationType                 `json:"type,omitempty"`
	CreateAccountOp              *CreateAccountOp              `json:"createAccountOp,omitempty"`
	PaymentOp                    *PaymentOp                    `json:"paymentOp,omitempty"`
	SetOptionsOp                 *SetOptionsOp                 `json:"setOptionsOp,omitempty"`
	ManageCoinsEmissionRequestOp *ManageCoinsEmissionRequestOp `json:"manageCoinsEmissionRequestOp,omitempty"`
	ReviewCoinsEmissionRequestOp *ReviewCoinsEmissionRequestOp `json:"reviewCoinsEmissionRequestOp,omitempty"`
	SetFeesOp                    *SetFeesOp                    `json:"setFeesOp,omitempty"`
	ManageAccountOp              *ManageAccountOp              `json:"manageAccountOp,omitempty"`
	ManageForfeitRequestOp       *ManageForfeitRequestOp       `json:"manageForfeitRequestOp,omitempty"`
	RecoverOp                    *RecoverOp                    `json:"recoverOp,omitempty"`
	ManageBalanceOp              *ManageBalanceOp              `json:"manageBalanceOp,omitempty"`
	ReviewPaymentRequestOp       *ReviewPaymentRequestOp       `json:"reviewPaymentRequestOp,omitempty"`
	ManageAssetOp                *ManageAssetOp                `json:"manageAssetOp,omitempty"`
	UploadPreemissionsOp         *UploadPreemissionsOp         `json:"uploadPreemissionsOp,omitempty"`
	SetLimitsOp                  *SetLimitsOp                  `json:"setLimitsOp,omitempty"`
	DirectDebitOp                *DirectDebitOp                `json:"directDebitOp,omitempty"`
	ManageAssetPairOp            *ManageAssetPairOp            `json:"manageAssetPairOp,omitempty"`
	ManageOfferOp                *ManageOfferOp                `json:"manageOfferOp,omitempty"`
	ManageInvoiceOp              *ManageInvoiceOp              `json:"manageInvoiceOp,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationBody) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationBody
func (u OperationBody) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypeCreateAccount:
		return "CreateAccountOp", true
	case OperationTypePayment:
		return "PaymentOp", true
	case OperationTypeSetOptions:
		return "SetOptionsOp", true
	case OperationTypeManageCoinsEmissionRequest:
		return "ManageCoinsEmissionRequestOp", true
	case OperationTypeReviewCoinsEmissionRequest:
		return "ReviewCoinsEmissionRequestOp", true
	case OperationTypeSetFees:
		return "SetFeesOp", true
	case OperationTypeManageAccount:
		return "ManageAccountOp", true
	case OperationTypeManageForfeitRequest:
		return "ManageForfeitRequestOp", true
	case OperationTypeRecover:
		return "RecoverOp", true
	case OperationTypeManageBalance:
		return "ManageBalanceOp", true
	case OperationTypeReviewPaymentRequest:
		return "ReviewPaymentRequestOp", true
	case OperationTypeManageAsset:
		return "ManageAssetOp", true
	case OperationTypeUploadPreemissions:
		return "UploadPreemissionsOp", true
	case OperationTypeSetLimits:
		return "SetLimitsOp", true
	case OperationTypeDirectDebit:
		return "DirectDebitOp", true
	case OperationTypeManageAssetPair:
		return "ManageAssetPairOp", true
	case OperationTypeManageOffer:
		return "ManageOfferOp", true
	case OperationTypeManageInvoice:
		return "ManageInvoiceOp", true
	}
	return "-", false
}

// NewOperationBody creates a new  OperationBody.
func NewOperationBody(aType OperationType, value interface{}) (result OperationBody, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypeCreateAccount:
		tv, ok := value.(CreateAccountOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountOp")
			return
		}
		result.CreateAccountOp = &tv
	case OperationTypePayment:
		tv, ok := value.(PaymentOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentOp")
			return
		}
		result.PaymentOp = &tv
	case OperationTypeSetOptions:
		tv, ok := value.(SetOptionsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsOp")
			return
		}
		result.SetOptionsOp = &tv
	case OperationTypeManageCoinsEmissionRequest:
		tv, ok := value.(ManageCoinsEmissionRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageCoinsEmissionRequestOp")
			return
		}
		result.ManageCoinsEmissionRequestOp = &tv
	case OperationTypeReviewCoinsEmissionRequest:
		tv, ok := value.(ReviewCoinsEmissionRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewCoinsEmissionRequestOp")
			return
		}
		result.ReviewCoinsEmissionRequestOp = &tv
	case OperationTypeSetFees:
		tv, ok := value.(SetFeesOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetFeesOp")
			return
		}
		result.SetFeesOp = &tv
	case OperationTypeManageAccount:
		tv, ok := value.(ManageAccountOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAccountOp")
			return
		}
		result.ManageAccountOp = &tv
	case OperationTypeManageForfeitRequest:
		tv, ok := value.(ManageForfeitRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestOp")
			return
		}
		result.ManageForfeitRequestOp = &tv
	case OperationTypeRecover:
		tv, ok := value.(RecoverOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be RecoverOp")
			return
		}
		result.RecoverOp = &tv
	case OperationTypeManageBalance:
		tv, ok := value.(ManageBalanceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBalanceOp")
			return
		}
		result.ManageBalanceOp = &tv
	case OperationTypeReviewPaymentRequest:
		tv, ok := value.(ReviewPaymentRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewPaymentRequestOp")
			return
		}
		result.ReviewPaymentRequestOp = &tv
	case OperationTypeManageAsset:
		tv, ok := value.(ManageAssetOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetOp")
			return
		}
		result.ManageAssetOp = &tv
	case OperationTypeUploadPreemissions:
		tv, ok := value.(UploadPreemissionsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be UploadPreemissionsOp")
			return
		}
		result.UploadPreemissionsOp = &tv
	case OperationTypeSetLimits:
		tv, ok := value.(SetLimitsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetLimitsOp")
			return
		}
		result.SetLimitsOp = &tv
	case OperationTypeDirectDebit:
		tv, ok := value.(DirectDebitOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be DirectDebitOp")
			return
		}
		result.DirectDebitOp = &tv
	case OperationTypeManageAssetPair:
		tv, ok := value.(ManageAssetPairOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetPairOp")
			return
		}
		result.ManageAssetPairOp = &tv
	case OperationTypeManageOffer:
		tv, ok := value.(ManageOfferOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferOp")
			return
		}
		result.ManageOfferOp = &tv
	case OperationTypeManageInvoice:
		tv, ok := value.(ManageInvoiceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageInvoiceOp")
			return
		}
		result.ManageInvoiceOp = &tv
	}
	return
}

// MustCreateAccountOp retrieves the CreateAccountOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustCreateAccountOp() CreateAccountOp {
	val, ok := u.GetCreateAccountOp()

	if !ok {
		panic("arm CreateAccountOp is not set")
	}

	return val
}

// GetCreateAccountOp retrieves the CreateAccountOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetCreateAccountOp() (result CreateAccountOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateAccountOp" {
		result = *u.CreateAccountOp
		ok = true
	}

	return
}

// MustPaymentOp retrieves the PaymentOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustPaymentOp() PaymentOp {
	val, ok := u.GetPaymentOp()

	if !ok {
		panic("arm PaymentOp is not set")
	}

	return val
}

// GetPaymentOp retrieves the PaymentOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetPaymentOp() (result PaymentOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentOp" {
		result = *u.PaymentOp
		ok = true
	}

	return
}

// MustSetOptionsOp retrieves the SetOptionsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetOptionsOp() SetOptionsOp {
	val, ok := u.GetSetOptionsOp()

	if !ok {
		panic("arm SetOptionsOp is not set")
	}

	return val
}

// GetSetOptionsOp retrieves the SetOptionsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetOptionsOp() (result SetOptionsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetOptionsOp" {
		result = *u.SetOptionsOp
		ok = true
	}

	return
}

// MustManageCoinsEmissionRequestOp retrieves the ManageCoinsEmissionRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageCoinsEmissionRequestOp() ManageCoinsEmissionRequestOp {
	val, ok := u.GetManageCoinsEmissionRequestOp()

	if !ok {
		panic("arm ManageCoinsEmissionRequestOp is not set")
	}

	return val
}

// GetManageCoinsEmissionRequestOp retrieves the ManageCoinsEmissionRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageCoinsEmissionRequestOp() (result ManageCoinsEmissionRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageCoinsEmissionRequestOp" {
		result = *u.ManageCoinsEmissionRequestOp
		ok = true
	}

	return
}

// MustReviewCoinsEmissionRequestOp retrieves the ReviewCoinsEmissionRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustReviewCoinsEmissionRequestOp() ReviewCoinsEmissionRequestOp {
	val, ok := u.GetReviewCoinsEmissionRequestOp()

	if !ok {
		panic("arm ReviewCoinsEmissionRequestOp is not set")
	}

	return val
}

// GetReviewCoinsEmissionRequestOp retrieves the ReviewCoinsEmissionRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetReviewCoinsEmissionRequestOp() (result ReviewCoinsEmissionRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewCoinsEmissionRequestOp" {
		result = *u.ReviewCoinsEmissionRequestOp
		ok = true
	}

	return
}

// MustSetFeesOp retrieves the SetFeesOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetFeesOp() SetFeesOp {
	val, ok := u.GetSetFeesOp()

	if !ok {
		panic("arm SetFeesOp is not set")
	}

	return val
}

// GetSetFeesOp retrieves the SetFeesOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetFeesOp() (result SetFeesOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetFeesOp" {
		result = *u.SetFeesOp
		ok = true
	}

	return
}

// MustManageAccountOp retrieves the ManageAccountOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageAccountOp() ManageAccountOp {
	val, ok := u.GetManageAccountOp()

	if !ok {
		panic("arm ManageAccountOp is not set")
	}

	return val
}

// GetManageAccountOp retrieves the ManageAccountOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageAccountOp() (result ManageAccountOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAccountOp" {
		result = *u.ManageAccountOp
		ok = true
	}

	return
}

// MustManageForfeitRequestOp retrieves the ManageForfeitRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageForfeitRequestOp() ManageForfeitRequestOp {
	val, ok := u.GetManageForfeitRequestOp()

	if !ok {
		panic("arm ManageForfeitRequestOp is not set")
	}

	return val
}

// GetManageForfeitRequestOp retrieves the ManageForfeitRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageForfeitRequestOp() (result ManageForfeitRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageForfeitRequestOp" {
		result = *u.ManageForfeitRequestOp
		ok = true
	}

	return
}

// MustRecoverOp retrieves the RecoverOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustRecoverOp() RecoverOp {
	val, ok := u.GetRecoverOp()

	if !ok {
		panic("arm RecoverOp is not set")
	}

	return val
}

// GetRecoverOp retrieves the RecoverOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetRecoverOp() (result RecoverOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RecoverOp" {
		result = *u.RecoverOp
		ok = true
	}

	return
}

// MustManageBalanceOp retrieves the ManageBalanceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageBalanceOp() ManageBalanceOp {
	val, ok := u.GetManageBalanceOp()

	if !ok {
		panic("arm ManageBalanceOp is not set")
	}

	return val
}

// GetManageBalanceOp retrieves the ManageBalanceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageBalanceOp() (result ManageBalanceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageBalanceOp" {
		result = *u.ManageBalanceOp
		ok = true
	}

	return
}

// MustReviewPaymentRequestOp retrieves the ReviewPaymentRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustReviewPaymentRequestOp() ReviewPaymentRequestOp {
	val, ok := u.GetReviewPaymentRequestOp()

	if !ok {
		panic("arm ReviewPaymentRequestOp is not set")
	}

	return val
}

// GetReviewPaymentRequestOp retrieves the ReviewPaymentRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetReviewPaymentRequestOp() (result ReviewPaymentRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewPaymentRequestOp" {
		result = *u.ReviewPaymentRequestOp
		ok = true
	}

	return
}

// MustManageAssetOp retrieves the ManageAssetOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageAssetOp() ManageAssetOp {
	val, ok := u.GetManageAssetOp()

	if !ok {
		panic("arm ManageAssetOp is not set")
	}

	return val
}

// GetManageAssetOp retrieves the ManageAssetOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageAssetOp() (result ManageAssetOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetOp" {
		result = *u.ManageAssetOp
		ok = true
	}

	return
}

// MustUploadPreemissionsOp retrieves the UploadPreemissionsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustUploadPreemissionsOp() UploadPreemissionsOp {
	val, ok := u.GetUploadPreemissionsOp()

	if !ok {
		panic("arm UploadPreemissionsOp is not set")
	}

	return val
}

// GetUploadPreemissionsOp retrieves the UploadPreemissionsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetUploadPreemissionsOp() (result UploadPreemissionsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "UploadPreemissionsOp" {
		result = *u.UploadPreemissionsOp
		ok = true
	}

	return
}

// MustSetLimitsOp retrieves the SetLimitsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetLimitsOp() SetLimitsOp {
	val, ok := u.GetSetLimitsOp()

	if !ok {
		panic("arm SetLimitsOp is not set")
	}

	return val
}

// GetSetLimitsOp retrieves the SetLimitsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetLimitsOp() (result SetLimitsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetLimitsOp" {
		result = *u.SetLimitsOp
		ok = true
	}

	return
}

// MustDirectDebitOp retrieves the DirectDebitOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustDirectDebitOp() DirectDebitOp {
	val, ok := u.GetDirectDebitOp()

	if !ok {
		panic("arm DirectDebitOp is not set")
	}

	return val
}

// GetDirectDebitOp retrieves the DirectDebitOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetDirectDebitOp() (result DirectDebitOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DirectDebitOp" {
		result = *u.DirectDebitOp
		ok = true
	}

	return
}

// MustManageAssetPairOp retrieves the ManageAssetPairOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageAssetPairOp() ManageAssetPairOp {
	val, ok := u.GetManageAssetPairOp()

	if !ok {
		panic("arm ManageAssetPairOp is not set")
	}

	return val
}

// GetManageAssetPairOp retrieves the ManageAssetPairOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageAssetPairOp() (result ManageAssetPairOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetPairOp" {
		result = *u.ManageAssetPairOp
		ok = true
	}

	return
}

// MustManageOfferOp retrieves the ManageOfferOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageOfferOp() ManageOfferOp {
	val, ok := u.GetManageOfferOp()

	if !ok {
		panic("arm ManageOfferOp is not set")
	}

	return val
}

// GetManageOfferOp retrieves the ManageOfferOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageOfferOp() (result ManageOfferOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageOfferOp" {
		result = *u.ManageOfferOp
		ok = true
	}

	return
}

// MustManageInvoiceOp retrieves the ManageInvoiceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageInvoiceOp() ManageInvoiceOp {
	val, ok := u.GetManageInvoiceOp()

	if !ok {
		panic("arm ManageInvoiceOp is not set")
	}

	return val
}

// GetManageInvoiceOp retrieves the ManageInvoiceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageInvoiceOp() (result ManageInvoiceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageInvoiceOp" {
		result = *u.ManageInvoiceOp
		ok = true
	}

	return
}

// Operation is an XDR Struct defines as:
//
//   struct Operation
//    {
//        // sourceAccount is the account used to run the operation
//        // if not set, the runtime defaults to "sourceAccount" specified at
//        // the transaction level
//        AccountID* sourceAccount;
//
//        union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountOp createAccountOp;
//        case PAYMENT:
//            PaymentOp paymentOp;
//        case SET_OPTIONS:
//            SetOptionsOp setOptionsOp;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestOp manageCoinsEmissionRequestOp;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestOp reviewCoinsEmissionRequestOp;
//        case SET_FEES:
//            SetFeesOp setFeesOp;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountOp manageAccountOp;
//    	case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestOp manageForfeitRequestOp;
//    	case RECOVER:
//    		RecoverOp recoverOp;
//    	case MANAGE_BALANCE:
//    		ManageBalanceOp manageBalanceOp;
//    	case REVIEW_PAYMENT_REQUEST:
//    		ReviewPaymentRequestOp reviewPaymentRequestOp;
//        case MANAGE_ASSET:
//            ManageAssetOp manageAssetOp;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsOp uploadPreemissionsOp;
//        case SET_LIMITS:
//            SetLimitsOp setLimitsOp;
//        case DIRECT_DEBIT:
//            DirectDebitOp directDebitOp;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairOp manageAssetPairOp;
//    	case MANAGE_OFFER:
//    		ManageOfferOp manageOfferOp;
//        case MANAGE_INVOICE:
//            ManageInvoiceOp manageInvoiceOp;
//        }
//        body;
//    };
//
type Operation struct {
	SourceAccount *AccountId    `json:"sourceAccount,omitempty"`
	Body          OperationBody `json:"body,omitempty"`
}

// MemoType is an XDR Enum defines as:
//
//   enum MemoType
//    {
//        MEMO_NONE = 0,
//        MEMO_TEXT = 1,
//        MEMO_ID = 2,
//        MEMO_HASH = 3,
//        MEMO_RETURN = 4
//    };
//
type MemoType int32

const (
	MemoTypeMemoNone   MemoType = 0
	MemoTypeMemoText   MemoType = 1
	MemoTypeMemoId     MemoType = 2
	MemoTypeMemoHash   MemoType = 3
	MemoTypeMemoReturn MemoType = 4
)

var MemoTypeAll = []MemoType{
	MemoTypeMemoNone,
	MemoTypeMemoText,
	MemoTypeMemoId,
	MemoTypeMemoHash,
	MemoTypeMemoReturn,
}

var memoTypeMap = map[int32]string{
	0: "MemoTypeMemoNone",
	1: "MemoTypeMemoText",
	2: "MemoTypeMemoId",
	3: "MemoTypeMemoHash",
	4: "MemoTypeMemoReturn",
}

var memoTypeShortMap = map[int32]string{
	0: "memo_none",
	1: "memo_text",
	2: "memo_id",
	3: "memo_hash",
	4: "memo_return",
}

var memoTypeRevMap = map[string]int32{
	"MemoTypeMemoNone":   0,
	"MemoTypeMemoText":   1,
	"MemoTypeMemoId":     2,
	"MemoTypeMemoHash":   3,
	"MemoTypeMemoReturn": 4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MemoType
func (e MemoType) ValidEnum(v int32) bool {
	_, ok := memoTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MemoType) String() string {
	name, _ := memoTypeMap[int32(e)]
	return name
}

func (e MemoType) ShortString() string {
	name, _ := memoTypeShortMap[int32(e)]
	return name
}

func (e MemoType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *MemoType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := memoTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = MemoType(value)
	return nil
}

// Memo is an XDR Union defines as:
//
//   union Memo switch (MemoType type)
//    {
//    case MEMO_NONE:
//        void;
//    case MEMO_TEXT:
//        string text<28>;
//    case MEMO_ID:
//        uint64 id;
//    case MEMO_HASH:
//        Hash hash; // the hash of what to pull from the content server
//    case MEMO_RETURN:
//        Hash retHash; // the hash of the tx you are rejecting
//    };
//
type Memo struct {
	Type    MemoType `json:"type,omitempty"`
	Text    *string  `json:"text,omitempty" xdrmaxsize:"28"`
	Id      *Uint64  `json:"id,omitempty"`
	Hash    *Hash    `json:"hash,omitempty"`
	RetHash *Hash    `json:"retHash,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Memo) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Memo
func (u Memo) ArmForSwitch(sw int32) (string, bool) {
	switch MemoType(sw) {
	case MemoTypeMemoNone:
		return "", true
	case MemoTypeMemoText:
		return "Text", true
	case MemoTypeMemoId:
		return "Id", true
	case MemoTypeMemoHash:
		return "Hash", true
	case MemoTypeMemoReturn:
		return "RetHash", true
	}
	return "-", false
}

// NewMemo creates a new  Memo.
func NewMemo(aType MemoType, value interface{}) (result Memo, err error) {
	result.Type = aType
	switch MemoType(aType) {
	case MemoTypeMemoNone:
		// void
	case MemoTypeMemoText:
		tv, ok := value.(string)
		if !ok {
			err = fmt.Errorf("invalid value, must be string")
			return
		}
		result.Text = &tv
	case MemoTypeMemoId:
		tv, ok := value.(Uint64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint64")
			return
		}
		result.Id = &tv
	case MemoTypeMemoHash:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.Hash = &tv
	case MemoTypeMemoReturn:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.RetHash = &tv
	}
	return
}

// MustText retrieves the Text value from the union,
// panicing if the value is not set.
func (u Memo) MustText() string {
	val, ok := u.GetText()

	if !ok {
		panic("arm Text is not set")
	}

	return val
}

// GetText retrieves the Text value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetText() (result string, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Text" {
		result = *u.Text
		ok = true
	}

	return
}

// MustId retrieves the Id value from the union,
// panicing if the value is not set.
func (u Memo) MustId() Uint64 {
	val, ok := u.GetId()

	if !ok {
		panic("arm Id is not set")
	}

	return val
}

// GetId retrieves the Id value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetId() (result Uint64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Id" {
		result = *u.Id
		ok = true
	}

	return
}

// MustHash retrieves the Hash value from the union,
// panicing if the value is not set.
func (u Memo) MustHash() Hash {
	val, ok := u.GetHash()

	if !ok {
		panic("arm Hash is not set")
	}

	return val
}

// GetHash retrieves the Hash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hash" {
		result = *u.Hash
		ok = true
	}

	return
}

// MustRetHash retrieves the RetHash value from the union,
// panicing if the value is not set.
func (u Memo) MustRetHash() Hash {
	val, ok := u.GetRetHash()

	if !ok {
		panic("arm RetHash is not set")
	}

	return val
}

// GetRetHash retrieves the RetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetRetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RetHash" {
		result = *u.RetHash
		ok = true
	}

	return
}

// TimeBounds is an XDR Struct defines as:
//
//   struct TimeBounds
//    {
//        uint64 minTime;
//        uint64 maxTime; // 0 here means no maxTime
//    };
//
type TimeBounds struct {
	MinTime Uint64 `json:"minTime,omitempty"`
	MaxTime Uint64 `json:"maxTime,omitempty"`
}

// TransactionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionExt
func (u TransactionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionExt creates a new  TransactionExt.
func NewTransactionExt(v LedgerVersion, value interface{}) (result TransactionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// Transaction is an XDR Struct defines as:
//
//   struct Transaction
//    {
//        // account used to run the transaction
//        AccountID sourceAccount;
//
//        Salt salt;
//
//        // validity range (inclusive) for the last ledger close time
//        TimeBounds timeBounds;
//
//        Memo memo;
//
//        Operation operations<100>;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type Transaction struct {
	SourceAccount AccountId      `json:"sourceAccount,omitempty"`
	Salt          Salt           `json:"salt,omitempty"`
	TimeBounds    TimeBounds     `json:"timeBounds,omitempty"`
	Memo          Memo           `json:"memo,omitempty"`
	Operations    []Operation    `json:"operations,omitempty" xdrmaxsize:"100"`
	Ext           TransactionExt `json:"ext,omitempty"`
}

// TransactionEnvelope is an XDR Struct defines as:
//
//   struct TransactionEnvelope
//    {
//        Transaction tx;
//        DecoratedSignature signatures<20>;
//    };
//
type TransactionEnvelope struct {
	Tx         Transaction          `json:"tx,omitempty"`
	Signatures []DecoratedSignature `json:"signatures,omitempty" xdrmaxsize:"20"`
}

// OperationResultCode is an XDR Enum defines as:
//
//   enum OperationResultCode
//    {
//        opINNER = 0, // inner object result is valid
//
//        opBAD_AUTH = -1,      // too few valid signatures / wrong network
//        opNO_ACCOUNT = -2,    // source account was not found
//    	opNOT_ALLOWED = -3,   // operation is not allowed for this type of source account
//    	opACCOUNT_BLOCKED = -4, // account is blocked
//        opNO_COUNTERPARTY = -5,
//        opCOUNTERPARTY_BLOCKED = -6,
//        opCOUNTERPARTY_WRONG_TYPE = -7,
//    	opBAD_AUTH_EXTRA = -8
//    };
//
type OperationResultCode int32

const (
	OperationResultCodeOpInner                 OperationResultCode = 0
	OperationResultCodeOpBadAuth               OperationResultCode = -1
	OperationResultCodeOpNoAccount             OperationResultCode = -2
	OperationResultCodeOpNotAllowed            OperationResultCode = -3
	OperationResultCodeOpAccountBlocked        OperationResultCode = -4
	OperationResultCodeOpNoCounterparty        OperationResultCode = -5
	OperationResultCodeOpCounterpartyBlocked   OperationResultCode = -6
	OperationResultCodeOpCounterpartyWrongType OperationResultCode = -7
	OperationResultCodeOpBadAuthExtra          OperationResultCode = -8
)

var OperationResultCodeAll = []OperationResultCode{
	OperationResultCodeOpInner,
	OperationResultCodeOpBadAuth,
	OperationResultCodeOpNoAccount,
	OperationResultCodeOpNotAllowed,
	OperationResultCodeOpAccountBlocked,
	OperationResultCodeOpNoCounterparty,
	OperationResultCodeOpCounterpartyBlocked,
	OperationResultCodeOpCounterpartyWrongType,
	OperationResultCodeOpBadAuthExtra,
}

var operationResultCodeMap = map[int32]string{
	0:  "OperationResultCodeOpInner",
	-1: "OperationResultCodeOpBadAuth",
	-2: "OperationResultCodeOpNoAccount",
	-3: "OperationResultCodeOpNotAllowed",
	-4: "OperationResultCodeOpAccountBlocked",
	-5: "OperationResultCodeOpNoCounterparty",
	-6: "OperationResultCodeOpCounterpartyBlocked",
	-7: "OperationResultCodeOpCounterpartyWrongType",
	-8: "OperationResultCodeOpBadAuthExtra",
}

var operationResultCodeShortMap = map[int32]string{
	0:  "op_inner",
	-1: "op_bad_auth",
	-2: "op_no_account",
	-3: "op_not_allowed",
	-4: "op_account_blocked",
	-5: "op_no_counterparty",
	-6: "op_counterparty_blocked",
	-7: "op_counterparty_wrong_type",
	-8: "op_bad_auth_extra",
}

var operationResultCodeRevMap = map[string]int32{
	"OperationResultCodeOpInner":                 0,
	"OperationResultCodeOpBadAuth":               -1,
	"OperationResultCodeOpNoAccount":             -2,
	"OperationResultCodeOpNotAllowed":            -3,
	"OperationResultCodeOpAccountBlocked":        -4,
	"OperationResultCodeOpNoCounterparty":        -5,
	"OperationResultCodeOpCounterpartyBlocked":   -6,
	"OperationResultCodeOpCounterpartyWrongType": -7,
	"OperationResultCodeOpBadAuthExtra":          -8,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationResultCode
func (e OperationResultCode) ValidEnum(v int32) bool {
	_, ok := operationResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationResultCode) String() string {
	name, _ := operationResultCodeMap[int32(e)]
	return name
}

func (e OperationResultCode) ShortString() string {
	name, _ := operationResultCodeShortMap[int32(e)]
	return name
}

func (e OperationResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *OperationResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := operationResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = OperationResultCode(value)
	return nil
}

// OperationResultTr is an XDR NestedUnion defines as:
//
//   union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountResult createAccountResult;
//        case PAYMENT:
//            PaymentResult paymentResult;
//        case SET_OPTIONS:
//            SetOptionsResult setOptionsResult;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestResult manageCoinsEmissionRequestResult;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestResult reviewCoinsEmissionRequestResult;
//        case SET_FEES:
//            SetFeesResult setFeesResult;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountResult manageAccountResult;
//        case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestResult manageForfeitRequestResult;
//        case RECOVER:
//    		RecoverResult recoverResult;
//        case MANAGE_BALANCE:
//            ManageBalanceResult manageBalanceResult;
//        case REVIEW_PAYMENT_REQUEST:
//            ReviewPaymentRequestResult reviewPaymentRequestResult;
//        case MANAGE_ASSET:
//            ManageAssetResult manageAssetResult;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsResult uploadPreemissionsResult;
//        case SET_LIMITS:
//            SetLimitsResult setLimitsResult;
//        case DIRECT_DEBIT:
//            DirectDebitResult directDebitResult;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairResult manageAssetPairResult;
//    	case MANAGE_OFFER:
//    		ManageOfferResult manageOfferResult;
//    	case MANAGE_INVOICE:
//    		ManageInvoiceResult manageInvoiceResult;
//        }
//
type OperationResultTr struct {
	Type                             OperationType                     `json:"type,omitempty"`
	CreateAccountResult              *CreateAccountResult              `json:"createAccountResult,omitempty"`
	PaymentResult                    *PaymentResult                    `json:"paymentResult,omitempty"`
	SetOptionsResult                 *SetOptionsResult                 `json:"setOptionsResult,omitempty"`
	ManageCoinsEmissionRequestResult *ManageCoinsEmissionRequestResult `json:"manageCoinsEmissionRequestResult,omitempty"`
	ReviewCoinsEmissionRequestResult *ReviewCoinsEmissionRequestResult `json:"reviewCoinsEmissionRequestResult,omitempty"`
	SetFeesResult                    *SetFeesResult                    `json:"setFeesResult,omitempty"`
	ManageAccountResult              *ManageAccountResult              `json:"manageAccountResult,omitempty"`
	ManageForfeitRequestResult       *ManageForfeitRequestResult       `json:"manageForfeitRequestResult,omitempty"`
	RecoverResult                    *RecoverResult                    `json:"recoverResult,omitempty"`
	ManageBalanceResult              *ManageBalanceResult              `json:"manageBalanceResult,omitempty"`
	ReviewPaymentRequestResult       *ReviewPaymentRequestResult       `json:"reviewPaymentRequestResult,omitempty"`
	ManageAssetResult                *ManageAssetResult                `json:"manageAssetResult,omitempty"`
	UploadPreemissionsResult         *UploadPreemissionsResult         `json:"uploadPreemissionsResult,omitempty"`
	SetLimitsResult                  *SetLimitsResult                  `json:"setLimitsResult,omitempty"`
	DirectDebitResult                *DirectDebitResult                `json:"directDebitResult,omitempty"`
	ManageAssetPairResult            *ManageAssetPairResult            `json:"manageAssetPairResult,omitempty"`
	ManageOfferResult                *ManageOfferResult                `json:"manageOfferResult,omitempty"`
	ManageInvoiceResult              *ManageInvoiceResult              `json:"manageInvoiceResult,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationResultTr) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationResultTr
func (u OperationResultTr) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypeCreateAccount:
		return "CreateAccountResult", true
	case OperationTypePayment:
		return "PaymentResult", true
	case OperationTypeSetOptions:
		return "SetOptionsResult", true
	case OperationTypeManageCoinsEmissionRequest:
		return "ManageCoinsEmissionRequestResult", true
	case OperationTypeReviewCoinsEmissionRequest:
		return "ReviewCoinsEmissionRequestResult", true
	case OperationTypeSetFees:
		return "SetFeesResult", true
	case OperationTypeManageAccount:
		return "ManageAccountResult", true
	case OperationTypeManageForfeitRequest:
		return "ManageForfeitRequestResult", true
	case OperationTypeRecover:
		return "RecoverResult", true
	case OperationTypeManageBalance:
		return "ManageBalanceResult", true
	case OperationTypeReviewPaymentRequest:
		return "ReviewPaymentRequestResult", true
	case OperationTypeManageAsset:
		return "ManageAssetResult", true
	case OperationTypeUploadPreemissions:
		return "UploadPreemissionsResult", true
	case OperationTypeSetLimits:
		return "SetLimitsResult", true
	case OperationTypeDirectDebit:
		return "DirectDebitResult", true
	case OperationTypeManageAssetPair:
		return "ManageAssetPairResult", true
	case OperationTypeManageOffer:
		return "ManageOfferResult", true
	case OperationTypeManageInvoice:
		return "ManageInvoiceResult", true
	}
	return "-", false
}

// NewOperationResultTr creates a new  OperationResultTr.
func NewOperationResultTr(aType OperationType, value interface{}) (result OperationResultTr, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypeCreateAccount:
		tv, ok := value.(CreateAccountResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountResult")
			return
		}
		result.CreateAccountResult = &tv
	case OperationTypePayment:
		tv, ok := value.(PaymentResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentResult")
			return
		}
		result.PaymentResult = &tv
	case OperationTypeSetOptions:
		tv, ok := value.(SetOptionsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsResult")
			return
		}
		result.SetOptionsResult = &tv
	case OperationTypeManageCoinsEmissionRequest:
		tv, ok := value.(ManageCoinsEmissionRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageCoinsEmissionRequestResult")
			return
		}
		result.ManageCoinsEmissionRequestResult = &tv
	case OperationTypeReviewCoinsEmissionRequest:
		tv, ok := value.(ReviewCoinsEmissionRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewCoinsEmissionRequestResult")
			return
		}
		result.ReviewCoinsEmissionRequestResult = &tv
	case OperationTypeSetFees:
		tv, ok := value.(SetFeesResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetFeesResult")
			return
		}
		result.SetFeesResult = &tv
	case OperationTypeManageAccount:
		tv, ok := value.(ManageAccountResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAccountResult")
			return
		}
		result.ManageAccountResult = &tv
	case OperationTypeManageForfeitRequest:
		tv, ok := value.(ManageForfeitRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestResult")
			return
		}
		result.ManageForfeitRequestResult = &tv
	case OperationTypeRecover:
		tv, ok := value.(RecoverResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be RecoverResult")
			return
		}
		result.RecoverResult = &tv
	case OperationTypeManageBalance:
		tv, ok := value.(ManageBalanceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBalanceResult")
			return
		}
		result.ManageBalanceResult = &tv
	case OperationTypeReviewPaymentRequest:
		tv, ok := value.(ReviewPaymentRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewPaymentRequestResult")
			return
		}
		result.ReviewPaymentRequestResult = &tv
	case OperationTypeManageAsset:
		tv, ok := value.(ManageAssetResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetResult")
			return
		}
		result.ManageAssetResult = &tv
	case OperationTypeUploadPreemissions:
		tv, ok := value.(UploadPreemissionsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be UploadPreemissionsResult")
			return
		}
		result.UploadPreemissionsResult = &tv
	case OperationTypeSetLimits:
		tv, ok := value.(SetLimitsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetLimitsResult")
			return
		}
		result.SetLimitsResult = &tv
	case OperationTypeDirectDebit:
		tv, ok := value.(DirectDebitResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be DirectDebitResult")
			return
		}
		result.DirectDebitResult = &tv
	case OperationTypeManageAssetPair:
		tv, ok := value.(ManageAssetPairResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetPairResult")
			return
		}
		result.ManageAssetPairResult = &tv
	case OperationTypeManageOffer:
		tv, ok := value.(ManageOfferResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferResult")
			return
		}
		result.ManageOfferResult = &tv
	case OperationTypeManageInvoice:
		tv, ok := value.(ManageInvoiceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageInvoiceResult")
			return
		}
		result.ManageInvoiceResult = &tv
	}
	return
}

// MustCreateAccountResult retrieves the CreateAccountResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustCreateAccountResult() CreateAccountResult {
	val, ok := u.GetCreateAccountResult()

	if !ok {
		panic("arm CreateAccountResult is not set")
	}

	return val
}

// GetCreateAccountResult retrieves the CreateAccountResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetCreateAccountResult() (result CreateAccountResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateAccountResult" {
		result = *u.CreateAccountResult
		ok = true
	}

	return
}

// MustPaymentResult retrieves the PaymentResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustPaymentResult() PaymentResult {
	val, ok := u.GetPaymentResult()

	if !ok {
		panic("arm PaymentResult is not set")
	}

	return val
}

// GetPaymentResult retrieves the PaymentResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetPaymentResult() (result PaymentResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentResult" {
		result = *u.PaymentResult
		ok = true
	}

	return
}

// MustSetOptionsResult retrieves the SetOptionsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetOptionsResult() SetOptionsResult {
	val, ok := u.GetSetOptionsResult()

	if !ok {
		panic("arm SetOptionsResult is not set")
	}

	return val
}

// GetSetOptionsResult retrieves the SetOptionsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetOptionsResult() (result SetOptionsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetOptionsResult" {
		result = *u.SetOptionsResult
		ok = true
	}

	return
}

// MustManageCoinsEmissionRequestResult retrieves the ManageCoinsEmissionRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageCoinsEmissionRequestResult() ManageCoinsEmissionRequestResult {
	val, ok := u.GetManageCoinsEmissionRequestResult()

	if !ok {
		panic("arm ManageCoinsEmissionRequestResult is not set")
	}

	return val
}

// GetManageCoinsEmissionRequestResult retrieves the ManageCoinsEmissionRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageCoinsEmissionRequestResult() (result ManageCoinsEmissionRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageCoinsEmissionRequestResult" {
		result = *u.ManageCoinsEmissionRequestResult
		ok = true
	}

	return
}

// MustReviewCoinsEmissionRequestResult retrieves the ReviewCoinsEmissionRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustReviewCoinsEmissionRequestResult() ReviewCoinsEmissionRequestResult {
	val, ok := u.GetReviewCoinsEmissionRequestResult()

	if !ok {
		panic("arm ReviewCoinsEmissionRequestResult is not set")
	}

	return val
}

// GetReviewCoinsEmissionRequestResult retrieves the ReviewCoinsEmissionRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetReviewCoinsEmissionRequestResult() (result ReviewCoinsEmissionRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewCoinsEmissionRequestResult" {
		result = *u.ReviewCoinsEmissionRequestResult
		ok = true
	}

	return
}

// MustSetFeesResult retrieves the SetFeesResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetFeesResult() SetFeesResult {
	val, ok := u.GetSetFeesResult()

	if !ok {
		panic("arm SetFeesResult is not set")
	}

	return val
}

// GetSetFeesResult retrieves the SetFeesResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetFeesResult() (result SetFeesResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetFeesResult" {
		result = *u.SetFeesResult
		ok = true
	}

	return
}

// MustManageAccountResult retrieves the ManageAccountResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageAccountResult() ManageAccountResult {
	val, ok := u.GetManageAccountResult()

	if !ok {
		panic("arm ManageAccountResult is not set")
	}

	return val
}

// GetManageAccountResult retrieves the ManageAccountResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageAccountResult() (result ManageAccountResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAccountResult" {
		result = *u.ManageAccountResult
		ok = true
	}

	return
}

// MustManageForfeitRequestResult retrieves the ManageForfeitRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageForfeitRequestResult() ManageForfeitRequestResult {
	val, ok := u.GetManageForfeitRequestResult()

	if !ok {
		panic("arm ManageForfeitRequestResult is not set")
	}

	return val
}

// GetManageForfeitRequestResult retrieves the ManageForfeitRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageForfeitRequestResult() (result ManageForfeitRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageForfeitRequestResult" {
		result = *u.ManageForfeitRequestResult
		ok = true
	}

	return
}

// MustRecoverResult retrieves the RecoverResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustRecoverResult() RecoverResult {
	val, ok := u.GetRecoverResult()

	if !ok {
		panic("arm RecoverResult is not set")
	}

	return val
}

// GetRecoverResult retrieves the RecoverResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetRecoverResult() (result RecoverResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RecoverResult" {
		result = *u.RecoverResult
		ok = true
	}

	return
}

// MustManageBalanceResult retrieves the ManageBalanceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageBalanceResult() ManageBalanceResult {
	val, ok := u.GetManageBalanceResult()

	if !ok {
		panic("arm ManageBalanceResult is not set")
	}

	return val
}

// GetManageBalanceResult retrieves the ManageBalanceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageBalanceResult() (result ManageBalanceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageBalanceResult" {
		result = *u.ManageBalanceResult
		ok = true
	}

	return
}

// MustReviewPaymentRequestResult retrieves the ReviewPaymentRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustReviewPaymentRequestResult() ReviewPaymentRequestResult {
	val, ok := u.GetReviewPaymentRequestResult()

	if !ok {
		panic("arm ReviewPaymentRequestResult is not set")
	}

	return val
}

// GetReviewPaymentRequestResult retrieves the ReviewPaymentRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetReviewPaymentRequestResult() (result ReviewPaymentRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewPaymentRequestResult" {
		result = *u.ReviewPaymentRequestResult
		ok = true
	}

	return
}

// MustManageAssetResult retrieves the ManageAssetResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageAssetResult() ManageAssetResult {
	val, ok := u.GetManageAssetResult()

	if !ok {
		panic("arm ManageAssetResult is not set")
	}

	return val
}

// GetManageAssetResult retrieves the ManageAssetResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageAssetResult() (result ManageAssetResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetResult" {
		result = *u.ManageAssetResult
		ok = true
	}

	return
}

// MustUploadPreemissionsResult retrieves the UploadPreemissionsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustUploadPreemissionsResult() UploadPreemissionsResult {
	val, ok := u.GetUploadPreemissionsResult()

	if !ok {
		panic("arm UploadPreemissionsResult is not set")
	}

	return val
}

// GetUploadPreemissionsResult retrieves the UploadPreemissionsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetUploadPreemissionsResult() (result UploadPreemissionsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "UploadPreemissionsResult" {
		result = *u.UploadPreemissionsResult
		ok = true
	}

	return
}

// MustSetLimitsResult retrieves the SetLimitsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetLimitsResult() SetLimitsResult {
	val, ok := u.GetSetLimitsResult()

	if !ok {
		panic("arm SetLimitsResult is not set")
	}

	return val
}

// GetSetLimitsResult retrieves the SetLimitsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetLimitsResult() (result SetLimitsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetLimitsResult" {
		result = *u.SetLimitsResult
		ok = true
	}

	return
}

// MustDirectDebitResult retrieves the DirectDebitResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustDirectDebitResult() DirectDebitResult {
	val, ok := u.GetDirectDebitResult()

	if !ok {
		panic("arm DirectDebitResult is not set")
	}

	return val
}

// GetDirectDebitResult retrieves the DirectDebitResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetDirectDebitResult() (result DirectDebitResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DirectDebitResult" {
		result = *u.DirectDebitResult
		ok = true
	}

	return
}

// MustManageAssetPairResult retrieves the ManageAssetPairResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageAssetPairResult() ManageAssetPairResult {
	val, ok := u.GetManageAssetPairResult()

	if !ok {
		panic("arm ManageAssetPairResult is not set")
	}

	return val
}

// GetManageAssetPairResult retrieves the ManageAssetPairResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageAssetPairResult() (result ManageAssetPairResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetPairResult" {
		result = *u.ManageAssetPairResult
		ok = true
	}

	return
}

// MustManageOfferResult retrieves the ManageOfferResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageOfferResult() ManageOfferResult {
	val, ok := u.GetManageOfferResult()

	if !ok {
		panic("arm ManageOfferResult is not set")
	}

	return val
}

// GetManageOfferResult retrieves the ManageOfferResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageOfferResult() (result ManageOfferResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageOfferResult" {
		result = *u.ManageOfferResult
		ok = true
	}

	return
}

// MustManageInvoiceResult retrieves the ManageInvoiceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageInvoiceResult() ManageInvoiceResult {
	val, ok := u.GetManageInvoiceResult()

	if !ok {
		panic("arm ManageInvoiceResult is not set")
	}

	return val
}

// GetManageInvoiceResult retrieves the ManageInvoiceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageInvoiceResult() (result ManageInvoiceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageInvoiceResult" {
		result = *u.ManageInvoiceResult
		ok = true
	}

	return
}

// OperationResult is an XDR Union defines as:
//
//   union OperationResult switch (OperationResultCode code)
//    {
//    case opINNER:
//        union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountResult createAccountResult;
//        case PAYMENT:
//            PaymentResult paymentResult;
//        case SET_OPTIONS:
//            SetOptionsResult setOptionsResult;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestResult manageCoinsEmissionRequestResult;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestResult reviewCoinsEmissionRequestResult;
//        case SET_FEES:
//            SetFeesResult setFeesResult;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountResult manageAccountResult;
//        case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestResult manageForfeitRequestResult;
//        case RECOVER:
//    		RecoverResult recoverResult;
//        case MANAGE_BALANCE:
//            ManageBalanceResult manageBalanceResult;
//        case REVIEW_PAYMENT_REQUEST:
//            ReviewPaymentRequestResult reviewPaymentRequestResult;
//        case MANAGE_ASSET:
//            ManageAssetResult manageAssetResult;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsResult uploadPreemissionsResult;
//        case SET_LIMITS:
//            SetLimitsResult setLimitsResult;
//        case DIRECT_DEBIT:
//            DirectDebitResult directDebitResult;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairResult manageAssetPairResult;
//    	case MANAGE_OFFER:
//    		ManageOfferResult manageOfferResult;
//    	case MANAGE_INVOICE:
//    		ManageInvoiceResult manageInvoiceResult;
//        }
//        tr;
//    default:
//        void;
//    };
//
type OperationResult struct {
	Code OperationResultCode `json:"code,omitempty"`
	Tr   *OperationResultTr  `json:"tr,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationResult
func (u OperationResult) ArmForSwitch(sw int32) (string, bool) {
	switch OperationResultCode(sw) {
	case OperationResultCodeOpInner:
		return "Tr", true
	default:
		return "", true
	}
}

// NewOperationResult creates a new  OperationResult.
func NewOperationResult(code OperationResultCode, value interface{}) (result OperationResult, err error) {
	result.Code = code
	switch OperationResultCode(code) {
	case OperationResultCodeOpInner:
		tv, ok := value.(OperationResultTr)
		if !ok {
			err = fmt.Errorf("invalid value, must be OperationResultTr")
			return
		}
		result.Tr = &tv
	default:
		// void
	}
	return
}

// MustTr retrieves the Tr value from the union,
// panicing if the value is not set.
func (u OperationResult) MustTr() OperationResultTr {
	val, ok := u.GetTr()

	if !ok {
		panic("arm Tr is not set")
	}

	return val
}

// GetTr retrieves the Tr value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResult) GetTr() (result OperationResultTr, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Tr" {
		result = *u.Tr
		ok = true
	}

	return
}

// TransactionResultCode is an XDR Enum defines as:
//
//   enum TransactionResultCode
//    {
//        txSUCCESS = 0, // all operations succeeded
//
//        txFAILED = -1, // one of the operations failed (none were applied)
//
//        txTOO_EARLY = -2,         // ledger closeTime before minTime
//        txTOO_LATE = -3,          // ledger closeTime after maxTime
//        txMISSING_OPERATION = -4, // no operation was specified
//
//        txBAD_AUTH = -5,             // too few valid signatures / wrong network
//        txNO_ACCOUNT = -6,           // source account not found
//        txBAD_AUTH_EXTRA = -7,      // unused signatures attached to transaction
//        txINTERNAL_ERROR = -8,      // an unknown error occured
//    	txACCOUNT_BLOCKED = -9,     // account is blocked and cannot be source of tx
//        txDUPLICATION = -10         // if timing is stored
//    };
//
type TransactionResultCode int32

const (
	TransactionResultCodeTxSuccess          TransactionResultCode = 0
	TransactionResultCodeTxFailed           TransactionResultCode = -1
	TransactionResultCodeTxTooEarly         TransactionResultCode = -2
	TransactionResultCodeTxTooLate          TransactionResultCode = -3
	TransactionResultCodeTxMissingOperation TransactionResultCode = -4
	TransactionResultCodeTxBadAuth          TransactionResultCode = -5
	TransactionResultCodeTxNoAccount        TransactionResultCode = -6
	TransactionResultCodeTxBadAuthExtra     TransactionResultCode = -7
	TransactionResultCodeTxInternalError    TransactionResultCode = -8
	TransactionResultCodeTxAccountBlocked   TransactionResultCode = -9
	TransactionResultCodeTxDuplication      TransactionResultCode = -10
)

var TransactionResultCodeAll = []TransactionResultCode{
	TransactionResultCodeTxSuccess,
	TransactionResultCodeTxFailed,
	TransactionResultCodeTxTooEarly,
	TransactionResultCodeTxTooLate,
	TransactionResultCodeTxMissingOperation,
	TransactionResultCodeTxBadAuth,
	TransactionResultCodeTxNoAccount,
	TransactionResultCodeTxBadAuthExtra,
	TransactionResultCodeTxInternalError,
	TransactionResultCodeTxAccountBlocked,
	TransactionResultCodeTxDuplication,
}

var transactionResultCodeMap = map[int32]string{
	0:   "TransactionResultCodeTxSuccess",
	-1:  "TransactionResultCodeTxFailed",
	-2:  "TransactionResultCodeTxTooEarly",
	-3:  "TransactionResultCodeTxTooLate",
	-4:  "TransactionResultCodeTxMissingOperation",
	-5:  "TransactionResultCodeTxBadAuth",
	-6:  "TransactionResultCodeTxNoAccount",
	-7:  "TransactionResultCodeTxBadAuthExtra",
	-8:  "TransactionResultCodeTxInternalError",
	-9:  "TransactionResultCodeTxAccountBlocked",
	-10: "TransactionResultCodeTxDuplication",
}

var transactionResultCodeShortMap = map[int32]string{
	0:   "tx_success",
	-1:  "tx_failed",
	-2:  "tx_too_early",
	-3:  "tx_too_late",
	-4:  "tx_missing_operation",
	-5:  "tx_bad_auth",
	-6:  "tx_no_account",
	-7:  "tx_bad_auth_extra",
	-8:  "tx_internal_error",
	-9:  "tx_account_blocked",
	-10: "tx_duplication",
}

var transactionResultCodeRevMap = map[string]int32{
	"TransactionResultCodeTxSuccess":          0,
	"TransactionResultCodeTxFailed":           -1,
	"TransactionResultCodeTxTooEarly":         -2,
	"TransactionResultCodeTxTooLate":          -3,
	"TransactionResultCodeTxMissingOperation": -4,
	"TransactionResultCodeTxBadAuth":          -5,
	"TransactionResultCodeTxNoAccount":        -6,
	"TransactionResultCodeTxBadAuthExtra":     -7,
	"TransactionResultCodeTxInternalError":    -8,
	"TransactionResultCodeTxAccountBlocked":   -9,
	"TransactionResultCodeTxDuplication":      -10,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for TransactionResultCode
func (e TransactionResultCode) ValidEnum(v int32) bool {
	_, ok := transactionResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e TransactionResultCode) String() string {
	name, _ := transactionResultCodeMap[int32(e)]
	return name
}

func (e TransactionResultCode) ShortString() string {
	name, _ := transactionResultCodeShortMap[int32(e)]
	return name
}

func (e TransactionResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *TransactionResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := transactionResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = TransactionResultCode(value)
	return nil
}

// TransactionResultResult is an XDR NestedUnion defines as:
//
//   union switch (TransactionResultCode code)
//        {
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        default:
//            void;
//        }
//
type TransactionResultResult struct {
	Code    TransactionResultCode `json:"code,omitempty"`
	Results *[]OperationResult    `json:"results,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionResultResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionResultResult
func (u TransactionResultResult) ArmForSwitch(sw int32) (string, bool) {
	switch TransactionResultCode(sw) {
	case TransactionResultCodeTxSuccess:
		return "Results", true
	case TransactionResultCodeTxFailed:
		return "Results", true
	default:
		return "", true
	}
}

// NewTransactionResultResult creates a new  TransactionResultResult.
func NewTransactionResultResult(code TransactionResultCode, value interface{}) (result TransactionResultResult, err error) {
	result.Code = code
	switch TransactionResultCode(code) {
	case TransactionResultCodeTxSuccess:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	case TransactionResultCodeTxFailed:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	default:
		// void
	}
	return
}

// MustResults retrieves the Results value from the union,
// panicing if the value is not set.
func (u TransactionResultResult) MustResults() []OperationResult {
	val, ok := u.GetResults()

	if !ok {
		panic("arm Results is not set")
	}

	return val
}

// GetResults retrieves the Results value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionResultResult) GetResults() (result []OperationResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Results" {
		result = *u.Results
		ok = true
	}

	return
}

// TransactionResultExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionResultExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionResultExt
func (u TransactionResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionResultExt creates a new  TransactionResultExt.
func NewTransactionResultExt(v LedgerVersion, value interface{}) (result TransactionResultExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TransactionResult is an XDR Struct defines as:
//
//   struct TransactionResult
//    {
//        int64 feeCharged; // actual fee charged for the transaction
//
//        union switch (TransactionResultCode code)
//        {
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        default:
//            void;
//        }
//        result;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TransactionResult struct {
	FeeCharged Int64                   `json:"feeCharged,omitempty"`
	Result     TransactionResultResult `json:"result,omitempty"`
	Ext        TransactionResultExt    `json:"ext,omitempty"`
}

// RequestType is an XDR Enum defines as:
//
//   enum RequestType
//    {
//        REQUEST_TYPE_SALE = 0,
//        REQUEST_TYPE_WITHDRAWAL = 1,
//        REQUEST_TYPE_REDEEM = 2,
//        REQUEST_TYPE_PAYMENT = 3
//    };
//
type RequestType int32

const (
	RequestTypeRequestTypeSale       RequestType = 0
	RequestTypeRequestTypeWithdrawal RequestType = 1
	RequestTypeRequestTypeRedeem     RequestType = 2
	RequestTypeRequestTypePayment    RequestType = 3
)

var RequestTypeAll = []RequestType{
	RequestTypeRequestTypeSale,
	RequestTypeRequestTypeWithdrawal,
	RequestTypeRequestTypeRedeem,
	RequestTypeRequestTypePayment,
}

var requestTypeMap = map[int32]string{
	0: "RequestTypeRequestTypeSale",
	1: "RequestTypeRequestTypeWithdrawal",
	2: "RequestTypeRequestTypeRedeem",
	3: "RequestTypeRequestTypePayment",
}

var requestTypeShortMap = map[int32]string{
	0: "request_type_sale",
	1: "request_type_withdrawal",
	2: "request_type_redeem",
	3: "request_type_payment",
}

var requestTypeRevMap = map[string]int32{
	"RequestTypeRequestTypeSale":       0,
	"RequestTypeRequestTypeWithdrawal": 1,
	"RequestTypeRequestTypeRedeem":     2,
	"RequestTypeRequestTypePayment":    3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for RequestType
func (e RequestType) ValidEnum(v int32) bool {
	_, ok := requestTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e RequestType) String() string {
	name, _ := requestTypeMap[int32(e)]
	return name
}

func (e RequestType) ShortString() string {
	name, _ := requestTypeShortMap[int32(e)]
	return name
}

func (e RequestType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *RequestType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := requestTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = RequestType(value)
	return nil
}

// PaymentRequestEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentRequestEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentRequestEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentRequestEntryExt
func (u PaymentRequestEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentRequestEntryExt creates a new  PaymentRequestEntryExt.
func NewPaymentRequestEntryExt(v LedgerVersion, value interface{}) (result PaymentRequestEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentRequestEntry is an XDR Struct defines as:
//
//   struct PaymentRequestEntry
//    {
//        uint64 paymentID;
//        BalanceID sourceBalance;
//        BalanceID* destinationBalance;
//        int64 sourceSend;
//        int64 sourceSendUniversal;
//        int64 destinationReceive;
//
//        uint64 createdAt;
//
//        uint64* invoiceID;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentRequestEntry struct {
	PaymentId           Uint64                 `json:"paymentID,omitempty"`
	SourceBalance       BalanceId              `json:"sourceBalance,omitempty"`
	DestinationBalance  *BalanceId             `json:"destinationBalance,omitempty"`
	SourceSend          Int64                  `json:"sourceSend,omitempty"`
	SourceSendUniversal Int64                  `json:"sourceSendUniversal,omitempty"`
	DestinationReceive  Int64                  `json:"destinationReceive,omitempty"`
	CreatedAt           Uint64                 `json:"createdAt,omitempty"`
	InvoiceId           *Uint64                `json:"invoiceID,omitempty"`
	Ext                 PaymentRequestEntryExt `json:"ext,omitempty"`
}

// ManageAssetAction is an XDR Enum defines as:
//
//   enum ManageAssetAction
//    {
//        CREATE = 0,
//        UPDATE_POLICIES = 1
//    };
//
type ManageAssetAction int32

const (
	ManageAssetActionCreate         ManageAssetAction = 0
	ManageAssetActionUpdatePolicies ManageAssetAction = 1
)

var ManageAssetActionAll = []ManageAssetAction{
	ManageAssetActionCreate,
	ManageAssetActionUpdatePolicies,
}

var manageAssetActionMap = map[int32]string{
	0: "ManageAssetActionCreate",
	1: "ManageAssetActionUpdatePolicies",
}

var manageAssetActionShortMap = map[int32]string{
	0: "create",
	1: "update_policies",
}

var manageAssetActionRevMap = map[string]int32{
	"ManageAssetActionCreate":         0,
	"ManageAssetActionUpdatePolicies": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetAction
func (e ManageAssetAction) ValidEnum(v int32) bool {
	_, ok := manageAssetActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetAction) String() string {
	name, _ := manageAssetActionMap[int32(e)]
	return name
}

func (e ManageAssetAction) ShortString() string {
	name, _ := manageAssetActionShortMap[int32(e)]
	return name
}

func (e ManageAssetAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageAssetAction) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageAssetActionRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageAssetAction(value)
	return nil
}

// ManageAssetOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetOpExt
func (u ManageAssetOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetOpExt creates a new  ManageAssetOpExt.
func NewManageAssetOpExt(v LedgerVersion, value interface{}) (result ManageAssetOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetOp is an XDR Struct defines as:
//
//   struct ManageAssetOp
//    {
//        ManageAssetAction action;
//    	AssetCode code;
//
//        int32 policies;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetOp struct {
	Action   ManageAssetAction `json:"action,omitempty"`
	Code     AssetCode         `json:"code,omitempty"`
	Policies Int32             `json:"policies,omitempty"`
	Ext      ManageAssetOpExt  `json:"ext,omitempty"`
}

// ManageAssetResultCode is an XDR Enum defines as:
//
//   enum ManageAssetResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	NOT_FOUND = -1,           // failed to find asset with such code
//    	ALREADY_EXISTS = -2,
//        MALFORMED = -3
//    };
//
type ManageAssetResultCode int32

const (
	ManageAssetResultCodeSuccess       ManageAssetResultCode = 0
	ManageAssetResultCodeNotFound      ManageAssetResultCode = -1
	ManageAssetResultCodeAlreadyExists ManageAssetResultCode = -2
	ManageAssetResultCodeMalformed     ManageAssetResultCode = -3
)

var ManageAssetResultCodeAll = []ManageAssetResultCode{
	ManageAssetResultCodeSuccess,
	ManageAssetResultCodeNotFound,
	ManageAssetResultCodeAlreadyExists,
	ManageAssetResultCodeMalformed,
}

var manageAssetResultCodeMap = map[int32]string{
	0:  "ManageAssetResultCodeSuccess",
	-1: "ManageAssetResultCodeNotFound",
	-2: "ManageAssetResultCodeAlreadyExists",
	-3: "ManageAssetResultCodeMalformed",
}

var manageAssetResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "not_found",
	-2: "already_exists",
	-3: "malformed",
}

var manageAssetResultCodeRevMap = map[string]int32{
	"ManageAssetResultCodeSuccess":       0,
	"ManageAssetResultCodeNotFound":      -1,
	"ManageAssetResultCodeAlreadyExists": -2,
	"ManageAssetResultCodeMalformed":     -3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetResultCode
func (e ManageAssetResultCode) ValidEnum(v int32) bool {
	_, ok := manageAssetResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetResultCode) String() string {
	name, _ := manageAssetResultCodeMap[int32(e)]
	return name
}

func (e ManageAssetResultCode) ShortString() string {
	name, _ := manageAssetResultCodeShortMap[int32(e)]
	return name
}

func (e ManageAssetResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageAssetResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageAssetResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageAssetResultCode(value)
	return nil
}

// ManageAssetSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetSuccessExt
func (u ManageAssetSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetSuccessExt creates a new  ManageAssetSuccessExt.
func NewManageAssetSuccessExt(v LedgerVersion, value interface{}) (result ManageAssetSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetSuccess is an XDR Struct defines as:
//
//   struct ManageAssetSuccess
//    {
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetSuccess struct {
	Ext ManageAssetSuccessExt `json:"ext,omitempty"`
}

// ManageAssetResult is an XDR Union defines as:
//
//   union ManageAssetResult switch (ManageAssetResultCode code)
//    {
//    case SUCCESS:
//        ManageAssetSuccess success;
//    default:
//        void;
//    };
//
type ManageAssetResult struct {
	Code    ManageAssetResultCode `json:"code,omitempty"`
	Success *ManageAssetSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetResult
func (u ManageAssetResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageAssetResultCode(sw) {
	case ManageAssetResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageAssetResult creates a new  ManageAssetResult.
func NewManageAssetResult(code ManageAssetResultCode, value interface{}) (result ManageAssetResult, err error) {
	result.Code = code
	switch ManageAssetResultCode(code) {
	case ManageAssetResultCodeSuccess:
		tv, ok := value.(ManageAssetSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageAssetResult) MustSuccess() ManageAssetSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageAssetResult) GetSuccess() (result ManageAssetSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageAssetPairAction is an XDR Enum defines as:
//
//   enum ManageAssetPairAction
//    {
//        CREATE = 0,
//        UPDATE_PRICE = 1,
//        UPDATE_POLICIES = 2
//    };
//
type ManageAssetPairAction int32

const (
	ManageAssetPairActionCreate         ManageAssetPairAction = 0
	ManageAssetPairActionUpdatePrice    ManageAssetPairAction = 1
	ManageAssetPairActionUpdatePolicies ManageAssetPairAction = 2
)

var ManageAssetPairActionAll = []ManageAssetPairAction{
	ManageAssetPairActionCreate,
	ManageAssetPairActionUpdatePrice,
	ManageAssetPairActionUpdatePolicies,
}

var manageAssetPairActionMap = map[int32]string{
	0: "ManageAssetPairActionCreate",
	1: "ManageAssetPairActionUpdatePrice",
	2: "ManageAssetPairActionUpdatePolicies",
}

var manageAssetPairActionShortMap = map[int32]string{
	0: "create",
	1: "update_price",
	2: "update_policies",
}

var manageAssetPairActionRevMap = map[string]int32{
	"ManageAssetPairActionCreate":         0,
	"ManageAssetPairActionUpdatePrice":    1,
	"ManageAssetPairActionUpdatePolicies": 2,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetPairAction
func (e ManageAssetPairAction) ValidEnum(v int32) bool {
	_, ok := manageAssetPairActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetPairAction) String() string {
	name, _ := manageAssetPairActionMap[int32(e)]
	return name
}

func (e ManageAssetPairAction) ShortString() string {
	name, _ := manageAssetPairActionShortMap[int32(e)]
	return name
}

func (e ManageAssetPairAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageAssetPairAction) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageAssetPairActionRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageAssetPairAction(value)
	return nil
}

// ManageAssetPairOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetPairOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetPairOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetPairOpExt
func (u ManageAssetPairOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetPairOpExt creates a new  ManageAssetPairOpExt.
func NewManageAssetPairOpExt(v LedgerVersion, value interface{}) (result ManageAssetPairOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetPairOp is an XDR Struct defines as:
//
//   struct ManageAssetPairOp
//    {
//        ManageAssetPairAction action;
//    	AssetCode base;
//    	AssetCode quote;
//
//        int64 physicalPrice;
//
//    	int64 physicalPriceCorrection; // correction of physical price in percents. If physical price is set and restriction by physical price set, mininal price for offer for this pair will be physicalPrice * physicalPriceCorrection
//    	int64 maxPriceStep;
//
//    	int32 policies;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetPairOp struct {
	Action                  ManageAssetPairAction `json:"action,omitempty"`
	Base                    AssetCode             `json:"base,omitempty"`
	Quote                   AssetCode             `json:"quote,omitempty"`
	PhysicalPrice           Int64                 `json:"physicalPrice,omitempty"`
	PhysicalPriceCorrection Int64                 `json:"physicalPriceCorrection,omitempty"`
	MaxPriceStep            Int64                 `json:"maxPriceStep,omitempty"`
	Policies                Int32                 `json:"policies,omitempty"`
	Ext                     ManageAssetPairOpExt  `json:"ext,omitempty"`
}

// ManageAssetPairResultCode is an XDR Enum defines as:
//
//   enum ManageAssetPairResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	NOT_FOUND = -1,           // failed to find asset with such code
//    	ALREADY_EXISTS = -2,
//        MALFORMED = -3,
//    	INVALID_ASSET = -4,
//    	INVALID_ACTION = -5,
//    	INVALID_POLICIES = -6,
//    	ASSET_NOT_FOUND = -7
//    };
//
type ManageAssetPairResultCode int32

const (
	ManageAssetPairResultCodeSuccess         ManageAssetPairResultCode = 0
	ManageAssetPairResultCodeNotFound        ManageAssetPairResultCode = -1
	ManageAssetPairResultCodeAlreadyExists   ManageAssetPairResultCode = -2
	ManageAssetPairResultCodeMalformed       ManageAssetPairResultCode = -3
	ManageAssetPairResultCodeInvalidAsset    ManageAssetPairResultCode = -4
	ManageAssetPairResultCodeInvalidAction   ManageAssetPairResultCode = -5
	ManageAssetPairResultCodeInvalidPolicies ManageAssetPairResultCode = -6
	ManageAssetPairResultCodeAssetNotFound   ManageAssetPairResultCode = -7
)

var ManageAssetPairResultCodeAll = []ManageAssetPairResultCode{
	ManageAssetPairResultCodeSuccess,
	ManageAssetPairResultCodeNotFound,
	ManageAssetPairResultCodeAlreadyExists,
	ManageAssetPairResultCodeMalformed,
	ManageAssetPairResultCodeInvalidAsset,
	ManageAssetPairResultCodeInvalidAction,
	ManageAssetPairResultCodeInvalidPolicies,
	ManageAssetPairResultCodeAssetNotFound,
}

var manageAssetPairResultCodeMap = map[int32]string{
	0:  "ManageAssetPairResultCodeSuccess",
	-1: "ManageAssetPairResultCodeNotFound",
	-2: "ManageAssetPairResultCodeAlreadyExists",
	-3: "ManageAssetPairResultCodeMalformed",
	-4: "ManageAssetPairResultCodeInvalidAsset",
	-5: "ManageAssetPairResultCodeInvalidAction",
	-6: "ManageAssetPairResultCodeInvalidPolicies",
	-7: "ManageAssetPairResultCodeAssetNotFound",
}

var manageAssetPairResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "not_found",
	-2: "already_exists",
	-3: "malformed",
	-4: "invalid_asset",
	-5: "invalid_action",
	-6: "invalid_policies",
	-7: "asset_not_found",
}

var manageAssetPairResultCodeRevMap = map[string]int32{
	"ManageAssetPairResultCodeSuccess":         0,
	"ManageAssetPairResultCodeNotFound":        -1,
	"ManageAssetPairResultCodeAlreadyExists":   -2,
	"ManageAssetPairResultCodeMalformed":       -3,
	"ManageAssetPairResultCodeInvalidAsset":    -4,
	"ManageAssetPairResultCodeInvalidAction":   -5,
	"ManageAssetPairResultCodeInvalidPolicies": -6,
	"ManageAssetPairResultCodeAssetNotFound":   -7,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetPairResultCode
func (e ManageAssetPairResultCode) ValidEnum(v int32) bool {
	_, ok := manageAssetPairResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetPairResultCode) String() string {
	name, _ := manageAssetPairResultCodeMap[int32(e)]
	return name
}

func (e ManageAssetPairResultCode) ShortString() string {
	name, _ := manageAssetPairResultCodeShortMap[int32(e)]
	return name
}

func (e ManageAssetPairResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageAssetPairResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageAssetPairResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageAssetPairResultCode(value)
	return nil
}

// ManageAssetPairSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetPairSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetPairSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetPairSuccessExt
func (u ManageAssetPairSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetPairSuccessExt creates a new  ManageAssetPairSuccessExt.
func NewManageAssetPairSuccessExt(v LedgerVersion, value interface{}) (result ManageAssetPairSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetPairSuccess is an XDR Struct defines as:
//
//   struct ManageAssetPairSuccess
//    {
//    	int64 currentPrice;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetPairSuccess struct {
	CurrentPrice Int64                     `json:"currentPrice,omitempty"`
	Ext          ManageAssetPairSuccessExt `json:"ext,omitempty"`
}

// ManageAssetPairResult is an XDR Union defines as:
//
//   union ManageAssetPairResult switch (ManageAssetPairResultCode code)
//    {
//    case SUCCESS:
//        ManageAssetPairSuccess success;
//    default:
//        void;
//    };
//
type ManageAssetPairResult struct {
	Code    ManageAssetPairResultCode `json:"code,omitempty"`
	Success *ManageAssetPairSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetPairResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetPairResult
func (u ManageAssetPairResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageAssetPairResultCode(sw) {
	case ManageAssetPairResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageAssetPairResult creates a new  ManageAssetPairResult.
func NewManageAssetPairResult(code ManageAssetPairResultCode, value interface{}) (result ManageAssetPairResult, err error) {
	result.Code = code
	switch ManageAssetPairResultCode(code) {
	case ManageAssetPairResultCodeSuccess:
		tv, ok := value.(ManageAssetPairSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetPairSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageAssetPairResult) MustSuccess() ManageAssetPairSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageAssetPairResult) GetSuccess() (result ManageAssetPairSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// AssetPolicy is an XDR Enum defines as:
//
//   enum AssetPolicy
//    {
//    	TRANSFERABLE = 1,
//        EMITTABLE_PRIMARY = 2,
//        EMITTABLE_SECONDARY = 4
//    };
//
type AssetPolicy int32

const (
	AssetPolicyTransferable       AssetPolicy = 1
	AssetPolicyEmittablePrimary   AssetPolicy = 2
	AssetPolicyEmittableSecondary AssetPolicy = 4
)

var AssetPolicyAll = []AssetPolicy{
	AssetPolicyTransferable,
	AssetPolicyEmittablePrimary,
	AssetPolicyEmittableSecondary,
}

var assetPolicyMap = map[int32]string{
	1: "AssetPolicyTransferable",
	2: "AssetPolicyEmittablePrimary",
	4: "AssetPolicyEmittableSecondary",
}

var assetPolicyShortMap = map[int32]string{
	1: "transferable",
	2: "emittable_primary",
	4: "emittable_secondary",
}

var assetPolicyRevMap = map[string]int32{
	"AssetPolicyTransferable":       1,
	"AssetPolicyEmittablePrimary":   2,
	"AssetPolicyEmittableSecondary": 4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AssetPolicy
func (e AssetPolicy) ValidEnum(v int32) bool {
	_, ok := assetPolicyMap[v]
	return ok
}

// String returns the name of `e`
func (e AssetPolicy) String() string {
	name, _ := assetPolicyMap[int32(e)]
	return name
}

func (e AssetPolicy) ShortString() string {
	name, _ := assetPolicyShortMap[int32(e)]
	return name
}

func (e AssetPolicy) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *AssetPolicy) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := assetPolicyRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = AssetPolicy(value)
	return nil
}

// AssetEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AssetEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetEntryExt
func (u AssetEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAssetEntryExt creates a new  AssetEntryExt.
func NewAssetEntryExt(v LedgerVersion, value interface{}) (result AssetEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AssetEntry is an XDR Struct defines as:
//
//   struct AssetEntry
//    {
//        AssetCode code;
//        int32 policies;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AssetEntry struct {
	Code     AssetCode     `json:"code,omitempty"`
	Policies Int32         `json:"policies,omitempty"`
	Ext      AssetEntryExt `json:"ext,omitempty"`
}

// SetFeesOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type SetFeesOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetFeesOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetFeesOpExt
func (u SetFeesOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetFeesOpExt creates a new  SetFeesOpExt.
func NewSetFeesOpExt(v LedgerVersion, value interface{}) (result SetFeesOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetFeesOp is an XDR Struct defines as:
//
//   struct SetFeesOp
//        {
//            FeeEntry* fee;
//    		bool isDelete;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        };
//
type SetFeesOp struct {
	Fee      *FeeEntry    `json:"fee,omitempty"`
	IsDelete bool         `json:"isDelete,omitempty"`
	Ext      SetFeesOpExt `json:"ext,omitempty"`
}

// SetFeesResultCode is an XDR Enum defines as:
//
//   enum SetFeesResultCode
//        {
//            // codes considered as "success" for the operation
//            SUCCESS = 0,
//
//            // codes considered as "failure" for the operation
//            INVALID_AMOUNT = -1,      // amount is negative
//    		INVALID_FEE_TYPE = -2,     // operation type is invalid
//            ASSET_NOT_FOUND = -3,
//            INVALID_ASSET = -4,
//            MALFORMED = -5,
//    		MALFORMED_RANGE = -6,
//    		RANGE_OVERLAP = -7,
//    		NOT_FOUND = -8,
//    		SUB_TYPE_NOT_EXIST = -9
//        };
//
type SetFeesResultCode int32

const (
	SetFeesResultCodeSuccess         SetFeesResultCode = 0
	SetFeesResultCodeInvalidAmount   SetFeesResultCode = -1
	SetFeesResultCodeInvalidFeeType  SetFeesResultCode = -2
	SetFeesResultCodeAssetNotFound   SetFeesResultCode = -3
	SetFeesResultCodeInvalidAsset    SetFeesResultCode = -4
	SetFeesResultCodeMalformed       SetFeesResultCode = -5
	SetFeesResultCodeMalformedRange  SetFeesResultCode = -6
	SetFeesResultCodeRangeOverlap    SetFeesResultCode = -7
	SetFeesResultCodeNotFound        SetFeesResultCode = -8
	SetFeesResultCodeSubTypeNotExist SetFeesResultCode = -9
)

var SetFeesResultCodeAll = []SetFeesResultCode{
	SetFeesResultCodeSuccess,
	SetFeesResultCodeInvalidAmount,
	SetFeesResultCodeInvalidFeeType,
	SetFeesResultCodeAssetNotFound,
	SetFeesResultCodeInvalidAsset,
	SetFeesResultCodeMalformed,
	SetFeesResultCodeMalformedRange,
	SetFeesResultCodeRangeOverlap,
	SetFeesResultCodeNotFound,
	SetFeesResultCodeSubTypeNotExist,
}

var setFeesResultCodeMap = map[int32]string{
	0:  "SetFeesResultCodeSuccess",
	-1: "SetFeesResultCodeInvalidAmount",
	-2: "SetFeesResultCodeInvalidFeeType",
	-3: "SetFeesResultCodeAssetNotFound",
	-4: "SetFeesResultCodeInvalidAsset",
	-5: "SetFeesResultCodeMalformed",
	-6: "SetFeesResultCodeMalformedRange",
	-7: "SetFeesResultCodeRangeOverlap",
	-8: "SetFeesResultCodeNotFound",
	-9: "SetFeesResultCodeSubTypeNotExist",
}

var setFeesResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "invalid_amount",
	-2: "invalid_fee_type",
	-3: "asset_not_found",
	-4: "invalid_asset",
	-5: "malformed",
	-6: "malformed_range",
	-7: "range_overlap",
	-8: "not_found",
	-9: "sub_type_not_exist",
}

var setFeesResultCodeRevMap = map[string]int32{
	"SetFeesResultCodeSuccess":         0,
	"SetFeesResultCodeInvalidAmount":   -1,
	"SetFeesResultCodeInvalidFeeType":  -2,
	"SetFeesResultCodeAssetNotFound":   -3,
	"SetFeesResultCodeInvalidAsset":    -4,
	"SetFeesResultCodeMalformed":       -5,
	"SetFeesResultCodeMalformedRange":  -6,
	"SetFeesResultCodeRangeOverlap":    -7,
	"SetFeesResultCodeNotFound":        -8,
	"SetFeesResultCodeSubTypeNotExist": -9,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetFeesResultCode
func (e SetFeesResultCode) ValidEnum(v int32) bool {
	_, ok := setFeesResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetFeesResultCode) String() string {
	name, _ := setFeesResultCodeMap[int32(e)]
	return name
}

func (e SetFeesResultCode) ShortString() string {
	name, _ := setFeesResultCodeShortMap[int32(e)]
	return name
}

func (e SetFeesResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *SetFeesResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := setFeesResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = SetFeesResultCode(value)
	return nil
}

// SetFeesResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    				{
//    				case EMPTY_VERSION:
//    					void;
//    				}
//
type SetFeesResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetFeesResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetFeesResultSuccessExt
func (u SetFeesResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetFeesResultSuccessExt creates a new  SetFeesResultSuccessExt.
func NewSetFeesResultSuccessExt(v LedgerVersion, value interface{}) (result SetFeesResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetFeesResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    				// reserved for future use
//    				union switch (LedgerVersion v)
//    				{
//    				case EMPTY_VERSION:
//    					void;
//    				}
//    				ext;
//    			}
//
type SetFeesResultSuccess struct {
	Ext SetFeesResultSuccessExt `json:"ext,omitempty"`
}

// SetFeesResult is an XDR Union defines as:
//
//   union SetFeesResult switch (SetFeesResultCode code)
//        {
//            case SUCCESS:
//                struct {
//    				// reserved for future use
//    				union switch (LedgerVersion v)
//    				{
//    				case EMPTY_VERSION:
//    					void;
//    				}
//    				ext;
//    			} success;
//            default:
//                void;
//        };
//
type SetFeesResult struct {
	Code    SetFeesResultCode     `json:"code,omitempty"`
	Success *SetFeesResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetFeesResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetFeesResult
func (u SetFeesResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetFeesResultCode(sw) {
	case SetFeesResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewSetFeesResult creates a new  SetFeesResult.
func NewSetFeesResult(code SetFeesResultCode, value interface{}) (result SetFeesResult, err error) {
	result.Code = code
	switch SetFeesResultCode(code) {
	case SetFeesResultCodeSuccess:
		tv, ok := value.(SetFeesResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetFeesResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u SetFeesResult) MustSuccess() SetFeesResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SetFeesResult) GetSuccess() (result SetFeesResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageBalanceAction is an XDR Enum defines as:
//
//   enum ManageBalanceAction
//    {
//        CREATE = 0,
//        DELETE = 1
//    };
//
type ManageBalanceAction int32

const (
	ManageBalanceActionCreate ManageBalanceAction = 0
	ManageBalanceActionDelete ManageBalanceAction = 1
)

var ManageBalanceActionAll = []ManageBalanceAction{
	ManageBalanceActionCreate,
	ManageBalanceActionDelete,
}

var manageBalanceActionMap = map[int32]string{
	0: "ManageBalanceActionCreate",
	1: "ManageBalanceActionDelete",
}

var manageBalanceActionShortMap = map[int32]string{
	0: "create",
	1: "delete",
}

var manageBalanceActionRevMap = map[string]int32{
	"ManageBalanceActionCreate": 0,
	"ManageBalanceActionDelete": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageBalanceAction
func (e ManageBalanceAction) ValidEnum(v int32) bool {
	_, ok := manageBalanceActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageBalanceAction) String() string {
	name, _ := manageBalanceActionMap[int32(e)]
	return name
}

func (e ManageBalanceAction) ShortString() string {
	name, _ := manageBalanceActionShortMap[int32(e)]
	return name
}

func (e ManageBalanceAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageBalanceAction) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageBalanceActionRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageBalanceAction(value)
	return nil
}

// ManageBalanceOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageBalanceOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBalanceOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBalanceOpExt
func (u ManageBalanceOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageBalanceOpExt creates a new  ManageBalanceOpExt.
func NewManageBalanceOpExt(v LedgerVersion, value interface{}) (result ManageBalanceOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageBalanceOp is an XDR Struct defines as:
//
//   struct ManageBalanceOp
//    {
//        BalanceID balanceID;
//        ManageBalanceAction action;
//        AccountID destination;
//        AssetCode asset;
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageBalanceOp struct {
	BalanceId   BalanceId           `json:"balanceID,omitempty"`
	Action      ManageBalanceAction `json:"action,omitempty"`
	Destination AccountId           `json:"destination,omitempty"`
	Asset       AssetCode           `json:"asset,omitempty"`
	Ext         ManageBalanceOpExt  `json:"ext,omitempty"`
}

// ManageBalanceResultCode is an XDR Enum defines as:
//
//   enum ManageBalanceResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,       // invalid destination
//        NOT_FOUND = -2,
//        DESTINATION_NOT_FOUND = -3,
//        ALREADY_EXISTS = -4,
//        ASSET_NOT_FOUND = -5,
//        INVALID_ASSET = -6
//    };
//
type ManageBalanceResultCode int32

const (
	ManageBalanceResultCodeSuccess             ManageBalanceResultCode = 0
	ManageBalanceResultCodeMalformed           ManageBalanceResultCode = -1
	ManageBalanceResultCodeNotFound            ManageBalanceResultCode = -2
	ManageBalanceResultCodeDestinationNotFound ManageBalanceResultCode = -3
	ManageBalanceResultCodeAlreadyExists       ManageBalanceResultCode = -4
	ManageBalanceResultCodeAssetNotFound       ManageBalanceResultCode = -5
	ManageBalanceResultCodeInvalidAsset        ManageBalanceResultCode = -6
)

var ManageBalanceResultCodeAll = []ManageBalanceResultCode{
	ManageBalanceResultCodeSuccess,
	ManageBalanceResultCodeMalformed,
	ManageBalanceResultCodeNotFound,
	ManageBalanceResultCodeDestinationNotFound,
	ManageBalanceResultCodeAlreadyExists,
	ManageBalanceResultCodeAssetNotFound,
	ManageBalanceResultCodeInvalidAsset,
}

var manageBalanceResultCodeMap = map[int32]string{
	0:  "ManageBalanceResultCodeSuccess",
	-1: "ManageBalanceResultCodeMalformed",
	-2: "ManageBalanceResultCodeNotFound",
	-3: "ManageBalanceResultCodeDestinationNotFound",
	-4: "ManageBalanceResultCodeAlreadyExists",
	-5: "ManageBalanceResultCodeAssetNotFound",
	-6: "ManageBalanceResultCodeInvalidAsset",
}

var manageBalanceResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "malformed",
	-2: "not_found",
	-3: "destination_not_found",
	-4: "already_exists",
	-5: "asset_not_found",
	-6: "invalid_asset",
}

var manageBalanceResultCodeRevMap = map[string]int32{
	"ManageBalanceResultCodeSuccess":             0,
	"ManageBalanceResultCodeMalformed":           -1,
	"ManageBalanceResultCodeNotFound":            -2,
	"ManageBalanceResultCodeDestinationNotFound": -3,
	"ManageBalanceResultCodeAlreadyExists":       -4,
	"ManageBalanceResultCodeAssetNotFound":       -5,
	"ManageBalanceResultCodeInvalidAsset":        -6,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageBalanceResultCode
func (e ManageBalanceResultCode) ValidEnum(v int32) bool {
	_, ok := manageBalanceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageBalanceResultCode) String() string {
	name, _ := manageBalanceResultCodeMap[int32(e)]
	return name
}

func (e ManageBalanceResultCode) ShortString() string {
	name, _ := manageBalanceResultCodeShortMap[int32(e)]
	return name
}

func (e ManageBalanceResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageBalanceResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageBalanceResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageBalanceResultCode(value)
	return nil
}

// ManageBalanceSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageBalanceSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBalanceSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBalanceSuccessExt
func (u ManageBalanceSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageBalanceSuccessExt creates a new  ManageBalanceSuccessExt.
func NewManageBalanceSuccessExt(v LedgerVersion, value interface{}) (result ManageBalanceSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageBalanceSuccess is an XDR Struct defines as:
//
//   struct ManageBalanceSuccess {
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageBalanceSuccess struct {
	Ext ManageBalanceSuccessExt `json:"ext,omitempty"`
}

// ManageBalanceResult is an XDR Union defines as:
//
//   union ManageBalanceResult switch (ManageBalanceResultCode code)
//    {
//    case SUCCESS:
//        ManageBalanceSuccess success;
//    default:
//        void;
//    };
//
type ManageBalanceResult struct {
	Code    ManageBalanceResultCode `json:"code,omitempty"`
	Success *ManageBalanceSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBalanceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBalanceResult
func (u ManageBalanceResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageBalanceResultCode(sw) {
	case ManageBalanceResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageBalanceResult creates a new  ManageBalanceResult.
func NewManageBalanceResult(code ManageBalanceResultCode, value interface{}) (result ManageBalanceResult, err error) {
	result.Code = code
	switch ManageBalanceResultCode(code) {
	case ManageBalanceResultCodeSuccess:
		tv, ok := value.(ManageBalanceSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBalanceSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageBalanceResult) MustSuccess() ManageBalanceSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageBalanceResult) GetSuccess() (result ManageBalanceSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ReviewPaymentRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type ReviewPaymentRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewPaymentRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewPaymentRequestOpExt
func (u ReviewPaymentRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewPaymentRequestOpExt creates a new  ReviewPaymentRequestOpExt.
func NewReviewPaymentRequestOpExt(v LedgerVersion, value interface{}) (result ReviewPaymentRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewPaymentRequestOp is an XDR Struct defines as:
//
//   struct ReviewPaymentRequestOp
//    {
//        uint64 paymentID;
//
//    	bool accept;
//        string256* rejectReason;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type ReviewPaymentRequestOp struct {
	PaymentId    Uint64                    `json:"paymentID,omitempty"`
	Accept       bool                      `json:"accept,omitempty"`
	RejectReason *String256                `json:"rejectReason,omitempty"`
	Ext          ReviewPaymentRequestOpExt `json:"ext,omitempty"`
}

// ReviewPaymentRequestResultCode is an XDR Enum defines as:
//
//   enum ReviewPaymentRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	NOT_FOUND = -1,           // failed to find Recovery request with such ID
//        LINE_FULL = -2
//    };
//
type ReviewPaymentRequestResultCode int32

const (
	ReviewPaymentRequestResultCodeSuccess  ReviewPaymentRequestResultCode = 0
	ReviewPaymentRequestResultCodeNotFound ReviewPaymentRequestResultCode = -1
	ReviewPaymentRequestResultCodeLineFull ReviewPaymentRequestResultCode = -2
)

var ReviewPaymentRequestResultCodeAll = []ReviewPaymentRequestResultCode{
	ReviewPaymentRequestResultCodeSuccess,
	ReviewPaymentRequestResultCodeNotFound,
	ReviewPaymentRequestResultCodeLineFull,
}

var reviewPaymentRequestResultCodeMap = map[int32]string{
	0:  "ReviewPaymentRequestResultCodeSuccess",
	-1: "ReviewPaymentRequestResultCodeNotFound",
	-2: "ReviewPaymentRequestResultCodeLineFull",
}

var reviewPaymentRequestResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "not_found",
	-2: "line_full",
}

var reviewPaymentRequestResultCodeRevMap = map[string]int32{
	"ReviewPaymentRequestResultCodeSuccess":  0,
	"ReviewPaymentRequestResultCodeNotFound": -1,
	"ReviewPaymentRequestResultCodeLineFull": -2,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ReviewPaymentRequestResultCode
func (e ReviewPaymentRequestResultCode) ValidEnum(v int32) bool {
	_, ok := reviewPaymentRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ReviewPaymentRequestResultCode) String() string {
	name, _ := reviewPaymentRequestResultCodeMap[int32(e)]
	return name
}

func (e ReviewPaymentRequestResultCode) ShortString() string {
	name, _ := reviewPaymentRequestResultCodeShortMap[int32(e)]
	return name
}

func (e ReviewPaymentRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ReviewPaymentRequestResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := reviewPaymentRequestResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ReviewPaymentRequestResultCode(value)
	return nil
}

// PaymentState is an XDR Enum defines as:
//
//   enum PaymentState
//    {
//        PENDING = 0,
//        PROCESSED = 1,
//        REJECTED = 2
//    };
//
type PaymentState int32

const (
	PaymentStatePending   PaymentState = 0
	PaymentStateProcessed PaymentState = 1
	PaymentStateRejected  PaymentState = 2
)

var PaymentStateAll = []PaymentState{
	PaymentStatePending,
	PaymentStateProcessed,
	PaymentStateRejected,
}

var paymentStateMap = map[int32]string{
	0: "PaymentStatePending",
	1: "PaymentStateProcessed",
	2: "PaymentStateRejected",
}

var paymentStateShortMap = map[int32]string{
	0: "pending",
	1: "processed",
	2: "rejected",
}

var paymentStateRevMap = map[string]int32{
	"PaymentStatePending":   0,
	"PaymentStateProcessed": 1,
	"PaymentStateRejected":  2,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PaymentState
func (e PaymentState) ValidEnum(v int32) bool {
	_, ok := paymentStateMap[v]
	return ok
}

// String returns the name of `e`
func (e PaymentState) String() string {
	name, _ := paymentStateMap[int32(e)]
	return name
}

func (e PaymentState) ShortString() string {
	name, _ := paymentStateShortMap[int32(e)]
	return name
}

func (e PaymentState) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *PaymentState) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := paymentStateRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = PaymentState(value)
	return nil
}

// ReviewPaymentResponseExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type ReviewPaymentResponseExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewPaymentResponseExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewPaymentResponseExt
func (u ReviewPaymentResponseExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewPaymentResponseExt creates a new  ReviewPaymentResponseExt.
func NewReviewPaymentResponseExt(v LedgerVersion, value interface{}) (result ReviewPaymentResponseExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewPaymentResponse is an XDR Struct defines as:
//
//   struct ReviewPaymentResponse {
//        PaymentState state;
//
//        uint64* relatedInvoiceID;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type ReviewPaymentResponse struct {
	State            PaymentState             `json:"state,omitempty"`
	RelatedInvoiceId *Uint64                  `json:"relatedInvoiceID,omitempty"`
	Ext              ReviewPaymentResponseExt `json:"ext,omitempty"`
}

// ReviewPaymentRequestResult is an XDR Union defines as:
//
//   union ReviewPaymentRequestResult switch (ReviewPaymentRequestResultCode code)
//    {
//    case SUCCESS:
//        ReviewPaymentResponse reviewPaymentResponse;
//    default:
//        void;
//    };
//
type ReviewPaymentRequestResult struct {
	Code                  ReviewPaymentRequestResultCode `json:"code,omitempty"`
	ReviewPaymentResponse *ReviewPaymentResponse         `json:"reviewPaymentResponse,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewPaymentRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewPaymentRequestResult
func (u ReviewPaymentRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ReviewPaymentRequestResultCode(sw) {
	case ReviewPaymentRequestResultCodeSuccess:
		return "ReviewPaymentResponse", true
	default:
		return "", true
	}
}

// NewReviewPaymentRequestResult creates a new  ReviewPaymentRequestResult.
func NewReviewPaymentRequestResult(code ReviewPaymentRequestResultCode, value interface{}) (result ReviewPaymentRequestResult, err error) {
	result.Code = code
	switch ReviewPaymentRequestResultCode(code) {
	case ReviewPaymentRequestResultCodeSuccess:
		tv, ok := value.(ReviewPaymentResponse)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewPaymentResponse")
			return
		}
		result.ReviewPaymentResponse = &tv
	default:
		// void
	}
	return
}

// MustReviewPaymentResponse retrieves the ReviewPaymentResponse value from the union,
// panicing if the value is not set.
func (u ReviewPaymentRequestResult) MustReviewPaymentResponse() ReviewPaymentResponse {
	val, ok := u.GetReviewPaymentResponse()

	if !ok {
		panic("arm ReviewPaymentResponse is not set")
	}

	return val
}

// GetReviewPaymentResponse retrieves the ReviewPaymentResponse value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ReviewPaymentRequestResult) GetReviewPaymentResponse() (result ReviewPaymentResponse, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "ReviewPaymentResponse" {
		result = *u.ReviewPaymentResponse
		ok = true
	}

	return
}

// Hash is an XDR Typedef defines as:
//
//   typedef opaque Hash[32];
//
type Hash [32]byte

// Uint256 is an XDR Typedef defines as:
//
//   typedef opaque uint256[32];
//
type Uint256 [32]byte

// Uint32 is an XDR Typedef defines as:
//
//   typedef unsigned int uint32;
//
type Uint32 uint32

// Int32 is an XDR Typedef defines as:
//
//   typedef int int32;
//
type Int32 int32

// Uint64 is an XDR Typedef defines as:
//
//   typedef unsigned hyper uint64;
//
type Uint64 uint64

// Int64 is an XDR Typedef defines as:
//
//   typedef hyper int64;
//
type Int64 int64

// CryptoKeyType is an XDR Enum defines as:
//
//   enum CryptoKeyType
//    {
//        KEY_TYPE_ED25519 = 0
//    };
//
type CryptoKeyType int32

const (
	CryptoKeyTypeKeyTypeEd25519 CryptoKeyType = 0
)

var CryptoKeyTypeAll = []CryptoKeyType{
	CryptoKeyTypeKeyTypeEd25519,
}

var cryptoKeyTypeMap = map[int32]string{
	0: "CryptoKeyTypeKeyTypeEd25519",
}

var cryptoKeyTypeShortMap = map[int32]string{
	0: "key_type_ed25519",
}

var cryptoKeyTypeRevMap = map[string]int32{
	"CryptoKeyTypeKeyTypeEd25519": 0,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CryptoKeyType
func (e CryptoKeyType) ValidEnum(v int32) bool {
	_, ok := cryptoKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e CryptoKeyType) String() string {
	name, _ := cryptoKeyTypeMap[int32(e)]
	return name
}

func (e CryptoKeyType) ShortString() string {
	name, _ := cryptoKeyTypeShortMap[int32(e)]
	return name
}

func (e CryptoKeyType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *CryptoKeyType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := cryptoKeyTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = CryptoKeyType(value)
	return nil
}

// PublicKeyType is an XDR Enum defines as:
//
//   enum PublicKeyType
//    {
//    	PUBLIC_KEY_TYPE_ED25519 = 0
//    };
//
type PublicKeyType int32

const (
	PublicKeyTypePublicKeyTypeEd25519 PublicKeyType = 0
)

var PublicKeyTypeAll = []PublicKeyType{
	PublicKeyTypePublicKeyTypeEd25519,
}

var publicKeyTypeMap = map[int32]string{
	0: "PublicKeyTypePublicKeyTypeEd25519",
}

var publicKeyTypeShortMap = map[int32]string{
	0: "public_key_type_ed25519",
}

var publicKeyTypeRevMap = map[string]int32{
	"PublicKeyTypePublicKeyTypeEd25519": 0,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PublicKeyType
func (e PublicKeyType) ValidEnum(v int32) bool {
	_, ok := publicKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e PublicKeyType) String() string {
	name, _ := publicKeyTypeMap[int32(e)]
	return name
}

func (e PublicKeyType) ShortString() string {
	name, _ := publicKeyTypeShortMap[int32(e)]
	return name
}

func (e PublicKeyType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *PublicKeyType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := publicKeyTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = PublicKeyType(value)
	return nil
}

// PublicKey is an XDR Union defines as:
//
//   union PublicKey switch (CryptoKeyType type)
//    {
//    case KEY_TYPE_ED25519:
//        uint256 ed25519;
//    };
//
type PublicKey struct {
	Type    CryptoKeyType `json:"type,omitempty"`
	Ed25519 *Uint256      `json:"ed25519,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PublicKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u PublicKey) ArmForSwitch(sw int32) (string, bool) {
	switch CryptoKeyType(sw) {
	case CryptoKeyTypeKeyTypeEd25519:
		return "Ed25519", true
	}
	return "-", false
}

// NewPublicKey creates a new  PublicKey.
func NewPublicKey(aType CryptoKeyType, value interface{}) (result PublicKey, err error) {
	result.Type = aType
	switch CryptoKeyType(aType) {
	case CryptoKeyTypeKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u PublicKey) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PublicKey) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// LedgerVersion is an XDR Enum defines as:
//
//   enum LedgerVersion {
//    	EMPTY_VERSION = 0
//    };
//
type LedgerVersion int32

const (
	LedgerVersionEmptyVersion LedgerVersion = 0
)

var LedgerVersionAll = []LedgerVersion{
	LedgerVersionEmptyVersion,
}

var ledgerVersionMap = map[int32]string{
	0: "LedgerVersionEmptyVersion",
}

var ledgerVersionShortMap = map[int32]string{
	0: "empty_version",
}

var ledgerVersionRevMap = map[string]int32{
	"LedgerVersionEmptyVersion": 0,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerVersion
func (e LedgerVersion) ValidEnum(v int32) bool {
	_, ok := ledgerVersionMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerVersion) String() string {
	name, _ := ledgerVersionMap[int32(e)]
	return name
}

func (e LedgerVersion) ShortString() string {
	name, _ := ledgerVersionShortMap[int32(e)]
	return name
}

func (e LedgerVersion) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *LedgerVersion) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := ledgerVersionRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = LedgerVersion(value)
	return nil
}

// Signature is an XDR Typedef defines as:
//
//   typedef opaque Signature<64>;
//
type Signature []byte

// SignatureHint is an XDR Typedef defines as:
//
//   typedef opaque SignatureHint[4];
//
type SignatureHint [4]byte

// NodeId is an XDR Typedef defines as:
//
//   typedef PublicKey NodeID;
//
type NodeId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u NodeId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u NodeId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewNodeId creates a new  NodeId.
func NewNodeId(aType CryptoKeyType, value interface{}) (result NodeId, err error) {
	u, err := NewPublicKey(aType, value)
	result = NodeId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u NodeId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u NodeId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// Curve25519Secret is an XDR Struct defines as:
//
//   struct Curve25519Secret
//    {
//            opaque key[32];
//    };
//
type Curve25519Secret struct {
	Key [32]byte `json:"key,omitempty"`
}

// Curve25519Public is an XDR Struct defines as:
//
//   struct Curve25519Public
//    {
//            opaque key[32];
//    };
//
type Curve25519Public struct {
	Key [32]byte `json:"key,omitempty"`
}

// HmacSha256Key is an XDR Struct defines as:
//
//   struct HmacSha256Key
//    {
//            opaque key[32];
//    };
//
type HmacSha256Key struct {
	Key [32]byte `json:"key,omitempty"`
}

// HmacSha256Mac is an XDR Struct defines as:
//
//   struct HmacSha256Mac
//    {
//            opaque mac[32];
//    };
//
type HmacSha256Mac struct {
	Mac [32]byte `json:"mac,omitempty"`
}

// AccountId is an XDR Typedef defines as:
//
//   typedef PublicKey AccountID;
//
type AccountId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u AccountId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewAccountId creates a new  AccountId.
func NewAccountId(aType CryptoKeyType, value interface{}) (result AccountId, err error) {
	u, err := NewPublicKey(aType, value)
	result = AccountId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u AccountId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// BalanceId is an XDR Typedef defines as:
//
//   typedef PublicKey BalanceID;
//
type BalanceId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BalanceId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u BalanceId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewBalanceId creates a new  BalanceId.
func NewBalanceId(aType CryptoKeyType, value interface{}) (result BalanceId, err error) {
	u, err := NewPublicKey(aType, value)
	result = BalanceId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u BalanceId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BalanceId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// Thresholds is an XDR Typedef defines as:
//
//   typedef opaque Thresholds[4];
//
type Thresholds [4]byte

// String32 is an XDR Typedef defines as:
//
//   typedef string string32<32>;
//
type String32 string

// XDRMaxSize implements the Sized interface for String32
func (e String32) XDRMaxSize() int {
	return 32
}

// String64 is an XDR Typedef defines as:
//
//   typedef string string64<64>;
//
type String64 string

// XDRMaxSize implements the Sized interface for String64
func (e String64) XDRMaxSize() int {
	return 64
}

// String256 is an XDR Typedef defines as:
//
//   typedef string string256<256>;
//
type String256 string

// XDRMaxSize implements the Sized interface for String256
func (e String256) XDRMaxSize() int {
	return 256
}

// Longstring is an XDR Typedef defines as:
//
//   typedef string longstring<>;
//
type Longstring string

// AssetCode is an XDR Typedef defines as:
//
//   typedef string AssetCode<16>;
//
type AssetCode string

// XDRMaxSize implements the Sized interface for AssetCode
func (e AssetCode) XDRMaxSize() int {
	return 16
}

// Salt is an XDR Typedef defines as:
//
//   typedef uint64 Salt;
//
type Salt Uint64

// DataValue is an XDR Typedef defines as:
//
//   typedef opaque DataValue<64>;
//
type DataValue []byte

// OperationType is an XDR Enum defines as:
//
//   enum OperationType
//    {
//        CREATE_ACCOUNT = 0,
//        PAYMENT = 1,
//        SET_OPTIONS = 2,
//        MANAGE_COINS_EMISSION_REQUEST = 3,
//        REVIEW_COINS_EMISSION_REQUEST = 4,
//        SET_FEES = 5,
//    	MANAGE_ACCOUNT = 6,
//        MANAGE_FORFEIT_REQUEST = 7,
//        RECOVER = 8,
//        MANAGE_BALANCE = 9,
//        REVIEW_PAYMENT_REQUEST = 10,
//        MANAGE_ASSET = 11,
//        UPLOAD_PREEMISSIONS = 12,
//        SET_LIMITS = 13,
//        DIRECT_DEBIT = 14,
//    	MANAGE_ASSET_PAIR = 15,
//    	MANAGE_OFFER = 16,
//        MANAGE_INVOICE = 17
//    };
//
type OperationType int32

const (
	OperationTypeCreateAccount              OperationType = 0
	OperationTypePayment                    OperationType = 1
	OperationTypeSetOptions                 OperationType = 2
	OperationTypeManageCoinsEmissionRequest OperationType = 3
	OperationTypeReviewCoinsEmissionRequest OperationType = 4
	OperationTypeSetFees                    OperationType = 5
	OperationTypeManageAccount              OperationType = 6
	OperationTypeManageForfeitRequest       OperationType = 7
	OperationTypeRecover                    OperationType = 8
	OperationTypeManageBalance              OperationType = 9
	OperationTypeReviewPaymentRequest       OperationType = 10
	OperationTypeManageAsset                OperationType = 11
	OperationTypeUploadPreemissions         OperationType = 12
	OperationTypeSetLimits                  OperationType = 13
	OperationTypeDirectDebit                OperationType = 14
	OperationTypeManageAssetPair            OperationType = 15
	OperationTypeManageOffer                OperationType = 16
	OperationTypeManageInvoice              OperationType = 17
)

var OperationTypeAll = []OperationType{
	OperationTypeCreateAccount,
	OperationTypePayment,
	OperationTypeSetOptions,
	OperationTypeManageCoinsEmissionRequest,
	OperationTypeReviewCoinsEmissionRequest,
	OperationTypeSetFees,
	OperationTypeManageAccount,
	OperationTypeManageForfeitRequest,
	OperationTypeRecover,
	OperationTypeManageBalance,
	OperationTypeReviewPaymentRequest,
	OperationTypeManageAsset,
	OperationTypeUploadPreemissions,
	OperationTypeSetLimits,
	OperationTypeDirectDebit,
	OperationTypeManageAssetPair,
	OperationTypeManageOffer,
	OperationTypeManageInvoice,
}

var operationTypeMap = map[int32]string{
	0:  "OperationTypeCreateAccount",
	1:  "OperationTypePayment",
	2:  "OperationTypeSetOptions",
	3:  "OperationTypeManageCoinsEmissionRequest",
	4:  "OperationTypeReviewCoinsEmissionRequest",
	5:  "OperationTypeSetFees",
	6:  "OperationTypeManageAccount",
	7:  "OperationTypeManageForfeitRequest",
	8:  "OperationTypeRecover",
	9:  "OperationTypeManageBalance",
	10: "OperationTypeReviewPaymentRequest",
	11: "OperationTypeManageAsset",
	12: "OperationTypeUploadPreemissions",
	13: "OperationTypeSetLimits",
	14: "OperationTypeDirectDebit",
	15: "OperationTypeManageAssetPair",
	16: "OperationTypeManageOffer",
	17: "OperationTypeManageInvoice",
}

var operationTypeShortMap = map[int32]string{
	0:  "create_account",
	1:  "payment",
	2:  "set_options",
	3:  "manage_coins_emission_request",
	4:  "review_coins_emission_request",
	5:  "set_fees",
	6:  "manage_account",
	7:  "manage_forfeit_request",
	8:  "recover",
	9:  "manage_balance",
	10: "review_payment_request",
	11: "manage_asset",
	12: "upload_preemissions",
	13: "set_limits",
	14: "direct_debit",
	15: "manage_asset_pair",
	16: "manage_offer",
	17: "manage_invoice",
}

var operationTypeRevMap = map[string]int32{
	"OperationTypeCreateAccount":              0,
	"OperationTypePayment":                    1,
	"OperationTypeSetOptions":                 2,
	"OperationTypeManageCoinsEmissionRequest": 3,
	"OperationTypeReviewCoinsEmissionRequest": 4,
	"OperationTypeSetFees":                    5,
	"OperationTypeManageAccount":              6,
	"OperationTypeManageForfeitRequest":       7,
	"OperationTypeRecover":                    8,
	"OperationTypeManageBalance":              9,
	"OperationTypeReviewPaymentRequest":       10,
	"OperationTypeManageAsset":                11,
	"OperationTypeUploadPreemissions":         12,
	"OperationTypeSetLimits":                  13,
	"OperationTypeDirectDebit":                14,
	"OperationTypeManageAssetPair":            15,
	"OperationTypeManageOffer":                16,
	"OperationTypeManageInvoice":              17,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationType
func (e OperationType) ValidEnum(v int32) bool {
	_, ok := operationTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationType) String() string {
	name, _ := operationTypeMap[int32(e)]
	return name
}

func (e OperationType) ShortString() string {
	name, _ := operationTypeShortMap[int32(e)]
	return name
}

func (e OperationType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *OperationType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := operationTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = OperationType(value)
	return nil
}

// DecoratedSignature is an XDR Struct defines as:
//
//   struct DecoratedSignature
//    {
//        SignatureHint hint;  // last 4 bytes of the public key, used as a hint
//        Signature signature; // actual signature
//    };
//
type DecoratedSignature struct {
	Hint      SignatureHint `json:"hint,omitempty"`
	Signature Signature     `json:"signature,omitempty"`
}

// ManageInvoiceOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageInvoiceOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageInvoiceOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageInvoiceOpExt
func (u ManageInvoiceOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageInvoiceOpExt creates a new  ManageInvoiceOpExt.
func NewManageInvoiceOpExt(v LedgerVersion, value interface{}) (result ManageInvoiceOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageInvoiceOp is an XDR Struct defines as:
//
//   struct ManageInvoiceOp
//    {
//        BalanceID receiverBalance;
//    	AccountID sender;
//        int64 amount; // if set to 0, delete the invoice
//
//        // 0=create a new invoice, otherwise edit an existing invoice
//        uint64 invoiceID;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageInvoiceOp struct {
	ReceiverBalance BalanceId          `json:"receiverBalance,omitempty"`
	Sender          AccountId          `json:"sender,omitempty"`
	Amount          Int64              `json:"amount,omitempty"`
	InvoiceId       Uint64             `json:"invoiceID,omitempty"`
	Ext             ManageInvoiceOpExt `json:"ext,omitempty"`
}

// ManageInvoiceResultCode is an XDR Enum defines as:
//
//   enum ManageInvoiceResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,
//        BALANCE_NOT_FOUND = -2,
//    	INVOICE_OVERFLOW = -3,
//
//        NOT_FOUND = -4,
//        TOO_MANY_INVOICES = -5,
//        CAN_NOT_DELETE_IN_PROGRESS = -6
//    };
//
type ManageInvoiceResultCode int32

const (
	ManageInvoiceResultCodeSuccess                ManageInvoiceResultCode = 0
	ManageInvoiceResultCodeMalformed              ManageInvoiceResultCode = -1
	ManageInvoiceResultCodeBalanceNotFound        ManageInvoiceResultCode = -2
	ManageInvoiceResultCodeInvoiceOverflow        ManageInvoiceResultCode = -3
	ManageInvoiceResultCodeNotFound               ManageInvoiceResultCode = -4
	ManageInvoiceResultCodeTooManyInvoices        ManageInvoiceResultCode = -5
	ManageInvoiceResultCodeCanNotDeleteInProgress ManageInvoiceResultCode = -6
)

var ManageInvoiceResultCodeAll = []ManageInvoiceResultCode{
	ManageInvoiceResultCodeSuccess,
	ManageInvoiceResultCodeMalformed,
	ManageInvoiceResultCodeBalanceNotFound,
	ManageInvoiceResultCodeInvoiceOverflow,
	ManageInvoiceResultCodeNotFound,
	ManageInvoiceResultCodeTooManyInvoices,
	ManageInvoiceResultCodeCanNotDeleteInProgress,
}

var manageInvoiceResultCodeMap = map[int32]string{
	0:  "ManageInvoiceResultCodeSuccess",
	-1: "ManageInvoiceResultCodeMalformed",
	-2: "ManageInvoiceResultCodeBalanceNotFound",
	-3: "ManageInvoiceResultCodeInvoiceOverflow",
	-4: "ManageInvoiceResultCodeNotFound",
	-5: "ManageInvoiceResultCodeTooManyInvoices",
	-6: "ManageInvoiceResultCodeCanNotDeleteInProgress",
}

var manageInvoiceResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "malformed",
	-2: "balance_not_found",
	-3: "invoice_overflow",
	-4: "not_found",
	-5: "too_many_invoices",
	-6: "can_not_delete_in_progress",
}

var manageInvoiceResultCodeRevMap = map[string]int32{
	"ManageInvoiceResultCodeSuccess":                0,
	"ManageInvoiceResultCodeMalformed":              -1,
	"ManageInvoiceResultCodeBalanceNotFound":        -2,
	"ManageInvoiceResultCodeInvoiceOverflow":        -3,
	"ManageInvoiceResultCodeNotFound":               -4,
	"ManageInvoiceResultCodeTooManyInvoices":        -5,
	"ManageInvoiceResultCodeCanNotDeleteInProgress": -6,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageInvoiceResultCode
func (e ManageInvoiceResultCode) ValidEnum(v int32) bool {
	_, ok := manageInvoiceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageInvoiceResultCode) String() string {
	name, _ := manageInvoiceResultCodeMap[int32(e)]
	return name
}

func (e ManageInvoiceResultCode) ShortString() string {
	name, _ := manageInvoiceResultCodeShortMap[int32(e)]
	return name
}

func (e ManageInvoiceResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageInvoiceResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageInvoiceResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageInvoiceResultCode(value)
	return nil
}

// ManageInvoiceSuccessResultExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageInvoiceSuccessResultExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageInvoiceSuccessResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageInvoiceSuccessResultExt
func (u ManageInvoiceSuccessResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageInvoiceSuccessResultExt creates a new  ManageInvoiceSuccessResultExt.
func NewManageInvoiceSuccessResultExt(v LedgerVersion, value interface{}) (result ManageInvoiceSuccessResultExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageInvoiceSuccessResult is an XDR Struct defines as:
//
//   struct ManageInvoiceSuccessResult
//    {
//    	uint64 invoiceID;
//    	AssetCode asset;
//    	BalanceID senderBalance;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageInvoiceSuccessResult struct {
	InvoiceId     Uint64                        `json:"invoiceID,omitempty"`
	Asset         AssetCode                     `json:"asset,omitempty"`
	SenderBalance BalanceId                     `json:"senderBalance,omitempty"`
	Ext           ManageInvoiceSuccessResultExt `json:"ext,omitempty"`
}

// ManageInvoiceResult is an XDR Union defines as:
//
//   union ManageInvoiceResult switch (ManageInvoiceResultCode code)
//    {
//    case SUCCESS:
//        ManageInvoiceSuccessResult success;
//    default:
//        void;
//    };
//
type ManageInvoiceResult struct {
	Code    ManageInvoiceResultCode     `json:"code,omitempty"`
	Success *ManageInvoiceSuccessResult `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageInvoiceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageInvoiceResult
func (u ManageInvoiceResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageInvoiceResultCode(sw) {
	case ManageInvoiceResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageInvoiceResult creates a new  ManageInvoiceResult.
func NewManageInvoiceResult(code ManageInvoiceResultCode, value interface{}) (result ManageInvoiceResult, err error) {
	result.Code = code
	switch ManageInvoiceResultCode(code) {
	case ManageInvoiceResultCodeSuccess:
		tv, ok := value.(ManageInvoiceSuccessResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageInvoiceSuccessResult")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageInvoiceResult) MustSuccess() ManageInvoiceSuccessResult {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageInvoiceResult) GetSuccess() (result ManageInvoiceSuccessResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// CoinsEmissionRequestEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CoinsEmissionRequestEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CoinsEmissionRequestEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CoinsEmissionRequestEntryExt
func (u CoinsEmissionRequestEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCoinsEmissionRequestEntryExt creates a new  CoinsEmissionRequestEntryExt.
func NewCoinsEmissionRequestEntryExt(v LedgerVersion, value interface{}) (result CoinsEmissionRequestEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CoinsEmissionRequestEntry is an XDR Struct defines as:
//
//   struct CoinsEmissionRequestEntry
//    {
//    	uint64 requestID;
//        string64 reference;
//        BalanceID receiver;
//    	AccountID issuer;
//        int64 amount;
//        AssetCode asset;
//    	bool isApproved;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CoinsEmissionRequestEntry struct {
	RequestId  Uint64                       `json:"requestID,omitempty"`
	Reference  String64                     `json:"reference,omitempty"`
	Receiver   BalanceId                    `json:"receiver,omitempty"`
	Issuer     AccountId                    `json:"issuer,omitempty"`
	Amount     Int64                        `json:"amount,omitempty"`
	Asset      AssetCode                    `json:"asset,omitempty"`
	IsApproved bool                         `json:"isApproved,omitempty"`
	Ext        CoinsEmissionRequestEntryExt `json:"ext,omitempty"`
}

// CoinsEmissionEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CoinsEmissionEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CoinsEmissionEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CoinsEmissionEntryExt
func (u CoinsEmissionEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCoinsEmissionEntryExt creates a new  CoinsEmissionEntryExt.
func NewCoinsEmissionEntryExt(v LedgerVersion, value interface{}) (result CoinsEmissionEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CoinsEmissionEntry is an XDR Struct defines as:
//
//   struct CoinsEmissionEntry
//    {
//    	string64 serialNumber;
//        int64 amount;
//        AssetCode asset;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CoinsEmissionEntry struct {
	SerialNumber String64              `json:"serialNumber,omitempty"`
	Amount       Int64                 `json:"amount,omitempty"`
	Asset        AssetCode             `json:"asset,omitempty"`
	Ext          CoinsEmissionEntryExt `json:"ext,omitempty"`
}

// AccountLimitsEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AccountLimitsEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountLimitsEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountLimitsEntryExt
func (u AccountLimitsEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAccountLimitsEntryExt creates a new  AccountLimitsEntryExt.
func NewAccountLimitsEntryExt(v LedgerVersion, value interface{}) (result AccountLimitsEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AccountLimitsEntry is an XDR Struct defines as:
//
//   struct AccountLimitsEntry
//    {
//        AccountID accountID;
//        Limits limits;
//
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AccountLimitsEntry struct {
	AccountId AccountId             `json:"accountID,omitempty"`
	Limits    Limits                `json:"limits,omitempty"`
	Ext       AccountLimitsEntryExt `json:"ext,omitempty"`
}

// InvoiceReferenceExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type InvoiceReferenceExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InvoiceReferenceExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InvoiceReferenceExt
func (u InvoiceReferenceExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewInvoiceReferenceExt creates a new  InvoiceReferenceExt.
func NewInvoiceReferenceExt(v LedgerVersion, value interface{}) (result InvoiceReferenceExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// InvoiceReference is an XDR Struct defines as:
//
//   struct InvoiceReference {
//        uint64 invoiceID;
//        bool accept;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type InvoiceReference struct {
	InvoiceId Uint64              `json:"invoiceID,omitempty"`
	Accept    bool                `json:"accept,omitempty"`
	Ext       InvoiceReferenceExt `json:"ext,omitempty"`
}

// FeeDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type FeeDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeDataExt
func (u FeeDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewFeeDataExt creates a new  FeeDataExt.
func NewFeeDataExt(v LedgerVersion, value interface{}) (result FeeDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// FeeData is an XDR Struct defines as:
//
//   struct FeeData {
//        int64 paymentFee;
//        int64 fixedFee;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type FeeData struct {
	PaymentFee Int64      `json:"paymentFee,omitempty"`
	FixedFee   Int64      `json:"fixedFee,omitempty"`
	Ext        FeeDataExt `json:"ext,omitempty"`
}

// PaymentFeeDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentFeeDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentFeeDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentFeeDataExt
func (u PaymentFeeDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentFeeDataExt creates a new  PaymentFeeDataExt.
func NewPaymentFeeDataExt(v LedgerVersion, value interface{}) (result PaymentFeeDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentFeeData is an XDR Struct defines as:
//
//   struct PaymentFeeData {
//        FeeData sourceFee;
//        FeeData destinationFee;
//        bool sourcePaysForDest;    // if true source account pays fee, else destination
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentFeeData struct {
	SourceFee         FeeData           `json:"sourceFee,omitempty"`
	DestinationFee    FeeData           `json:"destinationFee,omitempty"`
	SourcePaysForDest bool              `json:"sourcePaysForDest,omitempty"`
	Ext               PaymentFeeDataExt `json:"ext,omitempty"`
}

// PaymentOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentOpExt
func (u PaymentOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentOpExt creates a new  PaymentOpExt.
func NewPaymentOpExt(v LedgerVersion, value interface{}) (result PaymentOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentOp is an XDR Struct defines as:
//
//   struct PaymentOp
//    {
//        BalanceID sourceBalanceID;
//        BalanceID destinationBalanceID;
//        int64 amount;          // amount they end up with
//
//        PaymentFeeData feeData;
//
//        string256 subject;
//        string64 reference;
//
//        InvoiceReference* invoiceReference;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentOp struct {
	SourceBalanceId      BalanceId         `json:"sourceBalanceID,omitempty"`
	DestinationBalanceId BalanceId         `json:"destinationBalanceID,omitempty"`
	Amount               Int64             `json:"amount,omitempty"`
	FeeData              PaymentFeeData    `json:"feeData,omitempty"`
	Subject              String256         `json:"subject,omitempty"`
	Reference            String64          `json:"reference,omitempty"`
	InvoiceReference     *InvoiceReference `json:"invoiceReference,omitempty"`
	Ext                  PaymentOpExt      `json:"ext,omitempty"`
}

// PaymentResultCode is an XDR Enum defines as:
//
//   enum PaymentResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0, // payment successfuly completed
//
//        // codes considered as "failure" for the operation
//        MALFORMED = -1,       // bad input
//        UNDERFUNDED = -2,     // not enough funds in source account
//        LINE_FULL = -3,       // destination would go above their limit
//    	FEE_MISMATCHED = -4,   // fee is not equal to expected fee
//        BALANCE_NOT_FOUND = -5, // destination balance not found
//        BALANCE_ACCOUNT_MISMATCHED = -6,
//        BALANCE_ASSETS_MISMATCHED = -7,
//    	SRC_BALANCE_NOT_FOUND = -8, // source balance not found
//        REFERENCE_DUPLICATION = -9,
//        STATS_OVERFLOW = -10,
//        LIMITS_EXCEEDED = -11,
//        NOT_ALLOWED_BY_ASSET_POLICY = -12,
//        INVOICE_NOT_FOUND = -13,
//        INVOICE_WRONG_AMOUNT = -14,
//        INVOICE_BALANCE_MISMATCH = -15,
//        INVOICE_ACCOUNT_MISMATCH = -16,
//        INVOICE_ALREADY_PAID = -17
//    };
//
type PaymentResultCode int32

const (
	PaymentResultCodeSuccess                  PaymentResultCode = 0
	PaymentResultCodeMalformed                PaymentResultCode = -1
	PaymentResultCodeUnderfunded              PaymentResultCode = -2
	PaymentResultCodeLineFull                 PaymentResultCode = -3
	PaymentResultCodeFeeMismatched            PaymentResultCode = -4
	PaymentResultCodeBalanceNotFound          PaymentResultCode = -5
	PaymentResultCodeBalanceAccountMismatched PaymentResultCode = -6
	PaymentResultCodeBalanceAssetsMismatched  PaymentResultCode = -7
	PaymentResultCodeSrcBalanceNotFound       PaymentResultCode = -8
	PaymentResultCodeReferenceDuplication     PaymentResultCode = -9
	PaymentResultCodeStatsOverflow            PaymentResultCode = -10
	PaymentResultCodeLimitsExceeded           PaymentResultCode = -11
	PaymentResultCodeNotAllowedByAssetPolicy  PaymentResultCode = -12
	PaymentResultCodeInvoiceNotFound          PaymentResultCode = -13
	PaymentResultCodeInvoiceWrongAmount       PaymentResultCode = -14
	PaymentResultCodeInvoiceBalanceMismatch   PaymentResultCode = -15
	PaymentResultCodeInvoiceAccountMismatch   PaymentResultCode = -16
	PaymentResultCodeInvoiceAlreadyPaid       PaymentResultCode = -17
)

var PaymentResultCodeAll = []PaymentResultCode{
	PaymentResultCodeSuccess,
	PaymentResultCodeMalformed,
	PaymentResultCodeUnderfunded,
	PaymentResultCodeLineFull,
	PaymentResultCodeFeeMismatched,
	PaymentResultCodeBalanceNotFound,
	PaymentResultCodeBalanceAccountMismatched,
	PaymentResultCodeBalanceAssetsMismatched,
	PaymentResultCodeSrcBalanceNotFound,
	PaymentResultCodeReferenceDuplication,
	PaymentResultCodeStatsOverflow,
	PaymentResultCodeLimitsExceeded,
	PaymentResultCodeNotAllowedByAssetPolicy,
	PaymentResultCodeInvoiceNotFound,
	PaymentResultCodeInvoiceWrongAmount,
	PaymentResultCodeInvoiceBalanceMismatch,
	PaymentResultCodeInvoiceAccountMismatch,
	PaymentResultCodeInvoiceAlreadyPaid,
}

var paymentResultCodeMap = map[int32]string{
	0:   "PaymentResultCodeSuccess",
	-1:  "PaymentResultCodeMalformed",
	-2:  "PaymentResultCodeUnderfunded",
	-3:  "PaymentResultCodeLineFull",
	-4:  "PaymentResultCodeFeeMismatched",
	-5:  "PaymentResultCodeBalanceNotFound",
	-6:  "PaymentResultCodeBalanceAccountMismatched",
	-7:  "PaymentResultCodeBalanceAssetsMismatched",
	-8:  "PaymentResultCodeSrcBalanceNotFound",
	-9:  "PaymentResultCodeReferenceDuplication",
	-10: "PaymentResultCodeStatsOverflow",
	-11: "PaymentResultCodeLimitsExceeded",
	-12: "PaymentResultCodeNotAllowedByAssetPolicy",
	-13: "PaymentResultCodeInvoiceNotFound",
	-14: "PaymentResultCodeInvoiceWrongAmount",
	-15: "PaymentResultCodeInvoiceBalanceMismatch",
	-16: "PaymentResultCodeInvoiceAccountMismatch",
	-17: "PaymentResultCodeInvoiceAlreadyPaid",
}

var paymentResultCodeShortMap = map[int32]string{
	0:   "success",
	-1:  "malformed",
	-2:  "underfunded",
	-3:  "line_full",
	-4:  "fee_mismatched",
	-5:  "balance_not_found",
	-6:  "balance_account_mismatched",
	-7:  "balance_assets_mismatched",
	-8:  "src_balance_not_found",
	-9:  "reference_duplication",
	-10: "stats_overflow",
	-11: "limits_exceeded",
	-12: "not_allowed_by_asset_policy",
	-13: "invoice_not_found",
	-14: "invoice_wrong_amount",
	-15: "invoice_balance_mismatch",
	-16: "invoice_account_mismatch",
	-17: "invoice_already_paid",
}

var paymentResultCodeRevMap = map[string]int32{
	"PaymentResultCodeSuccess":                  0,
	"PaymentResultCodeMalformed":                -1,
	"PaymentResultCodeUnderfunded":              -2,
	"PaymentResultCodeLineFull":                 -3,
	"PaymentResultCodeFeeMismatched":            -4,
	"PaymentResultCodeBalanceNotFound":          -5,
	"PaymentResultCodeBalanceAccountMismatched": -6,
	"PaymentResultCodeBalanceAssetsMismatched":  -7,
	"PaymentResultCodeSrcBalanceNotFound":       -8,
	"PaymentResultCodeReferenceDuplication":     -9,
	"PaymentResultCodeStatsOverflow":            -10,
	"PaymentResultCodeLimitsExceeded":           -11,
	"PaymentResultCodeNotAllowedByAssetPolicy":  -12,
	"PaymentResultCodeInvoiceNotFound":          -13,
	"PaymentResultCodeInvoiceWrongAmount":       -14,
	"PaymentResultCodeInvoiceBalanceMismatch":   -15,
	"PaymentResultCodeInvoiceAccountMismatch":   -16,
	"PaymentResultCodeInvoiceAlreadyPaid":       -17,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PaymentResultCode
func (e PaymentResultCode) ValidEnum(v int32) bool {
	_, ok := paymentResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e PaymentResultCode) String() string {
	name, _ := paymentResultCodeMap[int32(e)]
	return name
}

func (e PaymentResultCode) ShortString() string {
	name, _ := paymentResultCodeShortMap[int32(e)]
	return name
}

func (e PaymentResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *PaymentResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := paymentResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = PaymentResultCode(value)
	return nil
}

// PaymentResponseExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentResponseExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentResponseExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentResponseExt
func (u PaymentResponseExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentResponseExt creates a new  PaymentResponseExt.
func NewPaymentResponseExt(v LedgerVersion, value interface{}) (result PaymentResponseExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentResponse is an XDR Struct defines as:
//
//   struct PaymentResponse {
//        AccountID destination;
//        uint64 paymentID;
//        AssetCode asset;
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentResponse struct {
	Destination AccountId          `json:"destination,omitempty"`
	PaymentId   Uint64             `json:"paymentID,omitempty"`
	Asset       AssetCode          `json:"asset,omitempty"`
	Ext         PaymentResponseExt `json:"ext,omitempty"`
}

// PaymentResult is an XDR Union defines as:
//
//   union PaymentResult switch (PaymentResultCode code)
//    {
//    case SUCCESS:
//        PaymentResponse paymentResponse;
//    default:
//        void;
//    };
//
type PaymentResult struct {
	Code            PaymentResultCode `json:"code,omitempty"`
	PaymentResponse *PaymentResponse  `json:"paymentResponse,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentResult
func (u PaymentResult) ArmForSwitch(sw int32) (string, bool) {
	switch PaymentResultCode(sw) {
	case PaymentResultCodeSuccess:
		return "PaymentResponse", true
	default:
		return "", true
	}
}

// NewPaymentResult creates a new  PaymentResult.
func NewPaymentResult(code PaymentResultCode, value interface{}) (result PaymentResult, err error) {
	result.Code = code
	switch PaymentResultCode(code) {
	case PaymentResultCodeSuccess:
		tv, ok := value.(PaymentResponse)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentResponse")
			return
		}
		result.PaymentResponse = &tv
	default:
		// void
	}
	return
}

// MustPaymentResponse retrieves the PaymentResponse value from the union,
// panicing if the value is not set.
func (u PaymentResult) MustPaymentResponse() PaymentResponse {
	val, ok := u.GetPaymentResponse()

	if !ok {
		panic("arm PaymentResponse is not set")
	}

	return val
}

// GetPaymentResponse retrieves the PaymentResponse value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PaymentResult) GetPaymentResponse() (result PaymentResponse, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "PaymentResponse" {
		result = *u.PaymentResponse
		ok = true
	}

	return
}

// RecoverOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type RecoverOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RecoverOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RecoverOpExt
func (u RecoverOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewRecoverOpExt creates a new  RecoverOpExt.
func NewRecoverOpExt(v LedgerVersion, value interface{}) (result RecoverOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// RecoverOp is an XDR Struct defines as:
//
//   struct RecoverOp
//    {
//        AccountID account;
//        PublicKey oldSigner;
//        PublicKey newSigner;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type RecoverOp struct {
	Account   AccountId    `json:"account,omitempty"`
	OldSigner PublicKey    `json:"oldSigner,omitempty"`
	NewSigner PublicKey    `json:"newSigner,omitempty"`
	Ext       RecoverOpExt `json:"ext,omitempty"`
}

// RecoverResultCode is an XDR Enum defines as:
//
//   enum RecoverResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//
//        MALFORMED = -1,
//        OLD_SIGNER_NOT_FOUND = -2,
//        SIGNER_ALREADY_EXISTS = -3
//    };
//
type RecoverResultCode int32

const (
	RecoverResultCodeSuccess             RecoverResultCode = 0
	RecoverResultCodeMalformed           RecoverResultCode = -1
	RecoverResultCodeOldSignerNotFound   RecoverResultCode = -2
	RecoverResultCodeSignerAlreadyExists RecoverResultCode = -3
)

var RecoverResultCodeAll = []RecoverResultCode{
	RecoverResultCodeSuccess,
	RecoverResultCodeMalformed,
	RecoverResultCodeOldSignerNotFound,
	RecoverResultCodeSignerAlreadyExists,
}

var recoverResultCodeMap = map[int32]string{
	0:  "RecoverResultCodeSuccess",
	-1: "RecoverResultCodeMalformed",
	-2: "RecoverResultCodeOldSignerNotFound",
	-3: "RecoverResultCodeSignerAlreadyExists",
}

var recoverResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "malformed",
	-2: "old_signer_not_found",
	-3: "signer_already_exists",
}

var recoverResultCodeRevMap = map[string]int32{
	"RecoverResultCodeSuccess":             0,
	"RecoverResultCodeMalformed":           -1,
	"RecoverResultCodeOldSignerNotFound":   -2,
	"RecoverResultCodeSignerAlreadyExists": -3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for RecoverResultCode
func (e RecoverResultCode) ValidEnum(v int32) bool {
	_, ok := recoverResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e RecoverResultCode) String() string {
	name, _ := recoverResultCodeMap[int32(e)]
	return name
}

func (e RecoverResultCode) ShortString() string {
	name, _ := recoverResultCodeShortMap[int32(e)]
	return name
}

func (e RecoverResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *RecoverResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := recoverResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = RecoverResultCode(value)
	return nil
}

// RecoverResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type RecoverResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RecoverResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RecoverResultSuccessExt
func (u RecoverResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewRecoverResultSuccessExt creates a new  RecoverResultSuccessExt.
func NewRecoverResultSuccessExt(v LedgerVersion, value interface{}) (result RecoverResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// RecoverResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type RecoverResultSuccess struct {
	Ext RecoverResultSuccessExt `json:"ext,omitempty"`
}

// RecoverResult is an XDR Union defines as:
//
//   union RecoverResult switch (RecoverResultCode code)
//    {
//    case SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type RecoverResult struct {
	Code    RecoverResultCode     `json:"code,omitempty"`
	Success *RecoverResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RecoverResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RecoverResult
func (u RecoverResult) ArmForSwitch(sw int32) (string, bool) {
	switch RecoverResultCode(sw) {
	case RecoverResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewRecoverResult creates a new  RecoverResult.
func NewRecoverResult(code RecoverResultCode, value interface{}) (result RecoverResult, err error) {
	result.Code = code
	switch RecoverResultCode(code) {
	case RecoverResultCodeSuccess:
		tv, ok := value.(RecoverResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be RecoverResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u RecoverResult) MustSuccess() RecoverResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u RecoverResult) GetSuccess() (result RecoverResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageForfeitRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageForfeitRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageForfeitRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageForfeitRequestOpExt
func (u ManageForfeitRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageForfeitRequestOpExt creates a new  ManageForfeitRequestOpExt.
func NewManageForfeitRequestOpExt(v LedgerVersion, value interface{}) (result ManageForfeitRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageForfeitRequestOp is an XDR Struct defines as:
//
//   struct ManageForfeitRequestOp
//    {
//        BalanceID balance;
//        int64 amount;
//    	int64 totalFee;
//        string details<>;
//    	AccountID reviewer;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//
//    };
//
type ManageForfeitRequestOp struct {
	Balance  BalanceId                 `json:"balance,omitempty"`
	Amount   Int64                     `json:"amount,omitempty"`
	TotalFee Int64                     `json:"totalFee,omitempty"`
	Details  string                    `json:"details,omitempty"`
	Reviewer AccountId                 `json:"reviewer,omitempty"`
	Ext      ManageForfeitRequestOpExt `json:"ext,omitempty"`
}

// ManageForfeitRequestResultCode is an XDR Enum defines as:
//
//   enum ManageForfeitRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	UNDERFUNDED = -1,
//        INVALID_AMOUNT = -2,
//        LINE_FULL = -3,
//        BALANCE_MISMATCH = -4,
//        STATS_OVERFLOW = -5,
//        LIMITS_EXCEEDED = -6,
//        REVIEWER_NOT_FOUND = -7,
//        INVALID_DETAILS = -8,
//    	FEE_MISMATCH = -9 // fee is not equal to expected fee
//    };
//
type ManageForfeitRequestResultCode int32

const (
	ManageForfeitRequestResultCodeSuccess          ManageForfeitRequestResultCode = 0
	ManageForfeitRequestResultCodeUnderfunded      ManageForfeitRequestResultCode = -1
	ManageForfeitRequestResultCodeInvalidAmount    ManageForfeitRequestResultCode = -2
	ManageForfeitRequestResultCodeLineFull         ManageForfeitRequestResultCode = -3
	ManageForfeitRequestResultCodeBalanceMismatch  ManageForfeitRequestResultCode = -4
	ManageForfeitRequestResultCodeStatsOverflow    ManageForfeitRequestResultCode = -5
	ManageForfeitRequestResultCodeLimitsExceeded   ManageForfeitRequestResultCode = -6
	ManageForfeitRequestResultCodeReviewerNotFound ManageForfeitRequestResultCode = -7
	ManageForfeitRequestResultCodeInvalidDetails   ManageForfeitRequestResultCode = -8
	ManageForfeitRequestResultCodeFeeMismatch      ManageForfeitRequestResultCode = -9
)

var ManageForfeitRequestResultCodeAll = []ManageForfeitRequestResultCode{
	ManageForfeitRequestResultCodeSuccess,
	ManageForfeitRequestResultCodeUnderfunded,
	ManageForfeitRequestResultCodeInvalidAmount,
	ManageForfeitRequestResultCodeLineFull,
	ManageForfeitRequestResultCodeBalanceMismatch,
	ManageForfeitRequestResultCodeStatsOverflow,
	ManageForfeitRequestResultCodeLimitsExceeded,
	ManageForfeitRequestResultCodeReviewerNotFound,
	ManageForfeitRequestResultCodeInvalidDetails,
	ManageForfeitRequestResultCodeFeeMismatch,
}

var manageForfeitRequestResultCodeMap = map[int32]string{
	0:  "ManageForfeitRequestResultCodeSuccess",
	-1: "ManageForfeitRequestResultCodeUnderfunded",
	-2: "ManageForfeitRequestResultCodeInvalidAmount",
	-3: "ManageForfeitRequestResultCodeLineFull",
	-4: "ManageForfeitRequestResultCodeBalanceMismatch",
	-5: "ManageForfeitRequestResultCodeStatsOverflow",
	-6: "ManageForfeitRequestResultCodeLimitsExceeded",
	-7: "ManageForfeitRequestResultCodeReviewerNotFound",
	-8: "ManageForfeitRequestResultCodeInvalidDetails",
	-9: "ManageForfeitRequestResultCodeFeeMismatch",
}

var manageForfeitRequestResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "underfunded",
	-2: "invalid_amount",
	-3: "line_full",
	-4: "balance_mismatch",
	-5: "stats_overflow",
	-6: "limits_exceeded",
	-7: "reviewer_not_found",
	-8: "invalid_details",
	-9: "fee_mismatch",
}

var manageForfeitRequestResultCodeRevMap = map[string]int32{
	"ManageForfeitRequestResultCodeSuccess":          0,
	"ManageForfeitRequestResultCodeUnderfunded":      -1,
	"ManageForfeitRequestResultCodeInvalidAmount":    -2,
	"ManageForfeitRequestResultCodeLineFull":         -3,
	"ManageForfeitRequestResultCodeBalanceMismatch":  -4,
	"ManageForfeitRequestResultCodeStatsOverflow":    -5,
	"ManageForfeitRequestResultCodeLimitsExceeded":   -6,
	"ManageForfeitRequestResultCodeReviewerNotFound": -7,
	"ManageForfeitRequestResultCodeInvalidDetails":   -8,
	"ManageForfeitRequestResultCodeFeeMismatch":      -9,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageForfeitRequestResultCode
func (e ManageForfeitRequestResultCode) ValidEnum(v int32) bool {
	_, ok := manageForfeitRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageForfeitRequestResultCode) String() string {
	name, _ := manageForfeitRequestResultCodeMap[int32(e)]
	return name
}

func (e ManageForfeitRequestResultCode) ShortString() string {
	name, _ := manageForfeitRequestResultCodeShortMap[int32(e)]
	return name
}

func (e ManageForfeitRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageForfeitRequestResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageForfeitRequestResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageForfeitRequestResultCode(value)
	return nil
}

// ManageForfeitRequestResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//                {
//                case EMPTY_VERSION:
//                    void;
//                }
//
type ManageForfeitRequestResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageForfeitRequestResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageForfeitRequestResultSuccessExt
func (u ManageForfeitRequestResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageForfeitRequestResultSuccessExt creates a new  ManageForfeitRequestResultSuccessExt.
func NewManageForfeitRequestResultSuccessExt(v LedgerVersion, value interface{}) (result ManageForfeitRequestResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageForfeitRequestResultSuccess is an XDR NestedStruct defines as:
//
//   struct
//            {
//                uint64 paymentID;
//
//                union switch (LedgerVersion v)
//                {
//                case EMPTY_VERSION:
//                    void;
//                }
//                ext;
//            }
//
type ManageForfeitRequestResultSuccess struct {
	PaymentId Uint64                               `json:"paymentID,omitempty"`
	Ext       ManageForfeitRequestResultSuccessExt `json:"ext,omitempty"`
}

// ManageForfeitRequestResult is an XDR Union defines as:
//
//   union ManageForfeitRequestResult switch (ManageForfeitRequestResultCode code)
//    {
//        case SUCCESS:
//            struct
//            {
//                uint64 paymentID;
//
//                union switch (LedgerVersion v)
//                {
//                case EMPTY_VERSION:
//                    void;
//                }
//                ext;
//            } success;
//        default:
//            void;
//    };
//
type ManageForfeitRequestResult struct {
	Code    ManageForfeitRequestResultCode     `json:"code,omitempty"`
	Success *ManageForfeitRequestResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageForfeitRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageForfeitRequestResult
func (u ManageForfeitRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageForfeitRequestResultCode(sw) {
	case ManageForfeitRequestResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageForfeitRequestResult creates a new  ManageForfeitRequestResult.
func NewManageForfeitRequestResult(code ManageForfeitRequestResultCode, value interface{}) (result ManageForfeitRequestResult, err error) {
	result.Code = code
	switch ManageForfeitRequestResultCode(code) {
	case ManageForfeitRequestResultCodeSuccess:
		tv, ok := value.(ManageForfeitRequestResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageForfeitRequestResult) MustSuccess() ManageForfeitRequestResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageForfeitRequestResult) GetSuccess() (result ManageForfeitRequestResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// SignerType is an XDR Enum defines as:
//
//   enum SignerType
//    {
//    	READER = 1,                  // can only read data from API and Horizon
//    	NOT_VERIFIED_ACC_MANAGER = 2,// can manage not verified account and block/unblock general
//    	GENERAL_ACC_MANAGER = 4,     // allowed to create account, block/unblock, change limits for particular general account
//    	DIRECT_DEBIT_OPERATOR = 8, // allowed to perform direct debit operation
//    	ASSET_MANAGER = 16, // allowed to create assets/asset pairs and update policies, set fees
//    	ASSET_RATE_MANAGER = 32, // allowed to set physical asset price
//    	BALANCE_MANAGER = 64, // allowed to create balances, spend assets from balances
//    	EMISSION_MANAGER = 128, // allowed to make emission requests, review emission, upload preemission
//    	INVOICE_MANAGER = 256, // allowed to create payment requests to other accounts
//    	PAYMENT_OPERATOR = 512, // allowed to review payment requests
//    	LIMITS_MANAGER = 1024, // allowed to change limits
//    	ACCOUNT_MANAGER = 2048, // allowed to add/delete signers and trust
//    	COMMISSION_BALANCE_MANAGER  = 4096,// allowed to spend from commission balances
//    	OPERATIONAL_BALANCE_MANAGER = 8192 // allowed to spend from operational balances
//    };
//
type SignerType int32

const (
	SignerTypeReader                    SignerType = 1
	SignerTypeNotVerifiedAccManager     SignerType = 2
	SignerTypeGeneralAccManager         SignerType = 4
	SignerTypeDirectDebitOperator       SignerType = 8
	SignerTypeAssetManager              SignerType = 16
	SignerTypeAssetRateManager          SignerType = 32
	SignerTypeBalanceManager            SignerType = 64
	SignerTypeEmissionManager           SignerType = 128
	SignerTypeInvoiceManager            SignerType = 256
	SignerTypePaymentOperator           SignerType = 512
	SignerTypeLimitsManager             SignerType = 1024
	SignerTypeAccountManager            SignerType = 2048
	SignerTypeCommissionBalanceManager  SignerType = 4096
	SignerTypeOperationalBalanceManager SignerType = 8192
)

var SignerTypeAll = []SignerType{
	SignerTypeReader,
	SignerTypeNotVerifiedAccManager,
	SignerTypeGeneralAccManager,
	SignerTypeDirectDebitOperator,
	SignerTypeAssetManager,
	SignerTypeAssetRateManager,
	SignerTypeBalanceManager,
	SignerTypeEmissionManager,
	SignerTypeInvoiceManager,
	SignerTypePaymentOperator,
	SignerTypeLimitsManager,
	SignerTypeAccountManager,
	SignerTypeCommissionBalanceManager,
	SignerTypeOperationalBalanceManager,
}

var signerTypeMap = map[int32]string{
	1:    "SignerTypeReader",
	2:    "SignerTypeNotVerifiedAccManager",
	4:    "SignerTypeGeneralAccManager",
	8:    "SignerTypeDirectDebitOperator",
	16:   "SignerTypeAssetManager",
	32:   "SignerTypeAssetRateManager",
	64:   "SignerTypeBalanceManager",
	128:  "SignerTypeEmissionManager",
	256:  "SignerTypeInvoiceManager",
	512:  "SignerTypePaymentOperator",
	1024: "SignerTypeLimitsManager",
	2048: "SignerTypeAccountManager",
	4096: "SignerTypeCommissionBalanceManager",
	8192: "SignerTypeOperationalBalanceManager",
}

var signerTypeShortMap = map[int32]string{
	1:    "reader",
	2:    "not_verified_acc_manager",
	4:    "general_acc_manager",
	8:    "direct_debit_operator",
	16:   "asset_manager",
	32:   "asset_rate_manager",
	64:   "balance_manager",
	128:  "emission_manager",
	256:  "invoice_manager",
	512:  "payment_operator",
	1024: "limits_manager",
	2048: "account_manager",
	4096: "commission_balance_manager",
	8192: "operational_balance_manager",
}

var signerTypeRevMap = map[string]int32{
	"SignerTypeReader":                    1,
	"SignerTypeNotVerifiedAccManager":     2,
	"SignerTypeGeneralAccManager":         4,
	"SignerTypeDirectDebitOperator":       8,
	"SignerTypeAssetManager":              16,
	"SignerTypeAssetRateManager":          32,
	"SignerTypeBalanceManager":            64,
	"SignerTypeEmissionManager":           128,
	"SignerTypeInvoiceManager":            256,
	"SignerTypePaymentOperator":           512,
	"SignerTypeLimitsManager":             1024,
	"SignerTypeAccountManager":            2048,
	"SignerTypeCommissionBalanceManager":  4096,
	"SignerTypeOperationalBalanceManager": 8192,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SignerType
func (e SignerType) ValidEnum(v int32) bool {
	_, ok := signerTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SignerType) String() string {
	name, _ := signerTypeMap[int32(e)]
	return name
}

func (e SignerType) ShortString() string {
	name, _ := signerTypeShortMap[int32(e)]
	return name
}

func (e SignerType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *SignerType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := signerTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = SignerType(value)
	return nil
}

// SignerExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type SignerExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SignerExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SignerExt
func (u SignerExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSignerExt creates a new  SignerExt.
func NewSignerExt(v LedgerVersion, value interface{}) (result SignerExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// Signer is an XDR Struct defines as:
//
//   struct Signer
//    {
//        AccountID pubKey;
//        uint32 weight; // really only need 1byte
//    	uint32 signerType;
//    	uint32 identity;
//    	string256 name;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type Signer struct {
	PubKey     AccountId `json:"pubKey,omitempty"`
	Weight     Uint32    `json:"weight,omitempty"`
	SignerType Uint32    `json:"signerType,omitempty"`
	Identity   Uint32    `json:"identity,omitempty"`
	Name       String256 `json:"name,omitempty"`
	Ext        SignerExt `json:"ext,omitempty"`
}

// TrustEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TrustEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustEntryExt
func (u TrustEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTrustEntryExt creates a new  TrustEntryExt.
func NewTrustEntryExt(v LedgerVersion, value interface{}) (result TrustEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TrustEntry is an XDR Struct defines as:
//
//   struct TrustEntry
//    {
//        AccountID allowedAccount;
//        BalanceID balanceToUse;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TrustEntry struct {
	AllowedAccount AccountId     `json:"allowedAccount,omitempty"`
	BalanceToUse   BalanceId     `json:"balanceToUse,omitempty"`
	Ext            TrustEntryExt `json:"ext,omitempty"`
}

// LimitsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LimitsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LimitsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LimitsExt
func (u LimitsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLimitsExt creates a new  LimitsExt.
func NewLimitsExt(v LedgerVersion, value interface{}) (result LimitsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// Limits is an XDR Struct defines as:
//
//   struct Limits
//    {
//        int64 dailyOut;
//    	int64 weeklyOut;
//    	int64 monthlyOut;
//        int64 annualOut;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//
//    };
//
type Limits struct {
	DailyOut   Int64     `json:"dailyOut,omitempty"`
	WeeklyOut  Int64     `json:"weeklyOut,omitempty"`
	MonthlyOut Int64     `json:"monthlyOut,omitempty"`
	AnnualOut  Int64     `json:"annualOut,omitempty"`
	Ext        LimitsExt `json:"ext,omitempty"`
}

// AccountPolicies is an XDR Enum defines as:
//
//   enum AccountPolicies
//    {
//    	NO_PERMISSIONS = 0,
//    	ALLOW_TO_CREATE_USER_VIA_API = 1
//    };
//
type AccountPolicies int32

const (
	AccountPoliciesNoPermissions           AccountPolicies = 0
	AccountPoliciesAllowToCreateUserViaApi AccountPolicies = 1
)

var AccountPoliciesAll = []AccountPolicies{
	AccountPoliciesNoPermissions,
	AccountPoliciesAllowToCreateUserViaApi,
}

var accountPoliciesMap = map[int32]string{
	0: "AccountPoliciesNoPermissions",
	1: "AccountPoliciesAllowToCreateUserViaApi",
}

var accountPoliciesShortMap = map[int32]string{
	0: "no_permissions",
	1: "allow_to_create_user_via_api",
}

var accountPoliciesRevMap = map[string]int32{
	"AccountPoliciesNoPermissions":           0,
	"AccountPoliciesAllowToCreateUserViaApi": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AccountPolicies
func (e AccountPolicies) ValidEnum(v int32) bool {
	_, ok := accountPoliciesMap[v]
	return ok
}

// String returns the name of `e`
func (e AccountPolicies) String() string {
	name, _ := accountPoliciesMap[int32(e)]
	return name
}

func (e AccountPolicies) ShortString() string {
	name, _ := accountPoliciesShortMap[int32(e)]
	return name
}

func (e AccountPolicies) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *AccountPolicies) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := accountPoliciesRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = AccountPolicies(value)
	return nil
}

// AccountType is an XDR Enum defines as:
//
//   enum AccountType
//    {
//    	OPERATIONAL = 1,       // operational account of the system
//    	GENERAL = 2,           // general account can perform payments, setoptions, be source account for tx, etc.
//    	COMMISSION = 3,        // commission account
//    	MASTER = 4,            // master account
//        NOT_VERIFIED = 5
//    };
//
type AccountType int32

const (
	AccountTypeOperational AccountType = 1
	AccountTypeGeneral     AccountType = 2
	AccountTypeCommission  AccountType = 3
	AccountTypeMaster      AccountType = 4
	AccountTypeNotVerified AccountType = 5
)

var AccountTypeAll = []AccountType{
	AccountTypeOperational,
	AccountTypeGeneral,
	AccountTypeCommission,
	AccountTypeMaster,
	AccountTypeNotVerified,
}

var accountTypeMap = map[int32]string{
	1: "AccountTypeOperational",
	2: "AccountTypeGeneral",
	3: "AccountTypeCommission",
	4: "AccountTypeMaster",
	5: "AccountTypeNotVerified",
}

var accountTypeShortMap = map[int32]string{
	1: "operational",
	2: "general",
	3: "commission",
	4: "master",
	5: "not_verified",
}

var accountTypeRevMap = map[string]int32{
	"AccountTypeOperational": 1,
	"AccountTypeGeneral":     2,
	"AccountTypeCommission":  3,
	"AccountTypeMaster":      4,
	"AccountTypeNotVerified": 5,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AccountType
func (e AccountType) ValidEnum(v int32) bool {
	_, ok := accountTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e AccountType) String() string {
	name, _ := accountTypeMap[int32(e)]
	return name
}

func (e AccountType) ShortString() string {
	name, _ := accountTypeShortMap[int32(e)]
	return name
}

func (e AccountType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *AccountType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := accountTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = AccountType(value)
	return nil
}

// BlockReasons is an XDR Enum defines as:
//
//   enum BlockReasons
//    {
//    	RECOVERY_REQUEST = 1,
//    	KYC_UPDATE = 2,
//    	SUSPICIOUS_BEHAVIOR = 3
//    };
//
type BlockReasons int32

const (
	BlockReasonsRecoveryRequest    BlockReasons = 1
	BlockReasonsKycUpdate          BlockReasons = 2
	BlockReasonsSuspiciousBehavior BlockReasons = 3
)

var BlockReasonsAll = []BlockReasons{
	BlockReasonsRecoveryRequest,
	BlockReasonsKycUpdate,
	BlockReasonsSuspiciousBehavior,
}

var blockReasonsMap = map[int32]string{
	1: "BlockReasonsRecoveryRequest",
	2: "BlockReasonsKycUpdate",
	3: "BlockReasonsSuspiciousBehavior",
}

var blockReasonsShortMap = map[int32]string{
	1: "recovery_request",
	2: "kyc_update",
	3: "suspicious_behavior",
}

var blockReasonsRevMap = map[string]int32{
	"BlockReasonsRecoveryRequest":    1,
	"BlockReasonsKycUpdate":          2,
	"BlockReasonsSuspiciousBehavior": 3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BlockReasons
func (e BlockReasons) ValidEnum(v int32) bool {
	_, ok := blockReasonsMap[v]
	return ok
}

// String returns the name of `e`
func (e BlockReasons) String() string {
	name, _ := blockReasonsMap[int32(e)]
	return name
}

func (e BlockReasons) ShortString() string {
	name, _ := blockReasonsShortMap[int32(e)]
	return name
}

func (e BlockReasons) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *BlockReasons) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := blockReasonsRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = BlockReasons(value)
	return nil
}

// AccountEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AccountEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountEntryExt
func (u AccountEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAccountEntryExt creates a new  AccountEntryExt.
func NewAccountEntryExt(v LedgerVersion, value interface{}) (result AccountEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AccountEntry is an XDR Struct defines as:
//
//   struct AccountEntry
//    {
//        AccountID accountID;      // master public key for this account
//
//        // fields used for signatures
//        // thresholds stores unsigned bytes: [weight of master|low|medium|high]
//        Thresholds thresholds;
//
//        Signer signers<>; // possible signers for this account
//        Limits* limits;
//
//    	uint32 blockReasons;
//        AccountType accountType; // type of the account
//
//        // Referral marketing
//        AccountID* referrer;     // parent account
//        int64 shareForReferrer; // share of fee to pay parent
//
//        int32 policies;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AccountEntry struct {
	AccountId        AccountId       `json:"accountID,omitempty"`
	Thresholds       Thresholds      `json:"thresholds,omitempty"`
	Signers          []Signer        `json:"signers,omitempty"`
	Limits           *Limits         `json:"limits,omitempty"`
	BlockReasons     Uint32          `json:"blockReasons,omitempty"`
	AccountType      AccountType     `json:"accountType,omitempty"`
	Referrer         *AccountId      `json:"referrer,omitempty"`
	ShareForReferrer Int64           `json:"shareForReferrer,omitempty"`
	Policies         Int32           `json:"policies,omitempty"`
	Ext              AccountEntryExt `json:"ext,omitempty"`
}

// ManageTrustAction is an XDR Enum defines as:
//
//   enum ManageTrustAction
//    {
//        TRUST_ADD = 0,
//        TRUST_REMOVE = 1
//    };
//
type ManageTrustAction int32

const (
	ManageTrustActionTrustAdd    ManageTrustAction = 0
	ManageTrustActionTrustRemove ManageTrustAction = 1
)

var ManageTrustActionAll = []ManageTrustAction{
	ManageTrustActionTrustAdd,
	ManageTrustActionTrustRemove,
}

var manageTrustActionMap = map[int32]string{
	0: "ManageTrustActionTrustAdd",
	1: "ManageTrustActionTrustRemove",
}

var manageTrustActionShortMap = map[int32]string{
	0: "trust_add",
	1: "trust_remove",
}

var manageTrustActionRevMap = map[string]int32{
	"ManageTrustActionTrustAdd":    0,
	"ManageTrustActionTrustRemove": 1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageTrustAction
func (e ManageTrustAction) ValidEnum(v int32) bool {
	_, ok := manageTrustActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageTrustAction) String() string {
	name, _ := manageTrustActionMap[int32(e)]
	return name
}

func (e ManageTrustAction) ShortString() string {
	name, _ := manageTrustActionShortMap[int32(e)]
	return name
}

func (e ManageTrustAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ManageTrustAction) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := manageTrustActionRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ManageTrustAction(value)
	return nil
}

// TrustDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type TrustDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustDataExt
func (u TrustDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTrustDataExt creates a new  TrustDataExt.
func NewTrustDataExt(v LedgerVersion, value interface{}) (result TrustDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TrustData is an XDR Struct defines as:
//
//   struct TrustData {
//        TrustEntry trust;
//        ManageTrustAction action;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type TrustData struct {
	Trust  TrustEntry        `json:"trust,omitempty"`
	Action ManageTrustAction `json:"action,omitempty"`
	Ext    TrustDataExt      `json:"ext,omitempty"`
}

// SetOptionsOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type SetOptionsOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsOpExt
func (u SetOptionsOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetOptionsOpExt creates a new  SetOptionsOpExt.
func NewSetOptionsOpExt(v LedgerVersion, value interface{}) (result SetOptionsOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetOptionsOp is an XDR Struct defines as:
//
//   struct SetOptionsOp
//    {
//        // account threshold manipulation
//        uint32* masterWeight; // weight of the master account
//        uint32* lowThreshold;
//        uint32* medThreshold;
//        uint32* highThreshold;
//
//        // Add, update or remove a signer for the account
//        // signer is deleted if the weight is 0
//        Signer* signer;
//
//        TrustData* trustData;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//
//    };
//
type SetOptionsOp struct {
	MasterWeight  *Uint32         `json:"masterWeight,omitempty"`
	LowThreshold  *Uint32         `json:"lowThreshold,omitempty"`
	MedThreshold  *Uint32         `json:"medThreshold,omitempty"`
	HighThreshold *Uint32         `json:"highThreshold,omitempty"`
	Signer        *Signer         `json:"signer,omitempty"`
	TrustData     *TrustData      `json:"trustData,omitempty"`
	Ext           SetOptionsOpExt `json:"ext,omitempty"`
}

// SetOptionsResultCode is an XDR Enum defines as:
//
//   enum SetOptionsResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        TOO_MANY_SIGNERS = -1, // max number of signers already reached
//        THRESHOLD_OUT_OF_RANGE = -2, // bad value for weight/threshold
//        BAD_SIGNER = -3,             // signer cannot be masterkey
//        BALANCE_NOT_FOUND = -4,
//        TRUST_MALFORMED = -5,
//    	TRUST_TOO_MANY = -6,
//    	INVALID_SIGNER_VERSION = -7 // if signer version is higher than ledger version
//    };
//
type SetOptionsResultCode int32

const (
	SetOptionsResultCodeSuccess              SetOptionsResultCode = 0
	SetOptionsResultCodeTooManySigners       SetOptionsResultCode = -1
	SetOptionsResultCodeThresholdOutOfRange  SetOptionsResultCode = -2
	SetOptionsResultCodeBadSigner            SetOptionsResultCode = -3
	SetOptionsResultCodeBalanceNotFound      SetOptionsResultCode = -4
	SetOptionsResultCodeTrustMalformed       SetOptionsResultCode = -5
	SetOptionsResultCodeTrustTooMany         SetOptionsResultCode = -6
	SetOptionsResultCodeInvalidSignerVersion SetOptionsResultCode = -7
)

var SetOptionsResultCodeAll = []SetOptionsResultCode{
	SetOptionsResultCodeSuccess,
	SetOptionsResultCodeTooManySigners,
	SetOptionsResultCodeThresholdOutOfRange,
	SetOptionsResultCodeBadSigner,
	SetOptionsResultCodeBalanceNotFound,
	SetOptionsResultCodeTrustMalformed,
	SetOptionsResultCodeTrustTooMany,
	SetOptionsResultCodeInvalidSignerVersion,
}

var setOptionsResultCodeMap = map[int32]string{
	0:  "SetOptionsResultCodeSuccess",
	-1: "SetOptionsResultCodeTooManySigners",
	-2: "SetOptionsResultCodeThresholdOutOfRange",
	-3: "SetOptionsResultCodeBadSigner",
	-4: "SetOptionsResultCodeBalanceNotFound",
	-5: "SetOptionsResultCodeTrustMalformed",
	-6: "SetOptionsResultCodeTrustTooMany",
	-7: "SetOptionsResultCodeInvalidSignerVersion",
}

var setOptionsResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "too_many_signers",
	-2: "threshold_out_of_range",
	-3: "bad_signer",
	-4: "balance_not_found",
	-5: "trust_malformed",
	-6: "trust_too_many",
	-7: "invalid_signer_version",
}

var setOptionsResultCodeRevMap = map[string]int32{
	"SetOptionsResultCodeSuccess":              0,
	"SetOptionsResultCodeTooManySigners":       -1,
	"SetOptionsResultCodeThresholdOutOfRange":  -2,
	"SetOptionsResultCodeBadSigner":            -3,
	"SetOptionsResultCodeBalanceNotFound":      -4,
	"SetOptionsResultCodeTrustMalformed":       -5,
	"SetOptionsResultCodeTrustTooMany":         -6,
	"SetOptionsResultCodeInvalidSignerVersion": -7,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetOptionsResultCode
func (e SetOptionsResultCode) ValidEnum(v int32) bool {
	_, ok := setOptionsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetOptionsResultCode) String() string {
	name, _ := setOptionsResultCodeMap[int32(e)]
	return name
}

func (e SetOptionsResultCode) ShortString() string {
	name, _ := setOptionsResultCodeShortMap[int32(e)]
	return name
}

func (e SetOptionsResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *SetOptionsResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := setOptionsResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = SetOptionsResultCode(value)
	return nil
}

// SetOptionsResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type SetOptionsResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsResultSuccessExt
func (u SetOptionsResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetOptionsResultSuccessExt creates a new  SetOptionsResultSuccessExt.
func NewSetOptionsResultSuccessExt(v LedgerVersion, value interface{}) (result SetOptionsResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetOptionsResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type SetOptionsResultSuccess struct {
	Ext SetOptionsResultSuccessExt `json:"ext,omitempty"`
}

// SetOptionsResult is an XDR Union defines as:
//
//   union SetOptionsResult switch (SetOptionsResultCode code)
//    {
//    case SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type SetOptionsResult struct {
	Code    SetOptionsResultCode     `json:"code,omitempty"`
	Success *SetOptionsResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsResult
func (u SetOptionsResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetOptionsResultCode(sw) {
	case SetOptionsResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewSetOptionsResult creates a new  SetOptionsResult.
func NewSetOptionsResult(code SetOptionsResultCode, value interface{}) (result SetOptionsResult, err error) {
	result.Code = code
	switch SetOptionsResultCode(code) {
	case SetOptionsResultCodeSuccess:
		tv, ok := value.(SetOptionsResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u SetOptionsResult) MustSuccess() SetOptionsResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SetOptionsResult) GetSuccess() (result SetOptionsResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ReviewCoinsEmissionRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ReviewCoinsEmissionRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewCoinsEmissionRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewCoinsEmissionRequestOpExt
func (u ReviewCoinsEmissionRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewCoinsEmissionRequestOpExt creates a new  ReviewCoinsEmissionRequestOpExt.
func NewReviewCoinsEmissionRequestOpExt(v LedgerVersion, value interface{}) (result ReviewCoinsEmissionRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewCoinsEmissionRequestOp is an XDR Struct defines as:
//
//   struct ReviewCoinsEmissionRequestOp
//    {
//    	CoinsEmissionRequestEntry request;  // request to be reviewed
//    	bool approve;
//    	string64 reason;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ReviewCoinsEmissionRequestOp struct {
	Request CoinsEmissionRequestEntry       `json:"request,omitempty"`
	Approve bool                            `json:"approve,omitempty"`
	Reason  String64                        `json:"reason,omitempty"`
	Ext     ReviewCoinsEmissionRequestOpExt `json:"ext,omitempty"`
}

// ReviewCoinsEmissionRequestResultCode is an XDR Enum defines as:
//
//   enum ReviewCoinsEmissionRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        INVALID_REASON = -1,        // reason must be null if approving
//    	NOT_FOUND = -2,             // failed to find emission request with such ID
//    	NOT_EQUAL = -3,             // stored emission request is not equal to request provided in op
//    	ALREADY_REVIEWED = -4,      // emission request have been already reviewed
//    	MALFORMED = -5,             // emission request is malformed
//        NOT_ENOUGH_PREEMISSIONS = -6,    // serial is already used in another review
//    	LINE_FULL = -9,             // balance will overflow
//        ASSET_NOT_FOUND = -10,
//        BALANCE_NOT_FOUND = -11,
//    	REFERENCE_DUPLICATION = -12
//    };
//
type ReviewCoinsEmissionRequestResultCode int32

const (
	ReviewCoinsEmissionRequestResultCodeSuccess               ReviewCoinsEmissionRequestResultCode = 0
	ReviewCoinsEmissionRequestResultCodeInvalidReason         ReviewCoinsEmissionRequestResultCode = -1
	ReviewCoinsEmissionRequestResultCodeNotFound              ReviewCoinsEmissionRequestResultCode = -2
	ReviewCoinsEmissionRequestResultCodeNotEqual              ReviewCoinsEmissionRequestResultCode = -3
	ReviewCoinsEmissionRequestResultCodeAlreadyReviewed       ReviewCoinsEmissionRequestResultCode = -4
	ReviewCoinsEmissionRequestResultCodeMalformed             ReviewCoinsEmissionRequestResultCode = -5
	ReviewCoinsEmissionRequestResultCodeNotEnoughPreemissions ReviewCoinsEmissionRequestResultCode = -6
	ReviewCoinsEmissionRequestResultCodeLineFull              ReviewCoinsEmissionRequestResultCode = -9
	ReviewCoinsEmissionRequestResultCodeAssetNotFound         ReviewCoinsEmissionRequestResultCode = -10
	ReviewCoinsEmissionRequestResultCodeBalanceNotFound       ReviewCoinsEmissionRequestResultCode = -11
	ReviewCoinsEmissionRequestResultCodeReferenceDuplication  ReviewCoinsEmissionRequestResultCode = -12
)

var ReviewCoinsEmissionRequestResultCodeAll = []ReviewCoinsEmissionRequestResultCode{
	ReviewCoinsEmissionRequestResultCodeSuccess,
	ReviewCoinsEmissionRequestResultCodeInvalidReason,
	ReviewCoinsEmissionRequestResultCodeNotFound,
	ReviewCoinsEmissionRequestResultCodeNotEqual,
	ReviewCoinsEmissionRequestResultCodeAlreadyReviewed,
	ReviewCoinsEmissionRequestResultCodeMalformed,
	ReviewCoinsEmissionRequestResultCodeNotEnoughPreemissions,
	ReviewCoinsEmissionRequestResultCodeLineFull,
	ReviewCoinsEmissionRequestResultCodeAssetNotFound,
	ReviewCoinsEmissionRequestResultCodeBalanceNotFound,
	ReviewCoinsEmissionRequestResultCodeReferenceDuplication,
}

var reviewCoinsEmissionRequestResultCodeMap = map[int32]string{
	0:   "ReviewCoinsEmissionRequestResultCodeSuccess",
	-1:  "ReviewCoinsEmissionRequestResultCodeInvalidReason",
	-2:  "ReviewCoinsEmissionRequestResultCodeNotFound",
	-3:  "ReviewCoinsEmissionRequestResultCodeNotEqual",
	-4:  "ReviewCoinsEmissionRequestResultCodeAlreadyReviewed",
	-5:  "ReviewCoinsEmissionRequestResultCodeMalformed",
	-6:  "ReviewCoinsEmissionRequestResultCodeNotEnoughPreemissions",
	-9:  "ReviewCoinsEmissionRequestResultCodeLineFull",
	-10: "ReviewCoinsEmissionRequestResultCodeAssetNotFound",
	-11: "ReviewCoinsEmissionRequestResultCodeBalanceNotFound",
	-12: "ReviewCoinsEmissionRequestResultCodeReferenceDuplication",
}

var reviewCoinsEmissionRequestResultCodeShortMap = map[int32]string{
	0:   "success",
	-1:  "invalid_reason",
	-2:  "not_found",
	-3:  "not_equal",
	-4:  "already_reviewed",
	-5:  "malformed",
	-6:  "not_enough_preemissions",
	-9:  "line_full",
	-10: "asset_not_found",
	-11: "balance_not_found",
	-12: "reference_duplication",
}

var reviewCoinsEmissionRequestResultCodeRevMap = map[string]int32{
	"ReviewCoinsEmissionRequestResultCodeSuccess":               0,
	"ReviewCoinsEmissionRequestResultCodeInvalidReason":         -1,
	"ReviewCoinsEmissionRequestResultCodeNotFound":              -2,
	"ReviewCoinsEmissionRequestResultCodeNotEqual":              -3,
	"ReviewCoinsEmissionRequestResultCodeAlreadyReviewed":       -4,
	"ReviewCoinsEmissionRequestResultCodeMalformed":             -5,
	"ReviewCoinsEmissionRequestResultCodeNotEnoughPreemissions": -6,
	"ReviewCoinsEmissionRequestResultCodeLineFull":              -9,
	"ReviewCoinsEmissionRequestResultCodeAssetNotFound":         -10,
	"ReviewCoinsEmissionRequestResultCodeBalanceNotFound":       -11,
	"ReviewCoinsEmissionRequestResultCodeReferenceDuplication":  -12,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ReviewCoinsEmissionRequestResultCode
func (e ReviewCoinsEmissionRequestResultCode) ValidEnum(v int32) bool {
	_, ok := reviewCoinsEmissionRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ReviewCoinsEmissionRequestResultCode) String() string {
	name, _ := reviewCoinsEmissionRequestResultCodeMap[int32(e)]
	return name
}

func (e ReviewCoinsEmissionRequestResultCode) ShortString() string {
	name, _ := reviewCoinsEmissionRequestResultCodeShortMap[int32(e)]
	return name
}

func (e ReviewCoinsEmissionRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ReviewCoinsEmissionRequestResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := reviewCoinsEmissionRequestResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ReviewCoinsEmissionRequestResultCode(value)
	return nil
}

// ReviewCoinsEmissionRequestResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ReviewCoinsEmissionRequestResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewCoinsEmissionRequestResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewCoinsEmissionRequestResultSuccessExt
func (u ReviewCoinsEmissionRequestResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewCoinsEmissionRequestResultSuccessExt creates a new  ReviewCoinsEmissionRequestResultSuccessExt.
func NewReviewCoinsEmissionRequestResultSuccessExt(v LedgerVersion, value interface{}) (result ReviewCoinsEmissionRequestResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewCoinsEmissionRequestResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		uint64 requestID;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type ReviewCoinsEmissionRequestResultSuccess struct {
	RequestId Uint64                                     `json:"requestID,omitempty"`
	Ext       ReviewCoinsEmissionRequestResultSuccessExt `json:"ext,omitempty"`
}

// ReviewCoinsEmissionRequestResult is an XDR Union defines as:
//
//   union ReviewCoinsEmissionRequestResult switch (ReviewCoinsEmissionRequestResultCode code)
//    {
//    case SUCCESS:
//    	struct {
//    		uint64 requestID;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type ReviewCoinsEmissionRequestResult struct {
	Code    ReviewCoinsEmissionRequestResultCode     `json:"code,omitempty"`
	Success *ReviewCoinsEmissionRequestResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewCoinsEmissionRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewCoinsEmissionRequestResult
func (u ReviewCoinsEmissionRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ReviewCoinsEmissionRequestResultCode(sw) {
	case ReviewCoinsEmissionRequestResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewReviewCoinsEmissionRequestResult creates a new  ReviewCoinsEmissionRequestResult.
func NewReviewCoinsEmissionRequestResult(code ReviewCoinsEmissionRequestResultCode, value interface{}) (result ReviewCoinsEmissionRequestResult, err error) {
	result.Code = code
	switch ReviewCoinsEmissionRequestResultCode(code) {
	case ReviewCoinsEmissionRequestResultCodeSuccess:
		tv, ok := value.(ReviewCoinsEmissionRequestResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewCoinsEmissionRequestResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ReviewCoinsEmissionRequestResult) MustSuccess() ReviewCoinsEmissionRequestResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ReviewCoinsEmissionRequestResult) GetSuccess() (result ReviewCoinsEmissionRequestResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// StatisticsEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type StatisticsEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u StatisticsEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of StatisticsEntryExt
func (u StatisticsEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewStatisticsEntryExt creates a new  StatisticsEntryExt.
func NewStatisticsEntryExt(v LedgerVersion, value interface{}) (result StatisticsEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// StatisticsEntry is an XDR Struct defines as:
//
//   struct StatisticsEntry
//    {
//    	AccountID accountID;
//
//    	int64 dailyOutcome;
//    	int64 weeklyOutcome;
//    	int64 monthlyOutcome;
//    	int64 annualOutcome;
//
//    	int64 updatedAt;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type StatisticsEntry struct {
	AccountId      AccountId          `json:"accountID,omitempty"`
	DailyOutcome   Int64              `json:"dailyOutcome,omitempty"`
	WeeklyOutcome  Int64              `json:"weeklyOutcome,omitempty"`
	MonthlyOutcome Int64              `json:"monthlyOutcome,omitempty"`
	AnnualOutcome  Int64              `json:"annualOutcome,omitempty"`
	UpdatedAt      Int64              `json:"updatedAt,omitempty"`
	Ext            StatisticsEntryExt `json:"ext,omitempty"`
}

// SetLimitsOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type SetLimitsOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetLimitsOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetLimitsOpExt
func (u SetLimitsOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetLimitsOpExt creates a new  SetLimitsOpExt.
func NewSetLimitsOpExt(v LedgerVersion, value interface{}) (result SetLimitsOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetLimitsOp is an XDR Struct defines as:
//
//   struct SetLimitsOp
//    {
//        AccountID* account;
//        AccountType* accountType;
//
//        Limits limits;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type SetLimitsOp struct {
	Account     *AccountId     `json:"account,omitempty"`
	AccountType *AccountType   `json:"accountType,omitempty"`
	Limits      Limits         `json:"limits,omitempty"`
	Ext         SetLimitsOpExt `json:"ext,omitempty"`
}

// SetLimitsResultCode is an XDR Enum defines as:
//
//   enum SetLimitsResultCode
//    {
//        // codes considered as "success" for the operation
//        SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        MALFORMED = -1
//    };
//
type SetLimitsResultCode int32

const (
	SetLimitsResultCodeSuccess   SetLimitsResultCode = 0
	SetLimitsResultCodeMalformed SetLimitsResultCode = -1
)

var SetLimitsResultCodeAll = []SetLimitsResultCode{
	SetLimitsResultCodeSuccess,
	SetLimitsResultCodeMalformed,
}

var setLimitsResultCodeMap = map[int32]string{
	0:  "SetLimitsResultCodeSuccess",
	-1: "SetLimitsResultCodeMalformed",
}

var setLimitsResultCodeShortMap = map[int32]string{
	0:  "success",
	-1: "malformed",
}

var setLimitsResultCodeRevMap = map[string]int32{
	"SetLimitsResultCodeSuccess":   0,
	"SetLimitsResultCodeMalformed": -1,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetLimitsResultCode
func (e SetLimitsResultCode) ValidEnum(v int32) bool {
	_, ok := setLimitsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetLimitsResultCode) String() string {
	name, _ := setLimitsResultCodeMap[int32(e)]
	return name
}

func (e SetLimitsResultCode) ShortString() string {
	name, _ := setLimitsResultCodeShortMap[int32(e)]
	return name
}

func (e SetLimitsResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *SetLimitsResultCode) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := setLimitsResultCodeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = SetLimitsResultCode(value)
	return nil
}

// SetLimitsResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type SetLimitsResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetLimitsResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetLimitsResultSuccessExt
func (u SetLimitsResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetLimitsResultSuccessExt creates a new  SetLimitsResultSuccessExt.
func NewSetLimitsResultSuccessExt(v LedgerVersion, value interface{}) (result SetLimitsResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetLimitsResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type SetLimitsResultSuccess struct {
	Ext SetLimitsResultSuccessExt `json:"ext,omitempty"`
}

// SetLimitsResult is an XDR Union defines as:
//
//   union SetLimitsResult switch (SetLimitsResultCode code)
//    {
//    case SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type SetLimitsResult struct {
	Code    SetLimitsResultCode     `json:"code,omitempty"`
	Success *SetLimitsResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetLimitsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetLimitsResult
func (u SetLimitsResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetLimitsResultCode(sw) {
	case SetLimitsResultCodeSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewSetLimitsResult creates a new  SetLimitsResult.
func NewSetLimitsResult(code SetLimitsResultCode, value interface{}) (result SetLimitsResult, err error) {
	result.Code = code
	switch SetLimitsResultCode(code) {
	case SetLimitsResultCodeSuccess:
		tv, ok := value.(SetLimitsResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetLimitsResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u SetLimitsResult) MustSuccess() SetLimitsResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SetLimitsResult) GetSuccess() (result SetLimitsResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// AssetPairPolicy is an XDR Enum defines as:
//
//   enum AssetPairPolicy
//    {
//    	TRADEABLE = 1, // if not set pair can not be traided
//    	PHYSICAL_PRICE_RESTRICTION = 2, // if set, then prices for new offers must be greater then physical price with correction
//    	CURRENT_PRICE_RESTRICTION = 4 // if set, then price for new offers must be in interval of (1 +- maxPriceStep)*currentPrice
//    };
//
type AssetPairPolicy int32

const (
	AssetPairPolicyTradeable                AssetPairPolicy = 1
	AssetPairPolicyPhysicalPriceRestriction AssetPairPolicy = 2
	AssetPairPolicyCurrentPriceRestriction  AssetPairPolicy = 4
)

var AssetPairPolicyAll = []AssetPairPolicy{
	AssetPairPolicyTradeable,
	AssetPairPolicyPhysicalPriceRestriction,
	AssetPairPolicyCurrentPriceRestriction,
}

var assetPairPolicyMap = map[int32]string{
	1: "AssetPairPolicyTradeable",
	2: "AssetPairPolicyPhysicalPriceRestriction",
	4: "AssetPairPolicyCurrentPriceRestriction",
}

var assetPairPolicyShortMap = map[int32]string{
	1: "tradeable",
	2: "physical_price_restriction",
	4: "current_price_restriction",
}

var assetPairPolicyRevMap = map[string]int32{
	"AssetPairPolicyTradeable":                1,
	"AssetPairPolicyPhysicalPriceRestriction": 2,
	"AssetPairPolicyCurrentPriceRestriction":  4,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AssetPairPolicy
func (e AssetPairPolicy) ValidEnum(v int32) bool {
	_, ok := assetPairPolicyMap[v]
	return ok
}

// String returns the name of `e`
func (e AssetPairPolicy) String() string {
	name, _ := assetPairPolicyMap[int32(e)]
	return name
}

func (e AssetPairPolicy) ShortString() string {
	name, _ := assetPairPolicyShortMap[int32(e)]
	return name
}

func (e AssetPairPolicy) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *AssetPairPolicy) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := assetPairPolicyRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = AssetPairPolicy(value)
	return nil
}

// AssetPairEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AssetPairEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetPairEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetPairEntryExt
func (u AssetPairEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAssetPairEntryExt creates a new  AssetPairEntryExt.
func NewAssetPairEntryExt(v LedgerVersion, value interface{}) (result AssetPairEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AssetPairEntry is an XDR Struct defines as:
//
//   struct AssetPairEntry
//    {
//        AssetCode base;
//    	AssetCode quote;
//
//        int64 currentPrice;
//        int64 physicalPrice;
//
//    	int64 physicalPriceCorrection; // correction of physical price in percents. If physical price is set and restriction by physical price set, mininal price for offer for this pair will be physicalPrice * physicalPriceCorrection
//    	int64 maxPriceStep; // max price step in percent. User is allowed to set offer with price < (1 - maxPriceStep)*currentPrice and > (1 + maxPriceStep)*currentPrice
//
//
//    	int32 policies;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AssetPairEntry struct {
	Base                    AssetCode         `json:"base,omitempty"`
	Quote                   AssetCode         `json:"quote,omitempty"`
	CurrentPrice            Int64             `json:"currentPrice,omitempty"`
	PhysicalPrice           Int64             `json:"physicalPrice,omitempty"`
	PhysicalPriceCorrection Int64             `json:"physicalPriceCorrection,omitempty"`
	MaxPriceStep            Int64             `json:"maxPriceStep,omitempty"`
	Policies                Int32             `json:"policies,omitempty"`
	Ext                     AssetPairEntryExt `json:"ext,omitempty"`
}

// ThresholdIndexes is an XDR Enum defines as:
//
//   enum ThresholdIndexes
//    {
//        MASTER_WEIGHT = 0,
//        LOW = 1,
//        MED = 2,
//        HIGH = 3
//    };
//
type ThresholdIndexes int32

const (
	ThresholdIndexesMasterWeight ThresholdIndexes = 0
	ThresholdIndexesLow          ThresholdIndexes = 1
	ThresholdIndexesMed          ThresholdIndexes = 2
	ThresholdIndexesHigh         ThresholdIndexes = 3
)

var ThresholdIndexesAll = []ThresholdIndexes{
	ThresholdIndexesMasterWeight,
	ThresholdIndexesLow,
	ThresholdIndexesMed,
	ThresholdIndexesHigh,
}

var thresholdIndexesMap = map[int32]string{
	0: "ThresholdIndexesMasterWeight",
	1: "ThresholdIndexesLow",
	2: "ThresholdIndexesMed",
	3: "ThresholdIndexesHigh",
}

var thresholdIndexesShortMap = map[int32]string{
	0: "master_weight",
	1: "low",
	2: "med",
	3: "high",
}

var thresholdIndexesRevMap = map[string]int32{
	"ThresholdIndexesMasterWeight": 0,
	"ThresholdIndexesLow":          1,
	"ThresholdIndexesMed":          2,
	"ThresholdIndexesHigh":         3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ThresholdIndexes
func (e ThresholdIndexes) ValidEnum(v int32) bool {
	_, ok := thresholdIndexesMap[v]
	return ok
}

// String returns the name of `e`
func (e ThresholdIndexes) String() string {
	name, _ := thresholdIndexesMap[int32(e)]
	return name
}

func (e ThresholdIndexes) ShortString() string {
	name, _ := thresholdIndexesShortMap[int32(e)]
	return name
}

func (e ThresholdIndexes) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *ThresholdIndexes) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := thresholdIndexesRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = ThresholdIndexes(value)
	return nil
}

// LedgerEntryType is an XDR Enum defines as:
//
//   enum LedgerEntryType
//    {
//        ACCOUNT = 0,
//    	COINS_EMISSION_REQUEST = 1,
//        FEE = 2,
//        COINS_EMISSION = 3,
//        BALANCE = 4,
//        PAYMENT_REQUEST = 5,
//        ASSET = 6,
//        REFERENCE_ENTRY = 7,
//        ACCOUNT_TYPE_LIMITS = 8,
//        STATISTICS = 9,
//        TRUST = 10,
//        ACCOUNT_LIMITS = 11,
//    	ASSET_PAIR = 12,
//    	OFFER_ENTRY = 13,
//        INVOICE = 14
//    };
//
type LedgerEntryType int32

const (
	LedgerEntryTypeAccount              LedgerEntryType = 0
	LedgerEntryTypeCoinsEmissionRequest LedgerEntryType = 1
	LedgerEntryTypeFee                  LedgerEntryType = 2
	LedgerEntryTypeCoinsEmission        LedgerEntryType = 3
	LedgerEntryTypeBalance              LedgerEntryType = 4
	LedgerEntryTypePaymentRequest       LedgerEntryType = 5
	LedgerEntryTypeAsset                LedgerEntryType = 6
	LedgerEntryTypeReferenceEntry       LedgerEntryType = 7
	LedgerEntryTypeAccountTypeLimits    LedgerEntryType = 8
	LedgerEntryTypeStatistics           LedgerEntryType = 9
	LedgerEntryTypeTrust                LedgerEntryType = 10
	LedgerEntryTypeAccountLimits        LedgerEntryType = 11
	LedgerEntryTypeAssetPair            LedgerEntryType = 12
	LedgerEntryTypeOfferEntry           LedgerEntryType = 13
	LedgerEntryTypeInvoice              LedgerEntryType = 14
)

var LedgerEntryTypeAll = []LedgerEntryType{
	LedgerEntryTypeAccount,
	LedgerEntryTypeCoinsEmissionRequest,
	LedgerEntryTypeFee,
	LedgerEntryTypeCoinsEmission,
	LedgerEntryTypeBalance,
	LedgerEntryTypePaymentRequest,
	LedgerEntryTypeAsset,
	LedgerEntryTypeReferenceEntry,
	LedgerEntryTypeAccountTypeLimits,
	LedgerEntryTypeStatistics,
	LedgerEntryTypeTrust,
	LedgerEntryTypeAccountLimits,
	LedgerEntryTypeAssetPair,
	LedgerEntryTypeOfferEntry,
	LedgerEntryTypeInvoice,
}

var ledgerEntryTypeMap = map[int32]string{
	0:  "LedgerEntryTypeAccount",
	1:  "LedgerEntryTypeCoinsEmissionRequest",
	2:  "LedgerEntryTypeFee",
	3:  "LedgerEntryTypeCoinsEmission",
	4:  "LedgerEntryTypeBalance",
	5:  "LedgerEntryTypePaymentRequest",
	6:  "LedgerEntryTypeAsset",
	7:  "LedgerEntryTypeReferenceEntry",
	8:  "LedgerEntryTypeAccountTypeLimits",
	9:  "LedgerEntryTypeStatistics",
	10: "LedgerEntryTypeTrust",
	11: "LedgerEntryTypeAccountLimits",
	12: "LedgerEntryTypeAssetPair",
	13: "LedgerEntryTypeOfferEntry",
	14: "LedgerEntryTypeInvoice",
}

var ledgerEntryTypeShortMap = map[int32]string{
	0:  "account",
	1:  "coins_emission_request",
	2:  "fee",
	3:  "coins_emission",
	4:  "balance",
	5:  "payment_request",
	6:  "asset",
	7:  "reference_entry",
	8:  "account_type_limits",
	9:  "statistics",
	10: "trust",
	11: "account_limits",
	12: "asset_pair",
	13: "offer_entry",
	14: "invoice",
}

var ledgerEntryTypeRevMap = map[string]int32{
	"LedgerEntryTypeAccount":              0,
	"LedgerEntryTypeCoinsEmissionRequest": 1,
	"LedgerEntryTypeFee":                  2,
	"LedgerEntryTypeCoinsEmission":        3,
	"LedgerEntryTypeBalance":              4,
	"LedgerEntryTypePaymentRequest":       5,
	"LedgerEntryTypeAsset":                6,
	"LedgerEntryTypeReferenceEntry":       7,
	"LedgerEntryTypeAccountTypeLimits":    8,
	"LedgerEntryTypeStatistics":           9,
	"LedgerEntryTypeTrust":                10,
	"LedgerEntryTypeAccountLimits":        11,
	"LedgerEntryTypeAssetPair":            12,
	"LedgerEntryTypeOfferEntry":           13,
	"LedgerEntryTypeInvoice":              14,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryType
func (e LedgerEntryType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryType) String() string {
	name, _ := ledgerEntryTypeMap[int32(e)]
	return name
}

func (e LedgerEntryType) ShortString() string {
	name, _ := ledgerEntryTypeShortMap[int32(e)]
	return name
}

func (e LedgerEntryType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *LedgerEntryType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := ledgerEntryTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = LedgerEntryType(value)
	return nil
}

// LedgerEntryData is an XDR NestedUnion defines as:
//
//   union switch (LedgerEntryType type)
//        {
//        case ACCOUNT:
//            AccountEntry account;
//    	case COINS_EMISSION_REQUEST:
//    		CoinsEmissionRequestEntry coinsEmissionRequest;
//        case FEE:
//            FeeEntry feeState;
//        case COINS_EMISSION:
//    		CoinsEmissionEntry coinsEmission;
//        case BALANCE:
//            BalanceEntry balance;
//        case PAYMENT_REQUEST:
//            PaymentRequestEntry paymentRequest;
//        case ASSET:
//            AssetEntry asset;
//        case REFERENCE_ENTRY:
//            ReferenceEntry payment;
//        case ACCOUNT_TYPE_LIMITS:
//            AccountTypeLimitsEntry accountTypeLimits;
//        case STATISTICS:
//            StatisticsEntry stats;
//        case TRUST:
//            TrustEntry trust;
//        case ACCOUNT_LIMITS:
//            AccountLimitsEntry accountLimits;
//    	case ASSET_PAIR:
//    		AssetPairEntry assetPair;
//    	case OFFER_ENTRY:
//    		OfferEntry offer;
//        case INVOICE:
//            InvoiceEntry invoice;
//        }
//
type LedgerEntryData struct {
	Type                 LedgerEntryType            `json:"type,omitempty"`
	Account              *AccountEntry              `json:"account,omitempty"`
	CoinsEmissionRequest *CoinsEmissionRequestEntry `json:"coinsEmissionRequest,omitempty"`
	FeeState             *FeeEntry                  `json:"feeState,omitempty"`
	CoinsEmission        *CoinsEmissionEntry        `json:"coinsEmission,omitempty"`
	Balance              *BalanceEntry              `json:"balance,omitempty"`
	PaymentRequest       *PaymentRequestEntry       `json:"paymentRequest,omitempty"`
	Asset                *AssetEntry                `json:"asset,omitempty"`
	Payment              *ReferenceEntry            `json:"payment,omitempty"`
	AccountTypeLimits    *AccountTypeLimitsEntry    `json:"accountTypeLimits,omitempty"`
	Stats                *StatisticsEntry           `json:"stats,omitempty"`
	Trust                *TrustEntry                `json:"trust,omitempty"`
	AccountLimits        *AccountLimitsEntry        `json:"accountLimits,omitempty"`
	AssetPair            *AssetPairEntry            `json:"assetPair,omitempty"`
	Offer                *OfferEntry                `json:"offer,omitempty"`
	Invoice              *InvoiceEntry              `json:"invoice,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryData) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryData
func (u LedgerEntryData) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeCoinsEmissionRequest:
		return "CoinsEmissionRequest", true
	case LedgerEntryTypeFee:
		return "FeeState", true
	case LedgerEntryTypeCoinsEmission:
		return "CoinsEmission", true
	case LedgerEntryTypeBalance:
		return "Balance", true
	case LedgerEntryTypePaymentRequest:
		return "PaymentRequest", true
	case LedgerEntryTypeAsset:
		return "Asset", true
	case LedgerEntryTypeReferenceEntry:
		return "Payment", true
	case LedgerEntryTypeAccountTypeLimits:
		return "AccountTypeLimits", true
	case LedgerEntryTypeStatistics:
		return "Stats", true
	case LedgerEntryTypeTrust:
		return "Trust", true
	case LedgerEntryTypeAccountLimits:
		return "AccountLimits", true
	case LedgerEntryTypeAssetPair:
		return "AssetPair", true
	case LedgerEntryTypeOfferEntry:
		return "Offer", true
	case LedgerEntryTypeInvoice:
		return "Invoice", true
	}
	return "-", false
}

// NewLedgerEntryData creates a new  LedgerEntryData.
func NewLedgerEntryData(aType LedgerEntryType, value interface{}) (result LedgerEntryData, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(AccountEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountEntry")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeCoinsEmissionRequest:
		tv, ok := value.(CoinsEmissionRequestEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be CoinsEmissionRequestEntry")
			return
		}
		result.CoinsEmissionRequest = &tv
	case LedgerEntryTypeFee:
		tv, ok := value.(FeeEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be FeeEntry")
			return
		}
		result.FeeState = &tv
	case LedgerEntryTypeCoinsEmission:
		tv, ok := value.(CoinsEmissionEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be CoinsEmissionEntry")
			return
		}
		result.CoinsEmission = &tv
	case LedgerEntryTypeBalance:
		tv, ok := value.(BalanceEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be BalanceEntry")
			return
		}
		result.Balance = &tv
	case LedgerEntryTypePaymentRequest:
		tv, ok := value.(PaymentRequestEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentRequestEntry")
			return
		}
		result.PaymentRequest = &tv
	case LedgerEntryTypeAsset:
		tv, ok := value.(AssetEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AssetEntry")
			return
		}
		result.Asset = &tv
	case LedgerEntryTypeReferenceEntry:
		tv, ok := value.(ReferenceEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReferenceEntry")
			return
		}
		result.Payment = &tv
	case LedgerEntryTypeAccountTypeLimits:
		tv, ok := value.(AccountTypeLimitsEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountTypeLimitsEntry")
			return
		}
		result.AccountTypeLimits = &tv
	case LedgerEntryTypeStatistics:
		tv, ok := value.(StatisticsEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be StatisticsEntry")
			return
		}
		result.Stats = &tv
	case LedgerEntryTypeTrust:
		tv, ok := value.(TrustEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be TrustEntry")
			return
		}
		result.Trust = &tv
	case LedgerEntryTypeAccountLimits:
		tv, ok := value.(AccountLimitsEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountLimitsEntry")
			return
		}
		result.AccountLimits = &tv
	case LedgerEntryTypeAssetPair:
		tv, ok := value.(AssetPairEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AssetPairEntry")
			return
		}
		result.AssetPair = &tv
	case LedgerEntryTypeOfferEntry:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	case LedgerEntryTypeInvoice:
		tv, ok := value.(InvoiceEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be InvoiceEntry")
			return
		}
		result.Invoice = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccount() AccountEntry {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccount() (result AccountEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustCoinsEmissionRequest() CoinsEmissionRequestEntry {
	val, ok := u.GetCoinsEmissionRequest()

	if !ok {
		panic("arm CoinsEmissionRequest is not set")
	}

	return val
}

// GetCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetCoinsEmissionRequest() (result CoinsEmissionRequestEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmissionRequest" {
		result = *u.CoinsEmissionRequest
		ok = true
	}

	return
}

// MustFeeState retrieves the FeeState value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustFeeState() FeeEntry {
	val, ok := u.GetFeeState()

	if !ok {
		panic("arm FeeState is not set")
	}

	return val
}

// GetFeeState retrieves the FeeState value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetFeeState() (result FeeEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeState" {
		result = *u.FeeState
		ok = true
	}

	return
}

// MustCoinsEmission retrieves the CoinsEmission value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustCoinsEmission() CoinsEmissionEntry {
	val, ok := u.GetCoinsEmission()

	if !ok {
		panic("arm CoinsEmission is not set")
	}

	return val
}

// GetCoinsEmission retrieves the CoinsEmission value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetCoinsEmission() (result CoinsEmissionEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmission" {
		result = *u.CoinsEmission
		ok = true
	}

	return
}

// MustBalance retrieves the Balance value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustBalance() BalanceEntry {
	val, ok := u.GetBalance()

	if !ok {
		panic("arm Balance is not set")
	}

	return val
}

// GetBalance retrieves the Balance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetBalance() (result BalanceEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Balance" {
		result = *u.Balance
		ok = true
	}

	return
}

// MustPaymentRequest retrieves the PaymentRequest value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustPaymentRequest() PaymentRequestEntry {
	val, ok := u.GetPaymentRequest()

	if !ok {
		panic("arm PaymentRequest is not set")
	}

	return val
}

// GetPaymentRequest retrieves the PaymentRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetPaymentRequest() (result PaymentRequestEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentRequest" {
		result = *u.PaymentRequest
		ok = true
	}

	return
}

// MustAsset retrieves the Asset value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAsset() AssetEntry {
	val, ok := u.GetAsset()

	if !ok {
		panic("arm Asset is not set")
	}

	return val
}

// GetAsset retrieves the Asset value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAsset() (result AssetEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Asset" {
		result = *u.Asset
		ok = true
	}

	return
}

// MustPayment retrieves the Payment value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustPayment() ReferenceEntry {
	val, ok := u.GetPayment()

	if !ok {
		panic("arm Payment is not set")
	}

	return val
}

// GetPayment retrieves the Payment value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetPayment() (result ReferenceEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Payment" {
		result = *u.Payment
		ok = true
	}

	return
}

// MustAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccountTypeLimits() AccountTypeLimitsEntry {
	val, ok := u.GetAccountTypeLimits()

	if !ok {
		panic("arm AccountTypeLimits is not set")
	}

	return val
}

// GetAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccountTypeLimits() (result AccountTypeLimitsEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountTypeLimits" {
		result = *u.AccountTypeLimits
		ok = true
	}

	return
}

// MustStats retrieves the Stats value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustStats() StatisticsEntry {
	val, ok := u.GetStats()

	if !ok {
		panic("arm Stats is not set")
	}

	return val
}

// GetStats retrieves the Stats value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetStats() (result StatisticsEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Stats" {
		result = *u.Stats
		ok = true
	}

	return
}

// MustTrust retrieves the Trust value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustTrust() TrustEntry {
	val, ok := u.GetTrust()

	if !ok {
		panic("arm Trust is not set")
	}

	return val
}

// GetTrust retrieves the Trust value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetTrust() (result TrustEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Trust" {
		result = *u.Trust
		ok = true
	}

	return
}

// MustAccountLimits retrieves the AccountLimits value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccountLimits() AccountLimitsEntry {
	val, ok := u.GetAccountLimits()

	if !ok {
		panic("arm AccountLimits is not set")
	}

	return val
}

// GetAccountLimits retrieves the AccountLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccountLimits() (result AccountLimitsEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountLimits" {
		result = *u.AccountLimits
		ok = true
	}

	return
}

// MustAssetPair retrieves the AssetPair value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAssetPair() AssetPairEntry {
	val, ok := u.GetAssetPair()

	if !ok {
		panic("arm AssetPair is not set")
	}

	return val
}

// GetAssetPair retrieves the AssetPair value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAssetPair() (result AssetPairEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetPair" {
		result = *u.AssetPair
		ok = true
	}

	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustOffer() OfferEntry {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetOffer() (result OfferEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// MustInvoice retrieves the Invoice value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustInvoice() InvoiceEntry {
	val, ok := u.GetInvoice()

	if !ok {
		panic("arm Invoice is not set")
	}

	return val
}

// GetInvoice retrieves the Invoice value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetInvoice() (result InvoiceEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Invoice" {
		result = *u.Invoice
		ok = true
	}

	return
}

// LedgerEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LedgerEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryExt
func (u LedgerEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerEntryExt creates a new  LedgerEntryExt.
func NewLedgerEntryExt(v LedgerVersion, value interface{}) (result LedgerEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerEntry is an XDR Struct defines as:
//
//   struct LedgerEntry
//    {
//        uint32 lastModifiedLedgerSeq; // ledger the LedgerEntry was last changed
//
//        union switch (LedgerEntryType type)
//        {
//        case ACCOUNT:
//            AccountEntry account;
//    	case COINS_EMISSION_REQUEST:
//    		CoinsEmissionRequestEntry coinsEmissionRequest;
//        case FEE:
//            FeeEntry feeState;
//        case COINS_EMISSION:
//    		CoinsEmissionEntry coinsEmission;
//        case BALANCE:
//            BalanceEntry balance;
//        case PAYMENT_REQUEST:
//            PaymentRequestEntry paymentRequest;
//        case ASSET:
//            AssetEntry asset;
//        case REFERENCE_ENTRY:
//            ReferenceEntry payment;
//        case ACCOUNT_TYPE_LIMITS:
//            AccountTypeLimitsEntry accountTypeLimits;
//        case STATISTICS:
//            StatisticsEntry stats;
//        case TRUST:
//            TrustEntry trust;
//        case ACCOUNT_LIMITS:
//            AccountLimitsEntry accountLimits;
//    	case ASSET_PAIR:
//    		AssetPairEntry assetPair;
//    	case OFFER_ENTRY:
//    		OfferEntry offer;
//        case INVOICE:
//            InvoiceEntry invoice;
//        }
//        data;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type LedgerEntry struct {
	LastModifiedLedgerSeq Uint32          `json:"lastModifiedLedgerSeq,omitempty"`
	Data                  LedgerEntryData `json:"data,omitempty"`
	Ext                   LedgerEntryExt  `json:"ext,omitempty"`
}

// EnvelopeType is an XDR Enum defines as:
//
//   enum EnvelopeType
//    {
//        SCP = 1,
//        TX = 2,
//        AUTH = 3
//    };
//
type EnvelopeType int32

const (
	EnvelopeTypeScp  EnvelopeType = 1
	EnvelopeTypeTx   EnvelopeType = 2
	EnvelopeTypeAuth EnvelopeType = 3
)

var EnvelopeTypeAll = []EnvelopeType{
	EnvelopeTypeScp,
	EnvelopeTypeTx,
	EnvelopeTypeAuth,
}

var envelopeTypeMap = map[int32]string{
	1: "EnvelopeTypeScp",
	2: "EnvelopeTypeTx",
	3: "EnvelopeTypeAuth",
}

var envelopeTypeShortMap = map[int32]string{
	1: "scp",
	2: "tx",
	3: "auth",
}

var envelopeTypeRevMap = map[string]int32{
	"EnvelopeTypeScp":  1,
	"EnvelopeTypeTx":   2,
	"EnvelopeTypeAuth": 3,
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EnvelopeType
func (e EnvelopeType) ValidEnum(v int32) bool {
	_, ok := envelopeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e EnvelopeType) String() string {
	name, _ := envelopeTypeMap[int32(e)]
	return name
}

func (e EnvelopeType) ShortString() string {
	name, _ := envelopeTypeShortMap[int32(e)]
	return name
}

func (e EnvelopeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

func (e *EnvelopeType) UnmarshalJSON(d []byte) error {
	var raw string
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	value, ok := envelopeTypeRevMap[raw]
	if !ok {
		return fmt.Errorf("unexpected json value: %s", raw)
	}

	*e = EnvelopeType(value)
	return nil
}

var fmtTest = fmt.Sprint("this is a dummy usage of fmt")
