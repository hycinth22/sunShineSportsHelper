package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/inkedawn/sunShineSportsHelper/utility"
)

var cmd string
var cmdFlags struct {
	silent bool

	schoolID int64
	user     string
	password string

	endTime         string
	rawRecord       bool
	ignoreCompleted bool
	distance        float64
	duration        time.Duration

	deviceName string
	modelType  string
	imei       string
	imsi       string
	userAgent  string
}

func parseFlag() {
	if len(os.Args) < 2 {
		fmt.Println("Arguments needed.")
		os.Exit(0)
	}
	cmd = os.Args[1]
	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	fmt.Println("Command:", cmd)
	flags.BoolVar(&cmdFlags.silent, "q", false, "quiet mode")

	flags.Int64Var(&cmdFlags.schoolID, "s", defaultSchoolID, "school ID")
	flags.StringVar(&cmdFlags.user, "u", "default", "account(stuNum)")
	flags.StringVar(&cmdFlags.password, "p", "", "password")

	flags.StringVar(&cmdFlags.endTime, "endTime", time.Now().Format(inputTimePattern), "upload test sport data")
	flags.BoolVar(&cmdFlags.rawRecord, "rawRecord", false, "upload rawRecord sport data")
	flags.BoolVar(&cmdFlags.ignoreCompleted, "ignoreCompleted", false, "continue to upload though completed")
	flags.Float64Var(&cmdFlags.distance, "distance", 0.0, "distance(精确到小数点后6位).")

	randomDuration := time.Duration(utility.RandRange(12, 20)) * time.Minute
	flags.DurationVar(&cmdFlags.duration, "duration", randomDuration, "time duration")

	flags.StringVar(&cmdFlags.deviceName, "device", "", "")
	flags.StringVar(&cmdFlags.modelType, "model", "", "")
	flags.StringVar(&cmdFlags.imei, "imei", "", "")
	flags.StringVar(&cmdFlags.imsi, "imsi", "", "")
	flags.StringVar(&cmdFlags.userAgent, "userAgent", "", "")
	flags.Parse(os.Args[2:])
}
