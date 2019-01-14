package resource

//RelationshipCollection  - represents collection of resource identifier objects to be included as relationship of
// sole resource
type RelationshipCollection struct {
	Data []Key `json:"data"`
}

// Add - adds new key to the data
func (c *RelationshipCollection) Add(key Key) {
	c.Data = append(c.Data, key)
}
