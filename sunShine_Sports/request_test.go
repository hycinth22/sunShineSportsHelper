package sunShine_Sports

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

var session *Session
var loginErr error

var fakeUserInfo = UserInfo{
	Id:            666,
	InClassName:   "haha",
	StudentName:   "666",
	StudentNumber: "6666",
	Sex:           "F",
	DistanceLimit: &DistanceParams{
		RandDistance:        Float64Range{2.6, 4.0},
		LimitSingleDistance: Float64Range{2.0, 4.0},
		LimitTotalDistance:  Float64Range{2.0, 5.0},
	},
}

func init() {
	session, loginErr = Login("021640302", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))))
}
func TestLogin(t *testing.T) {
	// 已在init()中执行函数功能函数，此处仅检测结果
	if loginErr != nil {
		if HTTPErr, ok := loginErr.(HTTPError); ok {
			t.Log(HTTPErr)
		}
		t.Fatalf("%v", loginErr)
	}
	t.Logf("%+v", session)
}

func TestGetSportResult(t *testing.T) {
	r, err := GetSportResult(session)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}

func TestSmartCreateRecords(t *testing.T) {
	records := SmartCreateRecords(fakeUserInfo, 5, time.Now())
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
