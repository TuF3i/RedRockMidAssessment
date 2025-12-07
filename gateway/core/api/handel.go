package api

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/service"
	"RedRockMidAssessment/core/utils/jwt"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"errors"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"
)

func RefreshTokensHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		//生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		//调用调用stu_service
		data, rsp := service.RefreshTokens(ctx, claims.UserID, claims.Role)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
		return
	}
}

func LogoutHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		//生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		//调用调用stu_service
		rsp := service.Logout(ctx, claims.UserID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
		return
	}
}

func LoginHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var userForm models.LoginForm
		//生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		//校验JSON
		if err := c.BindAndValidate(&userForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.ServerInternalError(err), nil))
			return
		}
		//调用调用stu_service
		data, rsp := service.Login(ctx, userForm)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
		return
	}
}

func RegisterHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var userForm models.Student
		//生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		//校验JSON
		if err := c.BindAndValidate(&userForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		//调用stu_service
		rsp := service.AddStudent(ctx, userForm)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
		return
	}
}

func GetStudentInfoForStuHandleFunc() app.HandlerFunc {
	// Permission: Student
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "student" { // 不可以拿admin来调用给学生的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 调用stu_service
		data, rsp := service.GetStuInfo(ctx, claims.UserID)
		if !errors.Is(rsp, response.OperationSuccess) {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			return
		}
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
		return
	}
}

func UpdateStudentInfoForStuHandleFunc() app.HandlerFunc {
	// Permission: student
	return func(ctx context.Context, c *app.RequestContext) {
		var updateData models.UpdateDataForStu
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "student" { // 不可以拿admin来调用给学生的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&updateData); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		//调用stu_service
		rsp := service.UpdateStuInfo(ctx, claims.UserID, updateData.UpdateColumns)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
		return
	}
}

func GetStudentListForAdminHandleFunc() app.HandlerFunc {
	// Permission: admin
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 获取Query参数
		page, err := strconv.Atoi(c.DefaultQuery("page", "1")) //查询结果的第几页
		if err != nil {
			core.Logger.Error(
				"Converting Error",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.ServerInternalError(err), nil))
			return
		}
		resNum, err := strconv.Atoi(c.DefaultQuery("resNum", "15")) //每页结果的条数
		if err != nil {
			core.Logger.Error(
				"Converting Error",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.ServerInternalError(err), nil))
			return
		}
		//调用stu_service
		data, rsp := service.GetStudentsList(ctx, page, resNum)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
		return
	}
}

func UpdateStudentInfoForAdminHandleFunc() app.HandlerFunc {
	// Permission: admin
	return func(ctx context.Context, c *app.RequestContext) {
		var updateData models.UpdateDataForStu
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&updateData); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用stu_service
		rsp := service.UpdateStuInfo(ctx, updateData.StudentID, updateData.UpdateColumns)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func AddStudentForAdminHandleFunc() app.HandlerFunc {
	// Permission: admin
	return func(ctx context.Context, c *app.RequestContext) {
		var userForm models.Student
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&userForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		//调用stu_service
		rsp := service.AddStudent(ctx, userForm)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
		return
	}
}

func DelStudentForAdminHandleFunc() app.HandlerFunc {
	// Permission: admin
	return func(ctx context.Context, c *app.RequestContext) {
		var userForm models.DelStudentForm
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&userForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		//调用stu_service
		rsp := service.DeleteStudent(ctx, userForm.StuID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func GetCourseInfoForStudentHandleFunc() app.HandlerFunc {
	// Permission: student
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "student" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 调用course_service
		data, rsp := service.GetCourseInfo(ctx)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
	}
}

func GetSelectedCourseForStudentHandleFunc() app.HandlerFunc {
	// Permission: student
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "student" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 调用course_service
		data, rsp := service.GetStuSelectedCourses(ctx, claims.UserID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
	}
}

func SubscribeCourseForStudentHandleFunc() app.HandlerFunc {
	// Permission: student
	return func(ctx context.Context, c *app.RequestContext) {
		var courseForm models.CourseForm
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "student" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&courseForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用course_service
		rsp := service.SubscribeCourse(ctx, claims.UserID, courseForm.ClassID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func DropCourseForStudentHandleFunc() app.HandlerFunc {
	// Permission: student
	return func(ctx context.Context, c *app.RequestContext) {
		var courseForm models.CourseForm
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "student" { // 不可以拿student来调用给admin的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&courseForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用course_service
		rsp := service.DropCourse(ctx, claims.UserID, courseForm.ClassID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func GetCourseInfoForAdminHandleFunc() app.HandlerFunc {
	// Permission: admin
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 调用course_service
		data, rsp := service.GetCourseInfoForAdmin(ctx)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
	}
}

func GetStuCourseSelectionInfoForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 获取Param参数
		stuID := c.Param("stuID")
		// 调用course_service
		data, rsp := service.GetStuSelectedCoursesForAdmin(ctx, stuID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, data))
	}
}

func AddStuCourseSelectionInfoForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var updateData models.UpdateCourseData
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&updateData); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用 course_service
		rsp := service.SubscribeCourseForAdmin(ctx, updateData.StuId, updateData.UpdateClassId)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func DelStuCourseSelectionInfoForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var updateData models.UpdateCourseData
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&updateData); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用 course_service
		rsp := service.DropCourseForAdmin(ctx, updateData.StuId, updateData.UpdateClassId)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func UpdateCourseInfoForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var updateData models.UpdateDataForCourse
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&updateData); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用course_service
		rsp := service.UpdateCourseInfoForAdmin(ctx, updateData.ClassID, updateData.UpdateColumns)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func UpdateCourseStockForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var updateData models.UpdateDataForCourseStock
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&updateData); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用course_service
		rsp := service.UpdateCourseStockForAdmin(ctx, updateData.ClassID, updateData.Stock)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func AddCourseForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var dataForm models.Course
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&dataForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用course_service
		rsp := service.AddCourseForAdmin(ctx, dataForm)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func DelCourseForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var dataForm models.DeleteDataForCourse
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		//校验JSON
		if err := c.BindAndValidate(&dataForm); err != nil {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.RevDataError, nil))
			return
		}
		// 调用course_service
		rsp := service.DelCourseForAdmin(ctx, dataForm.ClassID)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func StartCourseSelectionEventForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 调用course_service
		rsp := service.StartCourseSelection(ctx)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}

func StopCourseSelectionEventForAdminHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 解析JWT
		rawClaims, _ := c.Get("jwt_claims")
		claims := rawClaims.(jwt.CustomClaims)
		// 判断权限
		if claims.Role != "admin" { // 不可以拿admin来调用给student的接口，避免权限混乱
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.PermissionDenied, nil))
			return
		}
		// 调用course_service
		rsp := service.StopCourseSelection(ctx)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
	}
}
