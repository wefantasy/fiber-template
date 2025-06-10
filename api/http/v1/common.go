package v1

import (
	"app/code"
	"app/log"
	"app/model/input"
	"app/serv"
	"app/util/httputil"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CommonContro struct {
	userServ serv.UserServ
}

func NewCommonController(userServ serv.UserServ) BaseContro {
	return &CommonContro{
		userServ: userServ,
	}
}

func (o *CommonContro) Name() string {
	return "Common"
}

func (o *CommonContro) RegisterRoute(api fiber.Router) {
	api.Get("/ping", o.ping)
	api.Post("/login", o.login)
	api.Post("/register", o.register)
}

// @Summary			测试
// @Description		测试
// @Tags			common
// @Accept			json
// @Produce		json
// @Router			/ping	[get]
func (o *CommonContro) ping(c *fiber.Ctx) error {
	log.Info(httputil.GetRequestId(c), "pong")
	return c.Status(http.StatusOK).SendString("pong")
}

// @Summary	登录
// @Description	登录
// @Tags	common
// @Accept	json
// @Produce	json
// @Param	user	body	input.UserLogin	true	"登录信息"
// @Router	/login	[post]
func (o *CommonContro) login(c *fiber.Ctx) error {
	user := &input.UserLogin{}
	if err := c.BodyParser(user); err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	token, err := o.userServ.Login(c, user)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, token)
}

// @Summary	注册
// @Description	注册
// @Tags	common
// @Accept	json
// @Produce	json
// @Param	user	body	input.UserRegister	true	"用户信息"
// @Router	/register	[post]
func (o *CommonContro) register(c *fiber.Ctx) error {
	user := &input.UserRegister{}
	if err := c.BodyParser(user); err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	if *user.Username == "" || *user.Password == "" {
		log.Error(httputil.GetRequestId(c), "用户名或密码不能为空")
		return code.ParamError
	}
	err := o.userServ.Register(c, user)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, nil)
}
