package handlers

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"

	"gitlab.com/tokend/horizon/web_v2/resources"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

type getOperationsHandler struct {
	OperationQ history2.OperationQ
	Log        *logan.Entry
}

func GetOperations(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getOperationsHandler{
		OperationQ: history2.NewOperationQ(historyRepo),
		Log:        ctx.Log(r),
	}

	request, err := requests.NewGetOperations(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	result, err := handler.GetOperations(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get operations list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

func (h *getOperationsHandler) GetOperations(request *requests.GetOperations) (*regources.OperationListResponse, error) {
	q := h.OperationQ

	if request.Filters.Types != nil {
		q = q.FilterByOperationsTypes(request.Filters.Types)
	}

	if request.Filters.Source != nil {
		q = q.FilterByOperationSource(*request.Filters.Source)
	}

	opsAll, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load all operations")
	}

	if request.PageNumber != nil {
		q = q.PageOffset(pgdb.OffsetPageParams{
			Limit:      request.PageParams.Limit,
			Order:      request.PageParams.Order,
			PageNumber: *request.PageNumber,
		})
	} else {
		q = q.Page(request.PageParams)
	}

	historyOperations, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load operations")
	}

	result := regources.OperationListResponse{
		Data: make([]regources.Operation, 0, len(historyOperations)),
	}

	for _, historyOperation := range historyOperations {
		var operation regources.Operation
		operation = resources.NewOperation(historyOperation)

		if request.ShouldInclude(requests.IncludeTypeOperationsListOperationDetails) {
			result.Included.Add(resources.NewOperationDetails(historyOperation))
		}
		result.Data = append(result.Data, operation)
	}

	if request.PageNumber != nil {
		result.Links = request.GetOffsetLinks(pgdb.OffsetPageParams{
			Limit:      request.PageParams.Limit,
			Order:      request.PageParams.Order,
			PageNumber: *request.PageNumber,
		})

		err = result.PutMeta(requests.MetaPageParams{
			CurrentPage: *request.PageNumber,
			TotalPages:  (uint64(len(opsAll)) + request.PageParams.Limit - 1) / request.PageParams.Limit,
		})
	} else {
		if len(result.Data) > 0 {
			result.Links = request.GetCursorLinks(request.PageParams, result.Data[len(result.Data)-1].ID)
		} else {
			result.Links = request.GetCursorLinks(request.PageParams, "")
		}

		err = result.PutMeta(requests.MetaCursorParams{
			CurrentCursor: request.PageParams.Cursor,
			TotalPages:    (uint64(len(opsAll)) + request.PageParams.Limit - 1) / request.PageParams.Limit,
		})
	}

	if err != nil {
		return &result, errors.Wrap(err, "failed to put meta to response")
	}

	return &result, nil
}
