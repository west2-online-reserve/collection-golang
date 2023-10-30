package e

const (
	SUCCESS       = 200
	InvalidParams = 400
	ERROR         = 500

	ParseTokenFailed = 10001
	CheckTokenFailed = 10002

	JsonUnmarshalFailed    = 10003
	VerifyOtpFailed        = 10004
	GenerateTokenFailed    = 10005
	UpdateTotpStatusFailed = 10006
	GenerateOTPFailed      = 10007
	UpdateOTPSecretFailed  = 10008
	WriteFileFailed        = 10009
	SendEmailFailed        = 10010
)
