package utility

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func init()  {
	rand.Seed(time.Now().Unix())
}
func RandRange(min int, max int) int {

	return min + rand.Int()%(max-min+1)
}

func MD5String(raw string) string{
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}