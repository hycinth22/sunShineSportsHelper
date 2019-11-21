package subcommand

import (
	"fmt"

	jkwx "github.com/inkedawn/go-sunshinemotion/v3"
)

func GetTestRule(s *jkwx.Session) {
	rule, err := s.GetTestRule()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("TestRule: %+v", rule)
	}
}
