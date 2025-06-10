package repo

import (
	"app/code"
	"app/model"

	"github.com/gofiber/fiber/v2"
)

type BaseRepo interface {
	Migrate() code.Error
}

type UserRepo interface {
	Insert(*fiber.Ctx, *model.User) code.Error
	Delete(*fiber.Ctx, *model.User) code.Error
	Update(*fiber.Ctx, *model.User) code.Error
	Select(*fiber.Ctx, *model.User) ([]model.User, code.Error)
	SelectById(*fiber.Ctx, *model.User) (*model.User, code.Error)
	SelectByUsername(*fiber.Ctx, *model.User) (*model.User, code.Error)
	SelectWithPagination(*fiber.Ctx, *model.Pagination) code.Error
	SelectTotalCount(*fiber.Ctx) (int, code.Error)
	Migrate() code.Error
}
