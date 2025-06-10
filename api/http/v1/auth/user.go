package auth

import (
	v1 "app/api/http/v1"
	"app/code"
	"app/middleware"
	"app/model"
	"app/serv"
	"app/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserContro struct {
	userServ serv.UserServ
}

func NewUserController(userServ serv.UserServ) v1.BaseContro {
	return &UserContro{
		userServ: userServ,
	}
}
func (o *UserContro) RegisterRoute(api fiber.Router) {
	api.Post("/user", middleware.Jwt(), o.Insert)
	api.Delete("/user/:id", middleware.Jwt(), o.Delete)
	api.Put("/user", middleware.Jwt(), o.Update)
	api.Get("/user", middleware.Jwt(), o.Select)
	api.Get("/user/:id", middleware.Jwt(), o.SelectById)
	api.Get("/user/pagination/:size/:page", middleware.Jwt(), o.SelectWithPagination)
}

func (o *UserContro) Name() string {
	return "User"
}

// @Summary		新增用户
// @Description	新增用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			user	body		model.User	true	"用户信息"
// @Router			/user	[post]
func (o *UserContro) Insert(c *fiber.Ctx) error {
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	now := time.Now()
	user.DeletedAt = &now
	errCode := o.userServ.Insert(c, user)
	if errCode.IsNotNil() {
		return errCode
	}
	return util.JsonSuccess(c, "")
}

// @Summary		删除用户
// @Description	删除用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			id		path		int	true	"用户的 id"
// @Router			/user/:id	[delete]
func (o *UserContro) Delete(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.ParamsParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	now := time.Now()
	user.DeletedAt = &now
	if errCode := o.userServ.Delete(c, user); errCode.IsNotNil() {
		log.Error(util.GetRequestId(c), errCode)
		return code.ParamError
	}
	return util.JsonSuccess(c, nil)
}

// @Summary		更新用户
// @Description	更新用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			user	body		model.User	true	"用户信息"
// @Router			/user	[put]
func (o *UserContro) Update(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.BodyParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}

	errCode := o.userServ.Update(c, user)
	if errCode.IsNotNil() {
		return errCode
	}
	return util.JsonSuccess(c, nil)
}

// @Summary		查找用户
// @Description	查找用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Router			/user	[get]
func (o *UserContro) Select(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.QueryParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	users, errCode := o.userServ.Select(c, user)
	if errCode.IsNotNil() {
		return errCode
	}
	return util.JsonSuccess(c, users)
}

// @Summary		按id查找用户
// @Description	按id查找用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			id			path		int	true	"用户的 id"
// @Router			/user/{id}	[get]
func (o *UserContro) SelectById(c *fiber.Ctx) error {
	user := &model.User{}
	if err := c.ParamsParser(user); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	user, errCode := o.userServ.SelectById(c, user)
	if errCode.IsNotNil() {
		return errCode
	}
	return util.JsonSuccess(c, user)
}

// @Summary		分页查找用户
// @Description	分页查找用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			size							path		int	true	"分页大小"
// @Param			page							path		int	true	"查询页号"
// @Router			/user/pagination/{size}/{page}	[get]
func (o *UserContro) SelectWithPagination(c *fiber.Ctx) error {
	p := &model.Pagination{}
	if err := c.ParamsParser(p); err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.ParamError
	}
	p.Page = p.Page - 1
	errCode := o.userServ.SelectWithPagination(c, p)
	if errCode.IsNotNil() {
		return errCode
	}
	return util.JsonSuccess(c, p)
}
