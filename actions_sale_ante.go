package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type SaleAnteAction struct {
	Action
	q                    core.SaleAnteQI
	SaleID               string
	ParticipantBalanceID string
	ParticipantAccountID string
	Records              []core.SaleAnte
	Page                 hal.Page
}

func (action *SaleAnteAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadParticipantAccountID,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *SaleAnteAction) loadParams() {
	action.SaleID = action.GetString("sale_id")
	action.ParticipantBalanceID = action.GetString("participant_balance_id")
	action.Page.Filters = map[string]string{
		"sale_id":                action.SaleID,
		"participant_balance_id": action.ParticipantBalanceID,
	}
}

func (action *SaleAnteAction) checkAllowed() {
	action.IsAllowed(action.ParticipantAccountID)
}

func (action *SaleAnteAction) loadRecords() {
	action.q = action.CoreQ().SaleAntes()

	if action.SaleID != "" {
		action.q = action.q.ForSale(action.SaleID)
	}

	if action.ParticipantBalanceID != "" {
		action.q = action.q.ForBalance(action.ParticipantBalanceID)
	}

	if action.Err != nil {
		return
	}

	var err error

	action.Records, err = action.q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get sale antes from core DB")
		action.Err = &problem.ServerError
		return
	}
}

func (action *SaleAnteAction) loadPage() {
	for i := range action.Records {
		var res resource.SaleAnte
		res.Populate(&action.Records[i])
		action.Page.Add(&res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.PopulateLinks()
}

func (action *SaleAnteAction) loadParticipantAccountID() {
	if action.ParticipantBalanceID == "" {
		return
	}

	participantBalance, err := action.CoreQ().Balances().ByID(action.ParticipantBalanceID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get sale ante participant balance from core DB")
		action.Err = &problem.ServerError
		return
	}

	if participantBalance == nil {
		action.Err = &problem.BadRequest
		return
	}

	action.ParticipantAccountID = participantBalance.AccountID
}
