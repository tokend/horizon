package urlval

import (
	"fmt"
	"net/url"
	"reflect"

	"gitlab.com/distributed_lab/urlval/internal"
	"gitlab.com/distributed_lab/urlval/internal/betterreflect"
)

func Decode(values url.Values, dest interface{}) error {
	tokens, _ := internal.Tokenize(values)

	refdest := betterreflect.NewStruct(dest)
	for token, consumed := range tokens {
		if consumed {
			continue
		}
		if ok := decodeToken(token, refdest); !ok {
			delete(tokens, token)
			tokens[internal.Token{
				Type: internal.TokenTypeInvalid,
				Key:  token.Key,
			}] = true
		}
	}

	for token := range tokens {
		if token.Type == internal.TokenTypeInvalid {
			fmt.Println("UNKNOWN", token.Key)
		}
	}

	return nil
}

func decodeToken(token internal.Token, dest *betterreflect.Struct) bool {
	switch token.Type {
	case internal.TokenTypeFilter:
		for i := 0; i < dest.NumField(); i++ {
			tag := dest.Tag(i, "filter")
			if tag == token.Key {
				_ = dest.Set(i, token.Value)
				return true
			}
		}
	case internal.TokenTypeInclude:
		for i := 0; i < dest.NumField(); i++ {
			tag := dest.Tag(i, "include")
			if tag == token.Key {
				if dest.Type(i).Kind() != reflect.Bool {
					panic("bool expected for include")
				}
				_ = dest.Set(i, true)
				return true
			}
		}
	case internal.TokenTypePage:
		for i := 0; i < dest.NumField(); i++ {
			tag := dest.Tag(i, "page")
			if tag == token.Key {
				_ = dest.Set(i, token.Value)
				return true
			}
		}
	case internal.TokenTypeInvalid:
	default:
		panic("unknown token type")
	}
	return false
}
