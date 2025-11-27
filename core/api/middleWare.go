package api

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/utils/jwt"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"
)

func JWTAuthMiddleWare() app.HandlerFunc {
	// JWT验证中间件
	return func(ctx context.Context, c *app.RequestContext) {
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
				zap.String("middleWareID", "MiddleWare_JWTAuthMiddleWare_Error"),
				zap.String("detail", err.Error()),
			)
			c.JSON(consts.StatusOK, response.GenFinalResponse(response.ServerInternalError(err), nil))
			c.Abort()
			return
		}
		// 将Token写入上下文
		c.Set("jwt_claims", claim)
		return
	}
}
