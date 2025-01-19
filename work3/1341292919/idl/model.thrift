namespace go model

struct User{
    1: required i64 id            //用户id
    2: required string username,  //用户名
    3: required string password,  //密码
    4: required string created_at   //创建时间
    5: required string updated_at   //更新时间

}

struct Task{
    1: required i64 id            //用户id
    2: required string title,    //标题
    3: required string content,  //内容
    4: required i64 status,      //0-未完成 1-已完成
    5: required string created_at   //创建时间
    6: required string updated_at   //更新时间
    7: required string start_at     //开始时间
    8: required string end_at       //结束时间
}

struct TaskList{
    1: required list<Task> items,   //任务列表
    2: required i64 total,          //总数
}

struct BaseResp {
    1: required i64 code,          //请求返回的状态码
    2: required string msg,        //返回的消息
}