package config

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/swarmfund/horizon/assets"
)

const (
	MUST_NOT_BE_EMPTY string = "Must not be empty"
)

type Base struct {
	parent *Base
	prefix string
}

func NewBase(parent *Base, prefix string) *Base {
	return &Base{
		parent: parent,
		prefix: prefix,
	}
}

func (c *Base) invalidField(field, reason string) error {
	return fmt.Errorf("Field '%s' has invalid value. Error: %s", c.getEnvName(field), reason)
}

func (c *Base) getFieldName(baseName string) string {
	fieldName := c.getPrefix()
	if fieldName != "" {
		fieldName += "_"
	}

	fieldName += baseName
	return fieldName
}

func (c *Base) getString(name string) string {
	return viper.GetString(c.getFieldName(name))
}

func (c *Base) getNonEmptyString(name string) (string, error) {
	result := c.getString(name)
	if result == "" {
		return "", c.invalidField(name, MUST_NOT_BE_EMPTY)
	}

	return result, nil
}

func (c *Base) getEnvName(baseName string) string {
	return strings.ToUpper(c.getFieldName(baseName))
}

func (c *Base) bindEnv(baseName string) {
	viper.BindEnv(c.getFieldName(baseName), c.getEnvName(baseName))
}

func (c *Base) setDefault(baseName string, value interface{}) {
	viper.SetDefault(c.getFieldName(baseName), value)
}

func (c *Base) getPrefix() string {
	if c.parent != nil {
		parentPrefix := c.parent.getPrefix()
		if parentPrefix != "" {
			return parentPrefix + "_" + c.prefix
		}
	}
	return c.prefix
}

func (c *Base) getParsedURL(name string) (*url.URL, error) {
	rawUrl, err := c.getNonEmptyString(name)
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return nil, c.invalidField(name, err.Error())
	}

	return parsed, nil
}

func (c *Base) getOptionalParsedURL(name string) (*url.URL, error) {
	rawUrl := c.getString(name)
	if rawUrl == "" {
		// No URL is provided - it's OK value is optional
		return nil, nil
	}

	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return nil, c.invalidField(name, err.Error())
	}

	return parsed, nil
}

func (c *Base) getURLAsString(name string) (string, error) {
	result, err := c.getParsedURL(name)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (c *Base) getInt(name string) int {
	return viper.GetInt(c.getFieldName(name))
}

func (c *Base) getBool(name string) bool {
	return viper.GetBool(c.getFieldName(name))
}

func (c *Base) getTimeDuration(name string) (time.Duration, error) {
	rawDur := viper.GetString(c.getFieldName(name))
	return time.ParseDuration(rawDur)
}

func (c *Base) getOptionalTDuration(name string) (*time.Duration, error) {
	rawDur := viper.GetString(c.getFieldName(name))
	if rawDur == "" {
		return nil, nil
	}
	tDur, err := time.ParseDuration(rawDur)
	return &tDur, err
}
func (c *Base) getTemplate(name string) *template.Template {
	return assets.Templates.Lookup(name)
}
