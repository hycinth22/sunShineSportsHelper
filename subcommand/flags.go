package subcommand

import "time"

type CmdFlagsType struct {
	Silent bool

	SchoolID int64
	User     string
	Password string

	EndTime         string
	RawRecord       bool
	IgnoreCompleted bool
	Distance        float64
	Duration        time.Duration

	DeviceName string
	ModelType  string
	Screen     string
	Imei       string
	Imsi       string
	UserAgent  string
}
