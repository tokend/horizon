package janus

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/kit/janus/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var ErrAlreadyExists = errors.New("API already registered with different surname")

type Upstream struct {
	Target  string
	Surname string
}

type Janus struct {
	disabled bool
	upstream Upstream
	client   internal.Client
}

func NewNoOp() *Janus {
	return &Janus{
		disabled: true,
	}
}

func New(endpoint string, upstream Upstream) *Janus {
	return &Janus{
		upstream: upstream,
		client: internal.Client{
			endpoint,
		},
	}
}

// RegisterChi takes router and registers all endpoints in janus
func (j *Janus) RegisterChi(r chi.Router) error {
	if j.disabled {
		return nil
	}

	// walk the router without hitting janus
	services, err := internal.NewChi(r).Services()
	if err != nil {
		return errors.Wrap(err, "failed to walk chi router")
	}

	for _, candidate := range services {
		// check if service already exists
		remote, err := j.client.GetAPI(candidate.Name)
		if err != nil {
			return errors.Wrap(err, "failed to get remote service")
		}
		if remote != nil {
			if remote.Surname != j.upstream.Surname {
				return ErrAlreadyExists
			}

			// modify remote service

			// TODO check if upstream is duplicate

			remote.Proxy.Upstreams.Targets = append(
				remote.Proxy.Upstreams.Targets,
				internal.Target{Target: j.upstream.Target})

			if err := j.client.UpdateAPI(remote.Name, remote); err != nil {
				return errors.Wrap(err, "failed to update remote service")
			}
			continue
		}

		// add new service definition
		candidate.Surname = j.upstream.Surname
		candidate.Proxy.Upstreams = internal.Upstreams{
			Balancing: "weight",
			Targets:   []internal.Target{{Target: j.upstream.Target, Weight: 10}},
		}
		if err := j.client.AddAPI(&candidate); err != nil {
			return errors.Wrap(err, "failed to register service")
		}
	}

	return nil
}
