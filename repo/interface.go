package repo

import (
	"app/model"

	"github.com/gofiber/fiber/v2"
)

type BaseRepo interface {
}

type UserRepo interface {
	Insert(*fiber.Ctx, *model.User) error
	Delete(*fiber.Ctx, int) error
	Update(*fiber.Ctx, *model.User) error
	Select(*fiber.Ctx, *model.User) ([]model.User, error)
	SelectById(*fiber.Ctx, int) (*model.User, error)
	SelectByUsername(*fiber.Ctx, string) (*model.User, error)
	SelectWithPagination(*fiber.Ctx, *model.Pagination) error
	SelectTotalCount(*fiber.Ctx) (int, error)
}
