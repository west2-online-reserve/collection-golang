package consts

import "time"

const (
	PasswordCost = 12

	JwtSecret = "something-secret"

	AccessTokenExpireTime  = 24 * time.Hour
	RefreshTokenExpireTime = 24 * 10 * time.Hour
	AccessToken            = "access_token"
	RefreshToken           = "refresh_token"
)
