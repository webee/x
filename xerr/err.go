package xerr

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/webee/x/xlog"
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

// Newf 新建一个带格式的 Error
func Newf(code int, key string, format string, a ...interface{}) *Error {
	return &Error{
		code:    code,
		Key:     key,
		Message: fmt.Sprintf(format, a...),
	}
}

// Error makes it compatible with `error` interface.
func (e *Error) Error() string {
	return e.Key + ": " + e.Message
}

type ErrorFunc func(err error) (ok bool, code int, key string, msg string)

var log = xlog.Get()

// 定义错误
var (
	ErrNotFound     = New(404, "NotFound", "not found")
	ErrUnauthorized = New(401, "Unauthorized", "unauthorized")
	ErrForbidden    = New(403, "Forbidden", "forbidden")
)

var errorFuncs = []ErrorFunc{}

func RegisterErrorFunc(f ErrorFunc) {
	errorFuncs = append(errorFuncs, f)
}

// ErrorHandler customize echo's HTTP error handler.
func ErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "ServerError"
		msg  string
	)

	// 二话不说先打日志
	log.WithError(err).Errorf("handling error: %T", err)

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
	} else {
		for _, f := range errorFuncs {
			if ok, _code, _key, _msg := f(err); ok {
				code = _code
				key = _key
				msg = _msg
				break
			}
		}

		if !ok {
			if c.Echo().Debug {
				// 剩下的都是500 开了debug显示详细错误
				msg = err.Error()
			} else {
				// 500 不开debug 用标准错误描述 以防泄漏信息
				msg = http.StatusText(code)
			}
		}
	}

	// 判断 context 是否已经返回了
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, New(code, key, msg))
		}
		if err != nil {
			c.Logger().Error(err.Error())
		}
	}
}
