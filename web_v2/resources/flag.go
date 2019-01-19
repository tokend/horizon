package resources

import "github.com/cheekybits/genny/generic"

type flagValueType generic.Number

//ShortString - provides short name of the flag value
func (f flagValueType) ShortString() string {
	return ""
}
