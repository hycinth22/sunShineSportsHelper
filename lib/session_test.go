package lib

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

var session *Session
var loginErr error

func init() {
	session = CreateSession()
	loginErr = session.Login("021640302", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))))
}
func TestLogin(t *testing.T) {
	// 已在init()中执行函数功能函数，此处仅检测结果
	if loginErr != nil {
		t.Log(loginErr.Error())
		t.Fatalf("%v", loginErr)
	}
	t.Logf("%+v", session)
}

func TestGetSportResult(t *testing.T) {
	r, err := session.GetSportResult()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}

func TestSmartCreateRecords(t *testing.T) {
	records := SmartCreateRecords(&LimitParams{
		RandDistance:        Float64Range{2.6, 4.0},
		LimitSingleDistance: Float64Range{2.0, 4.0},
		LimitTotalDistance:  Float64Range{2.0, 5.0},
		MinuteDuration:      IntRange{11, 20},
	}, 5, time.Now())
	for _, r := range records {
		t.Logf("%+v", r)
	}
}

func TestGetXtcode(t *testing.T) {
	if r := GetXtcode(4290, "2018-06-02 11:13:40"); r != "5438d151" {
		t.Log(r)
		t.FailNow()
	}
}

func TestGetXtcodeV2(t *testing.T) {
	if r := GetXtcodeV2(4502, "2018-09-15 11:01:24.7", "2.520"); r != "61b1c85e" {
		t.Log(r)
		t.FailNow()
	}
}

