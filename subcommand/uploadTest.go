package subcommand

import (
	"fmt"
	"log"
	"os"
	"time"

	jkwx "github.com/inkedawn/go-sunshinemotion/v3"

	_const "github.com/inkedawn/sunShineSportsHelper/const"
	"github.com/inkedawn/sunShineSportsHelper/utility"
)

func UploadTestData(s *jkwx.Session, info jkwx.UserInfo, cmdFlags CmdFlagsType) {
	var (
		distance float64
		duration time.Duration
	)
	rule, err := s.GetTestRule()
	if err != nil {
		log.Fatal(err)
	}
	switch info.Sex {
	case "F":
		distance = rule.GirlDistance / 1000
		duration = time.Duration(utility.RandRange(12, 13))*time.Minute + time.Duration(utility.RandRange(0, 60))*time.Second
	case "M":
		distance = rule.ManDistance / 1000
		duration = time.Duration(utility.RandRange(15, 16))*time.Minute + time.Duration(utility.RandRange(0, 60))*time.Second
	default:
		log.Fatal("unknown sex")
	}
	endTime, err := time.ParseInLocation(_const.InputTimePattern, cmdFlags.EndTime, _const.TimeZoneCST)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	record := jkwx.CreateRecord(s.User.UserID, cmdFlags.SchoolID, distance, endTime, duration)
	if !confirm([]jkwx.Record{record}, cmdFlags.Silent) {
		return
	}
	err = s.UploadTestRecord(record)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("上传成功")
	}

}
