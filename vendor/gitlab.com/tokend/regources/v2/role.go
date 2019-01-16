package regources

type Role struct {
	ID      string                 `jsonapi:"primary,roles"`
	Details map[string]interface{} `jsonapi:"attr,details"`
	Rules   []*Rule                `jsonapi:"relation,rules,omitempty"`
}
