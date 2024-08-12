package postgres

import "fmt"

type UniqueFieldErr struct {
	Payload string
	Value   string
	Err     error
}

func (e *UniqueFieldErr) Error() string {
	return fmt.Sprintf("value: %s, already exist. %v", e.Value, e.Err)
}

func NewUniqueFieldErr(v, p string, err error) error {
	return &UniqueFieldErr{
		Value:   v,
		Payload: p,
		Err:     err,
	}
}
