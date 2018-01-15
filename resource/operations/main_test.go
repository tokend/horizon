package operations

import (
	"context"
	"testing"
	"time"

	"encoding/json"

	"fmt"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
)

func getOperation(opType xdr.OperationType, details string) history.Operation {
	ts, err := time.Parse(time.RFC3339, "2018-01-11T13:51:15Z")
	if err != nil {
		panic(err)
	}
	return history.Operation{
		TotalOrderID:     db2.TotalOrderID{ID: 231928242177},
		TransactionID:    231928242176,
		TransactionHash:  "73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd",
		ApplicationOrder: 1,
		Type:             opType,
		DetailsString:    null.NewString(details, true),
		LedgerCloseTime:  ts,
		SourceAccount:    "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
		State:            2,
		Identifier:       4,
	}
}

func getParticipants() []*history.Participant {
	participants := make([]*history.Participant, 3)
	for i := range participants {
		participants[i] = &history.Participant{
			OperationID: int64(1234567 + i),
			AccountID:   "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB63" + fmt.Sprintf("%v", i),
			BalanceID:   "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB63" + fmt.Sprintf("%v", i),
			Nickname:    "Nickname" + fmt.Sprintf("%v", i),
			Email:       "email" + fmt.Sprintf("%v", i) + "@test.com",
			Mobile:      "38012123456" + fmt.Sprintf("%v", i),
			Details:     []byte{1, 2, 3, 4, 5, byte(i)},
			UserType:    "general",
			Effects:     nil,
		}
	}
	return participants
}

