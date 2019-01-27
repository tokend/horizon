package internal

import (
	"fmt"
	"regexp"
	"strings"
)

// GetName create URL-friendly name for service
// All parameters will be replaced with "x" in order to avoid adding the same paths differing only in the name of the parameters
func GetName(endpoint, method string) string {
	//TODO check how root works
	if len(endpoint) == 1 {
		return fmt.Sprintf("root-%s", strings.ToLower(method)) // FIXME why?
	}
	t := endpoint[1:]
	r := regexp.MustCompile(`{([a-z\s-]+)}`)
	t = r.ReplaceAllString(t, "x")
	t = strings.Replace(t, "/*/", "/", -1) // FIXME why?
	t = strings.Replace(t, "/", "-", -1)
	t = strings.Replace(t, "_", "-", -1)
	methodName := strings.ToLower(method)
	if t[len(t)-1:] == "-" {
		return fmt.Sprintf("%s%s", methodName, t) // FIXME wat
	}
	return fmt.Sprintf("%s-%s", methodName, t)
}
