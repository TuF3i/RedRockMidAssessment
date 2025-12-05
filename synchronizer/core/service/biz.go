package service

import (
	"RedRockMidAssessment-Synchronizer/core/dao"
	"context"
)

func GetAllCourseID(ctx context.Context) interface{} {
	return dao.GetAllCourseID(ctx)
}

func UpdateSelectedStuList(ctx context.Context, courseID string) {
	// 获取已选课学生列表
	stuList := dao.CourseSelectedStu(ctx, courseID)
	if stuList == nil {
		return
	}
	// 类型断言
	list := stuList.([]string)
	for stu := range list {
		// 组装json消息

	}
}

func GetDroppedStu(ctx context.Context, courseID string) {
	return dao.CourseDroppedStu(ctx, courseID)
}

func GetCourseSelectedNum(ctx context.Context, courseID string) {
	return dao.GetCourseStock(ctx, courseID)
}
