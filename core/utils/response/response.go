package response

// 空错误
var Null = Response{Status: 0, Info: ""}

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

	EmptyToken   = Response{Status: 10016, Info: "empty token"}   // 空Token
	InvalidToken = Response{Status: 10016, Info: "invalid token"} // 无效Token

	//注册错误
	UserAlreadyExist      = Response{Status: 10003, Info: "user already exist"}       // 用户已存在
	GradeNotExist         = Response{Status: 10004, Info: "grade not exist"}          // 年级不存在
	ClassNotExist         = Response{Status: 10005, Info: "Class not exist"}          // 班级不存在
	StudentIDAlreadyExist = Response{Status: 10006, Info: "student id already exist"} // 学生ID已存在

	//选课错误
	CourseNotExist       = Response{Status: 10007, Info: "course not exist"}                   // 课程不存在
	CourseIsFull         = Response{Status: 10008, Info: "course capacity is full"}            // 选课人数已满
	CourseDoubleSelected = Response{Status: 10017, Info: "you have selected the course twice"} // 重复选课

	//退课错误
	RecordNotExist = Response{Status: 10016, Info: "Course Selection Record Not Exist"} // 选课记录不存在

	//数据校验错误
	InvalidUserName   = Response{Status: 10009, Info: "invalid user name"}   // 无效用户名
	InvalidPassword   = Response{Status: 10010, Info: "invalid password"}    // 无效密码
	InvalidSexSetting = Response{Status: 10010, Info: "invalid sex setting"} // 无效性别
	InvalidStudentID  = Response{Status: 10011, Info: "invalid student id"}  // 无效学生ID
	InvalidClass      = Response{Status: 10012, Info: "invalid class"}       // 无效班级
	InvalidGrade      = Response{Status: 10013, Info: "invalid grade"}       // 无效年级
	InvalidAge        = Response{Status: 10014, Info: "invalid age"}         // 无效年龄
	InvalidCourse     = Response{Status: 10016, Info: "invalid course ID"}   // 无效的课程ID

	EmptyData = Response{Status: 10015, Info: "empty update columns"} // 无效的字段信息
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
