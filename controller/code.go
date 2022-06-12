package controller


//定义一些可能出现的错误码

type ResCode int64
const(
	CodeFocusSuccess = 0
	CodeSuccess ResCode = 1000 +iota
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
	CodeSaveSuccess
	CodeInvalidErr
	CodeUserNotExist
	CodeNotAccordStandard
)

var codeMsgMap = map[ResCode]string{
	CodeFocusSuccess:			"操作成功",
	CodeSuccess:				"success",
	CodeServerBusy:				"服务繁忙",
	CodeNeedLogin:				"需要登录",
	CodeInvalidToken:			"无效Token",
	CodeSaveSuccess:         	"保存成功",
	CodeInvalidErr:				"参数错误",
	CodeUserNotExist:			"用户不存在",
	CodeNotAccordStandard: 		"必填字段未填或不符合数据规范",
}

// Msg 返回特定的错误提示信息
func (c ResCode)Msg()string{
	msg,ok := codeMsgMap[c]
	if !ok{
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}

