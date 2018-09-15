package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	UserID             int64
	TokenID            string
	UserExpirationTime time.Time
	UserInfo           UserInfo
	UserAgent          string
	LimitParams        *LimitParams
}
type httpError struct {
	msg     string
	httpErr error
}

func (e httpError) Error() string {
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
		return httpError{"HTTP Create Request Failed.", err}
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
		return httpError{"HTTP Send Request Failed! ", err}
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}

	var loginResult struct {
		Status             int64
		UserID             int64
		TokenID            string
		UserExpirationTime int64
		UserInfo           UserInfo
	}
	err = json.Unmarshal(respBytes, &loginResult)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	if loginResult.Status != 1 {
		return fmt.Errorf("resp status not ok. %d", loginResult.Status)
	}
	s.UserID, s.TokenID, s.UserExpirationTime, s.UserInfo = loginResult.UserID, loginResult.TokenID, time.Unix(loginResult.UserExpirationTime/1000, 0), loginResult.UserInfo
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
	return s.UploadData(record.Distance, record.BeginTime, record.EndTime, record.XTcode)
}
func (s *Session) UploadData(distance float64, beginTime time.Time, endTime time.Time, xtCode string) (e error) {
	req, err := http.NewRequest(http.MethodPost, uploadDataURL, strings.NewReader(url.Values{
		"results":   {toExchangeDistanceStr(distance)},
		"beginTime": {toExchangeTimeStr(beginTime)},
		"endTime":   {toExchangeTimeStr(endTime)},
		"isValid":   {"1"},
		"schoolId":  {schoolId},
		"xtCode":    {xtCode},
		"bz":        {""},
	}.Encode()))
	if err != nil {
		panic(httpError{"HTTP Create Request Failed.", err})
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("UserID", strconv.FormatInt(s.UserID, 10))
	req.Header.Set("TokenID", s.TokenID)
	req.Header.Set("app", "com.ccxyct.sunshinemotion")
	req.Header.Set("ver", "2.1.0")
	req.Header.Set("device", "Android,24,7.0")
	req.Header.Set("model", "Android")
	req.Header.Set("screen", "1080x1920")
	req.Header.Set("imei", "")
	req.Header.Set("imsi", "")
	req.Header.Set("crack", "0")
	req.Header.Set("latitude", "0.0")
	req.Header.Set("longitude", "0.0")
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(fmt.Errorf("HTTP Send Request Failed! %s", err.Error()))
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode)))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("HTTP Read Resp Failed! %s", err.Error()))
	}

	var uploadResult struct {
		Status       int
		ErrorMessage string
	}
	err = json.Unmarshal(respBytes, &uploadResult)
	if err != nil {
		panic(fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes)))
	}
	const successCode = 1
	if uploadResult.Status != successCode {
		return fmt.Errorf("server status %d , message: %s", uploadResult.Status, uploadResult.ErrorMessage)
	}
	return nil
}

type SportResult struct {
	LastTime  time.Time
	Qualified float64
	Distance  float64
}

func (s *Session) GetSportResult() (r *SportResult, e error) {
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return nil, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("UserID", strconv.FormatInt(s.UserID, 10))
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
	var httpSporstResult struct {
		Status       int
		ErrorMessage string
		LastTime     string  `json:"lastTime"`
		Qualified    float64 `json:"qualified"`
		Result       float64 `json:"result"`
		UserID       int64   `json:"userID"`
		Term         string  `json:"term"`
		Year         int     `json:"year"`
	}
	err = json.Unmarshal(respBytes, &httpSporstResult)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	const successCode = 1
	if httpSporstResult.Status != successCode {
		return nil, fmt.Errorf("server status %d , message: %s", httpSporstResult.Status, httpSporstResult.ErrorMessage)
	}
	r = new(SportResult)
	if httpSporstResult.LastTime != "" {
		r.LastTime, err = fromExchangeTimeStr(httpSporstResult.LastTime)
	} else {
		r.LastTime = time.Now()
	}
	if err != nil {
		log.Println(string(respBytes))
		panic(err)
	}
	r.Qualified = httpSporstResult.Qualified
	r.Distance = httpSporstResult.Result
	return r, nil
}
