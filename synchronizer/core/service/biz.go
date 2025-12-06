package service

import (
	"RedRockMidAssessment-Synchronizer/core/dao"
	"RedRockMidAssessment-Synchronizer/core/kafka"
	"RedRockMidAssessment-Synchronizer/core/utils/msg"
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
	s := stuList.([]string)
	// 遍历发送
	for _, stu := range s {
		// 组装json消息
		couMsg := msg.GenUpdateSelectedStuListMsg(ctx, stu, courseID)
		// kafka发送消息
		kafka.Writer(couMsg)
	}
}

func UpdateDroppedStuList(ctx context.Context, courseID string) {
	// 获取退课学生列表
	stuList := dao.CourseDroppedStu(ctx, courseID)
	if stuList == nil {
		return
	}
	// 类型断言
	s := stuList.([]string)
	// 遍历发送
	for _, stu := range s {
		// 组装json消息
		couMsg := msg.GenUpdateDroppedStuListMsg(ctx, stu, courseID)
		// kafka发送消息
		kafka.Writer(couMsg)
	}
}

func UpdateCourseSelectedNum(ctx context.Context, courseID string) {
	// 获取余量
	remaining := dao.GetCourseStock(ctx, courseID)
	if remaining == nil {
		return
	}
	// 类型断言
	r := remaining.(uint)
	// 组装json消息
	reMsg := msg.GenUpdateCourseSelectedNumMsg(ctx, courseID, r)
	// 发送json消息
	kafka.Writer(reMsg)
}
