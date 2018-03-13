package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	jkwx "./sunShine_Sports"
)

const sessionFile =  "sunShine_Sports.session"
func saveSession(s *jkwx.Session){
	f, err := os.Create(sessionFile)
	if err!=nil{
		panic(err)
	}
	if err := gob.NewEncoder(f).Encode(s); err != nil{
		panic(err)
	}
}
func readSession() *jkwx.Session {
	var s jkwx.Session
	f, err := os.Open(sessionFile)
	if err != nil{
		return nil
	}
	if err := gob.NewDecoder(f).Decode(&s); err != nil{
		fmt.Println(err.Error())
		return nil
	}
	nowTime := time.Now()
	expiredTime := time.Unix(s.UserExpirationTime/1000, 0)
	fmt.Println("Alread Login.")
	fmt.Println("nowTime", nowTime.Format(timePattern))
	fmt.Println("expiredTime", expiredTime.Format(timePattern))
	fmt.Println()
	if nowTime.After(expiredTime) {
		fmt.Println("Login Expired.")
		return nil
	}
	return &s
}