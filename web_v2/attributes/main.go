package attributes

type Model interface {
	MarshalSelf() ([]byte, error)
}
