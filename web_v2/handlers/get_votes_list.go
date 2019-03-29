package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
)

// GetVoteList - processes request to get the list of sales
func GetVoteList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getVoteListHandler{
		VotesQ: history2.NewVotesQ(historyRepo),
		PollsQ: history2.NewPollsQ(historyRepo),
		Log:    ctx.Log(r),
	}

	request, err := requests.NewGetVoteList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	poll, err := handler.PollsQ.GetByID(request.PollID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get poll for vote", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if poll == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, poll.OwnerID, poll.ResultProviderID) {
		return
	}

	result, err := handler.GetVoteList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get vote list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getVoteListHandler struct {
	VotesQ history2.VotesQ
	PollsQ history2.PollsQ
	Log    *logan.Entry
}

// GetVoteList returns the list of assets with related resources
func (h *getVoteListHandler) GetVoteList(request *requests.GetVoteList) (*regources.VotesResponse, error) {
	q := h.VotesQ.FilterByPollID(request.PollID)

	historyVotes, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get vote list")
	}

	response := &regources.VotesResponse{
		Data: make([]regources.Vote, 0, len(historyVotes)),
	}

	for _, historyVote := range historyVotes {
		vote := resources.NewVote(historyVote)

		response.Data = append(response.Data, vote)
	}
	h.PopulateLinks(response, request)
	return response, nil
}

func (h *getVoteListHandler) PopulateLinks(
	response *regources.VotesResponse, request *requests.GetVoteList,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}
