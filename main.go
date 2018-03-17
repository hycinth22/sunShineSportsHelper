package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"os"
	"time"

	jkwx "./sunShine_Sports"
	"./utility"
)

const timePattern = "2006-01-02 15:04:05"

var cmdFlags struct {
	help bool

	silent bool

	login    bool
	user     string
	password string

	status bool

	upload   bool
	distance float64
	duration time.Duration
}

func init() {
	flag.BoolVar(&cmdFlags.help, "h", false, "this help")
	flag.BoolVar(&cmdFlags.silent, "q", false, "quiet mode")

	flag.BoolVar(&cmdFlags.login, "login", false, "login into account")
	flag.StringVar(&cmdFlags.user, "u", defaultStuNum, "account(stuNum)")
	flag.StringVar(&cmdFlags.password, "p", "", "password")

	flag.BoolVar(&cmdFlags.status, "status", false, "view account status")

	flag.BoolVar(&cmdFlags.upload, "upload", false, "upload sport data")

	flag.Float64Var(&cmdFlags.distance, "distance", 3.500000, "distance(精确到小数点后6位)")

	randomDuration := time.Duration(utility.RandRange(12, 20))*time.Minute
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
		s := readSessionById(cmdFlags.user)
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

func printHelp() {
	flag.Usage()
	os.Exit(0)
}
func loginAccount() {
	var s *jkwx.Session
	s = readSessionById(cmdFlags.user)
	if s != nil{
		fmt.Println("Alread Login.")
		return
	}else{
		var err error
		s, err = jkwx.Login(cmdFlags.user, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.password))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	saveSession(s)
	showStatus(s)
}
func showStatus(s *jkwx.Session) {
	// TODO
	r, err := jkwx.GetSportResult(s)
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
	fmt.Println("-----------")
	fmt.Printf("LastTime：\t%s \n", r.LastTime)
	fmt.Printf("已跑距离：\t%07.6f 公里\n", r.Distance)
	fmt.Printf("达标距离：\t%07.6f 公里\n", r.Qualified)
	fmt.Println("-----------")
	// fmt.Printf("%+v", r)
}
func uploadData(s *jkwx.Session) {
	totalDistance := cmdFlags.distance
	records := jkwx.CreateRecords(s.UserInfo, totalDistance, time.Now())

	fmt.Println("--------------")
	fmt.Println("| 确认上传数据 |")
	fmt.Println("---------------")
	for i, record := range records{
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
		fmt.Scanf("%s", &confirm)
		fmt.Println("---------------")
		if confirm != "YES" {
			return
		}
	}

	allStatus := make([]int, len(records))
	for i, record := range records{
		var err error
		allStatus[i], err = jkwx.UploadRecord(s, record)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	fmt.Println("---------------")
	fmt.Println("上传结果：")
	hasError := false
	for _, status := range allStatus{
		if status == 1 {
			fmt.Println("OK.")
		} else {
			fmt.Printf("Status %d", s)
			hasError = true
		}
	}
	if !hasError{
		showStatus(s)
	}
	fmt.Println("---------------")
}