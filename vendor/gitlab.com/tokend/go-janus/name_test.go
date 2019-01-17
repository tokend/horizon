package janus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) {
	t.Run("with parameter and /*/", func(t *testing.T) {
		endpoint := "/users/*/{id}"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "users-x-get", name)
	})

	t.Run("with parameter", func(t *testing.T) {
		endpoint := "/users/{id}"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "users-x-get", name)
	})

	t.Run("simple", func(t *testing.T) {
		endpoint := "/users/"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "users-get", name)
	})

	t.Run("root", func(t *testing.T) {
		endpoint := "/"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "root-get", name)
	})

	t.Run("with two parameters and /*/", func(t *testing.T) {
		endpoint := "/users/*/{id}/another/{one}/*/"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "users-x-another-x-get", name)
	})

	t.Run("with two parameters", func(t *testing.T) {
		endpoint := "/users/{id}/another/{one}/"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "users-x-another-x-get", name)
	})

	t.Run("without parameters", func(t *testing.T) {
		endpoint := "/another/one/and/another/one/"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "another-one-and-another-one-get", name)
	})

	t.Run("without parameters, with /*/", func(t *testing.T) {
		endpoint := "/another/*/one/and/*/another/one/*/"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "another-one-and-another-one-get", name)
	})

	t.Run("with underscore", func(t *testing.T) {
		endpoint := "/another_one"
		method := "GET"
		name := GetName(endpoint, method)
		assert.Equal(t, "another-one-get", name)
	})
}
