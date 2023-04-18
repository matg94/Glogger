package glogger

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

func SomeError() Err {
	return NewError(errors.New("some go error")).
		Caused("failed to run example with returns").
		LogIfErr(SimpleLog)
}

func SomeFuncWithReturns() (string, error) {
	return "", errors.New("some failure")
}

func ExampleWithReturns() (string, Err) {
	res, err := SomeFuncWithReturns()
	if err != nil {
		return "", NewError(err).Caused("failed to run example with returns").LogIfErr(SimpleLog)
	}
	return res, EmptyErr()
}

func Test() {
	logger := CreateSimpleConsoleLogger(
		func(a ...any) {
			fmt.Println(a...)
		},
	)
	logger.Info(Str("test"))
}

func Example() {
	Test()

}
