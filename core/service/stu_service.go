package service

import (
	dao "RedRockMidAssessment/core/dao/mysql"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"RedRockMidAssessment/core/utils/verify"
	"context"
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

	/* MySQL写库 */
	return dao.InsertStudentIntoDB(ctx, userForm) // 直接上抛来自dao层的结果
}

func GetStuInfo(ctx context.Context, userID uint) (models.Student, response.Response) {
	/* 校验数据 */
	// 检查用户ID是否可用
	if !verify.VerifyUserID(userID) {
		return models.Student{}, response.InvalidStudentID
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
		case "password":
			v := value.(string)
			if !verify.VerifyPassword(v) {
				return response.InvalidPassword
			}
			field = append(field, column)
			dataList[column] = value
		case "stu_id":
			v := value.(uint)
			if !verify.VerifyUserID(v) {
				return response.InvalidStudentID
			}
			field = append(field, column)
			dataList[column] = value
		case "name":
			v := value.(string)
			if !verify.VerifyUserName(v) {
				return response.InvalidUserName
			}
			field = append(field, column)
			dataList[column] = value
		case "stu_class":
			v := value.(string)
			if !verify.VerifyStudentClass(v) {
				return response.InvalidClass
			}
			field = append(field, column)
			dataList[column] = value
		case "sex":
			v := value.(uint)
			if !verify.VerifySexSetting(v) {
				return response.InvalidSexSetting
			}
			field = append(field, column)
			dataList[column] = value
		case "grade":
			v := value.(uint)
			if !verify.VerifyGrade(v) {
				return response.InvalidGrade
			}
			field = append(field, column)
			dataList[column] = value
		case "age":
			v := value.(uint)
			if !verify.VerifyAge(v) {
				return response.InvalidAge
			}
			field = append(field, column)
			dataList[column] = value
		}
	} // for

	if len(field) == 0 {
		return response.EmptyData
	}

	return dao.UpdateStudentInfo(ctx, userID, field, dataList)
}
