package service

import (
	"RedRockMidAssessment/core/dao/redis"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
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
