package glogger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestGetBaseError(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	if err.GetBaseErr() != baseErr {
		t.Logf("expected base error to match but got %s", baseErr.Error())
		t.Fail()
	}
}

func TestErrorBasicCausedAt(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	wrapped := err.CausedAt("some new error", "first causedAt usage")

	if wrapped.GetBaseErr() != baseErr {
		t.Logf("expected base error to match but got %s", baseErr.Error())
		t.Fail()
	}

	if len(wrapped.GetStack()) != 1 {
		t.Logf("expected stack to contain 1 but got %d", len(wrapped.GetStack()))
		t.Fail()
	}

	newSubErr := SubErr{
		err: "some new error",
		at:  "first causedAt usage",
	}

	if wrapped.GetStack()[0] != newSubErr {
		t.Logf("expected stack to match but did not")
		t.Fail()
	}
}

func TestErrorBasicCaused(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	wrapped := err.Caused("some new error")

	if wrapped.GetBaseErr() != baseErr {
		t.Logf("expected base error to match but got %s", baseErr.Error())
		t.Fail()
	}

	if len(wrapped.GetStack()) != 1 {
		t.Logf("expected stack to contain 1 but got %d", len(wrapped.GetStack()))
		t.Fail()
	}

	newSubErr := SubErr{
		err: "some new error",
		at:  "",
	}

	if wrapped.GetStack()[0] != newSubErr {
		t.Logf("expected stack to match but did not")
		t.Fail()
	}
}

func TestErrorCodes(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	coded := err.Code(404)

	if coded.GetErrCode() != 404 {
		t.Logf("expected error code to match but got %d", coded.GetErrCode())
		t.Fail()
	}

}

func TestErrOK(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	if err.Ok() {
		t.Logf("expected error to not be ok but got: %t", err.Ok())
		t.Fail()
	}

	emptyErr := NewError(nil)
	t.Log(emptyErr.baseErr)

	if !emptyErr.Ok() {
		t.Logf("expected error to be ok but got: %t", emptyErr.Ok())
		t.Fail()
	}
}

func TestErrorMultiStack(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	wrapped_once := err.CausedAt("some new error", "some location")
	wrapped_twice := wrapped_once.Caused("some new error 2")

	if len(wrapped_twice.GetStack()) != 2 {
		t.Logf("expected stack to contain 2 but got %d", len(wrapped_twice.GetStack()))
		t.Fail()
	}

	once_subErr := SubErr{
		err: "some new error",
		at:  "some location",
	}

	if wrapped_twice.GetStack()[0] != once_subErr {
		t.Logf("expected first stack value to match but did not")
		t.Fail()
	}

	twice_subErr := SubErr{
		err: "some new error 2",
		at:  "",
	}

	if wrapped_twice.GetStack()[1] != twice_subErr {
		t.Logf("expected second stack value to match but did not")
		t.Fail()
	}
}

func TestPrintableStackTrace(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr)

	wrapped_once := err.CausedAt("some new error", "some location")
	wrapped_twice := wrapped_once.Caused("some new error 2")

	expectedStackTrace := `root cause: some error
caused: some new error
    at: some location
caused: some new error 2
`

	if wrapped_twice.Error() != expectedStackTrace {
		t.Logf("expected stack trace to match example but got %s", wrapped_twice.GetStackTrace())
		t.Fail()
	}
}

func TestJSONSerialization(t *testing.T) {
	// Set up an Error object for testing
	err := NewError(errors.New("root error")).
		Code(500).
		CausedAt("sub error 1", "location 1").
		CausedAt("sub error 2", "location 2")

	// Call the JSON method and check the output
	jsonStr, parseErr := err.JSON()
	t.Log(jsonStr)
	if !parseErr.Ok() {
		t.Errorf("Unexpected error returned: %v", parseErr.GetStackTrace())
	}

	// Unmarshal the JSON string into a JSONErr object
	var jsonErr JSONErr
	mErr := json.Unmarshal([]byte(jsonStr), &jsonErr)
	if mErr != nil {
		t.Errorf("Failed to unmarshal JSON string: %v", mErr)
	}

	// Check the JSONErr fields against the original Error object
	if jsonErr.Root != "root error" {
		t.Errorf("Unexpected root error: %s", jsonErr.Root)
	}
	if jsonErr.Code != 500 {
		t.Errorf("Unexpected error code: %d", jsonErr.Code)
	}
	if len(jsonErr.Caused) != 2 {
		t.Errorf("Unexpected number of sub errors: %d", len(jsonErr.Caused))
	}
	if jsonErr.Caused[0].Err != "sub error 1" {
		t.Errorf("Unexpected sub error message: %s", jsonErr.Caused[0].Err)
	}
	if jsonErr.Caused[0].Pos != 0 {
		t.Errorf("Unexpected sub error position: %d", jsonErr.Caused[0].Pos)
	}
	if jsonErr.Caused[0].At != "location 1" {
		t.Errorf("Unexpected sub error location: %s", jsonErr.Caused[0].At)
	}
	if jsonErr.Caused[1].Err != "sub error 2" {
		t.Errorf("Unexpected sub error message: %s", jsonErr.Caused[1].Err)
	}
	if jsonErr.Caused[1].Pos != 1 {
		t.Errorf("Unexpected sub error position: %d", jsonErr.Caused[1].Pos)
	}
	if jsonErr.Caused[1].At != "location 2" {
		t.Errorf("Unexpected sub error location: %s", jsonErr.Caused[1].At)
	}
}

func TestUserErrLoggedIfExists(t *testing.T) {
	baseErr := errors.New("some error")
	err := NewError(baseErr).UserErr("could not process").Caused("some issue")

	expectedUserErr := "could not process"

	buffer := new(bytes.Buffer)
	_ = err.LogIfErr(func(err Loggable) {
		fmt.Fprint(buffer, err.Error())
	})
	assertEqual(t, expectedUserErr, buffer.String())
	buffer.Reset()
}

func TestLogIfErr(t *testing.T) {
	baseErr := errors.New("error")
	err := NewError(baseErr)

	expectedStackTrace := "root cause: error\n"

	buffer := new(bytes.Buffer)
	_ = err.LogIfErr(func(err Loggable) {
		fmt.Fprint(buffer, err.Error())
	})
	assertEqual(t, expectedStackTrace, buffer.String())
	buffer.Reset()

	emptyErr := EmptyErr()
	_ = emptyErr.LogIfErr(func(err Loggable) {
		fmt.Fprint(buffer, err.Error())
	})
	assertEqual(t, "", buffer.String())
	buffer.Reset()
}
