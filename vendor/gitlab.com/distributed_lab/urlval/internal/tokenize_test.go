package internal

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractFilter(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{"filter[state]", "state"},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			ok, got := extractFilter(tc.in)
			require.Equal(t, tc.expected != "", ok)
			require.Equal(t, tc.expected, got)
		})
	}
}

func Test_tokenizeInclude(t *testing.T) {
	values := url.Values{
		"include": []string{"a,b,c"},
	}
	tokens := Tokens{}

	tokenizeIncludes(values, tokens)

	require.Contains(t, tokens, Token{Type: TokenTypeInclude, Key: "a", Value: ""})
	require.Contains(t, tokens, Token{Type: TokenTypeInclude, Key: "b", Value: ""})
	require.Contains(t, tokens, Token{Type: TokenTypeInclude, Key: "c", Value: ""})
}
