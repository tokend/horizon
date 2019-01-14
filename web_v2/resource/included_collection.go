package resource

import "encoding/json"

type resource interface {
	//GetKey - returns key of the resource
	GetKey() Key
}

// includedCollection - an array of resource objects that are related to the primary data and/or
// each other (“included resources”).
type includedCollection struct {
	includes map[Key]resource
}

// Add - adds new include into collection. If one already present - overrides it
func (c *includedCollection) Add(include resource) {
	if c.includes == nil {
		c.includes = make(map[Key]resource)
	}

	c.includes[include.GetKey()] = include
}

//MarshalJSON - marshals include collection as array of json objects
func (c *includedCollection) MarshalJSON() ([]byte, error) {
	uniqueEntries := make([]resource, 0, len(c.includes))
	for _, value := range c.includes {
		uniqueEntries = append(uniqueEntries, value)
	}

	return json.Marshal(uniqueEntries)
}
