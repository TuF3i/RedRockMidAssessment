package redis

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"context"

	"github.com/cloudwego/hertz/pkg/common/json"
	"go.uber.org/zap"
)

func GetAllCourseID(ctx context.Context) (interface{}, response.Response) {
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
		return nil, response.ServerInternalError(err)
	}

	return ids, response.OperationSuccess
}

func GetCourseDetails(ctx context.Context, courseID string) (interface{}, response.Response) {
	var data models.Course
	// 生成Key
	key := courseInfoKey(courseID)
	// 取出课程的信息
	raw, err := core.RedisConn.HGetAll(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get Course Info",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}
	if len(raw) == 0 {
		return models.Course{ClassID: courseID}, response.OperationSuccess // 保险加个，防止Key不存在导致查不出来
	}
	// 把查出来的map[string]string转成jsonByte
	buf, err := json.Marshal(raw)
	if err != nil {
		core.Logger.Error(
			"Convert map to jsonByte",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}
	// 把jsonByte重新转成json（保护数据类型）
	if err := json.Unmarshal(buf, &data); err != nil {
		core.Logger.Error(
			"Convert jsonByte to json",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}
	return data, response.OperationSuccess
}

func GetStuSelectedCourseID(ctx context.Context, userID string) (interface{}, response.Response) {
	// 生成Key
	key := studentSelectedCourseKey(userID)
	// 取出所有可选课程的ID
	ids, err := core.RedisConn.SMembers(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get Stu Selected CourseID",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return ids, response.OperationSuccess
}
