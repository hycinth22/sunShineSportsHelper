package sunShine_Sports

import (
	"crypto/md5"
	"fmt"
	"testing"
)

var session *Session
var loginErr error

func init(){
	session, loginErr = Login("021640302", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))))
}
func TestLogin(t *testing.T) {
	// 已在init()中执行函数功能函数，此处仅检测结果
	if loginErr != nil{
		t.Fatalf("%v", loginErr)
	}
	t.Logf("%+v", session)
}

func TestGetSportResult(t *testing.T) {
	r, err := GetSportResult(session)
	if err != nil{
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}