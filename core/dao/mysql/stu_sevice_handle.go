package dao

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"errors"

	"go.uber.org/zap"
)

func CheckIfStudentExist(ctx context.Context, studentID uint) (bool, response.Response) {
	// 查询学生是否存在
	tx := core.MysqlConn.Begin() // 开启数据库事务
	defer tx.Commit()            // 查询结束后提交

	result := tx.Where("stu_id = ?", studentID).Find(&models.Student{}) // 用StuID查询学生信息
	if err := result.Error; err != nil {
		core.Logger.Error(
			"Insert Student Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err) // 系统错误上抛
	}

	if result.RowsAffected == 0 { // 验证学生是否存在
		return false, response.OperationSuccess
	} else {
		return true, response.OperationSuccess
	}

}

func InsertStudentIntoDB(ctx context.Context, userForm models.Student) response.Response {
	// 插入学生
	ifExist, rsp := CheckIfStudentExist(ctx, userForm.StudentID) // 获取学生
	if !errors.Is(rsp, response.OperationSuccess) {              // 出现错误直接上抛
		return rsp
	}

	if ifExist { // 检测学生是否存在
		return response.StudentIDAlreadyExist
	}

	tx := core.MysqlConn.Begin() // 开启数据库事务
	if err := tx.Create(&userForm).Error; err != nil {
		tx.Rollback() // 错误回滚
		core.Logger.Error(
			"Insert Student Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	tx.Commit() // 成功提交
	return response.OperationSuccess
}

func GetStudentInfo(ctx context.Context, userID uint) (models.Student, response.Response) {
	var studentData models.Student

	ifExist, rsp := CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return studentData, rsp
	}

	if ifExist {
		return studentData, response.StudentIDAlreadyExist
	}

	tx := core.MysqlConn.Begin()
	if err := tx.Where("stu_id = ?", userID).First(&studentData).Error; err != nil {
		tx.Rollback()
		core.Logger.Error(
			"Get Student Info Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return studentData, response.ServerInternalError(err)
	}

	tx.Commit()
	return studentData, response.OperationSuccess
}

func UpdateStudentInfo(ctx context.Context, userID uint, field []string, dataList map[string]interface{}) response.Response {
	ifExist, rsp := CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return rsp
	}

	if !ifExist {
		return response.UserNotExiOrWrongStuID
	}

	tx := core.MysqlConn.Begin()
	if err := tx.Model(&models.Student{}).Where("stu_id = ?", userID).Select(field).Updates(dataList).Error; err != nil {
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
