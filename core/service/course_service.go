package service

import (
	"RedRockMidAssessment/core/dao/mysql"
	"RedRockMidAssessment/core/dao/redis"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"RedRockMidAssessment/core/utils/verify"
	"context"
	"errors"
	"reflect"
)

func GetCourseInfo(ctx context.Context) (interface{}, response.Response) {
	var courseInfo []models.Course
	// 取出所有课程的ID
	ids, rsp := redis.GetAllCourseID(ctx)
	if !errors.Is(rsp, response.OperationSuccess) {
		return nil, rsp
	}
	// 遍历获取所有课程
	idList := ids.([]string) // 加错误处理
	for _, id := range idList {
		info, rsp := redis.GetCourseDetails(ctx, id)
		if !errors.Is(rsp, response.OperationSuccess) {
			return nil, rsp
		}
		courseInfo = append(courseInfo, info.(models.Course))
	}

	return models.SelectableClasses{ // 组装data
		Info: courseInfo,
	}, response.OperationSuccess
}

func GetStuSelectedCourses(ctx context.Context, userID string) (interface{}, response.Response) {
	var courseInfo []models.Course
	// 取出已选课程的ID
	ids, rsp := redis.GetStuSelectedCourseID(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return nil, rsp
	}
	// 遍历获取所有课程
	idList := ids.([]string) // 加错误处理
	for _, id := range idList {
		info, rsp := redis.GetCourseDetails(ctx, id)
		if !errors.Is(rsp, response.OperationSuccess) {
			return nil, rsp
		}
		courseInfo = append(courseInfo, info.(models.Course))
	}
	return models.SelectedClasses{ // 组装data
		Info: courseInfo,
	}, response.OperationSuccess
}

func SubscribeCourse(ctx context.Context, userID string, courseID string) response.Response {
	// 检测CourseID是否有效
	if !verify.VerifyCourseID(courseID) {
		return response.InvalidCourse
	}
	// 检测课程是否存在
	ok, rsp := redis.CheckIfCourseExist(ctx, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	// 检测课程是否存在
	if !ok {
		return response.CourseNotExist
	}
	// 判断是否选择课程
	ok, rsp = redis.CheckIfCourseSelected(ctx, userID, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if ok {
		return response.CourseDoubleSelected
	}
	// 写入选课
	rsp = redis.SubscribeACourse(ctx, userID, courseID)
	return rsp
}

func DropCourse(ctx context.Context, userID string, courseID string) response.Response {
	// 检测CourseID是否有效
	if !verify.VerifyCourseID(courseID) {
		return response.InvalidCourse
	}
	// 检测课程是否存在
	ok, rsp := redis.CheckIfCourseExist(ctx, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.CourseNotExist
	}

	// 判断是否选择课程
	ok, rsp = redis.CheckIfCourseSelected(ctx, userID, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.RecordNotExist
	}

	// 执行退课
	rsp = redis.DropACourse(ctx, userID, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}

	return response.OperationSuccess
}

func UpdateCourseInfoForAdmin(ctx context.Context, courseID string, data []models.UpdateColumnsEntityForCourse) response.Response {
	field := make([]string, 0)
	dataList := make(map[string]interface{})
	for _, item := range data {
		v := reflect.ValueOf(item)
		t := reflect.TypeOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			t = t.Elem()
		}
		for i := 0; i < v.NumField(); i++ { // TODO
			value := v.Field(i).Interface()
			column := t.Field(i).Tag.Get("json")
			switch column {
			case "class_name": // 校验课程名称
				v := value.(string)
				if !verify.VerifyCourseName(v) {
					return response.InvalidCourseName
				}
				field = append(field, column)
				dataList[column] = value
			case "class_id":
				v := value.(string)
				if !verify.VerifyCourseID(v) {
					return response.InvalidCourse
				}
				field = append(field, column)
				dataList[column] = value
			case "class_location":
				v := value.(string)
				if !verify.VerifyCourseID(v) {
					return response.InvalidCourseLocation
				}
				field = append(field, column)
				dataList[column] = value
			case "class_time":
				v := value.(string)
				if !verify.VerifyCourseID(v) {
					return response.InvalidCourseName
				}
				field = append(field, column)
				dataList[column] = value
			case "class_teacher":
				v := value.(string)
				if !verify.VerifyCourseID(v) {
					return response.InvalidCourseTeacher
				}
				field = append(field, column)
				dataList[column] = value
			case "class_capcity":
				v := value.(string)
				if !verify.VerifyCourseID(v) {
					return response.InvalidCourseStock
				}
				field = append(field, column)
				dataList[column] = value
			}
		}
	} // for

	/* 空数据直接返回 */
	if len(field) == 0 {
		return response.EmptyData
	}

	/* 检查学生是否存在 */
	ifExist, rsp := mysql.CheckIfCourseExist(ctx, courseID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return rsp
	}

	if !ifExist {
		return response.CourseNotExist
	}

	/* 写MySQL */
	rsp = mysql.UpdateCourseInfo(ctx, courseID, field, dataList)
	return rsp
}

func GetCourseInfoForAdmin(ctx context.Context) (interface{}, response.Response) {
	// 获取所有课程的信息
	data, rsp := mysql.GetAllCourseInfo(ctx)
	if !errors.Is(rsp, response.OperationSuccess) {
		return nil, rsp
	}

	// 类型断言
	courseInfo := data.([]models.Course)

	return models.SelectableClasses{ // 组装data
		Info: courseInfo,
	}, response.OperationSuccess
}

func GetStuSelectedCoursesForAdmin(ctx context.Context, userID string) (interface{}, response.Response) {
	// 校验stuID是否可用
	if !verify.VerifyUserID(userID) {
		return nil, response.InvalidStudentID
	}
	// 获取所选课程
	data, rsp := mysql.GetStudentSelectedCourse(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return nil, rsp
	}

	return data, rsp
}

func SubscribeCourseForAdmin(ctx context.Context, userID string, courseID string) response.Response {
	// 校验stuID是否可用
	if !verify.VerifyUserID(userID) {
		return response.InvalidStudentID
	}
	// 校验courseID是否可用
	if !verify.VerifyCourseID(courseID) {
		return response.InvalidCourse
	}
	// 检查学生是否存在
	ok, rsp := mysql.CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.UserNotExiOrWrongStuID
	}
	// 检查课程是否存在
	ok, rsp = mysql.CheckIfCourseExist(ctx, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.CourseNotExist
	}
	// 检查选课记录是否存在
	ok, rsp = mysql.CheckIfCourseSelected(ctx, courseID, userID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if ok {
		return response.RecordAlreadyExist
	}
	// 执行写入操作
	rsp = mysql.AddCourseToStudent(ctx, courseID, userID)
	return rsp
}

func DropCourseForAdmin(ctx context.Context, userID string, courseID string) response.Response {
	// 校验stuID是否可用
	if !verify.VerifyUserID(userID) {
		return response.InvalidStudentID
	}
	// 校验courseID是否可用
	if !verify.VerifyCourseID(courseID) {
		return response.InvalidCourse
	}
	// 检查学生是否存在
	ok, rsp := mysql.CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.UserNotExiOrWrongStuID
	}
	// 检查课程是否存在
	ok, rsp = mysql.CheckIfCourseExist(ctx, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.CourseNotExist
	}
	// 检查选课记录是否存在
	ok, rsp = mysql.CheckIfCourseSelected(ctx, courseID, userID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.RecordNotExist
	}
	// 执行写入操作
	rsp = mysql.DelCourseToStudent(ctx, courseID, userID)
	return rsp
}

