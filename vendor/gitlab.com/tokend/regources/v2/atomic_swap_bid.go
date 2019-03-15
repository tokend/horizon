package regources

import (
	"time"
)

// AtomicSwapBidResponse - response for AtomicSwapBid handler
type AtomicSwapBidResponse struct {
	Data     AtomicSwapBid `json:"data"`
	Included Included      `json:"included"`
}

// AtomicSwapBidsResponse - response for AtomicSwapBid handler
type AtomicSwapBidsResponse struct {
	Links    *Links          `json:"links"`
	Data     []AtomicSwapBid `json:"data"`
	Included Included        `json:"included"`
}

// AtomicSwapBid - Resource object representing AtomicSwapBidEntry
type AtomicSwapBid struct {
	Key
	Attributes    AtomicSwapBidAttrs     `json:"attributes"`
	Relationships AtomicSwapBidRelations `json:"relationships"`
}

type AtomicSwapBidAttrs struct {
	AvailableAmount Amount    `json:"available_amount"`
	LockedAmount    Amount    `json:"locked_amount"`
	CreatedAt       time.Time `json:"created_at"`
	IsCanceled      bool      `json:"is_canceled"`
	Details         Details   `json:"details"`
}

type AtomicSwapBidRelations struct {
	BaseBalance *Relation           `json:"base_balance"`
	Owner       *Relation           `json:"owner"`
	QuoteAssets *RelationCollection `json:"quote_assets"`
}
