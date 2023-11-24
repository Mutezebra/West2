package consts

import "time"

const (
	PasswordCost = 12

	AccessTokenExpireTime  = 24 * 2 * time.Hour
	RefreshTokenExpireTime = 24 * 10 * time.Hour

	JwtSecret = "jwt-secret"
)
