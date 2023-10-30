package myutils

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOtp(userName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Mutezebra",
		AccountName: userName,
	})
	return key, err
}

func VerifyOtp(token, secret string) bool {
	valid := totp.Validate(token, secret)
	return valid
}
