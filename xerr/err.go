package xerr

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// 定义错误
var (
	ErrNotFound     = New(404, "NotFound", "没有找到相应记录")
	ErrAuthFailed   = New(401, "AuthFailed", "登录失败")
	ErrUnauthorized = New(401, "Unauthorized", "本接口只有登录用户才能调用")
	ErrForbidden    = New(403, "Forbidden", "权限不足")
)

// Error 对外输出的错误格式
type Error struct {
	code int
	// 错误代码，为英文字符串，前端可用此判断大的错误类型。
	Key string `json:"error"`
	// 错误消息，为详细错误描述，前端可选择性的展示此字段。
	Message string `json:"message"`
}

// New 新建一个 Error 对象
func New(code int, key string, msg string) *Error {
	return &Error{
		code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *Error) Error() string {
	return e.Key + ": " + e.Message
}

// ErrorHandler customize echo's HTTP error handler.
func ErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "ServerError"
		msg  string
	)
	// 二话不说先打日志
	c.Logger().Error(err.Error())

	if he, ok := err.(*Error); ok {
		// 我们自定的错误
		code = he.code
		key = he.Key
		msg = he.Message
	} else if ee, ok := err.(*echo.HTTPError); ok {
		// echo 框架的错误
		code = ee.Code
		key = http.StatusText(code)
		msg = fmt.Sprintf("%v", ee.Message)
	} else if ee, ok := err.(*cos.ErrorResponse); ok {
		// 腾讯云 cos 错误
		code = ee.Response.StatusCode
		key = ee.Code
		msg = ee.Message
	} else if err == gorm.ErrRecordNotFound {
		// 我们将 gorm 的没有找到直接返回 404
		code = http.StatusNotFound
		key = "NotFound"
		msg = "没有找到相应记录"
	} else if c.Echo().Debug {
		// 剩下的都是500 开了debug显示详细错误
		msg = err.Error()
	} else {
		// 500 不开debug 用标准错误描述 以防泄漏信息
		msg = http.StatusText(code)
	}

	// 判断 context 是否已经返回了
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				c.Logger().Error(err.Error())
			}
		} else {
			err := c.JSON(code, New(code, key, msg))
			if err != nil {
				c.Logger().Error(err.Error())
			}
		}
	}
}
