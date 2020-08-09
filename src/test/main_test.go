package test

import (
	"testing"
)

func TestMain(t *testing.T) {
	if true == false {
		t.Error("failure")
	}
}
