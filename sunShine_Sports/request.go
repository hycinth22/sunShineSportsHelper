package sunShine_Sports

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

	DistanceLimit *DistanceParams
}
type DistanceParams struct {
	// 随机区间（生成记录随机的单次距离区间）
	RandDistance Float64Range
	// 限制区间（目标系统限制的单次距离区间）
	LimitSingleDistance Float64Range
	// 限制区间（目标系统限制的总距离区间）
	LimitTotalDistance Float64Range
}

type Float64Range struct {
	Min float64
	Max float64
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
	s.UserInfo.DistanceLimit = GetDistanceParams(s)
	return s, nil
}

type Record struct {
	Distance  float64
	BeginTime time.Time
	EndTime   time.Time
}

func CreateRecord(distance float64, beforeTime time.Time, duration time.Duration) Record {
	return Record{Distance: distance,
		BeginTime: beforeTime.Add(-duration),
		EndTime:   beforeTime,
	}
}
func SmartCreateRecords(userInfo UserInfo, distance float64, beforeTime time.Time) []Record {
	records := make([]Record, 0, int(distance/3))
	remain := distance
	lastBeginTime := beforeTime
	for remain > 0 {
		var singleDistance float64
		// 范围取随机
		if remain > userInfo.DistanceLimit.RandDistance.Max {
			// 检查是否下一条可能丢弃较大的距离
			// 防止：剩下比较多，但却不满足最小限制距离，不能生成下一条记录
			if remain-userInfo.DistanceLimit.RandDistance.Max > userInfo.DistanceLimit.LimitSingleDistance.Min {
				// 正常取随机值
				singleDistance = float64(utility.RandRange(int(userInfo.DistanceLimit.RandDistance.Min*1000), int(userInfo.DistanceLimit.RandDistance.Max*1000))) / 1000
			} else {
				// 随机选择本条为最小限制距离，或者为下一条预留最小限制距离
				singleDistance = []float64{userInfo.DistanceLimit.LimitSingleDistance.Min, remain - userInfo.DistanceLimit.LimitSingleDistance.Min}[utility.RandRange(0, 1)]
			}
		} else if remain >= userInfo.DistanceLimit.LimitSingleDistance.Min && remain <= userInfo.DistanceLimit.LimitSingleDistance.Max {
			// 剩余的符合限制区间，直接使用剩余的生成最后一条记录
			singleDistance = remain
		} else {
			// 剩余较多，但不符合限制区间无法再生成一条合法记录，输出提醒
			if remain > 0.5 {
				fmt.Println("提醒：由于随机原则与区间限制的冲突，丢弃了较大的距离", remain, "公里，考虑重新设定距离值。")
			}
			break
		}

		// 小数部分随机化 -0.09 ~ 0.09
		tinyPart := float64(utility.RandRange(0, 99999)) / 1000000
		switch r := singleDistance + tinyPart; {
		case r < userInfo.DistanceLimit.LimitSingleDistance.Min:
			singleDistance = userInfo.DistanceLimit.LimitSingleDistance.Min
			/*case r > userInfo.DistanceLimit.LimitSingleDistance.Max:
			singleDistance = userInfo.DistanceLimit.LimitSingleDistance.Max
			*/
		default:
			singleDistance += tinyPart
		}

		// 检测结果合法性，由于TinyPart允许上下浮动0.1
		if singleDistance < userInfo.DistanceLimit.LimitSingleDistance.Min-0.1 || singleDistance > userInfo.DistanceLimit.LimitSingleDistance.Max+0.1 {
			// 丢弃不合法距离
			log.Println("Drop distance: ", singleDistance)
			continue
		}

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

		endTime := lastBeginTime.Add(-time.Duration(utility.RandRange(1, 30)) * time.Minute)
		endTime = endTime.Add(-time.Duration(utility.RandRange(1, 60)) * time.Second)
		beginTime := endTime.Add(-randomDuration)

		records = append(records, Record{
			Distance:  singleDistance,
			BeginTime: beginTime,
			EndTime:   endTime,
		})

		remain -= singleDistance - tinyPart
		lastBeginTime = beginTime
	}
	nRecord := len(records)
	reverse := make([]Record, nRecord)
	for i := 0; i < nRecord; i++ {
		reverse[i] = records[nRecord-i-1]
	}
	return reverse
}

func GetDistanceParams(s *Session) *DistanceParams {
	switch s.UserInfo.Sex {
	case "F":
		return &DistanceParams{
			RandDistance:        Float64Range{2.0, 3.0},
			LimitSingleDistance: Float64Range{1.0, 3.0},
			LimitTotalDistance:  Float64Range{1.0, 3.0},
		}
	case "M":
		return &DistanceParams{
			RandDistance:        Float64Range{2.6, 4.0},
			LimitSingleDistance: Float64Range{2.0, 4.0},
			LimitTotalDistance:  Float64Range{2.0, 5.0},
		}
	default:
		panic("Unknown Sex" + s.UserInfo.Sex)
	}
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
	LastTime  string  `json:"lastTime"`
	Year      int     `json:"year"`
	Qualified float64 `json:"qualified"`
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
