package serial

type Error struct {
	str string
	e   error
}

func (e Error) Error() string {
	return e.str
}

func (e Error) Unwrap() error {
	return e.e
}
