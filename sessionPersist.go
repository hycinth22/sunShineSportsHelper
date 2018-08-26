package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"time"

	jkwx "inkedawn/sunShineSportsHelper/lib"
	"inkedawn/sunShineSportsHelper/utility"
)

const sessionFileFormat = "sunShine_Sports_%s.session"

var (
	ErrSessionNotExist = errors.New("session does not exist")
)

func getSessionFilePath(s *jkwx.Session) string {
	return getSessionFilePathById(s.UserInfo.StudentNumber)
}

func getSessionFilePathById(stuNum string) string {
	return fmt.Sprintf(sessionFileFormat, stuNum)
}

func saveSession(s *jkwx.Session) {
	if s == nil {
		panic("try to save nil session")
	}
	f, err := os.Create(getSessionFilePath(s))
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		panic(err)
	}
	// TODO: 数据文件版本号
	if err := gob.NewEncoder(f).Encode(s); err != nil {
		panic(err)
	}
}

func readSession(stuNu string) (*jkwx.Session, error) {
	f, err := os.Open(getSessionFilePathById(stuNu))
	if f != nil {
		defer f.Close()
	}
	if os.IsNotExist(err) {
		return nil, ErrSessionNotExist
	}
	if err != nil {
		panic(err)
	}
	s := new(jkwx.Session)
	// TODO: 数据文件版本号
	if err := gob.NewDecoder(f).Decode(s); err != nil {
		panic(err)
	}
	if s.UserAgent == "" {
		fmt.Println("Upgrade session file from old version (before 2.0)")
		fmt.Println("Add UserAgent")
		s.UserAgent = utility.GetRandUserAgent()
		saveSession(s)
	}
	s.UpdateLimitParams()

	if time.Now().After(s.UserExpirationTime) {
		removeSession(stuNu)
		fmt.Println("Login Expired.")
		return nil, ErrSessionNotExist
	}
	return s, nil
}

func removeSession(stuNum string) error {
	err := os.Remove(getSessionFilePathById(stuNum))
	if os.IsNotExist(err) {
		return ErrSessionNotExist
	}
	return err
}
