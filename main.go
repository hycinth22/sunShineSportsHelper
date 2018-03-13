package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"os"
	"time"

	jkwx "./sunShine_Sports"
)

const timePattern = "2006-01-02 15:04:05"

var cmdFlags struct {
	help bool

	login    bool
	user     string
	password string

	status bool

	upload   bool
	distance float64
	duration time.Duration
}

func init() {
	flag.BoolVar(&cmdFlags.help, "help", false, "this help")

	flag.BoolVar(&cmdFlags.login, "login", false, "login into account")
	flag.StringVar(&cmdFlags.user, "user", "", "account(stuNum)")
	flag.StringVar(&cmdFlags.password, "password", "", "password")

	flag.BoolVar(&cmdFlags.status, "status", false, "view account status")

	flag.BoolVar(&cmdFlags.upload, "upload", false, "upload sport data")
	flag.Float64Var(&cmdFlags.distance, "distance", 3.000000 * (float64(randRange(8000, 12000))/10000), "distance(精确到小数点后6位)")
	flag.DurationVar(&cmdFlags.duration, "duration", time.Duration(randRange(15, 25))*time.Minute+time.Duration(randRange(0, 60))*time.Second, "time duration")
}

func main() {
	flag.Parse()
	switch {
	case cmdFlags.help:
		help()
	case cmdFlags.login:
		login()
	default:
		// need session
		s := readSession()
		if s == nil {
			fmt.Println("Need to login.")
			return
		}
		switch {
		case cmdFlags.status:
			status(s)
		case cmdFlags.upload:
			upload(s)
		}
	}
}

func help() {
	flag.Usage()
	os.Exit(0)
}
func login() {
	s, err := jkwx.Login(cmdFlags.user, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.password))))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	saveSession(s)
}
func status(s *jkwx.Session) {
	// TODO
	r, err := jkwx.GetSportResult(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("ID：\t", s.UserInfo.Id)
	fmt.Println("班级：\t", s.UserInfo.InClassName)
	fmt.Println("学号：\t", s.UserInfo.StudentNumber)
	fmt.Println("姓名：\t", s.UserInfo.StudentName)
	fmt.Println("-----------")
	fmt.Printf("LastTime：\t%s \n", r.LastTime)
	fmt.Printf("已跑距离：\t%07.6f 公里\n", r.Distance)
	fmt.Printf("达标距离：\t%07.6f 公里\n", r.Qualified)
	// fmt.Printf("%+v", r)
}
func upload(s *jkwx.Session) {
	distance := cmdFlags.distance
	duration := cmdFlags.duration
	endTime := time.Now().Add(-time.Duration(randRange(1, 10)) * time.Minute)
	beginTime := endTime.Add(-duration)
	v := cmdFlags.distance * 1000 / duration.Seconds()
	fmt.Println("起始时间：", beginTime.Format(timePattern))
	fmt.Println("结束时间：", endTime.Format(timePattern))
	fmt.Printf("将于%s内完成%.6f公里距离，速度约为%.2fm/s \n", duration, distance, v)

	fmt.Println("请输入YES确认")
	var confirm string
	fmt.Scanf("%s", &confirm)
	if confirm != "YES" {
		return
	}
	status, err := jkwx.UploadData(s, cmdFlags.distance, beginTime, endTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if status == 1 {
		fmt.Println("OK.")
	} else {
		fmt.Printf("Status %d", s)
	}
}
