package utility

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func RandRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Int()%(max-min+1)
}

func MD5String(raw string) string{
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}