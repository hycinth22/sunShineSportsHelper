package lib

type UserInfo struct {
	InClassID      int64  `json:"inClassID"`
	InClassName    string `json:"inClassName"`
	InCollegeID    int64  `json:"inCollegeID"`
	InCollegeName  string `json:"inCollegeName"`
	InSchoolID     int64  `json:"inSchoolID"`
	InSchoolName   string `json:"inSchoolName"`
	InSchoolNumber string `json:"inSchoolNumber"`
	NickName       string `json:"nickName"`
	StudentName    string `json:"studentName"`
	StudentNumber  string `json:"studentNumber"`
	IsTeacher      int    `json:"isTeacher"`
	Sex            string `json:"sex"`
	// PhoneNumber   string `json:"phoneNumber"`
	// UserRoleID    int
}
