package urlval

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval/internal/betterreflect"
)

func Encode(src interface{}) string {
	refsrc := betterreflect.NewStruct(src)

	values := url.Values{}
	var includes []string

	for i := 0; i < refsrc.NumField(); i++ {
		tag := refsrc.Tag(i, "include")
		if tag != "" {
			if refsrc.Value(i).Bool() {
				includes = append(includes, tag)
			}
		}

		tag = refsrc.Tag(i, "page")
		if tag != "" {
			values.Set(fmt.Sprintf("page[%s]", tag), cast.ToString(refsrc.Value(i).Uint()))
		}
	}

	values.Set("include", strings.Join(includes, ","))
	return values.Encode()
}
