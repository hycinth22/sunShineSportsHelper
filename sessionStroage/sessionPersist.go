package sessionStroage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"time"

	jkwx "github.com/inkedawn/go-sunshinemotion/v3"
)

const sessionFileFormat = "sunShine_Sports_%d_%s.session"

var (
	ErrSessionNotExist = errors.New("session does not exist")
)

type persistStruct struct {
	Session *jkwx.Session
	Info    jkwx.UserInfo
}

func getSessionFilePath(s *jkwx.Session) string {
	return GetSessionFilePathById(s.User.SchoolID, s.User.StuNum)
}

func GetSessionFilePathById(schoolID int64, stuNum string) string {
	return fmt.Sprintf(sessionFileFormat, schoolID, stuNum)
}

func SaveSession(s *jkwx.Session, info jkwx.UserInfo) {
	if s == nil {
		panic("try to save nil session")
	}
	f, err := os.Create(getSessionFilePath(s))
	if f != nil {
		defer func() {
			_ = f.Close()
		}()
	}
	if err != nil {
		panic(err)
	}
	// TODO: 数据文件版本号
	if err := gob.NewEncoder(f).Encode(persistStruct{
		Session: s,
		Info:    info,
	}); err != nil {
		panic(err)
	}
}

func ReadSession(schoolID int64, stuNum string) (*jkwx.Session, jkwx.UserInfo, error) {
	f, err := os.Open(GetSessionFilePathById(schoolID, stuNum))
	if f != nil {
		defer func() {
			_ = f.Close()
		}()
	}
	if os.IsNotExist(err) {
		return nil, jkwx.UserInfo{}, ErrSessionNotExist
	}
	if err != nil {
		panic(err)
	}
	s := persistStruct{}
	// TODO: 数据文件版本号
	if err := gob.NewDecoder(f).Decode(&s); err != nil {
		panic(err)
	}
	session := s.Session
	info := s.Info
	if time.Now().After(session.Token.ExpirationTime) {
		_ = removeSession(schoolID, stuNum)
		fmt.Println("Login Expired.")
		return nil, jkwx.UserInfo{}, ErrSessionNotExist
	}
	return session, info, nil
}

func removeSession(schoolID int64, stuNum string) error {
	err := os.Remove(GetSessionFilePathById(schoolID, stuNum))
	if os.IsNotExist(err) {
		return ErrSessionNotExist
	}
	return err
}
