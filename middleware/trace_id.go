package middleware

import (
	"app/util"
	"github.com/gofiber/fiber/v2"
)

func TraceId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceId := c.Get(util.TraceHeaderIdKey)
		if traceId == "" {
			traceId = util.RandTraceId()
		}
		// 注入到 c.Locals 中，方便在 Fiber handler 内部快速访问
		c.Locals(util.TraceIdKey, traceId)
		// 注入到标准的 context.Context 中，用于跨API边界传递
		ctx := util.NewRootContextWithTraceId(traceId)
		c.SetUserContext(ctx)
		// 在响应头中设置 TraceID，方便客户端追踪
		c.Set(util.TraceHeaderIdKey, traceId)
		return c.Next()
	}
}
