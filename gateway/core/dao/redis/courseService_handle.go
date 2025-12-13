package redis

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"errors"

	_ "embed"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// 加载Lua脚本

//go:embed script/subscribe.lua
var subscribe string

//go:embed script/drop.lua
var drop string

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
	stockKet := courseStockKey(courseID)
	// 取出课程的信息
	raw, err := core.RedisConn.Get(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get Course Info",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}
	// 取出课程容量
	remain, err := core.RedisConn.Get(ctx, stockKet).Int()
	if err != nil {
		core.Logger.Error(
			"Get Course Remaining",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}
	if len(raw) == 0 {
		return models.Course{ClassID: courseID}, response.OperationSuccess // 保险加个，防止Key不存在导致查不出来
	}
	//// 把查出来的map[string]string转成jsonByte
	//buf, err := json.Marshal(raw)
	//if err != nil {
	//	core.Logger.Error(
	//		"Convert map to jsonByte",
	//		zap.String("snowflake", ctx.Value("trace_id").(string)),
	//		zap.String("detail", err.Error()),
	//	)
	//	return nil, response.ServerInternalError(err)
	//}
	// 把jsonByte重新转成json（保护数据类型）
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		core.Logger.Error(
			"Convert jsonByte to json",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}
	// 重新填充已选人数
	data.ClassSelectedNum = data.ClassCapacity - uint(remain)
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

func SubscribeACourse(ctx context.Context, userID string, courseID string) response.Response {
	// 生成key
	keyForStu := courseUsersKey(courseID)
	keyForStock := courseStockKey(courseID)
	keyForStudentSelectedCourseKey := studentSelectedCourseKey(userID)
	keyForStuDropped := courseDroppedUsersKey(courseID)

	// 初始化脚本
	script := redis.NewScript(subscribe)
	// 执行脚本
	ok, err := script.Run(
		ctx,
		core.RedisConn,
		[]string{keyForStock, keyForStu, keyForStudentSelectedCourseKey, keyForStuDropped},
		[]string{userID, courseID},
	).Result()
	// 判断返回值
	if err != nil { // 先判断错误
		core.Logger.Error(
			"Subscribe Course",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	if ok.(int64) == 0 {
		return response.CourseIsFull
	}

	return response.OperationSuccess
}

func DropACourse(ctx context.Context, userID string, courseID string) response.Response {
	// 生成key
	keyForStu := courseUsersKey(courseID)
	keyForStock := courseStockKey(courseID)
	keyForStuDropped := courseDroppedUsersKey(courseID)
	keyForStudentSelectedCourseKey := studentSelectedCourseKey(userID)

	// 初始化脚本
	script := redis.NewScript(drop)
	// 执行脚本
	ok, err := script.Run(ctx,
		core.RedisConn,
		[]string{keyForStock, keyForStu, keyForStuDropped, keyForStudentSelectedCourseKey},
		[]string{userID, courseID},
	).Result()
	// 判断返回值
	if err != nil { // 先判断错误
		core.Logger.Error(
			"Drop Course",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	if ok.(int64) == 0 {
		return response.RecordNotExist
	}

	return response.OperationSuccess
}

func CheckIfCourseExist(ctx context.Context, courseID string) (bool, response.Response) {
	// 生成key
	key := courseIDsKey()
	// 读redis
	ok, err := core.RedisConn.SIsMember(ctx, key, courseID).Result()
	if err != nil {
		core.Logger.Error(
			"Check If Course Exist",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err)
	}
	// 判断是否存在
	return ok, response.OperationSuccess
}

func CheckIfCourseSelected(ctx context.Context, userID string, courseID string) (bool, response.Response) {
	// 生成key
	key := courseUsersKey(courseID)
	// 读redis
	ok, err := core.RedisConn.SIsMember(ctx, key, userID).Result()
	if err != nil {
		core.Logger.Error(
			"Check If Course Selected",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err)
	}
	// 判断是否存在
	return ok, response.OperationSuccess
}

func CheckIfCourseSelectionStarted(ctx context.Context) (bool, response.Response) { // 中间件直接调用
	// 生成key
	key := courseSelectionStatusKey()
	// 读Redis
	val, err := core.RedisConn.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, response.OperationSuccess
	}
	if err != nil {
		core.Logger.Error(
			"Check If Course Selection Start",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err)
	}
	// 判断是否开始选课
	if val == "1" {
		return true, response.OperationSuccess
	}

	return false, response.OperationSuccess
}

func UpdateCourseSelectionEventStatus(ctx context.Context, status bool) response.Response {
	// 1: 开始 0: 结束
	var Rstatus string
	// 生成key
	key := courseSelectionStatusKey()
	// 判断状态
	if status {
		Rstatus = "1"
	} else {
		Rstatus = "0"
	}
	// 写Redis
	err := core.RedisConn.Set(ctx, key, Rstatus, 0).Err()
	if err != nil {
		core.Logger.Error(
			"Update Course Selection Event Status",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	return response.OperationSuccess
}
