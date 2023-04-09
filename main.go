package main

import (
	"errors"
	"fmt"
)

func SomeExternalFunction() error {
	return errors.New("couldn't do some internet or something idk")
}

func MyOwnFunction() Err {
	err := SomeExternalFunction()
	if err != nil {
		return NewError(err).
			CausedAt("couldn't run my function", "SomeExternalFunction failed").
			Code(500)
	}
	return EmptyErr()
}

func ComplicatedFunction(something bool) (string, Err) {
	if something {
		return "something", EmptyErr()
	}
	return "", NewError(errors.New("some-error"))
}

func main() {
	logger := CreateFormattedLogger(
		"[blue][date][reset] -- [level]: [yellow][log]",
		func(a ...any) {
			fmt.Print(a...)
		},
	)

	err := MyOwnFunction()
	err.LogIfErr(logger.Error)

	err = MyOwnFunction()
	err.LogIfErr(logger.Error)
}
