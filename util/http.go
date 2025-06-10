package util

import (
	"app/code"
	"app/i18n"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	RequestId string      `json:"requestId" ` // 请求ID
	Code      code.Error  `json:"code" `      // 状态：""-正常成功
	Msg       string      `json:"msg" `       // 错误信息
	Data      interface{} `json:"data" `      // 返回数据
}

func JsonSuccess(c *fiber.Ctx, data interface{}) error {
	response := Response{}
	response.RequestId = GetRequestId(c)
	response.Code = ""
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
	response.RequestId = GetRequestId(c)
	response.Code = code.ParseError(err.Error())
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
	requestId += RandString(6)
	return requestId
}

func RandString(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
