package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ssmt "github.com/inkedawn/go-sunshinemotion/v3"

	_const "github.com/inkedawn/sunShineSportsHelper/const"
	"github.com/inkedawn/sunShineSportsHelper/sessionStroage"
	"github.com/inkedawn/sunShineSportsHelper/subcommand"
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
		subcommand.LoginAccount(cmdFlags)
		return
	}

	// need session
	s, info := tryResume()
	if s == nil {
		fmt.Println("Need to login.")
		return
	}
	checkAppVer(s)
	subcommand.ShowStatus(s, info)

	switch cmd {
	case "status":
		return
	case "getRoute":
		subcommand.GetRoute(s)
		return
	case "upload":
		subcommand.UploadData(s, info, cmdFlags)
		return
	case "uploadTest":
		subcommand.UploadTestData(s, cmdFlags)
		return
	case "testRule":
		subcommand.GetTestRule(s)
		return
	case "setDevice":
		subcommand.SetDevice(s, info, cmdFlags)
		return
	case "showDevice":
		subcommand.ShowDevice(s)
		return
	}
}

func tryResume() (*ssmt.Session, ssmt.UserInfo) {
	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Println("Resume failed.", err.Error())
			fmt.Println("** You can try to delete the session file ", sessionStroage.GetSessionFilePathById(cmdFlags.SchoolID, cmdFlags.User))
		}
	}()
	s, info, err := sessionStroage.ReadSession(cmdFlags.SchoolID, cmdFlags.User)
	if err == sessionStroage.ErrSessionNotExist {
		return nil, ssmt.UserInfo{}
	}
	if err != nil {
		panic(err)
	}
	if s != nil {
		fmt.Println("Use Existent Session.")
		fmt.Println("UserAgent", s.Device.UserAgent)
		fmt.Println("ExpiredTime", s.Token.ExpirationTime.Format(_const.DisplayTimePattern))
		fmt.Println()
	}
	return s, info
}

func printHelp() {
	flag.Usage()
	os.Exit(0)
}

func checkAppVer(s *ssmt.Session) {
	appInfo, err := s.GetAppInfo()
	if err != nil {
		fmt.Println("Warning: 获取APP版本失败")
		return
	}
	fmt.Println("Latest Ver: ", appInfo.VerNumber)
	fmt.Println("Lib Ver: ", ssmt.AppVersionID)
	if appInfo.VerNumber > ssmt.AppVersionID {
		fmt.Println("需要更新")
		os.Exit(1)
	}
	return
}
