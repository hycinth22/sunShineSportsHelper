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
	flag.StringVar(&cmdFlags.user, "user", defaultStuNum, "account(stuNum)")
	flag.StringVar(&cmdFlags.password, "password", "", "password")

	flag.BoolVar(&cmdFlags.status, "status", false, "view account status")

	flag.BoolVar(&cmdFlags.upload, "upload", false, "upload sport data")
	distanceRandomRatio :=  float64(utility.RandRange(9500, 11142))/10000 // 95%-111.42%
	flag.Float64Var(&cmdFlags.distance, "distance", 3.500000 *distanceRandomRatio, "distance(精确到小数点后6位)")

	randomDuration := time.Duration(utility.RandRange(12, 20))*time.Minute
	flag.DurationVar(&cmdFlags.duration, "duration", randomDuration, "time duration")
	flag.Parse()

	// TOOD: beginTime

	// TOOD: distance sperate


	// 小数部分随机化
	cmdFlags.distance += float64(utility.RandRange(-99999, 99999)) /1000000 // -0.09 ~ 0.09
	// 秒级随机化
	cmdFlags.duration += time.Duration(utility.RandRange(0, 60))*time.Second
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
	distance := cmdFlags.distance
	duration := cmdFlags.duration
	endTime := time.Now().Add(-time.Duration(utility.RandRange(1, 10)) * time.Minute)
	beginTime := endTime.Add(-duration)
	v := cmdFlags.distance * 1000 / duration.Seconds()

	fmt.Println("--------------")
	fmt.Println("| 确认上传数据 |")
	fmt.Println("---------------")
	fmt.Println("起始时间：", beginTime.Format(timePattern))
	fmt.Println("结束时间：", endTime.Format(timePattern))
	fmt.Printf("将于%s内完成%.6f公里距离，速度约为%.2fm/s \n", duration, distance, v)

	fmt.Println("请输入YES确认")
	var confirm string
	fmt.Scanf("%s", &confirm)
	fmt.Println("---------------")
	if confirm != "YES" {
		return
	}

	status, err := jkwx.UploadData(s, cmdFlags.distance, beginTime, endTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("---------------")
	fmt.Println("上传结果：")
	if status == 1 {
		fmt.Println("OK.")
		showStatus(s)
	} else {
		fmt.Printf("Status %d", s)
	}
	fmt.Println("---------------")
}