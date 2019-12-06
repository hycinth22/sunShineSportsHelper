package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_const "github.com/inkedawn/sunShineSportsHelper/const"
	"github.com/inkedawn/sunShineSportsHelper/subcommand"

	"github.com/inkedawn/sunShineSportsHelper/utility"
)

var cmd string
var cmdFlags subcommand.CmdFlagsType

const (
	defaultSchoolID = 60
)

func parseFlag() {
	if len(os.Args) < 2 {
		fmt.Println("Arguments needed.")
		os.Exit(0)
	}
	cmd = os.Args[1]
	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	fmt.Println("Command:", cmd)
	flags.BoolVar(&cmdFlags.Silent, "q", false, "quiet mode")

	flags.Int64Var(&cmdFlags.SchoolID, "s", defaultSchoolID, "school ID")
	flags.StringVar(&cmdFlags.User, "u", "default", "account(stuNum)")
	flags.StringVar(&cmdFlags.Password, "p", "", "password")

	flags.StringVar(&cmdFlags.EndTime, "endTime", time.Now().In(_const.TimeZoneCST).Format(_const.InputTimePattern), "upload test sport data")
	flags.BoolVar(&cmdFlags.RawRecord, "rawRecord", false, "upload rawRecord sport data")
	flags.BoolVar(&cmdFlags.IgnoreCompleted, "ignoreCompleted", false, "continue to upload though completed")
	flags.Float64Var(&cmdFlags.Distance, "distance", 0.0, "distance(精确到小数点后6位).")

	randomDuration := time.Duration(utility.RandRange(12, 20)) * time.Minute
	flags.DurationVar(&cmdFlags.Duration, "duration", randomDuration, "time duration")

	flags.StringVar(&cmdFlags.DeviceName, "device", "", "")
	flags.StringVar(&cmdFlags.ModelType, "model", "", "")
	flags.StringVar(&cmdFlags.Screen, "screen", "", "")
	flags.StringVar(&cmdFlags.Imei, "imei", "", "")
	flags.StringVar(&cmdFlags.Imsi, "imsi", "", "")
	flags.StringVar(&cmdFlags.UserAgent, "userAgent", "", "")
	_ = flags.Parse(os.Args[2:])
}
