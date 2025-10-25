package main

/* 子评论 */
type comment struct {
	Rpid    int64  // 评论唯一id
	Ctime   int64  // 时间戳
	Like    int    // 点赞数
	Message string // 内容
	User    user   // 用户信息
}

/* 主评论 */
type mainComment struct {
	comment

	SubComments []comment // 子评论
}

/* 用户 */
type user struct {
	Uid  int64
	Name string
	Sex  string
}
