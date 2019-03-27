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
}

func parseFlag() {
	if len(os.Args) < 2 {
		fmt.Println("Arguments needed.")
		os.Exit(0)
	}
	cmd = os.Args[1]
	fmt.Println("Command:", cmd)
	flag.BoolVar(&cmdFlags.silent, "q", false, "quiet mode")

	flag.Int64Var(&cmdFlags.schoolID, "s", defaultSchoolID, "school ID")
	flag.StringVar(&cmdFlags.user, "u", "default", "account(stuNum)")
	flag.StringVar(&cmdFlags.password, "p", "", "password")

	flag.StringVar(&cmdFlags.endTime, "endTime", time.Now().Format(inputTimePattern), "upload test sport data")
	flag.BoolVar(&cmdFlags.rawRecord, "rawRecord", false, "upload rawRecord sport data")
	flag.BoolVar(&cmdFlags.ignoreCompleted, "ignoreCompleted", false, "continue to upload though completed")
	flag.Float64Var(&cmdFlags.distance, "distance", 0.0, "distance(精确到小数点后6位).")

	randomDuration := time.Duration(utility.RandRange(12, 20)) * time.Minute
	flag.DurationVar(&cmdFlags.duration, "duration", randomDuration, "time duration")
	flag.Parse()
}
