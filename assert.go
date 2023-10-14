package try

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrAssertion struct {
	Msg      string
	Got      any
	Expected any
}

func (e ErrAssertion) Error() string {
	s := fmt.Sprintf("assert failed (expected %v, got %v)", e.Expected, e.Got)
	if e.Msg != "" {
		s = s + ": " + e.Msg
	}
	return s
}

// assert returns an *ErrAssertion if t is not true.
func assert(t bool, msg ...string) error {
	if !t {
		return &ErrAssertion{
			Msg:      strings.Join(msg, " "),
			Expected: true,
			Got:      false,
		}
	}
	return nil
}

// Assert panics with errCheck(*ErrAssertion{}) if t is false.
// The assert-panic should be handled by try.Handle.
// If msg is not empty, they will be joined as the error message.
//
// Assert is designed to make a quick exit of a function (by try.Check & try.Handle)
// when unexpected conditions (t == false) are met.
//
// In any case that the condition MAY be true, use the if statement
// instead of Assert (the implementation of the unexported try.assert
// is an example for this).
func Assert(t bool, msg ...string) {
	Check(assert(t, msg...))
}

// AssertEqual is a shorthand for Assert(got != expected, msg...).
func AssertEqual(got, expected any, msg ...string) {
	err := assert(got == expected, msg...)
	if err != nil {
		err.(*ErrAssertion).Expected = expected
		err.(*ErrAssertion).Got = got
		Check(err)
	}
}

// AssertDeepEqual is a shorthand for Assert(reflect.DeepEqual(got, expected), msg...).
func AssertDeepEqual(got, expected any, msg ...string) {
	err := assert(reflect.DeepEqual(got, expected), msg...)
	if err != nil {
		err.(*ErrAssertion).Expected = expected
		err.(*ErrAssertion).Got = got
		Check(err)
	}
}
