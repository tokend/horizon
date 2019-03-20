package problems

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/signcontrol"
)

// NotAllowed will try to guess details of error and populate problem accordingly.
func NotAllowed(errs ...error) *jsonapi.ErrorObject {
	// errs is optional for backward compatibility
	if len(errs) == 0 {
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusUnauthorized),
			Status: fmt.Sprintf("%d", http.StatusUnauthorized),
		}
	}

	if len(errs) != 1 {
		panic(errors.New("unexpected number of errors passed"))
	}

	switch cause := errors.Cause(errs[0]); cause.(type) {
	case *signcontrol.Error:
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Detail: "Request signature was invalid in some way",
			Meta: &map[string]interface{}{
				"reason": cause.Error(),
			},
		}
	default:
		// default sink for *doorman.Error and everything else,
		// since 401 is our safe default
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusUnauthorized),
			Status: fmt.Sprintf("%d", http.StatusUnauthorized),
		}
	}

}
