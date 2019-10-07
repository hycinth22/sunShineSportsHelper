package subcommand

import (
	"fmt"

	jkwx "github.com/inkedawn/go-sunshinemotion/v3"
)

func ShowDevice(s *jkwx.Session) {
	fmt.Printf("Device: %+v", *s.Device)
}
