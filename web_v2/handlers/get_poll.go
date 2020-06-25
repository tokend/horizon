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

func GetPoll(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetPoll(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getPollHandler{
		VotesQ: history2.NewVotesQ(ctx.HistoryRepo(r)),
		PollsQ: history2.NewPollsQ(ctx.HistoryRepo(r)),
		Log:    ctx.Log(r),
	}

	result, err := handler.getPoll(request)
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

	if request.ShouldIncludeAny(
		requests.IncludeTypePollParticipation,
		requests.IncludeTypePollParticipationVotes,
	) {
		if !isAllowed(r, w, &result.Data.Relationships.Owner.Data.ID, &result.Data.Relationships.ResultProvider.Data.ID) {
			return
		}
	}

	ape.Render(w, result)
}

type getPollHandler struct {
	PollsQ history2.PollsQ
	VotesQ history2.VotesQ
	Log    *logan.Entry
}

// GetSale returns sale with related resources
func (h *getPollHandler) getPoll(request *requests.GetPoll) (*regources.PollResponse, error) {

	record, err := h.PollsQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get poll")
	}

	if record == nil {
		return nil, nil
	}

	resource := resources.NewPoll(*record)
	response := &regources.PollResponse{
		Data: resource,
	}

	if !request.ShouldIncludeAny(
		requests.IncludeTypePollParticipation,
		requests.IncludeTypePollParticipationVotes,
	) {
		return response, nil
	}

	votes, err := h.VotesQ.FilterByPollID(request.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get votes for poll")
	}

	if request.ShouldInclude(requests.IncludeTypePollParticipation) {
		outcomeKey := resources.NewParticipation(record.ID, votes)
		response.Included.Add(&outcomeKey)
	}

	if request.ShouldInclude(requests.IncludeTypePollParticipationVotes) {
		for _, v := range votes {
			vote := resources.NewVote(v)
			response.Included.Add(&vote)
		}
	}

	return response, nil
}
