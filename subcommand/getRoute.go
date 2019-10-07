package subcommand

import (
	"fmt"

	jkwx "github.com/inkedawn/go-sunshinemotion/v3"
)

func GetRoute(s *jkwx.Session) {
	r, err := s.GetRandRoute()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(r)
	}
}
