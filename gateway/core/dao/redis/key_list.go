package redis

import "fmt"

var (
	accessTokenKey  = func(userID string) string { return fmt.Sprintf("user:%v:accessToken", userID) }  // 生成accessToken的key
	refreshTokenKey = func(userID string) string { return fmt.Sprintf("user:%v:refreshToken", userID) } // 生成refreshToken的key

	courseInfoKey         = func(courseID string) string { return fmt.Sprintf("course:%v:info", courseID) }    // 生成courseInfo的Key
	courseStockKey        = func(courseID string) string { return fmt.Sprintf("course:%v:stock", courseID) }   // 生成courseStock的Key
	courseUsersKey        = func(courseID string) string { return fmt.Sprintf("course:%v:users", courseID) }   // 生成courseUsers的Key
	courseDroppedUsersKey = func(courseID string) string { return fmt.Sprintf("course:%v:dropped", courseID) } // 生成courseDroppedUsers的Key
	courseIDsKey          = func() string { return fmt.Sprintf("course:allID") }                               // 所有课程的ID

	studentSelectedCourseKey = func(userID string) string { return fmt.Sprintf("user:%v:selectedCourse", userID) } // 生成studentSelectedCourse的Key

	courseSelectionStatusKey = func() string { return fmt.Sprintf("courseSelection:status") } // 选课服务状态
)
