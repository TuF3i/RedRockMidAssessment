package union_init

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"encoding/json"
	"strconv"

	"go.uber.org/zap"
)

func slice2interface(s []string) []interface{} {
	out := make([]interface{}, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

func PreloadCache(ctx context.Context) response.Response {
	/* 缓存预热 */
	// 初始化选课状态
	key := courseSelectionStatusKey()
	err := core.RedisConn.Set(context.Background(), key, "1", 0).Err()
	if err != nil {
		core.Logger.Error(
			"Preset Course Selection Status",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 从MySQL拿取所有的课程信息
	var data []models.Course
	tx := core.MysqlConn.Begin()
	if err := tx.Find(&data).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Course Info From MySQL",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 从Mysql拿取选课信息
	var selectData []models.Relation
	if err := tx.Find(&selectData).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Course Select Info From MySQL",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	tx.Commit()
	// 遍历填充数据
	var allIDs []string
	for _, item := range data {
		// 填充课程容量
		key := courseStockKey(item.ClassID)
		initial := item.ClassCapacity - item.ClassSelectedNum
		if err := core.RedisConn.SetNX(context.Background(), key, initial, 0).Err(); err != nil {
			core.Logger.Error(
				"Set Course Stock",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			return response.ServerInternalError(err)
		}
		// 填充课程信息
		key = courseInfoKey(item.ClassID)
		raw, err := json.Marshal(item)
		// 序列化JSON
		if err != nil {
			core.Logger.Error(
				"Set Course Info",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			return response.ServerInternalError(err)
		}
		// 写入Redis
		if err := core.RedisConn.Set(context.Background(), key, raw, 0).Err(); err != nil {
			core.Logger.Error(
				"Set Course Info",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			return response.ServerInternalError(err)
		}
		// 添加课程id
		allIDs = append(allIDs, item.ClassID)
	}
	// 将所有可选课程的ID写入Redis
	key = courseIDsKey()
	if err := core.RedisConn.SAdd(context.Background(), key, slice2interface(allIDs)...).Err(); err != nil {
		core.Logger.Error(
			"Set Course IDs",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 遍历填充选课信息
	for _, item := range selectData {
		if err := core.RedisConn.SAdd(context.Background(), courseUsersKey(item.CouID), item.StuID).Err(); err != nil {
			core.Logger.Error(
				"Set Course Selected Student Info",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			return response.ServerInternalError(err)
		}
		if err := core.RedisConn.SAdd(context.Background(), studentSelectedCourseKey(strconv.Itoa(int(item.StuID))), item.CouID).Err(); err != nil {
			core.Logger.Error(
				"Set Student Course Selection Info",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			return response.ServerInternalError(err)
		}
	}

	return response.OperationSuccess
}

func CleanUpCache(ctx context.Context) response.Response {
	// 清空数据库
	if err := core.RedisConn.FlushDB(ctx).Err(); err != nil {
		core.Logger.Error(
			"Flush Data Base",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}

	return response.OperationSuccess
}
