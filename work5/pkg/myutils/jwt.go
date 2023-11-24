package myutils

import (
	"five/consts"
	"five/pkg/log"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	UserName string
	UID      uint
	jwt.StandardClaims
}

// CheckToken 用于jwt中间件里检查token
func CheckToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	aClaims, err, aValid := ParseToken(aToken)
	if err != nil {
		log.LogrusObj.Println("atoken", err)
		return
	}
	rClaims, err, rValid := ParseToken(rToken)
	if err != nil {
		log.LogrusObj.Println("rtoken", err)
		return
	}
	// 如果aToken和rToken都不过期就只更新aToken
	if rValid && aValid {
		newAToken, err = GenerateAccessToken(aClaims.UserName, aClaims.UID)
		newRToken = rToken
		return
	}
	// 如果aToken过期 但是rToken不过期就只更新aToken
	if rValid && !aValid {
		newAToken, err = GenerateAccessToken(aClaims.UserName, aClaims.UID)
		newRToken = rToken
		return
	}
	// 全更新
	newAToken, err = GenerateAccessToken(aClaims.UserName, aClaims.UID)
	if err != nil {
		return
	}
	newRToken, err = GenerateRefreshToken(rClaims.UserName, rClaims.UID)
	if err != nil {
		return
	}
	return
}

func GenerateToken(userName string, id uint) (aToken, rToken string, err error) {
	aToken, err = GenerateAccessToken(userName, id)
	if err != nil {
		return "", "", err
	}
	rToken, err = GenerateRefreshToken(userName, id)
	if err != nil {
		return "", "", err
	}
	return
}

func GenerateAccessToken(userName string, id uint) (aToken string, err error) {
	expire := time.Now().Add(consts.AccessTokenExpireTime).Unix()
	claims := &Claims{
		UserName: userName,
		UID:      id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			Issuer:    "five",
			Subject:   userName,
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(consts.JwtSecret))
	return
}

func GenerateRefreshToken(userName string, id uint) (rToken string, err error) {
	expire := time.Now().Add(consts.RefreshTokenExpireTime).Unix()
	claims := &Claims{
		UserName: userName,
		UID:      id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			Issuer:    "five",
			Subject:   userName,
		},
	}
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(consts.JwtSecret))
	return
}

// ParseToken 解析token并判断其有没有过期
func ParseToken(token string) (*Claims, error, bool) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JwtSecret), nil
	})
	if err != nil {
		return nil, err, false
	}
	claims, ok := tokenClaims.Claims.(*Claims)
	if ok && tokenClaims.Valid {
		return claims, nil, IsValid(tokenClaims)
	}
	return nil, err, false
}

func IsValid(token *jwt.Token) bool {
	return token.Valid
}