//This test check only JSON output, he didn't check state of operation
func TestNew(t *testing.T) {
	ctx := context.TODO()
	participants := getParticipants()

	cases := []struct {
		name      string
		operation history.Operation
		expected  string
	}{
		{
			name: "CreateAccount",
			operation: getOperation(xdr.OperationTypeCreateAccount, `{
				"funder": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
				"account": "GBCADUBIWUSO5HZF7YZUC42Z6DWSCJH2GQQUEJH6HD6OE4TSYU46F7XT",
				"account_type": 2
			}`),
			expected: `{
				 "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "create_account",
			  "type_i": 0,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "funder": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "account": "GBCADUBIWUSO5HZF7YZUC42Z6DWSCJH2GQQUEJH6HD6OE4TSYU46F7XT",
			  "account_type": 2
			}`,
		},
		{
			name: "Payment",
			operation: getOperation(xdr.OperationTypePayment, `{
				"to": "GANVIVPOJ2Q7DTIYJJJSP5X64BZYFBGPQO4EXMBEOY6LT5CRJZ6PGC27",
				"from": "GA2ZQVZKQJUF3B3KSNXGAWVV2PEFBD4KCDRSCSFWD2CCVSGZ35S6K4P5",
				"asset": "SUN",
				"amount": "2.0691",
				"subject": "Test Staging\nPlease work well",
				"reference": "",
				"to_balance": "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
				"from_balance": "BDF6UAXEOJLKTDRBCEUJRGVNNSLGZGRRGCXQYOIZ4F25AED57OGFEZIX",
				"source_fixed_fee": "0.0000",
				"source_payment_fee": "0.0000",
				"source_pays_for_dest": true,
				"destination_fixed_fee": "0.0000",
				"destination_payment_fee": "0.0000"
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "payment",
			  "type_i": 1,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "from": "GA2ZQVZKQJUF3B3KSNXGAWVV2PEFBD4KCDRSCSFWD2CCVSGZ35S6K4P5",
			  "to": "GANVIVPOJ2Q7DTIYJJJSP5X64BZYFBGPQO4EXMBEOY6LT5CRJZ6PGC27",
			  "from_balance": "BDF6UAXEOJLKTDRBCEUJRGVNNSLGZGRRGCXQYOIZ4F25AED57OGFEZIX",
			  "to_balance": "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
			  "amount": "2.0691",
			  "asset": "SUN",
			  "source_payment_fee": "0.0000",
			  "destination_payment_fee": "0.0000",
			  "source_fixed_fee": "0.0000",
			  "destination_fixed_fee": "0.0000",
			  "source_pays_for_dest": true,
			  "subject": "Test Staging\nPlease work well",
			  "qasset": ""
			}`,
		},
		{
			name: "SetOptions",
			operation: getOperation(xdr.OperationTypeSetOptions, `{
				  "home_domain": "test.com",
				  "inflation_dest": "0.0000",
				  "master_key_weight": 1,
				  "signer_key": "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
				  "signer_weight": 1,
				  "signer_type": 1,
				  "signer_identity": 1,
				  "set_flags": [
					1,
					1,
					1
				  ],
				  "set_flags_s": [
					"test0",
					"test1",
					"test2"
				  ],
				  "clear_flags": [
					1,
					1,
					1
				  ],
				  "clear_flags_s": [
					"test0",
					"test1",
					"test2"
				  ],
				  "low_threshold": 1,
				  "med_threshold": 1,
				  "high_threshold": 1
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "set_options",
			  "type_i": 2,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "home_domain": "test.com",
			  "inflation_dest": "0.0000",
			  "master_key_weight": 1,
			  "signer_key": "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
			  "signer_weight": 1,
			  "signer_type": 1,
			  "signer_identity": 1,
			  "set_flags": [
				1,
				1,
				1
			  ],
			  "set_flags_s": [
				"test0",
				"test1",
				"test2"
			  ],
			  "clear_flags": [
				1,
				1,
				1
			  ],
			  "clear_flags_s": [
				"test0",
				"test1",
				"test2"
			  ],
			  "low_threshold": 1,
			  "med_threshold": 1,
			  "high_threshold": 1
			}`,
		},
		{
			name: "SetFees",
			operation: getOperation(xdr.OperationTypeSetFees, `{"fee":{
				"asset_code": "SUN",
				"fixed_fee": "0.0000",
				"percent_fee": "0.0000",
				"fee_type": 1,
				"account_id": "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
				"account_type": 1,
				"subtype": 		1,
				"lower_bound": 123,
				"upper_bound": 12345
			}}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "set_fees",
			  "type_i": 5,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "fee": {
				"asset_code": "SUN",
				"fixed_fee": "0.0000",
				"percent_fee": "0.0000",
				"fee_type": 1,
				"account_id": "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
				"account_type": 1,
				"subtype": 1,
				"lower_bound": 123,
				"upper_bound": 12345
			  }
			}`,
		},
		{
			name: "ManageAccount",
			operation: getOperation(xdr.OperationTypeManageAccount, `{
				"account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
				"block_reasons_to_add": 1,
  				"BlockReasonsToRemove": 1
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "manage_account",
			  "type_i": 6,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
			  "block_reasons_to_add": 1
			}`,
		},
		{
			name: "ManageAccount2",
			operation: getOperation(xdr.OperationTypeManageAccount, `{
				"account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
				"block_reasons_to_remove": 1
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "manage_account",
			  "type_i": 6,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
			  "block_reasons_to_remove": 1
			}`,
		},
		{
			name: "CreateWithdrawalRequest",
			operation: getOperation(xdr.OperationTypeCreateWithdrawalRequest, `{
				"amount": "1000.00",
				"balance": "BANTYPGNNSC64NSULLOBI2MOEUHQXJTNPUIFMCM4N7JXRX5",
				"fee_fixed": "0.0000",
				"fee_percent": "0.0000",
				"external_details": {
					"a": "some external details"
				},
				"dest_asset": "SUN",
				"dest_amount": "1000.00"
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "create_withdrawal_request",
			  "type_i": 7,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "amount": "1000.00",
			  "balance": "BANTYPGNNSC64NSULLOBI2MOEUHQXJTNPUIFMCM4N7JXRX5",
			  "fee_fixed": "0.0000",
			  "fee_percent": "0.0000",
			  "external_details": {
				"a": "some external details"
			  },
			  "dest_asset": "SUN",
			  "dest_amount": "1000.00"
			}`,
		},
		{
			name: "ManageInvoice",
			operation: getOperation(xdr.OperationTypeManageInvoice, `{
				"amount": "1000.00",
				"receiver_balance": "GBPBGYUANKZJWTFREEKMHEGSXZDFZJ6KSEMTHTR3AK3XSB2W3Y2FOL2B",
				"sender": "GAAB7JPFE4MSSF6Y7JIPKFK5KNITTOJM7VS5OZWQKK3KETTBU74JEOFW",
				"invoice_id": 1,
				"asset": "SUN"
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "manage_invoice",
			  "type_i": 17,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "amount": "1000.00",
			  "receiver_balance": "GBPBGYUANKZJWTFREEKMHEGSXZDFZJ6KSEMTHTR3AK3XSB2W3Y2FOL2B",
			  "sender": "GAAB7JPFE4MSSF6Y7JIPKFK5KNITTOJM7VS5OZWQKK3KETTBU74JEOFW",
			  "invoice_id": 1,
			  "asset": "SUN"
			}`,
		},
		{
			name: "ManageOffer",
			operation: getOperation(xdr.OperationTypeManageOffer, `{
			  "fee": "0.0000",
			  "price": "2.5500",
			  "amount": "7.8431",
			  "is_buy": true,
			  "offer_id": 1,
			  "is_deleted": false
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "manage_offer",
			  "type_i": 16,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "is_buy": true,
			  "amount": "7.8431",
			  "price": "2.5500",
			  "fee": "0.0000",
			  "offer_id": 1,
			  "is_deleted": false
			}`,
		},
		{
			name: "ManageAssetPair",
			operation: getOperation(xdr.OperationTypeManageAssetPair, `{
			  "base_asset": "ETH",
			  "quote_asset": "SUN",
			  "max_price_step": "0.0000",
			  "physical_price": "1326.0000",
			  "physical_price_correction": "0.0000"
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "manage_asset_pair",
			  "type_i": 15,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "base_asset": "ETH",
			  "quote_asset": "SUN",
			  "physical_price": "1326.0000",
			  "physical_price_correction": "0.0000",
			  "max_price_step": "0.0000"
			}`,
		},
		{
			name: "CreateIssuanceRequest",
			operation: getOperation(xdr.OperationTypeCreateIssuanceRequest, `{
			  "reference": "GAA6HKHWQWKWOPSQBIBLJZWWRSTHI7PBFINXYIURHZB7SVDU7AXSQU6K",
			  "amount": "1000.00",
			  "asset": "SUN",
			  "fee_fixed": "0.0000",
			  "fee_percent": "0.0000",
			  "external_details": {
				"a": "some external details"
			  }
			}`),
			expected: `{
			  "_links": {
				"self": {
				  "href": "/operations/231928242177"
				},
				"transaction": {
				  "href": "/transactions/73559b4bda9057acc6566da0e3f0e2a7eab6f7742df9ffe86a3a5cef6ef081cd"
				},
				"succeeds": {
				  "href": "/effects?order=desc&cursor=231928242177"
				},
				"precedes": {
				  "href": "/effects?order=asc&cursor=231928242177"
				}
			  },
			  "id": "231928242177",
			  "paging_token": "231928242177",
			  "transaction_id": "231928242176",
			  "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
			  "type": "create_issuance_request",
			  "type_i": 3,
			  "state_i": 2,
			  "state": "success",
			  "identifier": "4",
			  "ledger_close_time": "2018-01-11T13:51:15Z",
			  "participants": [
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB630",
				  "email": "email0@test.com",
				  "nickname": "Nickname0"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB631",
				  "email": "email1@test.com",
				  "nickname": "Nickname1"
				},
				{
				  "account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "balance_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB632",
				  "email": "email2@test.com",
				  "nickname": "Nickname2"
				}
			  ],
			  "reference": "GAA6HKHWQWKWOPSQBIBLJZWWRSTHI7PBFINXYIURHZB7SVDU7AXSQU6K",
			  "amount": "1000.00",
			  "asset": "SUN",
			  "fee_fixed": "0.0000",
			  "fee_percent": "0.0000",
			  "external_details": {
				"a": "some external details"
			  }
			}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := New(ctx, c.operation, participants, false)
			if err != nil {
				t.Fatal(err)
			}
			marshalRes, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, c.expected, string(marshalRes))
		})
	}
}
