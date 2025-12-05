package dao

import (
	"RedRockMidAssessment-Synchronizer/core"
	"context"
	"strconv"

	"go.uber.org/zap"
)

func GetAllCourseID(ctx context.Context) interface{} {
	// 生成Key
	key := courseIDsKey()
	// 取出所有可选课程的ID
	ids, err := core.RedisConn.SMembers(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get All CourseID",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return ids
}

func CourseSelectedStu(ctx context.Context, courseID string) interface{} {
	// 生成key
	key := courseUsersKey(courseID)
	// 取出选了该课程的学生ID
	ids, err := core.RedisConn.SMembers(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get Selected StuID",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return ids
}

func CourseDroppedStu(ctx context.Context, courseID string) interface{} {
	// 生成key
	key := courseDroppedUsersKey(courseID)
	// 取出选了该课程的学生ID
	ids, err := core.RedisConn.SMembers(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get Dropped StuID",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return ids
}

func GetCourseStock(ctx context.Context, courseID string) interface{} {
	// 生成key
	key := courseStockKey(courseID)
	//查询课程容量
	stoRaw, err := core.RedisConn.Get(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get Course Stock",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}
	// 类型转换
	stock, err := strconv.ParseUint(stoRaw, 10, 64)
	if err != nil {
		core.Logger.Error(
			"Convert String To Uint",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return uint(stock)
}
