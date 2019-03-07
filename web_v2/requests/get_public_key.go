package requests

import (
	"net/http"
)

// GetPublicKey - represents params to be specified by user for GetPublicKey handler
type GetPublicKey struct {
	*base
	ID string
}

// NewGetPublicKey - returns new instance of GetPublicKeyAccountList request
func NewGetPublicKey(r *http.Request) (*GetPublicKey, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id := b.getString("id")

	return &GetPublicKey{
		base: b,
		ID:   id,
	}, nil
}
