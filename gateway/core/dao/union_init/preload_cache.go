package union_init

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"context"

	"github.com/fatih/structs"
)

func slice2interface(s []string) []interface{} {
	out := make([]interface{}, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

func PreloadCache() error {
	/* 缓存预热 */
	// 初始化选课状态
	key := courseSelectionStatusKey()
	err := core.RedisConn.Set(context.Background(), key, "0", 0).Err()
	if err != nil {
		return err
	}
	// 从MySQL拿取所有的课程信息
	var data []models.Course
	tx := core.MysqlConn.Begin()
	if err := tx.Find(&data).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// 遍历填充数据
	var allIDs []string
	for _, item := range data {
		// 填充课程容量
		key := courseStockKey(item.ClassID)
		initial := item.ClassCapacity
		if err := core.RedisConn.SetNX(context.Background(), key, initial, 0).Err(); err != nil {
			return err
		}
		// 填充课程信息
		dataMap := structs.Map(item)
		key = courseInfoKey(item.ClassID)
		if err := core.RedisConn.HSet(context.Background(), key, dataMap).Err(); err != nil {
			return err
		}
		// 添加课程id
		allIDs = append(allIDs, item.ClassID)
	}
	// 将所有可选课程的ID写入Redis
	key = courseIDsKey()
	if err := core.RedisConn.SAdd(context.Background(), key, slice2interface(allIDs)...).Err(); err != nil {
		return err
	}

	return nil
}
