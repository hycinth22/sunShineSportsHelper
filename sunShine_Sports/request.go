package sunShine_Sports

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

	"../utility"
)

type Session struct {
	UserID             int
	TokenID            string
	UserExpirationTime int64
	UserInfo           UserInfo
	UserAgent          string
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
	ua                  = DefaultUserAgent
)

func SetUserAgent(newUA string) {
	ua = newUA
}

func Login(stuNum string, phoneNum string, passwordHash string) (s *Session, e error) {
	s = &Session{UserID: 0, TokenID: "", UserAgent: DefaultUserAgent}

	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(url.Values{
		"stuNum":   {stuNum},
		"phoneNum": {phoneNum},
		"passWd":   {passwordHash},
		"schoolId": {schoolId},
		"stuId":    {"1"},
		"token":    {""},
	}.Encode()))
	if err != nil {
		return nil, HTTPError{"HTTP Create Request Failed.", err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ua)
	req.Header.Set("UserID", "0")
	req.Header.Set("crack", "0")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, HTTPError{"HTTP Send Request Failed! ", err}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
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
		return nil, fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	if respMsg.Status == 0 {
		return nil, ErrIncorrectAccount
	}
	if respMsg.Status != 1 {
		return nil, fmt.Errorf("resp status not ok. %d", respMsg.Status)
	}
	s.UserID, s.TokenID, s.UserExpirationTime, s.UserInfo = respMsg.UserID, respMsg.TokenID, respMsg.UserExpirationTime, respMsg.UserInfo
	return s, nil
}

type Record struct {
	Distance  float64
	BeginTime time.Time
	EndTime   time.Time
}

func CreateRecords(userInfo UserInfo, distance float64, beforeTime time.Time) []Record {
	records := make([]Record, 0, int(distance/3))
	remain := distance
	lastBeginTime := beforeTime
	for remain > 0 {
		var singleDistance float64
		// 距离随机化
		//distanceRandomRatio :=  float64(utility.RandRange(9500, 11142))/10000 // 距离波动化比例 95%-111.42%
		// 范围取随机
		var minDistance float64
		var maxDistance float64
		switch userInfo.Sex {
		case "F":
			minDistance = 2.09
			maxDistance = 2.9
		case "M":
			minDistance = 2.59
			maxDistance = 3.9
		default:
			panic("Unknown Sex" + userInfo.Sex)
		}
		if remain < minDistance{
			break
		}
		if remain > maxDistance{
			singleDistance = float64(utility.RandRange(int(minDistance *1000), int(maxDistance*1000))) / 1000 // 取主要距离
		}else{
			singleDistance = remain
		}

		singleDistance += float64(utility.RandRange(-99999, 99999)) / 1000000 // 小数部分随机化 -0.09 ~ 0.09

		var randomDuration time.Duration
		// 时间间隔随机化
		// 参数设定：min>minDis*3, max<maxDis*10
		var minMinuteDuration int
		var maxMinuteDuration int
		switch userInfo.Sex {
		case "F":
			// 11-20min
			minMinuteDuration = 11
			maxMinuteDuration = 20
		case "M":
			// 14-25min
			minMinuteDuration = 14
			maxMinuteDuration = 25
		default:
			panic("Unknown Sex" + userInfo.Sex)
		}
		randomDuration = time.Duration(utility.RandRange(minMinuteDuration, maxMinuteDuration)) * time.Minute

		randomDuration += time.Duration(utility.RandRange(0, 60)) * time.Second // 时间间隔秒级随机化
		endTime := lastBeginTime.Add(-time.Duration(utility.RandRange(1, 10)) * time.Minute)
		beginTime := endTime.Add(-randomDuration)

		records = append(records, Record{
			Distance:  singleDistance,
			BeginTime: beginTime,
			EndTime:   endTime,
		})

		remain -= singleDistance
		lastBeginTime = beginTime
	}
	nRecord := len(records)
	reverse := make([]Record, nRecord)
	for i := 0; i < nRecord; i++ {
		reverse[i] = records[nRecord-i-1]
	}
	return reverse
}

func UploadRecord(session *Session, record Record) (status int, e error) {
	return UploadData(session, record.Distance, record.BeginTime, record.EndTime)
}

func UploadData(session *Session, distance float64, beginTime time.Time, endTime time.Time) (status int, e error) {
	const timePattern = "2006-01-02 15:04:05"
	req, err := http.NewRequest(http.MethodPost, uploadDataURL, strings.NewReader(url.Values{
		"results":   {fmt.Sprintf("%07.6f", distance)},
		"beginTime": {beginTime.Format(timePattern)},
		"endTime":   {endTime.Format(timePattern)},
		"isValid":   {"1"},
		"schoolId":  {schoolId},
		"bz":        {""},
	}.Encode()))
	if err != nil {
		return -1, HTTPError{"HTTP Create Request Failed.", err}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ua)
	req.Header.Set("UserID", strconv.Itoa(session.UserID))
	req.Header.Set("TokenID", session.TokenID)
	req.Header.Set("crack", "0")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return -1, fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}

	var respMsg struct {
		Status int
	}
	err = json.Unmarshal(respBytes, &respMsg)
	if err != nil {
		return -1, fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}

	return respMsg.Status, nil
}

type SportResult struct {
	Distance  float64 `json:"result"`
	LastTime  string  `json:"lastTime`
	Year      int     `json:"year`
	Qualified float64 `json:"qualified`
}

func GetSportResult(session *Session) (r *SportResult, e error) {
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return nil, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ua)
	req.Header.Set("UserID", strconv.Itoa(session.UserID))
	req.Header.Set("TokenID", session.TokenID)
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
