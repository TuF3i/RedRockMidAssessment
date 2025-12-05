package models

type Commander struct {
	Role string      `json:"role"`
	Msg  interface{} `json:"msg"`
}

type CourseMsg struct {
	Operation string `json:"operation"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"stu_id"`
}

type SelectedNum struct {
	CourseID    string `json:"course_id"`
	SelectedNum uint   `json:"selected_num"`
}
