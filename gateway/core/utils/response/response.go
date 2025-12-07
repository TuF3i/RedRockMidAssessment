package response

// 操作成功响应
var OperationSuccess = Response{Status: 20000, Info: "Operation Success"}

// api数据错误
var RevDataError = Response{Status: 70000, Info: "Cant Validate Data"}

// 业务错误集合
var (
	//登录错误
	UserNotExiOrWrongStuID = Response{Status: 10002, Info: "user not exist or wrong student id"} // 用户名错误或不存在
	WrongPassword          = Response{Status: 10001, Info: "wrong password"}                     // 密码错误
	PermissionDenied       = Response{Status: 10000, Info: "permission denied"}                  // 权限不足

	// 会话错误
	EmptyToken   = Response{Status: 10016, Info: "empty token"}   // 空Token
	InvalidToken = Response{Status: 10016, Info: "invalid token"} // 无效Token

	//注册错误
	StudentIDAlreadyExist = Response{Status: 10006, Info: "student id already exist"} // 学生ID已存在

	//选课错误
	CourseNotExist                   = Response{Status: 10007, Info: "course not exist"}                      // 课程不存在
	CourseIsFull                     = Response{Status: 10008, Info: "course capacity is full"}               // 选课人数已满
	CourseDoubleSelected             = Response{Status: 10017, Info: "you have selected the course twice"}    // 重复选课
	RecordAlreadyExist               = Response{Status: 10022, Info: "course selection record already exist"} // 选课记录已存在
	CourseSelectionEventNotStart     = Response{Status: 10024, Info: "course selection event not start"}      // 选课未开始
	CourseSelectionEventAlreadyStart = Response{Status: 10025, Info: "course selection event already start"}  // 选课已开始

	//退课错误
	RecordNotExist = Response{Status: 10016, Info: "Course Selection Record Not Exist"} // 选课记录不存在

	//添加课程错误
	CourseAlreadyExist = Response{Status: 10023, Info: "course already exist"} // 课程已选

	//数据校验错误
	InvalidUserName       = Response{Status: 10009, Info: "invalid user name"}       // 无效用户名
	InvalidPassword       = Response{Status: 10010, Info: "invalid password"}        // 无效密码
	InvalidSexSetting     = Response{Status: 10010, Info: "invalid sex setting"}     // 无效性别
	InvalidStudentID      = Response{Status: 10011, Info: "invalid student id"}      // 无效学生ID
	InvalidClass          = Response{Status: 10012, Info: "invalid class"}           // 无效班级
	InvalidGrade          = Response{Status: 10013, Info: "invalid grade"}           // 无效年级
	InvalidAge            = Response{Status: 10014, Info: "invalid age"}             // 无效年龄
	InvalidCourse         = Response{Status: 10016, Info: "invalid course id"}       // 无效的课程ID
	InvalidCourseName     = Response{Status: 10017, Info: "invalid course name"}     // 无效课程名称
	InvalidCourseLocation = Response{Status: 10018, Info: "invalid course location"} // 无效上课地点
	InvalidCourseTime     = Response{Status: 10019, Info: "invalid course time"}     // 无效上课时间
	InvalidCourseTeacher  = Response{Status: 10020, Info: "invalid course teacher"}  // 无效教师名称
	InvalidCourseStock    = Response{Status: 10021, Info: "invalid course stock"}    // 无效课程容量

	// 空数据错误
	EmptyData = Response{Status: 10015, Info: "empty update columns"} // 无效的字段信息

	// 管理员修改信息错误
	NoChangingData = Response{Status: 10026, Info: "you can not change the data until the selection event finished"} // 管理员无法在选课期间更改学生选课的如何信息
)

func GenFinalResponse(response Response, data interface{}) FinalResponse {
	return FinalResponse{Status: response.Status, Info: response.Info, Data: data}
}

// 业务层错误封装
type Response struct {
	Status uint   `json:"status"`
	Info   string `json:"info"`
}

func (r Response) Error() string {
	return r.Info
}

// 请求层错误封装
type FinalResponse struct {
	Status uint        `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

// 服务器内部错误封装
func ServerInternalError(err error) Response {
	return Response{
		Status: 500,
		Info:   err.Error(),
	}
}
