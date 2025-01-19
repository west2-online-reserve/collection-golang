namespace go user

include "model.thrift"

struct RegisterRequest{
    1: required string username,
    2: required string password,
}

struct RegisterResponse{
    1: model.BaseResp base,
    2: model.User data,
}

struct LoginRequest{
    1: required string username,
    2: required string password,
}

struct LoginResponse{
    1: model.BaseResp base,
    2: model.User data,
}

service UserService {
    RegisterResponse Register (1: RegisterRequest req) (api.post="/user/register"),
    LoginResponse Login(1: LoginRequest req) (api.post="/user/login"),
}

