package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
定义返回响应的格式：
{
	"code":	10001,
	"msg": "错误提示信息",
	"data":	"返回的数据"
}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 有许多不同的响应， 比如错误响应， 成功响应
func ResponseError(c *gin.Context, code ResCode) {
	rd := ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}

	c.JSON(http.StatusOK, rd)
}

// 就是说可能你要返回某个错误的时候， 你要返回的错误信息相比默认的那个更加具体。 比如错误进来了， 还有一个验证器错误， 将错误进行翻译
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	c.JSON(http.StatusOK, rd)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}

	c.JSON(http.StatusOK, rd)
}
