package history

import (
	"testing"

	"encoding/json"

	"fmt"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/go/xdr"
)

func TestOperation_Details(t *testing.T) {
	//todo add CreateAccount Marshal test
	cases := []struct {
		name            string
		details         OperationDetails
		expectedDetails string
	}{
		{
			name: "Payment",
			details: OperationDetails{
				Type: xdr.OperationTypePayment,
				Payment: &PaymentDetails{
					BasePayment: BasePayment{
						From:                  "GA2ZQVZKQJUF3B3KSNXGAWVV2PEFBD4KCDRSCSFWD2CCVSGZ35S6K4P5",
						To:                    "GANVIVPOJ2Q7DTIYJJJSP5X64BZYFBGPQO4EXMBEOY6LT5CRJZ6PGC27",
						FromBalance:           "BDF6UAXEOJLKTDRBCEUJRGVNNSLGZGRRGCXQYOIZ4F25AED57OGFEZIX",
						ToBalance:             "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
						Amount:                "2.0691",
						Asset:                 "SUN",
						SourcePaymentFee:      "0.0000",
						DestinationPaymentFee: "0.0000",
						SourceFixedFee:        "0.0000",
						DestinationFixedFee:   "0.0000",
						SourcePaysForDest:     true,
					},
					Subject:    "Test Staging\nPlease work well",
					Reference:  "",
					QuoteAsset: "SUN",
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 1,
				"string": "payment"
			  },
			  "payment": {
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
				"qasset": "SUN"
			  }
			}`,
		},
		{
			name: "SetOptions",
			details: OperationDetails{
				Type: xdr.OperationTypeSetOptions,
				SetOptions: &SetOptionsDetails{
					HomeDomain:                      "test.com",
					InflationDest:                   "0.0000",
					MasterKeyWeight:                 1,
					SignerKey:                       "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
					SignerWeight:                    1,
					SignerType:                      1,
					SignerIdentity:                  1,
					SetFlags:                        []int{1, 1, 1},
					SetFlagsS:                       []string{"test0", "test1", "test2"},
					ClearFlags:                      []int{1, 1, 1},
					ClearFlagsS:                     []string{"test0", "test1", "test2"},
					LowThreshold:                    1,
					MedThreshold:                    1,
					HighThreshold:                   1,
					LimitsUpdateRequestDocumentHash: "07997422f6829dbd8f625520133e9c93afc67d673a00c9a23a2de51cb1848271",
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 2,
				"string": "set_options"
			  },
			  "set_options": {
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
				"high_threshold": 1,
				"limits_update_request_document_hash": "07997422f6829dbd8f625520133e9c93afc67d673a00c9a23a2de51cb1848271"
			  }
			}`,
		},
		{
			name: "SetFees",
			details: OperationDetails{
				Type: xdr.OperationTypeSetFees,
				SetFees: &SetFeesDetails{
					Fee: &FeeDetails{
						AssetCode:   "SUN",
						FixedFee:    "0.0000",
						PercentFee:  "0.0000",
						FeeType:     1,
						AccountID:   "BA2UC6DJILEGPIHAPQFAVPGGGA7BF5PDJLB6WXHIOYO3RJZ3QIPRTEN7",
						AccountType: 1,
						Subtype:     1,
						LowerBound:  123,
						UpperBound:  12345,
					},
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 5,
				"string": "set_fees"
			  },
			  "set_fees": {
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
			  }
			}`,
		},
		{
			name: "ManageAccount",
			details: OperationDetails{
				Type: xdr.OperationTypeManageAccount,
				ManageAccount: &ManageAccountDetails{
					Account:           "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
					BlockReasonsToAdd: 1,
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 6,
				"string": "manage_account"
			  },
			  "manage_account": {
				"account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
				"block_reasons_to_add": 1
			  }
			}`,
		},
		{
			name: "ManageAccount2",
			details: OperationDetails{
				Type: xdr.OperationTypeManageAccount,
				ManageAccount: &ManageAccountDetails{
					Account:              "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
					BlockReasonsToRemove: 1,
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 6,
				"string": "manage_account"
			  },
			  "manage_account": {
				"account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB",
				"block_reasons_to_remove": 1
			  }
			}`,
		},
		{
			name: "CreateWithdrawalRequest",
			details: OperationDetails{
				Type: xdr.OperationTypeCreateWithdrawalRequest,
				CreateWithdrawalRequest: &CreateWithdrawalRequestDetails{
					Amount:     "1000.00",
					Balance:    "BANTYPGNNSC64NSULLOBI2MOEUHQXJTNPUIFMCM4N7JXRX5",
					FeeFixed:   "0.0000",
					FeePercent: "0.0000",
					ExternalDetails: map[string]interface{}{
						"a": "some external details",
					},
					DestAsset:  "SUN",
					DestAmount: "1000.00",
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 7,
				"string": "create_withdrawal_request"
			  },
			  "create_withdrawal_request": {
				"amount": "1000.00",
				"balance": "BANTYPGNNSC64NSULLOBI2MOEUHQXJTNPUIFMCM4N7JXRX5",
				"fee_fixed": "0.0000",
				"fee_percent": "0.0000",
				"external_details": {
				  "a": "some external details"
				},
				"dest_asset": "SUN",
				"dest_amount": "1000.00"
			  }
			}`,
		},
		{
			name: "ManageBalance",
			details: OperationDetails{
				Type: xdr.OperationTypeManageBalance,
				ManageBalance: &ManageBalanceDetails{
					Destination: "BANTYPGNNSC64NSULLOBI2MOEUHQXJTNPUIFMCM4N7JXRX5",
					Action:      123,
				},
			},
			expectedDetails: `{
			  "type": {
				"int": 9,
				"string": "manage_balance"
			  },
			  "manage_balance": {
				"destination": "BANTYPGNNSC64NSULLOBI2MOEUHQXJTNPUIFMCM4N7JXRX5",
				"action": 123
			  }
			}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jsonDetails, err := json.Marshal(c.details)
			if err != nil {
				panic(err)
			}

			fmt.Println(c.name, string(jsonDetails))

			assert.JSONEq(t, c.expectedDetails, string(jsonDetails))

		})
	}
}
