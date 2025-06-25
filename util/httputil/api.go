package httputil

import (
	"app/i18n"
	"app/util"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	RequestId string `json:"requestId" ` // 请求ID
	Msg       string `json:"msg" `       // success或错误信息
	Data      any    `json:"data" `      // 返回数据
}

func JsonSuccess(c *fiber.Ctx, data any) error {
	response := Response{}
	strings.TrimSpace(GetRequestId(c))
	response.RequestId = strings.TrimSpace(GetRequestId(c))
	response.Msg = "success"
	response.Data = data
	return c.Status(http.StatusOK).JSON(response)
}

func JsonErrorParse(c *fiber.Ctx, err error) error {
	errMsg := i18n.LocalizeWithCtx(c, err.Error())
	if errMsg == "" {
		return err
	}
	response := Response{}
	response.RequestId = strings.TrimSpace(GetRequestId(c))
	response.Msg = errMsg
	response.Data = nil
	return c.Status(http.StatusOK).JSON(response)
}

func GetRequestId(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	rid := c.Locals(fiber.HeaderXRequestID)
	if rid == "" {
		rid = GenerateRequestId()
		c.Locals(fiber.HeaderXRequestID, rid)
	}
	c.Locals(fiber.HeaderXRequestID, rid)
	return c.Locals(fiber.HeaderXRequestID).(string) + " "
}

func GenerateRequestId() string {
	requestId := time.Now().Format("20060102150405")
	requestId += util.RandString(6)
	return requestId
}
