package e

var MsgFlags = map[int]string{
	// common
	SUCCESS:           "ok",
	ERROR:             "operation failed",
	GetUserInfoFailed: "Get user info failed",
	InvalidParams:     "Params is invalid.",

	// user
	UserExists:          "User exists.",
	UserNotExists:       "User not exists.",
	CreateUserFailed:    "Create user failed.",
	SetPasswordFailed:   "Set password failed.",
	CheckPasswordFailed: "Check password failed.",

	// token
	GenerateTokenFailed: "Generate token failed.",
	CheckTokenFailed:    "Check token failed.",
	ParseTokenFailed:    "Parse token failed.",

	// task
	CreateTaskFailed: "Create task failed.",
	FindTaskFailed:   "Find task failed",
	UpdateTaskFailed: "Update task failed.",
	DeleteTaskFailed: "Delete task failed",
	SearchTaskFailed: "Search task failed",
	ListTaskFailed:   "List task failed",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
