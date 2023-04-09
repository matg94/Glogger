package main

import "testing"

func TestLogString(t *testing.T) {
	err := "some error"
	logStr := Str(err)

	if logStr.Error() != err {
		t.Logf("expected logstring.Error() to match err but got %s", logStr.Error())
		t.Fail()
	}
}
