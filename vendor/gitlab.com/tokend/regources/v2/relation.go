package regources

type Relation struct {
	Links *Links `json:"links,omitempty"`
	Data  *Key   `json:"data,omitempty"`
}

type RelationCollection struct {
	Links *Links `json:"links,omitempty"`
	Data  []Key  `json:"data,omitempty"`
}
