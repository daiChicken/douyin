package controller

//定义一些可能出现的错误码

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
	CodeSaveSuccess
	CodeInvalidParam
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:      "success",
	CodeServerBusy:   "服务繁忙",
	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效Token",
	CodeSaveSuccess:  "保存成功",
	CodeInvalidParam: "请求参数有误",
}

// Msg 返回特定的错误提示信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
