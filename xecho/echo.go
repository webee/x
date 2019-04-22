package xecho

import (
	"github.com/labstack/echo/v4"
)

// QueryParamExists 检查查询参数是否存在
func QueryParamExists(c echo.Context, name string) bool {
	return len(c.QueryParams()[name]) > 0
}
