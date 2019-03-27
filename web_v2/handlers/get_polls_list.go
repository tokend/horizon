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

// GetPollList - processes request to get the list of sales
func GetPollList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getPollListHandler{
		PollsQ: history2.NewPollsQ(historyRepo),
		VotesQ: history2.NewVotesQ(historyRepo),
		Log:    ctx.Log(r),
	}

	request, err := requests.NewGetPollList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetPollList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get poll list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getPollListHandler struct {
	PollsQ history2.PollsQ
	VotesQ history2.VotesQ
	Log    *logan.Entry
}

// GetPollList returns the list of assets with related resources
func (h *getPollListHandler) GetPollList(request *requests.GetPollList) (*regources.PollsResponse, error) {
	q := h.PollsQ

	if request.ShouldFilter(requests.FilterTypePollListOwner) {
		q = q.FilterByOwner(request.Filters.Owner)
	}

	if request.ShouldFilter(requests.FilterTypePollListResultProvider) {
		q = q.FilterByResultProvider(request.Filters.ResultProvider)
	}

	if request.ShouldFilter(requests.FilterTypePollListPermissionType) {
		q = q.FilterByPermissionType(request.Filters.PermissionType)
	}

	if request.ShouldFilter(requests.FilterTypePollListVoteConfirmation) {
		q = q.FilterByVoteConfirmation(request.Filters.VoteConfirmation)
	}

	if request.ShouldFilter(requests.FilterTypePollListMaxEndTime) {
		q = q.FilterByMaxEndTime(*request.Filters.MaxEndTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListMaxStartTime) {
		q = q.FilterByMaxStartTime(*request.Filters.MaxStartTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListMinStartTime) {
		q = q.FilterByMinStartTime(*request.Filters.MinStartTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListMinEndTime) {
		q = q.FilterByMinEndTime(*request.Filters.MinEndTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListState) {
		q = q.FilterByState(request.Filters.State)
	}

	if request.ShouldFilter(requests.FilterTypePollListPollType) {
		q = q.FilterByPollType(request.Filters.PollType)
	}

	historyPolls, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get poll list ")
	}

	response := &regources.PollsResponse{
		Data:  make([]regources.Poll, 0, len(historyPolls)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, historyPoll := range historyPolls {
		votes, err := h.VotesQ.FilterByPollID(historyPoll.ID).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get votes for poll", logan.F{
				"poll_id": historyPoll.ID,
			})
		}
		poll := resources.NewPoll(historyPoll, votes)

		response.Data = append(response.Data, poll)
	}

	return response, nil
}
