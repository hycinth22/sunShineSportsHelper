package lib

import (
	"testing"
)

func TestGetXtcode(t *testing.T) {
	if r := GetXTcode(4290, "2018-06-02 11:13:40"); r != "5438d151" {
		t.Log(r)
		t.FailNow()
	}
}

func TestGetXtcodeV2codeV2(t *testing.T) {
	if r := GetXTcodeV2(4502, "2018-09-15 11:01:24.7", "2.520"); r != "61b1c85e" {
		t.Log(r)
		t.FailNow()
	}
}
