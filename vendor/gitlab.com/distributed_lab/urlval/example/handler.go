package example

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/urlval"
)

type ListLotsRequest struct {
	Search          *string `url:"search"`
	ModelFilter     *string `filter:"model"`
	MakerFilter     *string `filter:"maker"`
	PlatformInclude bool    `include:"platform"`
	// pagination parameters are specified explicitly and not provided as
	// ready to use structures for couple of reasons:
	// * specification is agnostic to pagination strategy and only SHOULDs
	//   us to use `page` field
	// * allows us to be spec compliant and fail if customer does not support
	//   `size` or any other subset of "common" parameters
	// * allows to use custom types for extended semantics, like
	//   `page[number]=(random|last|whatever)`
	PageNumber uint64 `page:"number"`
	PageSize   uint64 `page:"size"`
	// responsibility for handling both page[order] and sort params is on your service,
	// because different implementations can require different behaviours. (maybe it's temporary)
	PageOrder string      `page:"order"`
	Sort      urlval.Sort `sort:"a,b,-c"`
}

func NewListLotsRequest(r *http.Request) (*ListLotsRequest, error) {
	var request ListLotsRequest
	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return nil, errors.Wrap(err, "failed decode query")
	}
	return &request, nil // TODO: validation is up to you
}

func ListLots(w http.ResponseWriter, r *http.Request) {
	request, err := NewListLotsRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.MakerFilter != nil {
		// filter by maker
	}

	if request.PlatformInclude {
		// include platform
	}

	// render links
	_, err = urlval.Encode(request) // self
	request.PageNumber += 1
	_, err = urlval.Encode(request) // next
}
