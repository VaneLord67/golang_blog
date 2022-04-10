package common

type ResultCode struct {
	Code    int
	Message string
}

// ResultCode 枚举
var SUCCESS = ResultCode{1, "操作成功"}
var FAIL = ResultCode{-1, "操作失败"}
var PASSWORD_WRONG = ResultCode{Code: 100, Message: "密码错误"}
var USER_NOT_EXISTS = ResultCode{Code: 101, Message: "用户不存在"}
var TOKEN_CREATE_ERROR = ResultCode{Code: 102, Message: "创建jwt出错"}
