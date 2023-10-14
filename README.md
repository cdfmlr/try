# try

`try` is a Go package that provides a Check and Handle mechanism for error handling that
avoids `if err != nil { return err }`s overshadowing the main logic.

> The abuse of this package may cause unexpected crashes. I recommend you to read the source code (it's short!) before using it.

Install:

```bash
go get github.com/cdfmlr/try
```

Usage:

```go
package example

import "github.com/cdfmlr/try"

func mayFail() (result Result, err error) {
    defer try.Handle(&err)

    a, err := doA()
    try.Check(err, try.Log())
    
    try.Assert(a > 0, "a must be positive")

    b, err := doB()
    try.Check(err)

    return Result{A: a, B: b}, nil
}
```

## How it works

The Check-Handle mechanism is implemented by panic-recover, which is not
designed for normal control flow. `try.Check` panics with a special error type,
and `try.Handle` recovers from the panics and converts them back to normal errors.

```go
defer try.Handle(&err)
err = doSomething()
try.Check(err, try.Log())

// is equivalent to

err = doSomething()
if err != nil {
    slog.Error("doSomething failed.", "err", err)
    return err
}
```

## How to use it

In next subsections, I will try to explain how this package should be used.
But it's too long and boring to read. 
I do recommend you to read the source code instead. 
It's much shorter and easier to understand, compared to my poor English.

### Check

`try.Check` allows `HandlerFunc`s to be passed in, which will be called
before it panics. The `HandlerFunc`s can be used to log the error (see try.Log),
and any other custom error handling logic.
`HandlerFunc` works like middleware in web frameworks or
the `Result::or_else` in Rust.

A `HandlerFunc` returns a nil error
indicates the error is handled, all further `HandlerFunc`s will be skipped,
and no panic will be raised.

### Handle

`try.Handle` requires a pointer to an error value as its argument, so that it
can modify the error value to the recovered error. A good practice is to
define the function with a named return value `err error`, 
and pass `&err` to `try.Handle`.

Because `recover()` can only be called inside a deferred function (and not be called in any nested function),
`try.Handle` must be used with `defer`. And it must not be wrapped.

### Assert

`try.Assert` is a simple wrapper of `try.Check`, which panics with a `try.AssertionError`
if the condition is not met. I guess it's useful in some cases.

Also, do `defer try.Handle(&err)` before using asserts.

### Log

I used to log errors like this:

```go
err = doSomething()
if err != nil {
	slog.Error("doSomething failed.", "err", err)
    return err
}
```

This causes the `if err != nil` even more verbose and makes the main logic
harder to read. So I wrote an `HandlerFunc` to do the logging with `try.Check`:

```go
err = doSomething()
try.Check(err, try.Log(try.WithMsg("doSomething failed.")))
```

Logging with (optional) custom message, logger and level is supported.

(Err, the functional options pattern makes it verbose in another way.
I tried a lot of ways to make it more concise, but failed.
Connect me if you have any idea.)

## Do not abuse it

I personally treat `Check` just as a syntax sugar for the `if err != nil`
idiom, and `defer try.Handle(&err)` is required to make it work.
Do not misunderstand it as a Try-Catch like other languages have.

A `try.Check` without a deferred `try.Handle` will cause the panic to escape,
which may finally crash the program. So, always use Check-Handle in pairs.

Though it is possible to Check an error in an inner function, and Handle it
in an outer function, it is not recommended. Check-Handle is not Try-Catch.
Go functions return error values instead of throwing exceptions. Do not
break the idiomatic paradigm.

If you are intended to use Check-Handle cross function boundaries anyway, at least
make sure that the most outer function of your package has a
`defer try.Handle(&err)` to recover from the inner panics.
Do never let the panics escape from your package, instead,
Return the errors to the callers through values.
See what stdlib `encoding/json` does for example.

## License

MIT License (c) 2023-present cdfmlr
