package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	UserID             int
	TokenID            string
	UserExpirationTime time.Time
	UserInfo           UserInfo
	UserAgent          string
	LimitParams        *LimitParams
}
type LimitParams struct {
	// 随机区间（生成记录随机的单次距离区间）
	RandDistance Float64Range
	// 限制区间（目标系统限制的单次距离区间）
	LimitSingleDistance Float64Range
	// 限制区间（目标系统限制的总距离区间）
	LimitTotalDistance Float64Range
	// 每条记录的时间区间
	MinuteDuration IntRange
}

type UserInfo struct {
	Id            int    `json:"id"`
	InClassID     int    `json:"inClassID"`
	InClassName   string `json:"inClassName"`
	InCollegeID   int    `json:"inCollegeID"`
	InCollegeName string `json:"inCollegeName"`
	IsTeacher     int    `json:"isTeacher"`
	NickName      string `json:"nickName"`
	PhoneNumber   string `json:"phoneNumber"`
	Sex           string `json:"sex"`
	StudentName   string `json:"studentName"`
	StudentNumber string `json:"studentNumber"`
	// UserRoleID    int
}
type Float64Range struct {
	Min float64
	Max float64
}
type IntRange struct {
	Min int
	Max int
}
type HTTPError struct {
	msg     string
	httpErr error
}

func (e HTTPError) Error() string {
	return e.msg + "\n" + e.httpErr.Error()
}

const (
	server            = "http://www.ccxyct.com:8080"
	loginURL          = server + "/sunShine_Sports/loginSport.action"
	uploadDataURL     = server + "/sunShine_Sports/xtUploadData.action"
	getSportResultURL = server + "/sunShine_Sports/xtGetSportResult.action"
	DefaultUserAgent  = "Dalvik/2.1.0 (Linux; U; Android 7.0)"

	schoolId = "60"
)

var (
	ErrIncorrectAccount = errors.New("account or password is INCORRECT")
	ErrIllegalData      = errors.New("illegal data")
)

func CreateSession() *Session {
	return &Session{UserID: 0, TokenID: "", UserAgent: DefaultUserAgent}
}

func (s *Session) Login(stuNum string, phoneNum string, passwordHash string) (e error) {
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(url.Values{
		"stuNum":   {stuNum},
		"phoneNum": {phoneNum},
		"passWd":   {passwordHash},
		"schoolId": {schoolId},
		"stuId":    {"1"},
		"token":    {""},
	}.Encode()))
	if err != nil {
		return HTTPError{"HTTP Create Request Failed.", err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("UserID", "0")
	req.Header.Set("crack", "0")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return HTTPError{"HTTP Send Request Failed! ", err}
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}

	var respMsg struct {
		Status             int
		Date               string
		UserInfo           UserInfo
		TokenID            string
		UserExpirationTime int64
		UserID             int
	}
	err = json.Unmarshal(respBytes, &respMsg)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	if respMsg.Status == 0 {
		return ErrIncorrectAccount
	}
	if respMsg.Status != 1 {
		return fmt.Errorf("resp status not ok. %d", respMsg.Status)
	}
	s.UserID, s.TokenID, s.UserExpirationTime, s.UserInfo = respMsg.UserID, respMsg.TokenID, time.Unix(respMsg.UserExpirationTime/1000, 0), respMsg.UserInfo
	s.UpdateLimitParams()
	return nil
}

func (s *Session) UpdateLimitParams() {
	// 参数设定：
	// MinuteDuration: min>minDis*3, max<maxDis*10
	switch s.UserInfo.Sex {
	case "F":
		s.LimitParams = &LimitParams{
			RandDistance:        Float64Range{2.0, 3.0},
			LimitSingleDistance: Float64Range{1.0, 3.0},
			LimitTotalDistance:  Float64Range{1.0, 3.0},
			MinuteDuration:      IntRange{11, 20},
		}
	case "M":
		s.LimitParams = &LimitParams{
			RandDistance:        Float64Range{2.6, 4.0},
			LimitSingleDistance: Float64Range{2.0, 4.0},
			LimitTotalDistance:  Float64Range{2.0, 5.0},
			MinuteDuration:      IntRange{14, 25},
		}

	default:
		panic("Unknown Sex" + s.UserInfo.Sex)
	}
}
func (s *Session) UploadRecord(record Record) (e error) {
	return s.UploadData(record.Distance, record.BeginTime, record.EndTime)
}

func (s *Session) UploadData(distance float64, beginTime time.Time, endTime time.Time) (e error) {
	req, err := http.NewRequest(http.MethodPost, uploadDataURL, strings.NewReader(url.Values{
		"results":   {fmt.Sprintf("%07.6f", distance)},
		"beginTime": {getTimeStr(beginTime)},
		"endTime":   {getTimeStr(endTime)},
		"isValid":   {"1"},
		"schoolId":  {schoolId},
		"xtCode":    {GetXtcode(s.UserInfo.Id, getTimeStr(beginTime))},
		"bz":        {""},
	}.Encode()))
	if err != nil {
		return HTTPError{"HTTP Create Request Failed.", err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("UserID", strconv.Itoa(s.UserID))
	req.Header.Set("TokenID", s.TokenID)
	req.Header.Set("app", "com.ccxyct.sunshinemotion")
	req.Header.Set("ver", "2.0.1")
	req.Header.Set("device", "Android,24,7.0")
	req.Header.Set("model", "Android")
	req.Header.Set("screen", "1080x1920")
	//req.Header.Set("imei", "000000000000000")
	//req.Header.Set("imsi", "000000000000000")
	req.Header.Set("crack", "0")
	req.Header.Set("latitude", "0.0")
	req.Header.Set("longitude", "0.0")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}

	var respMsg struct {
		Status int
	}
	err = json.Unmarshal(respBytes, &respMsg)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	switch respMsg.Status {
	case 10001:
		return ErrIllegalData
	case 1:
		return nil // success
	default:
		return fmt.Errorf("server return unknown status %d ", respMsg.Status)
	}
}

type SportResult struct {
	Distance  float64 `json:"result"`
	LastTime  string  `json:"lastTime"`
	Year      int     `json:"year"`
	Qualified float64 `json:"qualified"`
}

func (s *Session) GetSportResult() (r *SportResult, e error) {
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return nil, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("UserID", strconv.Itoa(s.UserID))
	req.Header.Set("TokenID", s.TokenID)
	req.Header.Set("crack", "0")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}
	var respMsg SportResult
	err = json.Unmarshal(respBytes, &respMsg)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	return &respMsg, nil
}
