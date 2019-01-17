package regources

type Signer struct {
	ID       string                 `jsonapi:"primary,signers"`
	Weight   int                    `jsonapi:"attr,weight"`
	Role     *Role                  `jsonapi:"relation,role,omitempty"`
	Identity int                    `jsonapi:"attr,identity"`
	Details  map[string]interface{} `jsonapi:"attr,details"`
}
