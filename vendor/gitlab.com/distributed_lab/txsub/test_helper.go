package txsub

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// StaticMockServer is a test helper that records it's last request
type StaticMockServer struct {
	*httptest.Server
	LastRequest *http.Request
}

// NewStaticMockServer creates a new mock server that always responds with
// `response`
func NewStaticMockServer(response string) *StaticMockServer {
	result := &StaticMockServer{}
	result.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result.LastRequest = r
		fmt.Fprintln(w, response)
	}))

	return result
}
