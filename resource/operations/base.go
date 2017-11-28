package operations

import (
	"fmt"
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/httpx"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource/base"
	"golang.org/x/net/context"
)

// PagingToken implements hal.Pageable
func (this Base) PagingToken() string {
	return this.PT
}

// Populate fills out this resource using `row` as the source.
func (this *Base) Populate(
	ctx context.Context, row history.Operation, participants []*history.Participant, public bool,
) error {
	this.ID = fmt.Sprintf("%d", row.ID)
	this.PT = row.PagingToken()
	this.TransactionID = fmt.Sprintf("%d", row.TransactionID)
	this.SourceAccount = row.SourceAccount
	this.populateType(row)
	this.LedgerCloseTime = row.LedgerCloseTime
	this.Participants = make([]base.Participant, len(participants))
	this.State = int32(row.State)
	this.Identifier = strconv.FormatInt(row.Identifier, 10)
	for i := range participants {
		err := this.Participants[i].Populate(participants[i], row.Type, public)
		if err != nil {
			return err
		}
	}

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	self := fmt.Sprintf("/operations/%d", row.ID)
	this.Links.Self = lb.Link(self)
	this.Links.Succeeds = lb.Linkf(nil, "/effects?order=desc&cursor=%s", this.PT)
	this.Links.Precedes = lb.Linkf(nil, "/effects?order=asc&cursor=%s", this.PT)
	this.Links.Transaction = lb.Linkf(nil, "/transactions/%s", row.TransactionHash)

	if public {
		this.SourceAccount = ""
	}
	return nil
}

func (this *Base) populateType(row history.Operation) {
	var ok bool
	this.TypeI = int32(row.Type)
	this.Type, ok = TypeNames[row.Type]

	if !ok {
		this.Type = "unknown"
	}
}
