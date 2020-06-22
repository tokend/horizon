package internal

import (
	"net/url"
	"regexp"
	"strings"
)

type TokenType int32

const (
	TokenTypeFilter TokenType = 1 + iota
	TokenTypeInclude
	TokenTypePage
	TokenTypeInvalid
)

type Token struct {
	Type  TokenType
	Key   string
	Value string
}

// bool for consumed
type Tokens map[Token]bool

func Tokenize(values url.Values) (Tokens, error) {
	tokens := Tokens{}

	tokenizeIncludes(values, tokens)
	tokenizeFilters(values, tokens)
	tokenizePagination(values, tokens)

	for k := range values {
		tokens[Token{
			Type: TokenTypeInvalid,
			Key:  k,
		}] = false
	}

	return tokens, nil
}

func tokenizeIncludes(values url.Values, tokens Tokens) {
	includes := values.Get("include")
	if includes == "" {
		return
	}
	for _, include := range strings.Split(includes, ",") {
		tokens[Token{
			Type: TokenTypeInclude,
			Key:  include,
		}] = false
	}
	values.Del("include")
}

func tokenizeFilters(values url.Values, tokens Tokens) {
	for k, v := range values {
		// TODO: check v length
		ok, key := extractFilter(k)
		if ok {
			tokens[Token{
				Type:  TokenTypeFilter,
				Key:   key,
				Value: v[0],
			}] = false
			values.Del(k)
		}
	}
}

func extractFilter(s string) (bool, string) {
	r := regexp.MustCompile(`^filter\[([^\]]+)\]$`)
	match := r.FindStringSubmatch(s)
	if len(match) != 2 {
		return false, ""
	}
	return true, match[1]
}

func tokenizePagination(values url.Values, tokens Tokens) {
	for k, v := range values {
		// TODO: check v length
		ok, key := extractPage(k)
		if ok {
			tokens[Token{
				Type:  TokenTypePage,
				Key:   key,
				Value: v[0],
			}] = false
			values.Del(k)
		}
	}
}

func extractPage(s string) (bool, string) {
	r := regexp.MustCompile(`^page\[([^\]]+)\]$`)
	match := r.FindStringSubmatch(s)
	if len(match) != 2 {
		return false, ""
	}
	return true, match[1]
}
