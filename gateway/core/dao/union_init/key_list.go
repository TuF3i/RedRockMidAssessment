package union_init

import "fmt"

var (
	courseSelectionStatusKey = func() string { return fmt.Sprintf("courseSelection:status") }                   // 选课服务状态
	courseIDsKey             = func() string { return fmt.Sprintf("course:allID") }                             // 所有课程的ID
	courseStockKey           = func(courseID string) string { return fmt.Sprintf("course:%v:stock", courseID) } // 生成courseStock的Key
	courseInfoKey            = func(courseID string) string { return fmt.Sprintf("course:%v:info", courseID) }  // 生成courseInfo的Key
)
