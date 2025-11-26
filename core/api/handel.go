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

func RegisterHandleFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var userForm models.Student

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

// TODO LoginHandleFunc

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
		// 将JWT中的UserID转为uint
		num, err := strconv.ParseUint(claims.UserID, 10, 32)
		core.Logger.Error(
			"Converting Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		// 调用stu_service
		data, rsp := service.GetStuInfo(ctx, uint(num))
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
		var updateData models.UpdateData
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
		// 将JWT中的UserID转为uint
		num, err := strconv.ParseUint(claims.UserID, 10, 32)
		core.Logger.Error(
			"Converting Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		//调用stu_service
		rsp := service.UpdateStuInfo(ctx, uint(num), updateData)
		c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
		return
	}
}
