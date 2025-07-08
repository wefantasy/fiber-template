package httputil

import (
	"app/i18n"
	"app/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Response struct {
	TraceId string `json:"traceId"` // 请求ID
	Msg     string `json:"msg" `    // success或错误信息
	Data    any    `json:"data" `   // 返回数据
}

func JsonSuccess(c *fiber.Ctx, data any) error {
	response := Response{}
	response.TraceId = GetTraceId(c)
	response.Msg = "success"
	response.Data = data
	return c.Status(http.StatusOK).JSON(response)
}

func JsonErrorParse(c *fiber.Ctx, err error) error {
	errMsg := i18n.LocalizeWithCtx(c, err.Error())
	if errMsg == "" {
		errMsg = err.Error()
	}
	response := Response{}
	response.TraceId = GetTraceId(c)
	response.Msg = errMsg
	response.Data = nil
	return c.Status(http.StatusOK).JSON(response)
}

func GetTraceId(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	return c.Locals(util.TraceIdKey).(string)
}
