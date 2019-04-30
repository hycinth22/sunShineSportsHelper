package subcommand

import (
	"fmt"
	"github.com/inkedawn/go-sunshinemotion"
	_const "github.com/inkedawn/sunShineSportsHelper/const"

	"log"
	"os"
	"time"
)

const (
	defaultDistanceFemale = 2.5
	defaultDistanceMale   = 4.5
)

func UploadData(s *ssmt.Session, info ssmt.UserInfo, cmdFlags CmdFlagsType) {
	rawRecord := cmdFlags.RawRecord
	ignoreCompleted := cmdFlags.IgnoreCompleted
	totalDistance := cmdFlags.Distance
	if cmdFlags.Distance == 0.0 {
		switch info.Sex {
		case "F":
			totalDistance = defaultDistanceFemale
		case "M":
			totalDistance = defaultDistanceMale
		default:
			log.Panicln("Unknown sex", info.Sex)
		}
	}
	endTime, err := time.Parse(_const.InputTimePattern, cmdFlags.EndTime)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	duration := cmdFlags.Duration

	var records []ssmt.Record
	if !rawRecord {
		if !ignoreCompleted {
			r, err := s.GetSportResult()
			if err == nil && r.ActualDistance >= r.QualifiedDistance {
				fmt.Println("已达标，停止操作")

				return
			}
		}
		records = ssmt.SmartCreateRecordsBefore(s.User.SchoolID, s.User.UserID, ssmt.GetDefaultLimitParams(info.Sex), totalDistance, endTime)
	} else {
		records = []ssmt.Record{
			ssmt.CreateRecord(s.User.UserID, s.User.SchoolID, totalDistance, endTime, duration),
		}
	}

	if !confirm(records, cmdFlags.Silent) {
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
	ShowStatus(s, info)
	if !ignoreCompleted {
		r, err := s.GetSportResult()
		if err == nil && r.ActualDistance >= r.QualifiedDistance {
			fmt.Println("已达标")
		}
	}
	fmt.Println("---------------")
}

func confirm(records []ssmt.Record, silent bool) bool {
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
		fmt.Println("起始时间：", record.BeginTime.Format(_const.DisplayTimePattern))
		fmt.Println("结束时间：", record.EndTime.Format(_const.DisplayTimePattern))
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
