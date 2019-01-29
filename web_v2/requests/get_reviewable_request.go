package requests

import "net/http"

// GetReviewableRequest - represents params to be specified by user for Get Request handler
type GetReviewableRequest struct {
	*base
	ID uint64
}

//NewGetReviewableRequest - returns new instance of GetRequest request
func NewGetReviewableRequest(r *http.Request) (*GetReviewableRequest, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("id")
	if err != nil {
		return nil, err
	}

	return &GetReviewableRequest{
		base: b,
		ID:   id,
	}, nil
}
