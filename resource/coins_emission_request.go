package resource

import (
	"encoding/json"
	"strconv"
	"time"

	"gitlab.com/tokend/horizon/db2/history"
)

type CoinsEmissionRequest struct {
	PT           string        `json:"paging_token"`
	Receiver     string        `json:"receiver"`
	Issuer       string        `json:"issuer"`
	Reference    string        `json:"reference"`
	Amount       string        `json:"amount"`
	Asset        string        `json:"asset"`
	Approved     *bool         `json:"approved"`
	PreEmissions []preEmission `json:"pre_emissions"`
	Reason       *string       `json:"reason"`
	CreatedAt    *time.Time    `json:"created_at"`
	UpdatedAt    *time.Time    `json:"updated_at"`
	ExchangeName string        `json:"exchange_name"`
}

type preEmission struct {
	SerialNumber string `json:"serialNumber"`
	Amount       string `json:"amount"`
	Asset        string `json:"amount"`
}

// Populate fills out the resource's fields
func (request *CoinsEmissionRequest) Populate(row *history.CoinsEmissionRequest) error {
	request.PT = strconv.FormatInt(row.RequestID, 10)
	request.Receiver = row.Receiver
	request.Issuer = row.Issuer
	request.Amount = row.Amount
	request.Asset = row.Asset
	request.Reference = row.Reference
	request.Approved = row.Approved
	var preEmissions []preEmission
	if row.PreEmissions != nil {
		json.Unmarshal([]byte(*row.PreEmissions), &preEmissions)
		request.PreEmissions = preEmissions
	}
	request.Reason = row.Reason
	request.CreatedAt = row.CreatedAt
	request.UpdatedAt = row.UpdatedAt

	request.ExchangeName = row.ExchangeName
	return nil
}

// PagingToken implementation for hal.Pageable
func (request CoinsEmissionRequest) PagingToken() string {
	return request.PT
}
