package resource

import "gitlab.com/tokend/horizon/db2/history"

type PriceHistory struct {
	Prices []history.PricePoint `json:"prices"`
}
