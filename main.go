package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	jkwx "github.com/inkedawn/go-sunshinemotion"
	"github.com/inkedawn/sunShineSportsHelper/utility"
)

const (
	defaultDistanceFemale = 2.5
	defaultDistanceMale   = 4.5

	displayTimePattern = "2006-01-02 15:04:05"
	inputTimePattern   = displayTimePattern
	defaultSchoolID    = 60
)

var closeLog = true

func init() {
	if //noinspection GoBoolExpressions
	closeLog {
		log.SetOutput(ioutil.Discard)
	}
	parseFlag()
}

func main() {
	switch cmd {
	case "help":
		printHelp()
		return
	case "login":
		loginAccount()
		return
	}

	// need session
	s, info := tryResume()
	if s == nil {
		fmt.Println("Need to login.")
		return
	}
	checkAppVer(s)
	showStatus(s, info)

	switch cmd {
	case "status":
		return
	case "getRoute":
		getRoute(s)
		return
	case "upload":
		uploadData(s, info)
		return
	case "uploadTest":
		uploadTestData(s)
		return
	case "testRule":
		getTestRule(s)
		return
	}
}

func tryResume() (*jkwx.Session, jkwx.UserInfo) {
	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Println("Resume failed.", err.Error())
			fmt.Println("** You can try to delete the session file ", getSessionFilePathById(cmdFlags.schoolID, cmdFlags.user))
		}
	}()
	s, info, err := readSession(cmdFlags.schoolID, cmdFlags.user)
	if err == ErrSessionNotExist {
		return nil, jkwx.UserInfo{}
	}
	if err != nil {
		panic(err)
	}
	if s != nil {
		fmt.Println("Use Existent Session.")
		fmt.Println("UserAgent", s.Device.UserAgent)
		fmt.Println("ExpiredTime", s.Token.ExpirationTime.Format(displayTimePattern))
		fmt.Println()
	}
	return s, info
}

func printHelp() {
	flag.Usage()
	os.Exit(0)
}

func checkAppVer(s *jkwx.Session) {
	appInfo, err := s.GetAppInfo()
	if err != nil {
		fmt.Println("Warning: 获取APP版本失败")
		return
	}
	fmt.Println("Latest Ver: ", appInfo.VerNumber)
	fmt.Println("Lib Ver: ", jkwx.AppVersionID)
	if appInfo.VerNumber > jkwx.AppVersionID {
		fmt.Println("需要更新")
		os.Exit(1)
	}
	return
}

func loginAccount() {
	s, info, err := readSession(cmdFlags.schoolID, cmdFlags.user)
	if err != nil && err != ErrSessionNotExist {
		panic(err)
	}
	if s != nil {
		// old user
		fmt.Println("Already Login. Use Old Session(Keep Same Device).")
		info, err = s.Login(cmdFlags.schoolID, cmdFlags.user, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.password))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		// new user
		s = jkwx.CreateSession()
		info, err = s.Login(cmdFlags.schoolID, cmdFlags.user, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.password))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		s.Device.UserAgent = utility.GetRandUserAgent()
	}
	fmt.Printf("Device: %+v\n", *s.Device)
	fmt.Printf("Token: %+v\n", *s.Token)
	saveSession(s, info)
	showStatus(s, info)
}
func showStatus(s *jkwx.Session, info jkwx.UserInfo) {
	r, err := s.GetSportResult()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("-----------")
	fmt.Println("| 帐号信息 |")
	fmt.Println("-----------")
	fmt.Println("ID：\t", s.User.UserID)
	fmt.Println("SchoolID：\t", s.User.SchoolID)
	fmt.Println("StuNum：\t", s.User.StuNum)
	fmt.Println("-----------")
	fmt.Println("班级：\t", info.ClassName)
	fmt.Println("学号：\t", info.StudentNumber)
	fmt.Println("姓名：\t", info.StudentName)
	fmt.Println("性别：\t", info.Sex)
	fmt.Println("-----------")
	fmt.Printf("LastTime：\t%s \n", r.LastTime.Format(displayTimePattern))
	fmt.Printf("已跑距离：\t%05.3f 公里\n", r.ActualDistance)
	fmt.Printf("达标距离：\t%05.3f 公里\n", r.QualifiedDistance)
	fmt.Println("-----------")
	// fmt.Printf("%+v", r)
}

