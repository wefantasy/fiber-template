package auth

import (
	v1 "app/api/http/v1"
	"app/code"
	"app/log"
	"app/middleware"
	"app/model"
	"app/model/input"
	"app/serv"
	"app/util/httputil"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	api.Post("/user", middleware.JwtAuth(), o.Insert)
	api.Delete("/user/:id", middleware.JwtAuth(), o.Delete)
	api.Put("/user", middleware.JwtAuth(), o.Update)
	api.Get("/user", middleware.JwtAuth(), o.Select)
	api.Get("/user/:id", middleware.JwtAuth(), o.SelectById)
	api.Get("/user/pagination/:size/:page", middleware.JwtAuth(), o.SelectWithPagination)
}

func (o *UserContro) Name() string {
	return "User"
}

// Insert @Summary		新增用户
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
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	err := o.userServ.Insert(c, user)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, "")
}

// Delete @Summary		删除用户
// @Description	删除用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			id		path		int	true	"用户的 id"
// @Router			/user/{id}	[delete]
func (o *UserContro) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return code.ParamError
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	if err := o.userServ.Delete(c, id); err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	return httputil.JsonSuccess(c, nil)
}

// Update @Summary		更新用户
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
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}

	err := o.userServ.Update(c, user)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, nil)
}

// Select @Summary		查找用户
// @Description	查找用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param	username	query	string	false	"用户账户"
// @Router			/user	[get]
func (o *UserContro) Select(c *fiber.Ctx) error {
	userFilter := &input.UserFilter{}
	if err := c.QueryParser(userFilter); err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	users, err := o.userServ.Select(c, userFilter)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, users)
}

// SelectById @Summary		按id查找用户
// @Description	按id查找用户
// @Tags			user
// @Accept			json
// @Produce		json
// @Param	Authorization	header	string	true	"Authentication header" default(Bearer xxxx)
// @Param			id			path		int	true	"用户的 id"
// @Router			/user/{id}	[get]
func (o *UserContro) SelectById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return code.ParamError
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	userOutput, err := o.userServ.SelectById(c, id)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, userOutput)
}

// SelectWithPagination @Summary		分页查找用户
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
		log.Error(httputil.GetRequestId(c), err)
		return code.ParamError
	}
	err := o.userServ.SelectWithPagination(c, p)
	if err != nil {
		return err
	}
	return httputil.JsonSuccess(c, p)
}
