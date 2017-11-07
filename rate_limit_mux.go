package horizon

import (
	"net/http"

	"math"
	"strconv"

	"bullioncoin.githost.io/development/api/log"
	"bullioncoin.githost.io/development/go/signcontrol"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"github.com/zenazn/goji/web"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
	"gopkg.in/throttled/throttled.v2/store/redigostore"
)

type RateLimitedMux struct {
	*web.Mux

	limiter throttled.RateLimiter
}

func NewRateLimitedMux(app *App) (*RateLimitedMux, error) {
	var err error
	var store throttled.GCRAStore

	if app.redis != nil {
		store, err = redigostore.New(app.redis, "ratelimit:", 0)
	} else {
		store, err = memstore.New(2 << 10)
	}
	if err != nil {
		log.WithField("service", "rate-limiter").WithError(err).Fatal("failed to init rate limiter")
	}

	quota := throttled.RateQuota{throttled.PerMin(levelCritical * 5), levelCritical * 2}

	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return nil, err
	}

	return &RateLimitedMux{
		Mux:     web.New(),
		limiter: rateLimiter,
	}, nil
}

// kinda middleware but vsrato-parametrized style
func (m *RateLimitedMux) rateLimit(c web.C, w http.ResponseWriter, r *http.Request, limits []int, handler web.Handler) {
	if m.limiter == nil || r.Header.Get(IsAdminHeader) == IsAdminHeaderValue {
		handler.ServeHTTPC(c, w, r)
		return
	}

	limit := 0
	for _, v := range limits {
		limit += v
	}

	if limit == 0 {
		handler.ServeHTTPC(c, w, r)
		return
	}

	key := r.Header.Get(signcontrol.PublicKeyHeader)

	limited, context, err := m.limiter.RateLimit(key, limit)

	if err != nil {
		log.WithField("service", "rate-limiter").WithError(err).Error("failed to rate limit")
		handler.ServeHTTPC(c, w, r)
		return
	}

	if v := context.Limit; v >= 0 {
		w.Header().Add("X-RateLimit-Limit", strconv.Itoa(v))
	}

	if v := context.Remaining; v >= 0 {
		w.Header().Add("X-RateLimit-Remaining", strconv.Itoa(v))
	}

	if v := context.ResetAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("X-RateLimit-Reset", strconv.Itoa(vi))
	}

	if v := context.RetryAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		w.Header().Add("Retry-After", strconv.Itoa(vi))
	}

	if !limited {
		handler.ServeHTTPC(c, w, r)
	} else {
		problem.Render(nil, w, &problem.RateLimitExceeded)
		return
	}

}

func (m *RateLimitedMux) Get(pattern web.PatternType, handler web.Handler, limits ...int) {
	m.Mux.Get(pattern, func(c web.C, w http.ResponseWriter, r *http.Request) {
		m.rateLimit(c, w, r, limits, handler)
	})
}

func (m *RateLimitedMux) Post(pattern web.PatternType, handler web.Handler, limits ...int) {
	m.Mux.Post(pattern, func(c web.C, w http.ResponseWriter, r *http.Request) {
		m.rateLimit(c, w, r, limits, handler)
	})
}
