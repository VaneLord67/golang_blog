package common

type ResultCode struct {
	Code    int
	Message string
}

func (r ResultCode) Error() string {
	return r.Message
}

// ResultCode 枚举
var SUCCESS = ResultCode{1, "操作成功"}
var PERMISSON_DENIED = ResultCode{2, "权限不足"}
var FAIL = ResultCode{-1, "操作失败"}
var PASSWORD_WRONG = ResultCode{Code: 100, Message: "密码错误"}
var USER_NOT_EXISTS = ResultCode{Code: 101, Message: "用户不存在"}
var TOKEN_CREATE_ERROR = ResultCode{Code: 102, Message: "创建jwt出错"}
var USER_ALREADY_EXISTS = ResultCode{Code: 103, Message: "用户已存在"}
var TOKEN_PARSE_ERROR = ResultCode{Code: 104, Message: "jwt解析出错"}
var TOKEN_EXPIRED = ResultCode{Code: 105, Message: "jwt过期"}
var PARAMETER_PARSE_ERROR = ResultCode{Code: 106, Message: "参数解析异常"}
var SERVICE_FIND_ERROR = ResultCode{Code: 107, Message: "服务发现失败"}
var CAPTCHA_ERROR = ResultCode{Code: 108, Message: "验证码错误"}
var TOO_MANY_REQUESTS = ResultCode{Code: 429, Message: "请求过多"}
var GITHUB_TOKEN_GET_WRONG = ResultCode{Code: 109, Message: "获取Github token失败"}
var GITHUB_ACCOUNT_NOT_BIND = ResultCode{Code: 110, Message: "Github账号未绑定"}
var GITHUB_ACCOUNT_ALREADY_BIND = ResultCode{Code: 111, Message: "Github账号已绑定"}
var GITHUB_INFO_GET_ERROR = ResultCode{Code: 112, Message: "获取Github用户信息失败"}
