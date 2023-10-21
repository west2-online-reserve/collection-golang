namespace go task

include "model.thrift"

// 增
struct AddTaskRequest {
    1: required string title,   // 标题
    2: required string content, // 内容
    3: required i64 start_at,   // 开始时间
    4: required i64 end_at,     // 结束时间
}

struct AddTaskResponse {
    1: model.BaseResp base,
    2: model.Task data,
}


// 改
struct UpdateTaskRequest {
    1: required i64 id,         // 任务id， 0 表示全部
    6: required i64 status,     // 任务状态，0-未完成 1-已完成
}

struct UpdateTaskResponse {
    1: model.BaseResp base,
}

// 删除一条任务
struct DeleteTaskRequest {
    1: required i64 id,         // 任务id
}

struct DeleteTaskResponse {
    1: model.BaseResp base,
}

// 根据状态批量删除任务
struct DeleteTaskByStatusRequest {
    1: required i64 status,     // 任务状态，0-未完成 1-已完成 2-全部
}

struct DeleteTaskByStatusResponse {
    1: model.BaseResp base,
}

// 获取待办事项列表-根据状态
struct QueryTaskListByStatusRequest {
    1: required i64 page_size, // 每一页的数量
    2: required i64 page_num,  // 页码
    3: required i64 status,    // 任务状态 0-未完成 1-已完成
}

struct QueryTaskListByStatusResponse {
    1: model.BaseResp base,
    2: model.TaskList data,
}

// 获取待办事项列表-根据关键词
struct QueryTaskListByKeywordRequest {
    1: required i64 page_size, // 每一页的数量
    2: required i64 page_num,  // 页码
    3: required string keyword, // 关键词
}

struct QueryTaskListByKeywordResponse {
    1: model.BaseResp base,
    2: model.TaskList data,
}

service TaskService {
    AddTaskResponse AddTask(1:AddTaskRequest req) (api.post="/task/add"),

    UpdateTaskResponse UpdateTask(1:UpdateTaskRequest req) (api.put="/task/update"),

    DeleteTaskResponse DeleteTask(1:DeleteTaskRequest req) (api.delete="/task/delete"),
    DeleteTaskByStatusResponse DeleteTaskByStatus(1:DeleteTaskByStatusRequest req) (api.delete="/task/delete/status"),

    QueryTaskListByStatusResponse QueryTaskListByStatus(1:QueryTaskListByStatusRequest req) (api.get="/task/query/status"),
    QueryTaskListByKeywordResponse QueryTaskListByKeyword(1:QueryTaskListByKeywordRequest req) (api.get="/task/query/keyword"),
}