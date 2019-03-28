package doorman

var (
	// ErrNotAllowed constraint check failed
	ErrNotAllowed = &Error{"not allowed"}
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}
