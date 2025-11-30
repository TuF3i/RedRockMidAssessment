package service

import (
	"RedRockMidAssessment/core/dao/redis"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"RedRockMidAssessment/core/utils/verify"
	"context"
	"errors"
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
