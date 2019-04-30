package subcommand

import (
	"fmt"
	jkwx "github.com/inkedawn/go-sunshinemotion"
	_const "github.com/inkedawn/sunShineSportsHelper/const"
	"os"
	"time"
)

func UploadTestData(s *jkwx.Session, cmdFlags CmdFlagsType) {
	totalDistance := cmdFlags.Distance
	duration := cmdFlags.Duration
	endTime, err := time.Parse(_const.InputTimePattern, cmdFlags.EndTime)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	record := jkwx.CreateRecord(s.User.UserID, cmdFlags.SchoolID, totalDistance, endTime, duration)
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
