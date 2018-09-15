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
	// 已在init()中执行功能函数，此处仅检测结果
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
	records := SmartCreateRecords(0, &LimitParams{
		RandDistance:        Float64Range{2.6, 4.0},
		LimitSingleDistance: Float64Range{2.0, 4.0},
		LimitTotalDistance:  Float64Range{2.0, 5.0},
		MinuteDuration:      IntRange{11, 20},
	}, 5, time.Now())
	for _, r := range records {
		t.Logf("%+v", r)
	}
}

