package requests

import "net/http"

const (
	// IncludeTypePollVotes - defines if votes should be included in the response
	IncludeTypePollVotes = "votes"
)

var includeTypePollAll = map[string]struct{}{
	IncludeTypePollVotes: {},
}

type GetPoll struct {
	*base
	ID int64
}

func NewGetPoll(r *http.Request) (*GetPoll, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypePollAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}
	return &GetPoll{
		base: b,
		ID:   int64(id),
	}, nil
}
