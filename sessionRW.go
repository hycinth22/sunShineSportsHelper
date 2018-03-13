package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	jkwx "./sunShine_Sports"
)

const sessionFileFormat = "sunShine_Sports_%s.session"
const defaultStuNum = "default" //用于未输入user时的默认参数名

func getSessionFilePath(s *jkwx.Session) string {
	return getSessionFilePathById(s.UserInfo.StudentNumber)
}
func getSessionFilePathById(stuNum string) string {
	return fmt.Sprintf(sessionFileFormat, stuNum)
}
func saveSession(s *jkwx.Session) {
	f, err := os.Create(getSessionFilePath(s))
	if err != nil {
		panic(err)
	}
	if err := gob.NewEncoder(f).Encode(s); err != nil {
		panic(err)
	}
}
func readSession() *jkwx.Session {
	return readSessionById(defaultStuNum)
}
func readSessionById(stuNu string) *jkwx.Session {
	f, err := os.Open(getSessionFilePathById(stuNu))
	var s jkwx.Session
	if err != nil {
		return nil
	}
	if err := gob.NewDecoder(f).Decode(&s); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	nowTime := time.Now()
	expiredTime := time.Unix(s.UserExpirationTime/1000, 0)
	fmt.Println("Alread Login.")
	fmt.Println("nowTime", nowTime.Format(timePattern))
	fmt.Println("expiredTime", expiredTime.Format(timePattern))
	fmt.Println()
	if nowTime.After(expiredTime) {
		fmt.Println("Login Expired.")
		return nil
	}
	return &s
}