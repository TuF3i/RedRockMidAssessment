package service

import (
	"RedRockMidAssessment/core"
	dao "RedRockMidAssessment/core/dao/mysql"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/jwt"
	"RedRockMidAssessment/core/utils/md5"
	"RedRockMidAssessment/core/utils/response"
	"RedRockMidAssessment/core/utils/verify"
	"context"
	"errors"
	"reflect"
	"strconv"

	"go.uber.org/zap"
)

func Login(ctx context.Context, userForm models.LoginForm) (interface{}, response.Response) {
	userID := userForm.StuID
	password := userForm.Password
	/* 校验数据 */
	// 校验用户名是否可用
	if !verify.VerifyUserID(userForm.StuID) {
		return nil, response.InvalidStudentID
	}
	// 校验密码是否可用
	if !verify.VerifyPassword(userForm.Password) {
		return nil, response.InvalidPassword
	}

	/* 检查学生是否存在 */
	ifExist, rsp := dao.CheckIfStudentExist(ctx, userForm.StuID) // 获取学生
	if !errors.Is(rsp, response.OperationSuccess) {              // 出现错误直接上抛
		return nil, rsp
	}

	if ifExist { // 检测学生是否存在
		return nil, response.StudentIDAlreadyExist
	}

	/* MySQL读库 */
	data, rsp := dao.GetStudentInfo(ctx, userID)
	if data == nil {
		return nil, rsp
	}

	typedData := data.(models.Student)

	/* 校验密码 */
	if typedData.Password != md5.GenMD5(password) {
		return nil, response.WrongPassword
	}

	/* 校验角色 */
	role := func() string {
		if typedData.Role {
			return "admin"
		} else {
			return "student"
		}
	}()

	/* 生成令牌 */
	ID := strconv.FormatUint(uint64(userID), 10)
	accessToken, refreshToken, err := jwt.GenTokens(ID, role)
	if err != nil {
		core.Logger.Error(
			"Generate Token Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return models.LoginRsp{ // 构造data
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, rsp

}

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
	rsp = dao.InsertStudentIntoDB(ctx, userForm)
	return rsp // 直接上抛来自dao层的结果
}

func RefreshTokens(ctx context.Context, userID string, role string) (interface{}, response.Response) {
	/* 生成新Token */
	accessToken, refreshToken, err := jwt.GenTokens(userID, role)
	if err != nil {
		core.Logger.Error(
			"Generate Token Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return models.LoginRsp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, response.OperationSuccess
}

func GetStuInfo(ctx context.Context, userID uint) (interface{}, response.Response) {
	/* 校验数据 */
	// 检查用户ID是否可用
	if !verify.VerifyUserID(userID) {
		return nil, response.InvalidStudentID
	}
	/* 判断学生是否存在 */
	ifExist, rsp := dao.CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return nil, rsp
	}

	if !ifExist {
		return nil, response.UserNotExiOrWrongStuID
	}

	/* 查MySQL */
	data, rsp := dao.GetStudentInfo(ctx, userID)
	return data, rsp
}

func UpdateStuInfo(ctx context.Context, userID uint, data []models.UpdateColumnsEntity) response.Response {
	field := make([]string, 0)
	dataList := make(map[string]interface{})
	for _, item := range data {
		v := reflect.ValueOf(item)
		t := reflect.TypeOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			t = t.Elem()
		}
		for i := 0; i < v.NumField(); i++ {
			value := v.Field(i).Interface()
			column := t.Field(i).Tag.Get("json")
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
	rsp = dao.UpdateStudentInfo(ctx, userID, field, dataList)
	return rsp
}

func GetStudentsList(ctx context.Context, page int, resNum int) (interface{}, response.Response) {
	//var students models.Students
	/* 数据校验 */
	if page <= 0 {
		page = 1
	}
	if resNum <= 0 {
		resNum = 10
	}

	/* 计算偏移量 */
	offset := (page - 1) * resNum

	/* 查MySQL */
	data, total, rsp := dao.GetStudentList(ctx, resNum, offset, page)

	return models.Students{ // 组装data
		Total:        total,
		Page:         page,
		PageSize:     resNum,
		StudentsList: data.([]models.StudentsListEntity), // 类型断言
	}, rsp
}

func DeleteStudent(ctx context.Context, userID uint) response.Response {
	/* 检查学生是否存在 */
	ifExist, rsp := dao.CheckIfStudentExist(ctx, userID)
	if !errors.Is(rsp, response.OperationSuccess) { // 出现错误直接上抛
		return rsp
	}

	if !ifExist {
		return response.UserNotExiOrWrongStuID
	}

	/* 删数据库记录 */
	rsp = dao.DeleteStudent(ctx, userID)
	return rsp
}
