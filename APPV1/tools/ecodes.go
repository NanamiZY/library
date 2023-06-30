package tools

import "time"

const (
	OK                  = 0
	NotLogin            = 10001 //你还没登录
	UserInfoErr         = 10002 //用户信息错误
	DoErr               = 10003
	NotFound            = 10004 //信息不存在
	InternalServerError = 10005 //连接错误
	CaptchaError        = 10006 //验证码错误
	T                   = 30 * 24 * time.Hour
)

type HttpCode struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
