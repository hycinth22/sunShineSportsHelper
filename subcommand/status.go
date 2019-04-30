package subcommand

import (
	"fmt"
	"github.com/inkedawn/go-sunshinemotion"
	_const "github.com/inkedawn/sunShineSportsHelper/const"
)

func ShowStatus(s *ssmt.Session, info ssmt.UserInfo) {
	r, err := s.GetSportResult()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("-----------")
	fmt.Println("| 帐号信息 |")
	fmt.Println("-----------")
	fmt.Println("ID：\t", s.User.UserID)
	fmt.Println("SchoolID：\t", s.User.SchoolID)
	fmt.Println("StuNum：\t", s.User.StuNum)
	fmt.Println("-----------")
	fmt.Println("班级：\t", info.ClassName)
	fmt.Println("学号：\t", info.StudentNumber)
	fmt.Println("姓名：\t", info.StudentName)
	fmt.Println("性别：\t", info.Sex)
	fmt.Println("UserRoleID：\t", info.UserRoleID)
	fmt.Println("-----------")
	fmt.Printf("LastTime：\t%s \n", r.LastTime.Format(_const.DisplayTimePattern))
	fmt.Printf("已跑距离：\t%05.3f 公里\n", r.ActualDistance)
	fmt.Printf("达标距离：\t%05.3f 公里\n", r.QualifiedDistance)
	fmt.Println("-----------")
	// fmt.Printf("%+v", r)
}
