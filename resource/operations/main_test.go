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
			Effects:     nil, //TODO ask about it
		}
	}
	return participants
}

//This test check only JSON output, he didn't check state of operation
func TestNew(t *testing.T) {
	ctx := context.TODO()
	participants := getParticipants()

	operation := getOperation(xdr.OperationTypeCreateAccount, `{"funder": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636", "account": "GBCADUBIWUSO5HZF7YZUC42Z6DWSCJH2GQQUEJH6HD6OE4TSYU46F7XT", "account_type": 2}`)

	result, err := New(ctx, operation, participants, false)
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
	}`
	assert.JSONEq(t, expected, string(marshalRes))
}
