package redis

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/utils/response"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

func AddTokenToRedis(ctx context.Context, userID string, tokenType uint, uuid string, expireTime time.Duration) response.Response {
	// tokenType: 0为AccessToken, 1为RefreshToken
	var key string
	if tokenType == 0 {
		key = accessTokenKey(userID)
	} else if tokenType == 1 {
		key = refreshTokenKey(userID)
	} else {
		err := errors.New("undefined token type")
		core.Logger.Error(
			"Add Token To Redis Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}

	// 写入Token
	_, err := core.RedisConn.Set(ctx, key, uuid, expireTime).Result()
	if err != nil {
		core.Logger.Error(
			"Add Token To Redis Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	return response.OperationSuccess
}

func DelTokenFromRedis(ctx context.Context, userID string, tokenType uint) response.Response {
	// tokenType: 0为AccessToken, 1为RefreshToken
	var key string
	if tokenType == 0 {
		key = accessTokenKey(userID)
	} else if tokenType == 1 {
		key = refreshTokenKey(userID)
	} else {
		err := errors.New("undefined token type")
		core.Logger.Error(
			"Delete Token From Redis Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	// 写入Token
	_, err := core.RedisConn.Del(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Delete Token From Redis Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return response.ServerInternalError(err)
	}
	return response.OperationSuccess
}

func IfTokenExist(ctx context.Context, userID string, tokenType uint) (bool, response.Response) {
	// tokenType: 0为AccessToken, 1为RefreshToken
	var key string
	if tokenType == 0 {
		key = accessTokenKey(userID)
	} else if tokenType == 1 {
		key = refreshTokenKey(userID)
	} else {
		err := errors.New("undefined token type")
		core.Logger.Error(
			"Delete Token From Redis Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err)
	}
	// 检查是否存在
	val, err := core.RedisConn.Exists(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Check If The Key Exists",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err)
	}

	if val == 0 {
		return false, response.OperationSuccess
	}

	return true, response.OperationSuccess
}

func GetUUIDOfToken(ctx context.Context, userID string, tokenType uint) (interface{}, response.Response) {
	// tokenType: 0为AccessToken, 1为RefreshToken
	var key string
	if tokenType == 0 {
		key = accessTokenKey(userID)
	} else if tokenType == 1 {
		key = refreshTokenKey(userID)
	} else {
		err := errors.New("undefined token type")
		core.Logger.Error(
			"Delete Token From Redis Error",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return false, response.ServerInternalError(err)
	}
	// 取值
	val, err := core.RedisConn.Get(ctx, key).Result()
	if err != nil {
		core.Logger.Error(
			"Get UUID From Key",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil, response.ServerInternalError(err)
	}

	return val, response.OperationSuccess
}
