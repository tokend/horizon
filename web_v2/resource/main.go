package resource

import "gitlab.com/tokend/regources"

type LinksObject struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

// Mask - represent bit mask
type Mask struct {
	Flags []regources.Flag `json:"flags"`
	Mask  int32            `json:"mask"`
}
