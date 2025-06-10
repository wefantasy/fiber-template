package v1

import (
	"app/code"
	"app/middleware"
	"app/model"
	"app/serv"
	"app/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
	return code.ParamError
	return c.Status(http.StatusOK).SendString("pong")
}

// @Summary	登录
// @Description	登录
// @Tags	common
// @Accept	json
// @Produce	json
// @Param	user	body	model.User	true	"用户名"
// @Router	/login	[post]
func (o *CommonContro) login(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.BodyParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	if *user.Username == "" || *user.Password == "" {
		log.Error(util.GetRequestId(c), "用户名或密码不能为空")
		return code.ParamError
	}
	userDB, token, errCode := o.userServ.Login(c, user)
	if errCode.IsNotNil() {
		return errCode
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(middleware.JwtExpireTime),
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
	})
	data := map[string]interface{}{
		"user":  userDB,
		"token": token,
	}
	return util.JsonSuccess(c, data)
}

// @Summary	注册
// @Description	注册
// @Tags	common
// @Accept	json
// @Produce	json
// @Param	user	body	model.User	true	"用户名"
// @Router	/register	[post]
func (o *CommonContro) register(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.BodyParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	if *user.Username == "" || *user.Password == "" {
		log.Error(util.GetRequestId(c), "用户名或密码不能为空")
		return code.ParamError
	}
	errCode := o.userServ.Insert(c, user)
	if errCode.IsNotNil() {
		return errCode
	}
	return util.JsonSuccess(c, nil)
}
