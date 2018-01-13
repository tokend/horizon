package operations

import (
	"context"
	"testing"
	"time"

	"encoding/json"

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
		TransactionHash:  "hashy",
		ApplicationOrder: 1,
		Type:             opType,
		DetailsString:    null.NewString(details, true),
		LedgerCloseTime:  ts,
		SourceAccount:    "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
		State:            2,
		Identifier:       4,
	}
}

//This test check only JSON output, he didn't check state of operation
func TestNewCreateAccountOperation(t *testing.T) {
	ctx := context.TODO()
	// TODO helper for participants
	participant := []*history.Participant{}

	operation := getOperation(xdr.OperationTypeCreateAccount, `{"funder": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636", "account": "GBCADUBIWUSO5HZF7YZUC42Z6DWSCJH2GQQUEJH6HD6OE4TSYU46F7XT", "account_type": 2}`)

	result, err := New(ctx, operation, participant, false)
	if err != nil {
		t.Fatal(err)
	}

	marshalRes, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{
		  "_links": {
			"self": {
			  "href": "/operations/231928242177"
			},
			"transaction": {
			  "href": "/transactions/hashy"
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
		  "funder": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
		  "account": "GBCADUBIWUSO5HZF7YZUC42Z6DWSCJH2GQQUEJH6HD6OE4TSYU46F7XT",
		  "account_type": 2
	}`
	assert.JSONEq(t, expected, string(marshalRes))
}
