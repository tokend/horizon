package janus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Janus struct {
	URL     string
	Target  string
	Surname string
}

//TODO add retry
//TODO add service "surname" to definition and try to teach janus to use different service for one path with different method
//FIXME if at the end of the listen_path it is written "/" janus thinks it is different api
//TODO goji-style names and paths?
//TODO add logger?
func NewJanus(url, target, surname string) *Janus {
	return &Janus{
		URL:     url,
		Target:  target,
		Surname: surname,
	}
}

// DoRegister takes router and registers all endpoints in janus
func (j *Janus) DoRegister(r chi.Router, log *logan.Entry) error {
	walk := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		safeRoute := strings.Replace(route, "/*/", "/", -1)
		log.Info(fmt.Sprintf("janus creating service with endpoint: %s", safeRoute))

		err := j.NewAPI(safeRoute, method, j.Target, j.Surname)
		if err != nil {
			return errors.Wrap(err, "failed to add service")
		}
		return nil
	}
	err := chi.Walk(r, walk)
	if err != nil {
		return errors.Wrap(err, "walk return error")
	}
	return nil
}

// NewAPI register new service in janus
// if service already exist - updates it
func (j *Janus) NewAPI(endpoint, method, target, surname string) error {
	// check if service already exist
	janus, err := j.GetAPI(GetName(endpoint, method))
	if err != nil {
		return errors.Wrap(err, "failed to get service")
	}
	if janus != nil {
		if janus.Surname != surname {
			return errors.New("service already exists with another surname")
		}
		err := j.ModifyAPI(target, *janus)
		if err != nil {
			return errors.Wrap(err, "failed to modify api")
		}
		return nil
	}

	service := NewService(target, endpoint, method, surname)
	jsonStr, err := json.Marshal(service)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}
	return j.addAPI(jsonStr)
}

// NewService create service definition with given method, endpoint and target
// create URL-friendly name with given endpoint and sets roundrobin as default balancing
func NewService(target, endpoint, method, serviceSurname string) *Service {
	methods := []string{method}
	targets := []Target{{Target: target, Weight: 10}}
	return &Service{
		Name:    GetName(endpoint, method),
		Surname: serviceSurname,
		Active:  true,
		Proxy: Proxy{
			AppendPath: true,
			ListenPath: endpoint,
			Upstreams: Upstreams{
				Balancing: "weight",
				Targets:   targets},
			Methods: methods},
	}
}

// GetName receive endpoint and create URL-friendly name for service
// All parameters will be replaced with "x" in order to avoid adding the same paths differing only in the name of the parameters
func GetName(endpoint, method string) string {
	//TODO check how root works
	if len(endpoint) == 1 {
		return fmt.Sprintf("root-%s", strings.ToLower(method))
	}
	t := endpoint[1:]
	r := regexp.MustCompile(`{([a-z\s-]+)}`)
	t = r.ReplaceAllString(t, "x")
	t = strings.Replace(t, "/*/", "/", -1)
	t = strings.Replace(t, "/", "-", -1)
	t = strings.Replace(t, "_", "-", -1)
	methodName := strings.ToLower(method)
	if t[len(t)-1:] == "-" {
		return fmt.Sprintf("%s%s", t, methodName)
	}
	return fmt.Sprintf("%s-%s", t, methodName)
}

// GetAPI returns service by name
// return nil,nil if service with that name was not found
// expect prepared (URL-friendly) name of service (using the dash)
func (j *Janus) GetAPI(name string) (*Service, error) {
	return j.getAPI(name)
}

// ModifyAPI check that there are no such target in the service already and updates the service
func (j *Janus) ModifyAPI(target string, service Service) error {
	newTarget := true
	for _, s := range service.Proxy.Upstreams.Targets {
		if s.Target == target {
			newTarget = false
		}
	}
	if !newTarget {
		return nil
	}

	service.Proxy.Upstreams.Targets = append(service.Proxy.Upstreams.Targets, Target{Target: target})
	jsonStr, err := json.Marshal(service)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}

	return j.modifyAPI(service.Name, jsonStr)
}
