package resource

import "gitlab.com/distributed_lab/tokend/horizon/db2/history"

type PriceHistory struct {
	Prices []history.PricePoint `json:"prices"`
}
