package requests

import (
	"net/http"
)

type GetBalance struct {
	*base
	ID string
}

func NewGetBalance(r *http.Request) (*GetBalance, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := newBalanceID(b, "id")
	if err != nil {
		return nil, err
	}

	return &GetBalance{
		base: b,
		ID:   id,
	}, nil
}