func uploadData(s *jkwx.Session, info jkwx.UserInfo) {
	rawRecord := cmdFlags.rawRecord
	ignoreCompleted := cmdFlags.ignoreCompleted
	totalDistance := cmdFlags.distance
	if cmdFlags.distance == 0.0 {
		switch info.Sex {
		case "F":
			totalDistance = defaultDistanceFemale
		case "M":
			totalDistance = defaultDistanceMale
		default:
			log.Panicln("Unknown sex", info.Sex)
		}
	}
	endTime, err := time.Parse(inputTimePattern, cmdFlags.endTime)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	duration := cmdFlags.duration

	var records []jkwx.Record
	if !rawRecord {
		if !ignoreCompleted {
			r, err := s.GetSportResult()
			if err == nil && r.ActualDistance >= r.QualifiedDistance {
				fmt.Println("已达标，停止操作")

				return
			}
		}
		records = jkwx.SmartCreateRecordsBefore(s.User.SchoolID, s.User.UserID, jkwx.GetDefaultLimitParams(info.Sex), totalDistance, endTime)
	} else {
		records = []jkwx.Record{
			jkwx.CreateRecord(s.User.UserID, s.User.SchoolID, totalDistance, endTime, duration),
		}
	}

	if !confirm(records) {
		return
	}
	allErr := make([]error, len(records))
	for i, record := range records {
		allErr[i] = s.UploadRecord(record)
	}
	fmt.Println("---------------")
	fmt.Println("上传结果：")
	for _, err := range allErr {
		if err == nil {
			fmt.Println("OK.")
		} else {
			fmt.Println(err.Error())
		}
	}
	showStatus(s, info)
	if !ignoreCompleted {
		r, err := s.GetSportResult()
		if err == nil && r.ActualDistance >= r.QualifiedDistance {
			fmt.Println("已达标")
		}
	}
	fmt.Println("---------------")
}

func uploadTestData(s *jkwx.Session) {
	totalDistance := cmdFlags.distance
	duration := cmdFlags.duration
	endTime, err := time.Parse(inputTimePattern, cmdFlags.endTime)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	record := jkwx.CreateRecord(s.User.UserID, cmdFlags.schoolID, totalDistance, endTime, duration)
	if !confirm([]jkwx.Record{record}) {
		return
	}
	err = s.UploadTestRecord(record)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("上传成功")
	}

}

func confirm(records []jkwx.Record) bool {
	silent := cmdFlags.silent
	if silent {
		return true
	}

	fmt.Println("--------------")
	fmt.Println("| 确认上传数据 |")
	fmt.Println("---------------")
	for i, record := range records {
		distance := record.Distance
		duration := record.EndTime.Sub(record.BeginTime)
		v := record.Distance * 1000 / duration.Seconds()
		fmt.Println("第", i+1, "条")
		fmt.Println("起始时间：", record.BeginTime.Format(displayTimePattern))
		fmt.Println("结束时间：", record.EndTime.Format(displayTimePattern))
		fmt.Printf("用时%s内完成%.3f公里距离，速度约为%.2fm/s \n", duration, distance, v)
	}

	fmt.Println("请输入YES确认")
	var confirm string
	_, err := fmt.Scan(&confirm)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("---------------")
	return confirm == "YES"
}

func getRoute(s *jkwx.Session) {
	r, err := s.GetRandRoute()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(r)
	}
}

func getTestRule(s *jkwx.Session) {
	rule, err := s.GetTestRule()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(rule)
	}
}
