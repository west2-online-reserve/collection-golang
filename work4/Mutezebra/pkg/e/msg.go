package e

var MsgFlags = map[int]string{
	// Common
	SUCCESS: "operator success",
	ERROR:   "operator failed",

	ParseTokenFailed: "parse token failed",
	CheckTokenFailed: "check token failed",

	JsonUnmarshalFailed: "Json unmarshal failed",

	// User
	SetPasswordFailed:      "Set password failed",
	CreateUserFailed:       "Create user failed",
	UserExists:             "User have exists",
	UserDoNotExist:         "User don`t exist",
	CheckPasswordFailed:    "Check password failed",
	VerifyOtpFailed:        "Please input you otp",
	GenerateTokenFailed:    "Generate token failed",
	GetUserInfoFailed:      "Get user info failed",
	FindUserFailed:         "Find user failed",
	UpdateTotpStatusFailed: "Update totp status failed",
	GenerateOTPFailed:      "Generate otp failed",
	UpdateOTPSecretFailed:  "Update otp secret failed",
	WriteFileFailed:        "Write file failed",
	SendEmailFailed:        "Send email failed",
	UserInfoUpdateFailed:   "User info update failed",
	InvalidEmailFormat:     "Invalid email format",
	UnValidAvatar:          "This avatar is invalid",
	SaveAvatarFailed:       "Save avatar failed",
	UpdateAvatarFailed:     "Update avatar failed",
	FollowFailed:           "Follow failed",
	UnFollowFailed:         "Unfollow failed",
	HaveFollowed:           "Your have followed",
	GetFriendListFailed:    "Get friend list failed",
	GetFollowerListFailed:  "Get follower list failed",
	DeleteUserFailed:       "Delete user Failed",
	AvatarFileOpenFailed:   "Avatar file open failed",
	ReadAvatarFileFailed:   "Read avatar file failed",

	//video
	OpenVideoHeaderFailed:  "Open video header failed",
	ReadVideoFileFailed:    "Read video file failed",
	OpenFileFailed:         " Open filed failed",
	VideoWriteToFileFailed: "Video write to file failed",
	VideoUpdateFailed:      "Video update failed",
	InvalidVideo:           "invalid video",
	CloseFileFailed:        "close file failed",
	FindVideoFailed:        "Find video failed",
	CreateVideoFailed:      "Create video failed",
	DeleteVideoFailed:      "Delete video failed",
	OSSUploadVideoFailed:   "oss upload file failed",
	CachedVideoFailed:      "cached video failed",

	//comment
	VideoCommentCreateFailed:     "Video comment create failed",
	UpdateVideoCommentTreeFailed: "Update video comment tree failed",
	ReplyRecordNotExist:          "Reply record not exist",
	FindCommentRootFailed:        "Find comment root failed",

	// search
	GetSearchItemFailed: "Get item from redis failed",
	SearchUserFailed:    "Search user failed",
	SearchFailed:        "Search failed",
	CreateIndexFailed:   "Create index failed",
	CreateDocFailed:     "Create doc failed",
	VideoFilterFailed:   "Video filter failed",
	SearchDocFailed:     "Search doc failed",
	UpdateDocFailed:     "Update doc failed",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
