package consts

import "time"

const (
	PasswordCost = 12

	JwtSecret = "jwt-secret"

	AccessTokenExpireTime  = 24 * 2 * time.Hour
	RefreshTokenExpireTime = 24 * 10 * time.Hour
	AccessToken            = "access_token"
	RefreshToken           = "refresh_token"
	OtpCodeStorePath       = "./static/imgs/qr-code/"
	EmailSubject           = "这是您的验证码,请使用authenticator扫描"

	MaxAvatarSize = 3 * MB
)
