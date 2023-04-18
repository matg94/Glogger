package glogger

import (
	"encoding/json"
	"fmt"
)

type Err interface {
	GetBaseErr() error
	GetErrCode() int
	GetStack() []SubErr
	GetStackTrace() string
	Error() string
	JSON() (string, Err)
	Caused(string) Err
	CausedAt(string, string) Err
	Code(int) Err
	Ok() bool
	LogIfErr(func(Loggable)) Err
}

type SubErr struct {
	err string
	at  string
}

type Error struct {
	baseErr error
	code    int
	stack   []SubErr
}

type JSONErr struct {
	Root   string       `json:"root"`
	Code   int          `json:"code"`
	Caused []JSONSubErr `json:"errors"`
}

type JSONSubErr struct {
	Err string `json:"error"`
	Pos int    `json:"order"`
	At  string `json:"at"`
}

func NewError(base error) *Error {
	return &Error{
		baseErr: base,
	}
}

func EmptyErr() *Error {
	return &Error{}
}

func (e *Error) GetBaseErr() error {
	return e.baseErr
}

func (e *Error) GetErrCode() int {
	return e.code
}

func (e *Error) GetStack() []SubErr {
	return e.stack
}

func (e *Error) GetStackTrace() string {
	trace := fmt.Sprintf("root cause: %s\n", e.GetBaseErr().Error())
	for _, caused := range e.GetStack() {
		trace += fmt.Sprintf("caused: %s\n", caused.err)
		if caused.at != "" {
			trace += fmt.Sprintf("    at: %s\n", caused.at)
		}
	}
	return trace
}

func (e *Error) Error() string {
	return e.GetStackTrace()
}

func (e *Error) JSON() (string, Err) {
	jsonSubErrs := make([]JSONSubErr, len(e.stack))
	for i, subErr := range e.GetStack() {
		jsonSubErrs[i] = JSONSubErr{
			Err: subErr.err,
			Pos: i,
			At:  subErr.at,
		}
	}

	jsonErr := JSONErr{
		Root:   e.GetBaseErr().Error(),
		Code:   e.GetErrCode(),
		Caused: jsonSubErrs,
	}

	jsonBytes, err := json.Marshal(jsonErr)
	if err != nil {
		return "", NewError(err).CausedAt("failed to serialize error object", "EZrror JSONify implementation marshal")
	}

	return string(jsonBytes), EmptyErr()
}

func (e *Error) Caused(err string) Err {
	subErr := SubErr{
		err: err,
		at:  "",
	}
	e.stack = append(e.stack, subErr)
	return e
}

func (e *Error) CausedAt(err, at string) Err {
	subErr := SubErr{
		err: err,
		at:  at,
	}
	e.stack = append(e.stack, subErr)
	return e
}

func (e *Error) Code(code int) Err {
	e.code = code
	return e
}

func (e *Error) LogIfErr(log func(Loggable)) Err {
	if !e.Ok() {
		log(e)
	}
	return e
}

func (e *Error) Ok() bool {
	return e.baseErr == nil
}
