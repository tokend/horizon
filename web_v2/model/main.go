package model

type Model interface {
	MarshalSelf() ([]byte, error)
}
