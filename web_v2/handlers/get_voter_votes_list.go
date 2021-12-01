package handlers

import (
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/generator"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

func GetVoterVotesList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)

	handler := getVoteListHandler{
		VotesQ:        history2.NewVotesQ(historyRepo),
		PollsQ:        history2.NewPollsQ(historyRepo),
		AccountsQ:     core2.NewAccountsQ(coreRepo),
		LedgerHeaderQ: *core2.NewLedgerHeaderQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetVotersVotes(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, &request.VoterID) {
		return
	}

	result, err := handler.GetVoterVotesList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get votes list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

func (h *getVoteListHandler) GetVoterVotesList(request *requests.GetVoterVoteList) (*regources.VoteListResponse, error) {
	q := h.VotesQ.
		FilterByVoterID(request.VoterID).
		Page(*request.PageParams)

	historyVotes, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "cannot select votes")
	}

	response := regources.VoteListResponse{
		Data: make([]regources.Vote, 0, len(historyVotes)),
	}
	for _, vote := range historyVotes {
		ledgerSequence := generator.GetSeqFromInt64(vote.ID) // ledger sequence
		header, err := h.LedgerHeaderQ.GetBySequence(ledgerSequence)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get header of ledger sequence")
		}
		created := uint64(header.CloseTime)
		vote.VoteData.CreationTime = &created

		response.Data = append(response.Data, resources.NewVote(vote))
		if request.ShouldInclude(requests.IncludeTypeVoterVoteListPolls) {
			historyPoll, err := h.PollsQ.GetByID(vote.PollID)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get voter account")
			}

			poll := resources.NewPoll(*historyPoll)
			response.Included.Add(&poll)
		}
	}

	if len(response.Data) > 0 {
		if request.ShouldInclude(requests.IncludeTypeVotersVoteListAccount) {
			coreAccount, err := h.AccountsQ.GetByAddress(request.VoterID)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get voter account")
			}

			account := resources.NewAccount(*coreAccount, nil)
			response.Included.Add(&account)
		}

		response.Links = request.GetCursorLinks(*request.PageParams, fmt.Sprintf("%d", historyVotes[len(historyVotes)-1].ID))
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}
