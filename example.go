package main

import (
	"errors"
	"fmt"
)

func SimpleLog(log Loggable) {
	logger := CreateSimpleConsoleLogger(
		func(a ...any) {
			fmt.Print(a...)
		},
	)
	logger.Error(log)
}

func SomeFuncWithReturns(err bool) (string, error) {
	if err {
		return "", errors.New("some failure")
	}
	return "hello", nil
}

func ExampleWithReturns(shouldErr bool) (string, Err) {
	res, err := SomeFuncWithReturns(shouldErr)
	if err != nil {
		return "", NewError(err).Caused("failed to run example with returns").LogIfErr(SimpleLog)
	}
	return res + ", World!", EmptyErr()
}

func Example() {
	val, err := ExampleWithReturns(true)
	if !err.Ok() {
		return
	}
	fmt.Println(val)
}
