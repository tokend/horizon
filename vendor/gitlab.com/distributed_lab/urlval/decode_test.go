package urlval

import (
	"net/url"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"gitlab.com/distributed_lab/urlval/internal/betterreflect"
)

func mustValues(query string) url.Values {
	values, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}

	return values
}

func requireDecodePanics(t *testing.T, query string, dest interface{}) {
	require.Panics(t, func() {
		_ = Decode(mustValues(query), dest)
	})
}

func TestDecodeSupported(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		var request struct {
			A *string `filter:"a"`
		}
		err := DecodeSupported(mustValues("filter[a]=foo&filter[b]=bar"), &request)
		require.NoError(t, err)
		require.NotNil(t, request.A)
		require.Equal(t, "foo", *request.A)
	})
}

func TestDecodeFilters(t *testing.T) {
	t.Run("decodes to pointers to built-in types", func(t *testing.T) {
		t.Run("*string", func(t *testing.T) {
			var request struct {
				A *string `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=foo"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.Equal(t, "foo", *request.A)
		})
		t.Run("*int", func(t *testing.T) {
			var request struct {
				A *int `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=12"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, 12, *request.A)
		})
		t.Run("*int32", func(t *testing.T) {
			var request struct {
				A *int32 `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=12"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, 12, *request.A)
		})
		t.Run("*int64", func(t *testing.T) {
			var request struct {
				A *int64 `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=12"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, 12, *request.A)
		})
		t.Run("*bool", func(t *testing.T) {
			var request struct {
				A *bool `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=true"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.True(t, *request.A)
		})
	})
	t.Run("decodes to pointers to aliases to built-in types", func(t *testing.T) {
		t.Run("*string", func(t *testing.T) {
			type mytype string
			var request struct {
				A *mytype `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=foo"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, "foo", *request.A)
		})
		t.Run("*int", func(t *testing.T) {
			type mytype int
			var request struct {
				A *mytype `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=12"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, 12, *request.A)
		})
		t.Run("*int32", func(t *testing.T) {
			type mytype int32
			var request struct {
				A *mytype `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=12"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, 12, *request.A)
		})
		t.Run("*int64", func(t *testing.T) {
			type mytype int64
			var request struct {
				A *mytype `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=12"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, 12, *request.A)
		})
		t.Run("*bool", func(t *testing.T) {
			type mytype bool
			var request struct {
				A *mytype `filter:"a"`
			}
			err := Decode(mustValues("filter[a]=true"), &request)
			require.NoError(t, err)
			require.NotNil(t, request.A)
			require.EqualValues(t, true, *request.A)
		})
	})
	t.Run("decodes to a slice of strings", func(t *testing.T) {
		var request struct {
			A []string `filter:"a"`
		}
		err := Decode(mustValues("filter[a]=foo,bar,baz"), &request)
		require.NoError(t, err)
		require.NotNil(t, request.A)
		require.EqualValues(t, []string{"foo", "bar", "baz"}, request.A)
	})
	t.Run("decodes to a slice of ints", func(t *testing.T) {
		var request struct {
			A []int `filter:"a"`
		}
		err := Decode(mustValues("filter[a]=1,2,"), &request)
		require.NoError(t, err)
		require.NotNil(t, request.A)
		require.EqualValues(t, []int{1, 2}, request.A)
	})
	t.Run("properly handles missing values", func(t *testing.T) {
		t.Run("leaves dest nil when values are nil", func(t *testing.T) {
			type Request struct {
				A *string `filter:"a"`
				B *int    `filter:"b"`
			}

			var request Request
			err := Decode(mustValues(`filter[a]&filter[b]=`), &request)
			require.NoError(t, err)
			require.Nil(t, request.A)
			require.Nil(t, request.B)
		})
	})
	t.Run("returns proper errors", func(t *testing.T) {
		type Request struct {
			A *int    `filter:"a"`
			B *string `filter:"b"`
		}
		var request Request

		t.Run("when can't convert to destination type", func(t *testing.T) {
			err := Decode(mustValues(`filter[a]=foo&filter[b]=123`), &request)
			require.IsType(t, errBadRequest{}, err)
			message := err.(errBadRequest)["a"].Error()
			require.Contains(t, message, "unable to cast")
		})
		t.Run("when filter type is not supported", func(t *testing.T) {
			err := Decode(mustValues(`filter[a]=foo&filter[b]=123&filter[c]=123`), &request)
			message := err.(errBadRequest)["c"].Error()
			require.NotNil(t, err)
			require.Contains(t, message, errNotSupportedParameter.Error())
		})
	})
	t.Run("panics when filter is not a pointer", func(t *testing.T) {
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter bool `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter string `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter int `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter int16 `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter int32 `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter int64 `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter uint `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter uint16 `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter uint32 `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter uint64 `filter:"a"`
		}{})
		requireDecodePanics(t, "filter[a]=foo", &struct {
			Filter uint64 `filter:"a"`
		}{})
	})
	t.Run("sets default values when specified", func(t *testing.T) {
		type Request struct {
			FilterA *string `filter:"a" default:"foo"`
			FilterB *string `filter:"b" default:"bar"`
			FilterC *int    `filter:"c" default:"123"`
		}

		var request Request
		err := Decode(mustValues("filter[a]=&filter[c]"), &request)

		require.NoError(t, err)
		require.NotNil(t, request.FilterA)
		require.NotNil(t, request.FilterB)
		require.NotNil(t, request.FilterC)
		require.Equal(t, "foo", *request.FilterA)
		require.Equal(t, "bar", *request.FilterB)
		require.Equal(t, 123, *request.FilterC)
	})
}

func TestDecodeSorts(t *testing.T) {
	t.Run("decodes to a slice of Sort", func(t *testing.T) {
		var request struct {
			Sort []Sort `url:"sort"`
		}
		err := Decode(mustValues("sort=a,b,-c,c.d,-e.f"), &request)
		require.NoError(t, err)
		require.Equal(t, []Sort{"a", "b", "-c", "c.d", "-e.f"}, request.Sort)
	})
	t.Run("decodes to a slice of string", func(t *testing.T) {
		var request struct {
			Sort []string `url:"sort"`
		}
		err := Decode(mustValues("sort=a,b,-c,c.d,-e.f"), &request)
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b", "-c", "c.d", "-e.f"}, request.Sort)
	})
	t.Run("decodes to a slice of other string-compatible types", func(t *testing.T) {
		type String string
		var request struct {
			Sort []String `url:"sort"`
		}
		err := Decode(mustValues("sort=a,b,-c,c.d,-e.f"), &request)
		require.NoError(t, err)
		require.Equal(t, []String{"a", "b", "-c", "c.d", "-e.f"}, request.Sort)
	})
	t.Run("decodes to a custom typed slice of custom typed strings", func(t *testing.T) {
		type String string
		type StringSlice []String
		var request struct {
			Sort StringSlice `url:"sort"`
		}
		err := Decode(mustValues("sort=a,b,-c,c.d,-e.f"), &request)
		require.NoError(t, err)
		require.EqualValues(t, StringSlice{"a", "b", "-c", "c.d", "-e.f"}, request.Sort)
	})
	t.Run("decodes to Sort", func(t *testing.T) {
		var request struct {
			Sort Sort `url:"sort"`
		}
		err := Decode(mustValues("sort=a"), &request)
		require.NoError(t, err)
		require.Equal(t, Sort("a"), request.Sort)

		t.Run("returns error for multiple keys", func(t *testing.T) {
			var request struct {
				Sort Sort `url:"sort"`
			}
			err := Decode(mustValues("sort=a,-b"), &request)
			require.Contains(t, err.Error(), betterreflect.ErrSingleValueExpected.Error())
		})
	})
	t.Run("decodes to string", func(t *testing.T) {
		var request struct {
			Sort string `url:"sort"`
		}
		err := Decode(mustValues("sort=a"), &request)
		require.NoError(t, err)
		require.Equal(t, "a", request.Sort)
	})
	t.Run("decodes to a custom types string", func(t *testing.T) {
		type String string
		var request struct {
			Sort String `url:"sort"`
		}
		err := Decode(mustValues("sort=a"), &request)
		require.NoError(t, err)
		require.Equal(t, String("a"), request.Sort)
	})
	t.Run("panics if got type different from slice of string", func(t *testing.T) {
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort bool `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort int `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort int16 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort int32 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort int64 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort uint `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort uint16 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort uint32 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort uint64 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *string `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *bool `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *int `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *int16 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *int32 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *int64 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *uint `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *uint16 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *uint32 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *uint64 `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort []int `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort []bool `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *[]string `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *[]int `url:"sort"`
		}{})
		requireDecodeError(t, "sort=a,b,-c,c.d,-e.f", &struct {
			Sort *[]bool `url:"sort"`
		}{})
	})
	t.Run("leaves zero-length slice when param is not specified", func(t *testing.T) {
		t.Run("when param is not present at all", func(t *testing.T) {
			var request struct {
				Sort []string `url:"sort"`
			}
			err := Decode(mustValues(""), &request)
			require.NoError(t, err)
			require.Len(t, request.Sort, 0)
		})
		t.Run("when param is present but has empty value", func(t *testing.T) {
			var request struct {
				Sort []string `url:"sort"`
			}
			err := Decode(mustValues("sort="), &request)
			require.NoError(t, err)
			require.Len(t, request.Sort, 0)
		})
		t.Run("when param is present but has empty value (without '=')", func(t *testing.T) {
			var request struct {
				Sort []string `url:"sort"`
			}
			err := Decode(mustValues("sort"), &request)
			require.NoError(t, err)
			require.Len(t, request.Sort, 0)
		})
	})
	t.Run("sets default values param is not specified but default values is present", func(t *testing.T) {
		type String string
		type StringSlice []String
		var request struct {
			Sort StringSlice `url:"sort" default:"a,b,-c,c.d,-e.f"`
		}
		err := Decode(mustValues("sort="), &request)
		require.NoError(t, err)
		require.EqualValues(t, StringSlice{"a", "b", "-c", "c.d", "-e.f"}, request.Sort)
	})
}

func requireDecodeError(t *testing.T, s string, value interface{}) {
	require.Error(t, Decode(mustValues(s), value))
}

func TestDecodePages(t *testing.T) {
	var request struct {
		Order  string `page:"order"`
		Cursor string `page:"cursor"`
		Limit  uint64 `page:"limit"`
	}
	err := Decode(mustValues(`page[cursor]=foobar&page[limit]=10&page[order]=desc`), &request)
	require.NoError(t, err)
	require.Equal(t, "desc", request.Order)
	require.Equal(t, "foobar", request.Cursor)
	require.EqualValues(t, 10, request.Limit)

	t.Run("when order is a custom type", func(t *testing.T) {
		type Order string
		var request struct {
			Order  Order  `page:"order"`
			Cursor string `page:"cursor"`
			Limit  uint64 `page:"limit"`
		}
		err := Decode(mustValues(`page[cursor]=foobar&page[limit]=10&page[order]=desc`), &request)
		require.NoError(t, err)
		require.Equal(t, "foobar", request.Cursor)
		require.EqualValues(t, "desc", request.Order)
		require.EqualValues(t, 10, request.Limit)
	})

	t.Run("sets default values when params are missing but defaults are specified", func(t *testing.T) {
		type Order string
		var request struct {
			Order Order  `page:"order" default:"desc"`
			Limit uint64 `page:"limit" default:"24"`
		}
		err := Decode(mustValues(""), &request)
		require.NoError(t, err)
		require.EqualValues(t, "desc", request.Order)
		require.EqualValues(t, 24, request.Limit)
	})
}

func TestDecodeSearch(t *testing.T) {
	t.Run("decodes to pointer to string", func(t *testing.T) {
		var request struct {
			Search *string `url:"search"`
		}

		err := Decode(mustValues(`search=foo+bar`), &request)
		require.NoError(t, err)
		require.NotNil(t, request.Search)
		require.Equal(t, "foo bar", *request.Search)
	})
	t.Run("decodes to pointer to string alias", func(t *testing.T) {
		type mytype string
		var request struct {
			Search *mytype `url:"search"`
		}

		err := Decode(mustValues(`search=foo+bar`), &request)
		require.NoError(t, err)
		require.NotNil(t, request.Search)
		require.EqualValues(t, "foo bar", *request.Search)
	})
	t.Run("panics if field is not a pointer to string/string alias", func(t *testing.T) {
		requireDecodeError(t, "search=foo+bar", &struct {
			Search bool `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search int `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search int16 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search int32 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search int64 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search uint `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search uint16 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search uint32 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search uint64 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *bool `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *int `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *int16 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *int32 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *int64 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *uint `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *uint16 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *uint32 `url:"search"`
		}{})
		requireDecodeError(t, "search=foo+bar", &struct {
			Search *uint64 `url:"search"`
		}{})
	})
}

func TestDecodeIncludes(t *testing.T) {
	t.Run("decodes to bool", func(t *testing.T) {
		var request struct {
			IncludeA bool `include:"a"`
			IncludeB bool `include:"b"`
			IncludeC bool `include:"c"`
		}

		err := Decode(mustValues(`include=a,b`), &request)
		require.NoError(t, err)
		require.True(t, request.IncludeA)
		require.True(t, request.IncludeB)
		require.False(t, request.IncludeC)
	})
	t.Run("decodes when param has comma at the end", func(t *testing.T) {
		// may look strange, but had real bug with this
		var request struct {
			IncludeA bool `include:"a"`
			IncludeB bool `include:"b"`
			IncludeC bool `include:"c"`
		}

		err := Decode(mustValues(`include=a,b,`), &request)
		require.NoError(t, err)
		require.True(t, request.IncludeA)
		require.True(t, request.IncludeB)
		require.False(t, request.IncludeC)
	})
	t.Run("decodes empty include", func(t *testing.T) {
		// may look strange, but had real bug with this
		var request struct {
			IncludeA bool `include:"a"`
		}

		err := Decode(mustValues(`include=`), &request)
		require.NoError(t, err)
		require.False(t, request.IncludeA)
	})
	t.Run("panics if field is not a bool", func(t *testing.T) {
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *bool `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA string `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA int `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA int16 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA int32 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA int64 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA uint `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA uint16 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA uint32 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA uint64 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *string `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *int `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *int16 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *int32 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *int64 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *uint `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *uint16 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *uint32 `include:"a"`
		}{})
		requireDecodePanics(t, "include=a", &struct {
			IncludeA *uint64 `include:"a"`
		}{})
	})
}

func TestDecodeToNestedStructures(t *testing.T) {
	query := "page[number]=2&page[limit]=64&page[order]=desc&filter[author.first_name]=John&filter[author.last_name]=Doe"

	t.Run("supports nested structures", func(t *testing.T) {
		type PageParams struct {
			Number uint64 `page:"number"`
			Limit  uint64 `page:"limit"`
			Order  string `page:"order"`
		}

		type Filters struct {
			FirstName *string `filter:"author.first_name"`
			LastName  *string `filter:"author.last_name"`
		}

		type Request struct {
			Page    PageParams
			Filters Filters
		}

		var request Request
		err := Decode(mustValues(query), &request)

		require.NoError(t, err)
		require.EqualValues(t, 2, request.Page.Number)
		require.EqualValues(t, 64, request.Page.Limit)
		require.Equal(t, "desc", request.Page.Order)
		require.NotNil(t, request.Filters.FirstName)
		require.Equal(t, "John", *request.Filters.FirstName)
		require.NotNil(t, request.Filters.LastName)
		require.Equal(t, "Doe", *request.Filters.LastName)
	})

	t.Run("supports embed structures", func(t *testing.T) {
		type PageParams struct {
			Number uint64 `page:"number"`
			Limit  uint64 `page:"limit"`
			Order  string `page:"order"`
		}

		type Filters struct {
			FirstName *string `filter:"author.first_name"`
			LastName  *string `filter:"author.last_name"`
		}

		type Request struct {
			PageParams
			Filters
		}

		var request Request
		err := Decode(mustValues(query), &request)

		require.NoError(t, err)
		require.EqualValues(t, 2, request.Number)
		require.EqualValues(t, 64, request.Limit)
		require.Equal(t, "desc", request.Order)
		require.NotNil(t, *request.FirstName)
		require.Equal(t, "John", *request.FirstName)
		require.NotNil(t, *request.LastName)
		require.Equal(t, "Doe", *request.LastName)
	})

	t.Run("panics when embed structures has similar parameters", func(t *testing.T) {
		type FiltersA struct {
			FirstName string `filter:"author.first_name"`
		}

		type FiltersB struct {
			FirstName string `filter:"author.first_name"`
		}

		type Request struct {
			FiltersB
			Filters   FiltersA
			FirstName string `filter:"author.first_name"`
		}

		var request Request

		require.Panics(t, func() {
			_ = Decode(mustValues("filter[author.first_name]=John"), &request)
		})
	})
	t.Run("panics when embed structure has similar parameter with the root one", func(t *testing.T) {
		type Filters struct {
			FirstName string `filter:"author.first_name"`
		}

		type Request struct {
			AnotherFirstName string `filter:"author.first_name"`
			Filters
		}

		var request Request

		require.Panics(t, func() {
			_ = Decode(mustValues("filter[author.first_name]=John"), &request)
		})
	})
}

func TestDecodeErrors(t *testing.T) {
	testErrs := func(t *testing.T, query string, dest interface{}, expectedErr errBadRequest) {
		err := Decode(mustValues(query), dest)

		require.NotNil(t, err)
		require.IsType(t, errBadRequest{}, err)

		actualErr := err.(errBadRequest)

		require.Equal(t, expectedErr.Error(), actualErr.Error())
	}

	t.Run("when contains not JSON API param", func(t *testing.T) {
		testErrs(t, "my_param1=1", struct{}{}, errBadRequest{
			"my_param1": errNotSupportedParameter,
		})
	})

	t.Run("when contains valid, but not supported param", func(t *testing.T) {
		testErrs(t, "filter[name]=John", struct{}{}, errBadRequest{
			"name": errNotSupportedParameter,
		})
	})
}

type customAlias int32

func (c *customAlias) UnmarshalText(b []byte) error {
	switch string(b) {
	case "a":
		*c = 1
	case "b":
		*c = 2
	default:
		return errors.New("unknown custom type value")
	}
	return nil
}

type customStruct struct {
	str string
}

func (m *customStruct) UnmarshalText(text []byte) error {
	m.str = string(text)
	return nil
}

func TestDecodeUnmarshalTextTypes(t *testing.T) {
	t.Run("when custom type defines UnmarshalText method", func(t *testing.T) {
		var request struct {
			A *customAlias `filter:"a"`
		}
		err := Decode(mustValues("filter[a]=a"), &request)
		require.NoError(t, err)
		require.NotNil(t, request.A)
		require.EqualValues(t, customAlias(1), *request.A)
	})
	t.Run("when slice of custom type defines UnmarshalText method", func(t *testing.T) {
		var request struct {
			A []customAlias `filter:"a"`
		}
		err := Decode(mustValues("filter[a]=a,b"), &request)
		require.NoError(t, err)
		require.NotNil(t, request.A)
		require.EqualValues(t, []customAlias{1, 2}, request.A)
	})
	t.Run("when slice of custom type defines UnmarshalText method and it is sort", func(t *testing.T) {
		var request struct {
			A []customAlias `url:"sort"`
		}
		err := Decode(mustValues("sort=a,b"), &request)
		require.NoError(t, err)
		require.NotNil(t, request.A)
		require.EqualValues(t, []customAlias{1, 2}, request.A)
	})
	t.Run("decodes to a TextUnmarshaller struct", func(t *testing.T) {
		var request struct {
			A *customStruct `filter:"a"`
		}
		err := Decode(mustValues("filter[a]=pretty epic huh?"), &request)
		require.NoError(t, err)
		require.EqualValues(t, "pretty epic huh?", request.A.str)
	})
	t.Run("decodes to a slice of TextUnmarshaller struct", func(t *testing.T) {
		var request struct {
			A []customStruct `filter:"a"`
		}
		err := Decode(mustValues("filter[a]=pretty epic huh?,yes!?"), &request)
		require.NoError(t, err)
		require.EqualValues(t, []customStruct{{"pretty epic huh?"}, {"yes!?"}}, request.A)
	})
	t.Run("decodes to a slice of TextUnmarshaller struct", func(t *testing.T) {
		var request struct {
			A customStruct `url:"sort"`
		}
		err := Decode(mustValues("sort=text"), &request)
		require.NoError(t, err)
		require.EqualValues(t, "text", request.A.str)
	})
}

type SortStruct struct {
	Field    string
	OrderASC bool
}

var defaultSort = SortStruct{
	Field:    "id",
	OrderASC: true,
}

func (s *SortStruct) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*s = defaultSort
		return nil
	}

	strText := string(text)
	if strings.HasPrefix(strText, "-") {
		s.OrderASC = false
	}

	switch strings.TrimPrefix(strText, "-") {
	case "id":
		s.Field = "id"
	case "price":
		s.Field = "price"
	default:
		return errors.New("invalid value")
	}

	return nil
}

func (s *SortStruct) MarshalText() ([]byte, error) {
	sign := ""
	if !s.OrderASC {
		sign = "-"
	}
	return []byte(sign + s.Field), nil
}

func TestDecodeURLTag(t *testing.T) {
	type request struct {
		Sort SortStruct `url:"sort"`
	}
	t.Run("no query with the filter, default stays", func(t *testing.T) {
		r := request{
			Sort: defaultSort,
		}
		err := Decode(url.Values{}, &r)
		require.NoError(t, err)
		require.Equal(t, defaultSort, r.Sort)
	})
	t.Run("not pointer", func(t *testing.T) {
		r := request{
			Sort: defaultSort,
		}
		err := Decode(mustValues("sort=-price"), &r)
		require.NoError(t, err)
		require.Equal(t, SortStruct{
			Field:    "price",
			OrderASC: false,
		}, r.Sort)
	})
}
