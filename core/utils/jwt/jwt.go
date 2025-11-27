package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("Welcome_To_Red_Rock")
	refreshSecret = []byte("I_Want_To_Join_Red_Rock")
	accessTTL     = 15 * time.Minute   // 访问令牌有效期
	refreshTTL    = 7 * 24 * time.Hour // 刷新令牌有效期

	issuer = "StuClassMS.ApiNode"
)

type CustomClaims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"` // "admin" or "student"
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func GenTokens(userID string, role string) (accessToken string, refreshToken string, err error) {
	t := time.Now() //统一时间
	// AccessToken
	accessClaims := CustomClaims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  []string{"student", "admin"},
			ExpiresAt: jwt.NewNumericDate(t.Add(accessTTL)),
			NotBefore: jwt.NewNumericDate(t.Add(-5 * time.Second)), //此令牌在time.Now()之前一律无效
			IssuedAt:  jwt.NewNumericDate(t),                       // 签发时间
		},
	}
	//生成令牌
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret) //签名密钥
	if err != nil {
		return "", "", fmt.Errorf("gen accessToken: %v", err.Error())
	}

	// RefreshToken
	refreshClaims := CustomClaims{
		UserID: userID,
		Role:   role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  []string{"student", "admin"},
			ExpiresAt: jwt.NewNumericDate(t.Add(refreshTTL)),
			NotBefore: jwt.NewNumericDate(t.Add(-5 * time.Second)), //此令牌在time.Now()之前一律无效
			IssuedAt:  jwt.NewNumericDate(t),                       // 签发时间
		},
	}
	//生成令牌
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret) //签名密钥
	if err != nil {
		return "", "", fmt.Errorf("gen refershToken: %v", err.Error())
	}

	return accessToken, refreshToken, nil
}

func VerifyAccessToken(raw string) (*CustomClaims, error) {
	signMethodVerifier := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			// 检查签名方式是否合法以及签名方式是否为sha256
			return nil, fmt.Errorf("wrong signing method: %v", t.Header["alg"])
		}
		return accessSecret, nil
	}

	//解析 + 校验
	token, err := jwt.ParseWithClaims(
		raw,                           //原始Token
		&CustomClaims{},               //自定义Claim
		signMethodVerifier,            //签名方法校验函数
		jwt.WithLeeway(5*time.Second), //允许5秒内的时钟误差
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.Type != "access" {
		return nil, errors.New("not an access token")
	}

	return claims, nil
}

func VerifyRefreshToken(raw string) (*CustomClaims, error) {
	signMethodVerifier := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			// 检查签名方式是否合法以及签名方式是否为sha256
			return nil, fmt.Errorf("wrong signing method: %v", t.Header["alg"])
		}
		return refreshSecret, nil
	}

	//解析 + 校验
	token, err := jwt.ParseWithClaims(
		raw,                           //原始Token
		&CustomClaims{},               //自定义Claim
		signMethodVerifier,            //签名方法校验函数
		jwt.WithLeeway(5*time.Second), //允许5秒内的时钟误差
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("not an refresh token")
	}

	return claims, nil
}
