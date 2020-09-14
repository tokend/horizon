package urlval

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeFilters(t *testing.T) {
	t.Run("encodes from built-in types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			var request = struct {
				Foo string `filter:"foo"`
			}{Foo: "bar"}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=bar", encoded)
		})
		t.Run("int", func(t *testing.T) {
			var request = struct {
				Foo int `filter:"foo"`
			}{Foo: 1}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int32", func(t *testing.T) {
			var request = struct {
				Foo int32 `filter:"foo"`
			}{Foo: 1}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int64", func(t *testing.T) {
			var request = struct {
				Foo int64 `filter:"foo"`
			}{Foo: 1}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("bool", func(t *testing.T) {
			var request = struct {
				Foo bool `filter:"foo"`
			}{Foo: false}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=false", encoded)
		})
	})

	t.Run("encodes from aliases to built-in types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			type mytype string
			var request = struct {
				Foo mytype `filter:"foo"`
			}{Foo: "bar"}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=bar", encoded)
		})
		t.Run("int", func(t *testing.T) {
			type mytype int
			var request = struct {
				Foo mytype `filter:"foo"`
			}{Foo: 1}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int32", func(t *testing.T) {
			type mytype int32
			var request = struct {
				Foo mytype `filter:"foo"`
			}{Foo: 1}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int64", func(t *testing.T) {
			type mytype int64
			var request = struct {
				Foo mytype `filter:"foo"`
			}{Foo: 1}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("bool", func(t *testing.T) {
			type mytype bool
			var request = struct {
				Foo mytype `filter:"foo"`
			}{Foo: false}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=false", encoded)
		})
	})

	t.Run("encodes from pointers to built-in types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			var value = "bar"
			var request = struct {
				Foo *string `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=bar", encoded)
		})
		t.Run("int", func(t *testing.T) {
			var value = 1
			var request = struct {
				Foo *int `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int32", func(t *testing.T) {
			var value int32 = 1
			var request = struct {
				Foo *int32 `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int64", func(t *testing.T) {
			var value int64 = 1
			var request = struct {
				Foo *int64 `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("bool", func(t *testing.T) {
			var value = false
			var request = struct {
				Foo *bool `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=false", encoded)
		})
	})
	t.Run("encodes from pointers on aliases to built-in types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			type mytype string
			var value mytype = "bar"
			var request = struct {
				Foo *mytype `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=bar", encoded)
		})
		t.Run("int", func(t *testing.T) {
			type mytype int
			var value mytype = 1
			var request = struct {
				Foo *mytype `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int32", func(t *testing.T) {
			type mytype int32
			var value mytype = 1
			var request = struct {
				Foo *mytype `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("int64", func(t *testing.T) {
			type mytype int64
			var value mytype = 1
			var request = struct {
				Foo *mytype `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.EqualValues(t, "filter%5Bfoo%5D=1", encoded)
		})
		t.Run("bool", func(t *testing.T) {
			type mytype bool
			var value mytype = false
			var request = struct {
				Foo *mytype `filter:"foo"`
			}{Foo: &value}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "filter%5Bfoo%5D=false", encoded)
		})
	})
	t.Run("encodes from slice of strings", func(t *testing.T) {
		type Request struct {
			Foo []string `filter:"foo"`
		}

		request := Request{
			Foo: []string{"a", "b", "c"},
		}

		encoded, err := Encode(request)
		require.NoError(t, err)
		require.Equal(t, "filter%5Bfoo%5D=a%2Cb%2Cc", encoded)
	})
	t.Run("encodes primitive slices", func(t *testing.T) {
		type Request struct {
			Foo []int16 `filter:"foo"`
		}

		request := Request{
			Foo: []int16{1, 2, 3},
		}

		encoded, err := Encode(request)
		require.NoError(t, err)
		require.Equal(t, "filter%5Bfoo%5D=1%2C2%2C3", encoded)
	})
	t.Run("doesn't include nil filters at all", func(t *testing.T) {
		type Request struct {
			Foo *string `filter:"foo"`
			Bar *string `filter:"bar"`
		}
		request := Request{
			Foo: nil,
			Bar: nil,
		}

		encoded, err := Encode(request)
		require.NoError(t, err)
		require.Equal(t, "", encoded)
	})
}

func TestEncodeSorts(t *testing.T) {
	t.Run("properly encodes sorts", func(t *testing.T) {
		t.Run("from a slice of Sort", func(t *testing.T) {
			type Request struct {
				Sort []Sort `url:"sort"`
			}
			request := Request{
				Sort: []Sort{"a", "-b", "-c.d"},
			}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "sort=a%2C-b%2C-c.d", encoded)
		})
		t.Run("from a slice of strings", func(t *testing.T) {
			type Request struct {
				Sort []string `url:"sort"`
			}
			request := Request{
				Sort: []string{"a", "-b", "-c.d"},
			}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "sort=a%2C-b%2C-c.d", encoded)
		})
		t.Run("from a slice of other string-convertible type", func(t *testing.T) {
			type String string
			type Request struct {
				Sort []String `url:"sort"`
			}
			request := Request{
				Sort: []String{"a", "-b", "-c.d"},
			}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "sort=a%2C-b%2C-c.d", encoded)
		})
		t.Run("from a custom typed slice of custom typed strings", func(t *testing.T) {
			type String string
			type StringSlice []String
			type Request struct {
				Sort StringSlice `url:"sort"`
			}
			request := Request{
				Sort: StringSlice{"a", "-b", "-c.d"},
			}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "sort=a%2C-b%2C-c.d", encoded)
		})
	})
}

func TestEncodePages(t *testing.T) {
	type Request struct {
		Limit  uint64 `page:"limit"`
		Number uint64 `page:"number"`
		Order  string `page:"order"`
	}
	request := Request{
		Limit:  15,
		Number: 100,
		Order:  "desc",
	}

	encoded, err := Encode(request)
	require.NoError(t, err)
	require.Equal(t, "page%5Blimit%5D=15&page%5Bnumber%5D=100&page%5Border%5D=desc", encoded)

	t.Run("when order is a custom type", func(t *testing.T) {
		type Order string
		type Request struct {
			Limit  uint64 `page:"limit"`
			Number uint64 `page:"number"`
			Order  Order  `page:"order"`
		}
		request := Request{
			Limit:  15,
			Number: 100,
			Order:  "desc",
		}

		encoded, err := Encode(request)
		require.NoError(t, err)
		require.Equal(t, "page%5Blimit%5D=15&page%5Bnumber%5D=100&page%5Border%5D=desc", encoded)
	})
}

func TestEncodeSearch(t *testing.T) {
	t.Run("encodes properly", func(t *testing.T) {
		t.Run("from a string", func(t *testing.T) {
			type Request struct {
				Search string `url:"search"`
			}
			request := Request{Search: "foo bar"}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "search=foo+bar", encoded)
		})
		t.Run("from a *string", func(t *testing.T) {
			type Request struct {
				Search *string `url:"search"`
			}
			s := "foo bar"
			request := Request{Search: &s}

			encoded, err := Encode(request)
			require.NoError(t, err)
			require.Equal(t, "search=foo+bar", encoded)

			t.Run("isn't present in a query when field is nil", func(t *testing.T) {
				request := Request{Search: nil}

				encoded, err := Encode(request)
				require.NoError(t, err)
				require.Equal(t, "", encoded)
			})
		})
	})
}

func TestEncodeNestedStructs(t *testing.T) {
	type Page struct {
		Number uint64 `page:"number"`
		Limit  uint64 `page:"limit"`
		Order  string `page:"order"`
	}

	type Filters struct {
		FirstName string `filter:"author.first_name"`
		LastName  string `filter:"author.last_name"`
	}

	t.Run("encodes from nested structures", func(t *testing.T) {
		type Request struct {
			Page    Page
			Filters Filters
		}
		request := Request{
			Page: Page{
				Number: 2,
				Limit:  64,
				Order:  "desc",
			},
			Filters: Filters{
				FirstName: "John",
				LastName:  "Doe",
			},
		}

		encoded, err := Encode(request)
		require.NoError(t, err)
		require.NoError(t, err)
		require.Equal(t, "filter%5Bauthor.first_name%5D=John&filter%5Bauthor.last_name%5D=Doe&page%5Blimit%5D=64&page%5Bnumber%5D=2&page%5Border%5D=desc", encoded)
	})
	t.Run("encodes from embed structures", func(t *testing.T) {
		type Request struct {
			Page
			Filters
		}
		request := Request{
			Page: Page{
				Number: 2,
				Limit:  64,
				Order:  "desc",
			},
			Filters: Filters{
				FirstName: "John",
				LastName:  "Doe",
			},
		}

		encoded, err := Encode(request)
		require.NoError(t, err)
		require.Equal(t, "filter%5Bauthor.first_name%5D=John&filter%5Bauthor.last_name%5D=Doe&page%5Blimit%5D=64&page%5Bnumber%5D=2&page%5Border%5D=desc", encoded)
	})
}

func (c customAlias) MarshalText() ([]byte, error) {
	switch c {
	case 1:
		return []byte("a"), nil
	case 2:
		return []byte("b"), nil
	default:
		return nil, errors.New("unknown value")
	}
}

func TestEncodeCustomTypes(t *testing.T) {
	t.Run("custom type encoding", func(t *testing.T) {
		type Request struct {
			A customAlias `filter:"a"`
		}
		req := Request{A: 1}
		encoded, err := Encode(&req)
		require.NoError(t, err)
		require.Equal(t, "filter%5Ba%5D=a", encoded)
	})
	t.Run("custom type pointer encoding", func(t *testing.T) {
		type Request struct {
			A *customAlias `filter:"a"`
		}
		a := customAlias(1)
		req := Request{A: &a}

		encoded, err := Encode(&req)
		require.NoError(t, err)
		require.Equal(t, "filter%5Ba%5D=a", encoded)
	})
	t.Run("custom type null pointer encoding", func(t *testing.T) {
		type Request struct {
			A *customAlias `filter:"a"`
		}
		req := Request{A: nil}

		encoded, err := Encode(&req)
		require.NoError(t, err)
		require.Equal(t, "", encoded)
	})
	t.Run("custom type slice encoding", func(t *testing.T) {
		type Request struct {
			A []customAlias `filter:"a"`
		}
		req := Request{A: []customAlias{1, 2}}

		encoded, err := Encode(&req)
		require.NoError(t, err)
		require.Equal(t, "filter%5Ba%5D=a%2Cb", encoded)
	})
}

func TestEncodeInclude(t *testing.T) {
	t.Run("encode includes", func(t *testing.T) {
		type Request struct {
			IncludeA bool `include:"a"`
			IncludeB bool `include:"b"`
		}
		req := Request{IncludeA: true}
		encoded, err := Encode(&req)
		require.NoError(t, err)
		require.Equal(t, "include=a", encoded)
	})
}
