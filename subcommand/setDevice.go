package subcommand

import (
	jkwx "github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/sunShineSportsHelper/sessionStroage"
)

func SetDevice(s *jkwx.Session, info jkwx.UserInfo, cmdFlags CmdFlagsType) {
	if cmdFlags.DeviceName != "" {
		s.Device.DeviceName = cmdFlags.DeviceName
	}
	if cmdFlags.ModelType != "" {
		s.Device.ModelType = cmdFlags.ModelType
	}
	if cmdFlags.Screen != "" {
		s.Device.Screen = cmdFlags.Screen
	}
	if cmdFlags.Imei != "" {
		s.Device.IMEI = cmdFlags.Imei
	}
	if cmdFlags.Imsi != "" {
		s.Device.IMSI = cmdFlags.Imsi
	}
	if cmdFlags.UserAgent != "" {
		s.Device.UserAgent = cmdFlags.UserAgent
	}
	sessionStroage.SaveSession(s, info)
}
