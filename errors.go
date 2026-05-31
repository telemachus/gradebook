package gradebook

// InvalidClassError wraps class invariant violations found while parsing.
type InvalidClassError struct {
	Err error
}

func (e *InvalidClassError) Error() string {
	return e.Err.Error()
}

func (e *InvalidClassError) Unwrap() error {
	return e.Err
}
