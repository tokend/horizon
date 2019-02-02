package regources

type Flag struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Flags struct {
	Int    int    `json:"int"`
	Values []Flag `json:"flags"`
}

func FlagsFromMask(mask int, allFlags map[int]string) Flags {
	values := []Flag{}

	for value, name := range allFlags {
		if (value & mask) == value {
			values = append(values, Flag{
				Value: value,
				Name:  name,
			})
		}
	}

	return Flags{
		Int:    mask,
		Values: values,
	}
}
