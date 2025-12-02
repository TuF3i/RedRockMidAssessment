package mysql

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"context"

	"go.uber.org/zap"
)

func CheckIfCourseSelected(ctx context.Context, courseID string, userID string) (bool, response.Response) {
	// 查询课程是否存在
	tx := core.MysqlConn.Begin() // 开启数据库事务
	defer tx.Commit()            // 查询结束后提交

	result := tx.Where("class_id = ? AND student_id = ?", courseID, userID).Find(&models.Relation{}) // 用courseID查询学生信息
	if err := result.Error; err != nil {
		core.Logger.Error(
			"Check Course Selected Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err) // 系统错误上抛
	}

	if result.RowsAffected == 0 { // 验证课程是否存在
		return false, response.OperationSuccess
	} else {
		return true, response.OperationSuccess
	}
}

func CheckIfCourseExist(ctx context.Context, courseID string) (bool, response.Response) {
	// 查询课程是否存在
	tx := core.MysqlConn.Begin() // 开启数据库事务
	defer tx.Commit()            // 查询结束后提交

	result := tx.Where("class_id = ?", courseID).Find(&models.Course{}) // 用courseID查询学生信息
	if err := result.Error; err != nil {
		core.Logger.Error(
			"Check Course Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err) // 系统错误上抛
	}

	if result.RowsAffected == 0 { // 验证课程是否存在
		return false, response.OperationSuccess
	} else {
		return true, response.OperationSuccess
	}
}

func InsertCourseIntoDB(ctx context.Context, courseForm models.Course) response.Response {
	// 插入课程
	tx := core.MysqlConn.Begin() // 开启数据库事务
	if err := tx.Create(&courseForm).Error; err != nil {
		tx.Rollback() // 错误回滚
		core.Logger.Error(
			"Insert Course Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	tx.Commit() // 成功提交
	return response.OperationSuccess
}

func UpdateCourseInfo(ctx context.Context, courseID string, field []string, dataList map[string]interface{}) response.Response {
	tx := core.MysqlConn.Begin()
	if err := tx.Model(&models.Student{}).Where("class_id = ?", courseID).Select(field).Updates(dataList).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Update Student Info Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	tx.Commit()
	return response.OperationSuccess
}

func DeleteCourse(ctx context.Context, courseID string) response.Response {
	tx := core.MysqlConn.Begin()
	if err := tx.Where("class_id = ?", courseID).Delete(&models.Student{}).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Delete Course Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	tx.Commit()
	return response.OperationSuccess
}

func GetAllCourseInfo(ctx context.Context) (interface{}, response.Response) {
	var data []models.Course
	tx := core.MysqlConn.Begin()
	defer tx.Commit()
	if err := tx.Find(&data).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Course List Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return data, response.OperationSuccess
}

func GetStudentSelectedCourse(ctx context.Context, userID string) (interface{}, response.Response) {
	var data models.Student
	tx := core.MysqlConn.Begin()
	defer tx.Commit()

	if err := tx.Preload("Courses").Where("student_id = ?", userID).First(&data).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Student Selected Course Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return data.Courses, response.OperationSuccess
}

func AddCourseToStudent(ctx context.Context, courseID string, userID string) response.Response {
	var course models.Course
	var student models.Student

	//开启数据库事务
	tx := core.MysqlConn.Begin()
	// 查询学生
	if err := tx.Where("student_id = ?", userID).First(&student).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Add Course To Student",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 查询课程
	if err := tx.Where("class_id = ?", courseID).First(&course).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Add Course To Student",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 建立关联
	if err := tx.Model(&student).Association("Courses").Append(&course); err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Add Course To Student",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}

	tx.Commit()
	return response.OperationSuccess
}

func DelCourseToStudent(ctx context.Context, courseID string, userID string) response.Response {
	var course models.Course
	var student models.Student

	//开启数据库事务
	tx := core.MysqlConn.Begin()
	// 查询学生
	if err := tx.Where("student_id = ?", userID).First(&student).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Add Course To Student",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 查询课程
	if err := tx.Where("class_id = ?", courseID).First(&course).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Add Course To Student",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 建立关联
	if err := tx.Model(&student).Association("Courses").Delete(&course); err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Add Course To Student",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}

	tx.Commit()
	return response.OperationSuccess
}

func GetCourseInfo(ctx context.Context, courseID string) (interface{}, response.Response) {
	var data models.Course
	tx := core.MysqlConn.Begin()
	defer tx.Commit()
	if err := tx.Where("class_id = ?", courseID).Find(&data).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Course Info Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return data, response.OperationSuccess
}

func UpdateCourseStock(ctx context.Context, courseID string, stock uint) response.Response {
	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	// 更新课程容量
	if err := tx.Model(&models.Course{}).Where("class_id = ?", courseID).Update("class_capacity", stock).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Course Info Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}

	return response.OperationSuccess
}