func UpdateCourseStockForAdmin(ctx context.Context, courseID string, stock uint) response.Response {
	// 验证courseID是否可用
	if !verify.VerifyCourseID(courseID) {
		return response.InvalidCourse
	}
	// 验证stock是否合法
	if !verify.VerifyCourseStock(stock) {
		return response.InvalidCourseStock
	}
	// 验证课程是否存在
	ok, rsp := mysql.CheckIfCourseExist(ctx, courseID)
	if errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.CourseNotExist
	}
	// 获取库存数据
	data, rsp := mysql.GetCourseInfo(ctx, courseID)
	//stock := data.(models.Course).ClassCapacity
	selectedNum := data.(models.Course).ClassSelectedNum
	// 判断新stock是否可用
	if stock < selectedNum {
		return response.InvalidCourseStock
	}
	// 更改库存
	rsp = mysql.UpdateCourseStock(ctx, courseID, stock)
	return rsp
}

func AddCourseForAdmin(ctx context.Context, data models.Course) response.Response {
	//校验课程名称
	if !verify.VerifyCourseName(data.ClassName) {
		return response.InvalidCourseName
	}
	// 校验课程ID
	if !verify.VerifyCourseID(data.ClassID) {
		return response.InvalidCourse
	}
	// 校验上课地带
	if !verify.VerifyCourseLocation(data.ClassLocation) {
		return response.InvalidCourseLocation
	}
	// 校验上课时间
	if !verify.VerifyCourseTime(data.ClassTime) {
		return response.InvalidCourseTime
	}
	// 校验上课老师
	if !verify.VerifyCourseTeacher(data.ClassTeacher) {
		return response.InvalidCourseTeacher
	}
	// 校验课程容量
	if !verify.VerifyCourseStock(data.ClassCapacity) {
		return response.InvalidCourseStock
	}

	// 检查课程是否存在
	ok, rsp := mysql.CheckIfCourseExist(ctx, data.ClassID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if ok {
		return response.CourseAlreadyExist
	}

	// 添加课程
	rsp = mysql.InsertCourseIntoDB(ctx, data)
	return rsp
}

func DelCourseForAdmin(ctx context.Context, courseID string) response.Response {
	// 校验课程ID
	if !verify.VerifyCourseID(courseID) {
		return response.InvalidCourse
	}

	// 检查课程是否存在
	ok, rsp := mysql.CheckIfCourseExist(ctx, courseID)
	if !errors.Is(rsp, response.OperationSuccess) {
		return rsp
	}
	if !ok {
		return response.CourseNotExist
	}

	// 删除课程
	rsp = mysql.DeleteCourse(ctx, courseID)
	return rsp
}
