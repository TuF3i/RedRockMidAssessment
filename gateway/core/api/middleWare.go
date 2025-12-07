package api

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/dao/redis"
	"RedRockMidAssessment/core/utils/jwt"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"
)

func JWTAuthMiddleWare() app.HandlerFunc {
	// JWT验证中间件
	return func(ctx context.Context, c *app.RequestContext) {
		//生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 获取Token
		authorizationToken := string(c.GetHeader("Authorization"))
		// 检测Token是否为空
		if authorizationToken == "" {
			c.JSON(consts.StatusOK, response.EmptyToken)
			c.Abort()
			return
		}
		// 检测Token格式
		token := strings.TrimPrefix(authorizationToken, "Bearer ")
		if token == authorizationToken { // 没前缀
			c.JSON(consts.StatusOK, response.InvalidToken)
			c.Abort()
			return
		}
		// 提取Claim信息，并校验Token是否有效
		claim, err := jwt.VerifyAccessToken(token)
		if err != nil {
			core.Logger.Error(
				"Get Claim Error",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.ServerInternalError(err), nil))
			c.Abort()
			return
		}

		userID := claim.UserID
		uuid := claim.UUID
		// 从Redis中校验Token
		ok, rsp := redis.IfTokenExist(ctx, userID, 0) // 是否存在该AccessToken
		if !ok {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			c.Abort()
			return
		}
		data, rsp := redis.GetUUIDOfToken(ctx, userID, 0) // 获取Token的UUID
		if !errors.Is(rsp, response.OperationSuccess) {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			c.Abort()
			return
		}

		value, ok := data.(string)
		if value == "" || !ok { // 类型断言
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.InvalidToken, nil))
			c.Abort()
			return
		}

		if value != uuid { // 判断claim内uuid是否与redis内uuid相同
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.InvalidToken, nil))
			c.Abort()
			return
		}

		// 将Token写入上下文
		c.Set("jwt_claims", claim)
		c.Next(context.Background()) // 以后再改
	}
}

func JWTRefreshMiddleWare() app.HandlerFunc {
	// JWT刷新中间件
	return func(ctx context.Context, c *app.RequestContext) {
		//生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 获取Token
		authorizationToken := string(c.GetHeader("Authorization"))
		// 检测Token是否为空
		if authorizationToken == "" {
			c.JSON(consts.StatusOK, response.EmptyToken)
			c.Abort()
			return
		}
		// 检测Token格式
		token := strings.TrimPrefix(authorizationToken, "Bearer ")
		if token == authorizationToken { // 没前缀
			c.JSON(consts.StatusOK, response.InvalidToken)
			c.Abort()
			return
		}
		// 提取Claim信息，并校验Token是否有效
		claim, err := jwt.VerifyRefreshToken(token)
		if err != nil {
			core.Logger.Error(
				"Get Claim Error",
				zap.String("snowflake", ctx.Value("trace_id").(string)),
				zap.String("detail", err.Error()),
			)
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.ServerInternalError(err), nil))
			c.Abort()
			return
		}

		userID := claim.UserID
		uuid := claim.UUID
		// 从Redis中校验Token
		ok, rsp := redis.IfTokenExist(ctx, userID, 1) // 是否存在该RefreshToken
		if !ok {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			c.Abort()
			return
		}
		data, rsp := redis.GetUUIDOfToken(ctx, userID, 1) // 获取Token的UUID
		if !errors.Is(rsp, response.OperationSuccess) {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			c.Abort()
			return
		}

		value, ok := data.(string)
		if value == "" || !ok { // 类型断言
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.InvalidToken, nil))
			c.Abort()
			return
		}

		if value != uuid { // 判断claim内uuid是否与redis内uuid相同
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.InvalidToken, nil))
			c.Abort()
			return
		}

		// 将Token写入上下文
		c.Set("jwt_claims", claim)
		c.Next(context.Background()) // 以后再改
	}
}

func CheckIfCourseSelectionStartedMiddleWare() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 查Redis
		ok, rsp := redis.CheckIfCourseSelectionStarted(ctx)
		if !errors.Is(rsp, response.OperationSuccess) {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			c.Abort()
			return
		}
		// 判断选课是否开始
		if !ok {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.CourseSelectionEventNotStart, nil))
			c.Abort()
			return
		}
		c.Next(context.Background()) // 以后再改
	}
}

func CheckIfCourseSelectionStartedMiddleWareForAdmin() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成TraceID
		traceID := core.SnowFlake.TraceID()
		ctx = context.WithValue(ctx, "trace_id", traceID)
		// 查Redis
		ok, rsp := redis.CheckIfCourseSelectionStarted(ctx)
		if !errors.Is(rsp, response.OperationSuccess) {
			c.JSON(consts.StatusOK, response.GenFinalResponse(rsp, nil))
			c.Abort()
			return
		}
		// 判断选课是否开始
		if ok {
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.NoChangingData, nil))
			c.Abort()
			return
		}
		c.Next(context.Background()) // 以后再改
	}
}
