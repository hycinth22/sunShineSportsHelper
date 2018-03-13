package main

import (
	"math/rand"
	"time"
)

func randRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Int()%(max-min+1)
}
