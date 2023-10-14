package try

// Check if the error is nil, if not, handle the error with the handlers.
// If any handler returns a nil error, the error is considered handled and
// the function returns. Otherwise, Check panics with the error.
//
// With no handlers, Check is equivalent to:
//
//	if err != nil {
//	    panic(errCheck{err})
//	}
//
// NOTICE: There should be a defer Handle(&err) before calling Check(err error).
func Check(err error, handlers ...HandlerFunc) {
	for _, h := range handlers {
		if err == nil { // early stop
			return
		}
		err = h(err)
	}

	if err != nil { // throw the error
		panic(errCheck{Err: err})
	}
}

// Handle recover from the panicked error raised by Check(err error),
// fill the error to errOut if errOut is not nil.
//
// If the panic is not raised by Check, Handle panics again.
func Handle(errOut *error) {
	switch v := recover().(type) {
	case nil:
		return
	case errCheck: // raised by Check
		if errOut != nil {
			*errOut = v.Err
		}
	default:
		panic(v)
	}
}

// errCheck is used to identify if the panic is raised by Check.
type errCheck struct {
	Err error
}

func (e errCheck) Error() string {
	// never called except "Check without defer Handler" => panic
	return "try.errCheck: " + e.Err.Error() +
		" (Note: This panic may be raised by a call to try.Check without " +
		"a defer try.Handle(&err) before it, which may be a wrong usage."
}

// HandlerFunc is a function that handles the error
// when Check(err error) is called with a non-nil error.
//
// If the error is handled, the error should be returned as nil.
type HandlerFunc func(err error) error
