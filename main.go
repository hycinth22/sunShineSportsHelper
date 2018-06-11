package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	jkwx "./lib"
	"./utility"
)

const timePattern = "2006-01-02 15:04:05"

var cmdFlags struct {
	help bool

	silent bool

	login      bool
	forceLogin bool
	user       string
	password   string

	status bool

	upload          bool
	rawRecord       bool
	ignoreCompleted bool
	distance        float64
	duration        time.Duration
}

const (
	defaultDistanceFemale = 2.5
	defaultDistanceMale   = 4.5
)

func init() {
	flag.BoolVar(&cmdFlags.help, "h", false, "this help")
	flag.BoolVar(&cmdFlags.silent, "q", false, "quiet mode")

	flag.BoolVar(&cmdFlags.login, "login", false, "login into account")
	flag.BoolVar(&cmdFlags.forceLogin, "forceLogin", false, "login into account(not use existent session)")
	flag.StringVar(&cmdFlags.user, "u", "default", "account(stuNum)")
	flag.StringVar(&cmdFlags.password, "p", "", "password")

	flag.BoolVar(&cmdFlags.status, "status", false, "view account status")

	flag.BoolVar(&cmdFlags.upload, "upload", false, "upload sport data")
	flag.BoolVar(&cmdFlags.rawRecord, "rawRecord", false, "upload rawRecord sport data")
	flag.BoolVar(&cmdFlags.ignoreCompleted, "ignoreCompleted", false, "continue to upload though completed")
	flag.Float64Var(&cmdFlags.distance, "distance", 0.0, "distance(精确到小数点后6位).")

	randomDuration := time.Duration(utility.RandRange(12, 20)) * time.Minute
	flag.DurationVar(&cmdFlags.duration, "duration", randomDuration, "time duration")
	flag.Parse()
}

func main() {
	switch {
	case cmdFlags.help:
		printHelp()
	case cmdFlags.login:
		loginAccount()
	default:
		// need session
		s := tryResume()
		if s == nil {
			fmt.Println("Need to login.")
			return
		}
		showStatus(s)
		switch {
		case cmdFlags.upload:
			uploadData(s)
		}
	}
}

func tryResume() *jkwx.Session {
	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Println("Reusme failed.", err.Error())
			fmt.Println("** You can try to delete the session file ", getSessionFilePathById(cmdFlags.user))
		}
	}()
	s, _ := readSession(cmdFlags.user)
	if s != nil {
		fmt.Println("Use Existent Session.")
		fmt.Println("UserAgent", s.UserAgent)
		fmt.Println("expiredTime", s.UserExpirationTime.Format(timePattern))
		fmt.Println()
	}
	return s
}

func printHelp() {
	flag.Usage()
	os.Exit(0)
}
func loginAccount() {
	s, _ := readSession(cmdFlags.user)
	if s != nil {
		fmt.Println("Alread Login.")
		return
	} else {
		s = jkwx.CreateSession()
		err := s.Login(cmdFlags.user, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.password))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		s.UserAgent = utility.GetRandUserAgent()
	}
	saveSession(s)
	showStatus(s)
}
func showStatus(s *jkwx.Session) {
	// TODO
	r, err := s.GetSportResult()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("-----------")
	fmt.Println("| 帐号信息 |")
	fmt.Println("-----------")
	fmt.Println("ID：\t", s.UserInfo.Id)
	fmt.Println("班级：\t", s.UserInfo.InClassName)
	fmt.Println("学号：\t", s.UserInfo.StudentNumber)
	fmt.Println("姓名：\t", s.UserInfo.StudentName)
	fmt.Println("性别：\t", s.UserInfo.Sex)
	fmt.Println("-----------")
	fmt.Printf("LastTime：\t%s \n", r.LastTime)
	fmt.Printf("已跑距离：\t%07.6f 公里\n", r.Distance)
	fmt.Printf("达标距离：\t%07.6f 公里\n", r.Qualified)
	fmt.Println("-----------")
	// fmt.Printf("%+v", r)
}
func uploadData(s *jkwx.Session) {
	totalDistance := cmdFlags.distance
	if cmdFlags.distance == 0.0 {
		switch s.UserInfo.Sex {
		case "F":
			totalDistance = defaultDistanceFemale
		case "M":
			totalDistance = defaultDistanceMale
		default:
			log.Panicln("Unknown sex", s.UserInfo.Sex)
		}
	}
	ignoreCompleted := cmdFlags.ignoreCompleted

	var records []jkwx.Record
	if !cmdFlags.rawRecord {
		if totalDistance < s.LimitParams.LimitTotalDistance.Min || totalDistance > s.LimitParams.LimitTotalDistance.Max {
			fmt.Printf("超出限制的总距离（%f - %f）\n", s.LimitParams.LimitTotalDistance.Min, s.LimitParams.LimitTotalDistance.Max)
			return
		}

		if !ignoreCompleted {
			r, err := s.GetSportResult()
			if err == nil && r.Distance >= r.Qualified {
				fmt.Println("已达标，停止操作")

				return
			}
		}
		records = jkwx.SmartCreateRecords(s.LimitParams, totalDistance, time.Now())
	} else {
		records = []jkwx.Record{
			jkwx.CreateRecord(totalDistance, time.Now(), cmdFlags.duration),
		}
	}

	fmt.Println("--------------")
	fmt.Println("| 确认上传数据 |")
	fmt.Println("---------------")
	for i, record := range records {
		distance := record.Distance
		duration := record.EndTime.Sub(record.BeginTime)
		v := record.Distance * 1000 / duration.Seconds()
		fmt.Println("第", i+1, "条")
		fmt.Println("起始时间：", record.BeginTime.Format(timePattern))
		fmt.Println("结束时间：", record.EndTime.Format(timePattern))
		fmt.Printf("用时%s内完成%.6f公里距离，速度约为%.2fm/s \n", duration, distance, v)
	}

	if !cmdFlags.silent {
		fmt.Println("请输入YES确认")
		var confirm string
		fmt.Scan(&confirm)
		fmt.Println("---------------")
		if confirm != "YES" {
			return
		}
	}

	allErr := make([]error, len(records))
	for i, record := range records {
		allErr[i] = s.UploadRecord(record)
	}

	fmt.Println("---------------")
	fmt.Println("上传结果：")
	for _, err := range allErr {
		switch err {
		case nil:
			fmt.Println("OK.")
			showStatus(s)
			if !ignoreCompleted {
				r, err := s.GetSportResult()
				if err == nil && r.Distance > r.Qualified {
					fmt.Println("已达标")
				}
			}
		case jkwx.ErrIllegalData:
			fmt.Println("非法数据.")
			fallthrough
		default:
			fmt.Println(err.Error())
		}
	}
	fmt.Println("---------------")
}
