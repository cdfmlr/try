// Package try provides a Check and Handle mechanism for error handling that
// avoids `if err != nil { return err }`s overshadowing the main logic.
//
// I personally treat this package just as a syntax sugar for the `if err != nil`
// idiom. Do not misunderstand it as a Try-Catch like other languages have.
// The abuse of this package may cause unexpected crashes.
// I recommend you to read the source code (it's short!) before using it.
//
// Example:
//
//	func mayFail() (result Result, err error) {
//	    defer try.Handle(&err)
//
//	    a, err := doA()
//	    try.Check(err, try.Log())
//
//	    b, err := doB()
//	    try.Check(err)
//
//	    return Result{A: a, B: b}, nil
//	}
//
// Equal to:
//
//	func mayFail() (result Result, err error) {
//	    a, err := doA()
//	    if err != nil {
//	        slog.Error("doA failed", "err", err)
//	        return result, err
//	    }
//
//	    b, err := doB()
//	    if err != nil {
//	        return result, err
//	    }
//
//	    return Result{A: a, B: b}, nil
//	}
//
// The Check-Handle mechanism is implemented by panic-recover, which is not
// designed for normal control flow. Check panics with a special error type,
// and Handle recovers from the panics and converts them back to normal errors.
//
// A Check without a defer Handle will cause the panic to escape,
// which may finally crash the program. So, always use Check-Handle in pairs.
//
// Though it is possible to Check an error in an inner function, and Handle it
// in an outer function, it is not recommended. Check-Handle is not Try-Catch.
// Go functions return error values instead of throwing exceptions. Do not
// break the idiomatic paradigm.
//
// If you are intended to use Check-Handle cross function boundaries anyway, at least
// make sure that the most outer function of your package has a
// defer Handle(&err) to recover from the inner panics.
// Do never let the panics escape from your package, instead,
// Return the errors to the callers through values.
// See what stdlib encoding/json does for example.
package try
