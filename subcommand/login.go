package subcommand

import (
	"crypto/md5"
	"fmt"

	"github.com/inkedawn/go-sunshinemotion/v3"

	"github.com/inkedawn/sunShineSportsHelper/sessionStroage"
	"github.com/inkedawn/sunShineSportsHelper/utility"
)

func LoginAccount(cmdFlags CmdFlagsType) {
	s, info, err := sessionStroage.ReadSession(cmdFlags.SchoolID, cmdFlags.User)
	if err != nil && err != sessionStroage.ErrSessionNotExist {
		panic(err)
	}
	if s != nil {
		// old user
		fmt.Println("Already Login. Use Old Session(Keep Same Device).")
		info, err = s.Login(cmdFlags.SchoolID, cmdFlags.User, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.Password))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		// new user
		s = ssmt.CreateSession()
		info, err = s.Login(cmdFlags.SchoolID, cmdFlags.User, "123", fmt.Sprintf("%x", md5.Sum([]byte(cmdFlags.Password))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		s.Device.UserAgent = utility.GetRandUserAgent()
	}
	fmt.Printf("Device: %+v\n", *s.Device)
	fmt.Printf("Token: %+v\n", *s.Token)
	sessionStroage.SaveSession(s, info)
	ShowStatus(s, info)
}
