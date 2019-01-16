package regources

type Rule struct {
	ID          string                 `jsonapi:"primary,rules"`
	Resource    string                 `jsonapi:"attr,resource"`
	Action      string                 `jsonapi:"attr,action"`
	IsForbidden bool                   `jsonapi:"attr,is_forbidden"`
	Details     map[string]interface{} `jsonapi:"attr,details"`
}
