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

func GetVote(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetVote(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	handler := getVoteHandler{
		VotesQ: history2.NewVotesQ(ctx.HistoryRepo(r)),
		Log:    ctx.Log(r),
	}

	result, err := handler.getVote(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get poll", logan.F{
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
	Log    *logan.Entry
}

// GetSale returns sale with related resources
func (h *getVoteHandler) getVote(request *requests.GetVote) (*regources.VoteResponse, error) {

	record, err := h.VotesQ.
		FilterByPollID(request.PollID).
		FilterByVoter(request.VoterID).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get poll")
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
