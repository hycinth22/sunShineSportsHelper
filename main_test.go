package main

import (
	"testing"
	"time"
	"./utility"
)

func TestTime(t *testing.T)  {
	e := time.Now()
	b := e.Add(time.Duration(utility.RandRange(-35, -25))*time.Minute)
	t.Logf("%v\n%v", b.Format(timePattern), e.Format(timePattern))
}
