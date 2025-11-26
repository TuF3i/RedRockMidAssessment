package service

import (
	dao "RedRockMidAssessment/core/dao/mysql"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"RedRockMidAssessment/core/utils/verify"
	"context"
	"errors"
)

func AddStudent(ctx context.Context, userForm models.Student) response.Response {
	/* 校验数据 */
	// 检查用户名是否可用
	if !verify.VerifyUserName(userForm.Name) {
		return response.InvalidUserName
	}

	// 校验用户ID
	if !verify.VerifyUserID(userForm.StudentID) {
		return response.InvalidStudentID
	}

	// 校验班级字符串
	if !verify.VerifyStudentClass(userForm.StudentClass) {
		return response.InvalidClass
	}
	// 校验密码
	if !verify.VerifyPassword(userForm.Password) {
		return response.InvalidPassword
	}

	// 校验性别
	if !verify.VerifySexSetting(userForm.Sex) {
		return response.InvalidSexSetting
	}

	// 校验年级
	if !verify.VerifyGrade(userForm.Grade) {
		return response.InvalidGrade
	}

	// 校验年龄
	if !verify.VerifyAge(userForm.Age) {
		return response.InvalidAge
	}

	/* 检查学生是否存在 */
	ifExist, rsp := dao.CheckIfStudentExist(ctx, userForm.StudentID) // 获取学生
	if !errors.Is(rsp, response.OperationSuccess) {                  // 出现错误直接上抛
		return rsp
	}

	if ifExist { // 检测学生是否存在
		return response.StudentIDAlreadyExist
	}

	/* MySQL写库 */
	return dao.InsertStudentIntoDB(ctx, userForm) // 直接上抛来自dao层的结果
}

func GetStuInfo(ctx context.Context, userID uint) (models.Student, response.Response) {
	/* 校验数据 */
	// 检查用户ID是否可用
	if !verify.VerifyUserID(userID) {
		return models.Student{}, response.InvalidStudentID
	}
	/* 判断学生是否存在 */
	ifExist, rsp := dao.CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return models.Student{}, rsp
	}

	if !ifExist {
		return models.Student{}, response.UserNotExiOrWrongStuID
	}

	/* 查MySQL */
	return dao.GetStudentInfo(ctx, userID)
}

func UpdateStuInfo(ctx context.Context, userID uint, data models.UpdateData) response.Response {
	field := make([]string, 0)
	dataList := make(map[string]interface{}, 0)
	for _, item := range data.UpdateColumns {
		column := item.Field
		value := item.Value
		switch column {
		case "password": // 校验密码规范
			v := value.(string)
			if !verify.VerifyPassword(v) {
				return response.InvalidPassword
			}
			field = append(field, column)
			dataList[column] = value
		case "stu_id": // 校验学生ID规范
			v := value.(uint)
			if !verify.VerifyUserID(v) {
				return response.InvalidStudentID
			}
			field = append(field, column)
			dataList[column] = value
		case "name": // 校验名字规范
			v := value.(string)
			if !verify.VerifyUserName(v) {
				return response.InvalidUserName
			}
			field = append(field, column)
			dataList[column] = value
		case "stu_class": // 校验班级字符串规范
			v := value.(string)
			if !verify.VerifyStudentClass(v) {
				return response.InvalidClass
			}
			field = append(field, column)
			dataList[column] = value
		case "sex": // 校验性别设置规范
			v := value.(uint)
			if !verify.VerifySexSetting(v) {
				return response.InvalidSexSetting
			}
			field = append(field, column)
			dataList[column] = value
		case "grade": // 校验年级设置规范
			v := value.(uint)
			if !verify.VerifyGrade(v) {
				return response.InvalidGrade
			}
			field = append(field, column)
			dataList[column] = value
		case "age": // 校验年龄设置规范
			v := value.(uint)
			if !verify.VerifyAge(v) {
				return response.InvalidAge
			}
			field = append(field, column)
			dataList[column] = value
		}
	} // for

	/* 空数据直接返回 */
	if len(field) == 0 {
		return response.EmptyData
	}

	/* 检查学生是否存在 */
	ifExist, rsp := dao.CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return rsp
	}

	if !ifExist {
		return response.UserNotExiOrWrongStuID
	}

	/* 写MySQL */
	return dao.UpdateStudentInfo(ctx, userID, field, dataList)
}
