package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
	"net/http"
)

// GetMatchList - processes request to get the list of matches
func GetMatchList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getMatchListHandler{
		MatchQ: history2.NewSquashedMatchesQ(historyRepo),
	}

	request, err := requests.NewGetMatchList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetMatchList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get match list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getMatchListHandler struct {
	MatchQ history2.MatchQ
	Log    *logan.Entry
}

// GetMatchList returns list of matches with related resources
func (h *getMatchListHandler) GetMatchList(request *requests.GetMatchList) (*regources.MatchsResponse, error) {
	q := h.MatchQ

	if request.ShouldFilter(requests.FilterTypeMatchListBaseAsset) {
		q = q.FilterByBaseAsset(request.Filters.BaseAsset)
	}
	if request.ShouldFilter(requests.FilterTypeMatchListQuoteAsset) {
		q = q.FilterByQuoteAsset(request.Filters.QuoteAsset)
	}
	if request.ShouldFilter(requests.FilterTypeMatchListOrderBook) {
		q = q.FilterByOrderBookID(request.Filters.OrderBook)
	}

	coreMatches, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get matches from db")
	}

	response := regources.MatchsResponse{
		Data: make([]regources.Match, 0, len(coreMatches)),
		// TODO: includes
	}

	for _, coreMatch := range coreMatches {
		response.Data = append(response.Data, resources.NewMatch(coreMatch))
	}

	return &response, nil
}
