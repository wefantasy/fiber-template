package serv

import (
	"app/model"
	"app/model/input"
	"app/model/output"
	"github.com/gofiber/fiber/v2"
)

type BaseServ interface {
}

type UserServ interface {
	Insert(*fiber.Ctx, *model.User) error
	Delete(*fiber.Ctx, int) error
	Update(*fiber.Ctx, *model.User) error
	Select(*fiber.Ctx, *input.UserFilter) ([]output.UserOutput, error)
	SelectById(*fiber.Ctx, int) (*output.UserOutput, error)
	SelectWithPagination(*fiber.Ctx, *model.Pagination) error
	Login(*fiber.Ctx, *input.UserLogin) (string, error)
	Register(*fiber.Ctx, *input.UserRegister) error
}
