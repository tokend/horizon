package handlers

import (
	"net/http"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
)

func GetVote(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetVote(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	handler := getVoteHandler{
		VotesQ: history2.NewVotesQ(ctx.HistoryRepo(r)),
		PollsQ: history2.NewPollsQ(ctx.HistoryRepo(r)),
		Log:    ctx.Log(r),
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

	if !isAllowed(r, w, &request.VoterID, &poll.ResultProviderID, &poll.OwnerID) {
		return
	}

	result, err := handler.getVote(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get vote", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getVoteHandler struct {
	VotesQ history2.VotesQ
	PollsQ history2.PollsQ
	Log    *logan.Entry
}

// GetSale returns sale with related resources
func (h *getVoteHandler) getVote(request *requests.GetVote) (*regources.VoteResponse, error) {

	record, err := h.VotesQ.
		FilterByPollID(request.PollID).
		FilterByVoterID(request.VoterID).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get vote")
	}

	if record == nil {
		return nil, nil
	}

	resource := resources.NewVote(*record)
	response := &regources.VoteResponse{
		Data: resource,
	}

	return response, nil
}
