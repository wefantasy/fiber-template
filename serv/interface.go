package serv

import (
	"app/code"
	"app/model"

	"github.com/gofiber/fiber/v2"
)

type BaseServ interface {
}

type UserServ interface {
	Insert(*fiber.Ctx, *model.User) code.Error
	Delete(*fiber.Ctx, *model.User) code.Error
	Update(*fiber.Ctx, *model.User) code.Error
	Select(*fiber.Ctx, *model.User) ([]model.User, code.Error)
	SelectById(*fiber.Ctx, *model.User) (*model.User, code.Error)
	SelectWithPagination(*fiber.Ctx, *model.Pagination) code.Error
	Login(*fiber.Ctx, *model.User) (*model.User, string, code.Error)
}
