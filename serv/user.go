package serv

import (
	"app/code"
	"app/middleware"
	"app/model"
	"app/repo"
	"app/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type userServ struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserServ {
	return &userServ{
		userRepo: userRepo,
	}
}

func (o *userServ) Insert(c *fiber.Ctx, user *model.User) code.Error {
	password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.PasswordCryptFailed
	}
	user.Password = util.EnPointer(string(password))
	return o.userRepo.Insert(c, user)
}

func (o *userServ) Delete(c *fiber.Ctx, user *model.User) code.Error {
	return o.userRepo.Delete(c, user)
}

func (o *userServ) Update(c *fiber.Ctx, user *model.User) code.Error {
	if user.Password != nil && len(util.DePointer(user.Password)) > 0 {
		password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(util.GetRequestId(c), err)
			return code.PasswordCryptFailed
		}
		user.Password = util.EnPointer(string(password))
	}

	return o.userRepo.Update(c, user)
}

func (o *userServ) Select(c *fiber.Ctx, user *model.User) ([]model.User, code.Error) {
	return o.userRepo.Select(c, user)
}

func (o *userServ) SelectById(c *fiber.Ctx, user *model.User) (*model.User, code.Error) {
	return o.userRepo.SelectById(c, user)
}

func (o *userServ) SelectWithPagination(c *fiber.Ctx, p *model.Pagination) code.Error {
	return o.userRepo.SelectWithPagination(c, p)
}

func (o *userServ) Login(c *fiber.Ctx, user *model.User) (*model.User, string, code.Error) {
	userDB, errCode := o.userRepo.SelectByUsername(c, user)
	if errCode.IsNotNil() {
		return nil, "", errCode
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*userDB.Password), []byte(*user.Password)); err != nil {
		log.Error(util.GetRequestId(c), err)
		return nil, "", code.UsernameOrPasswordFailed
	}
	token, err := middleware.GenerateJwt(*user.Username)
	if err != nil {
		return nil, "", code.TokenGenerateFailed
	}
	userDB.Password = nil
	return userDB, token, code.Nil
}
