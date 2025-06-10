package serv

import (
	"app/middleware"
	"app/model"
	"app/model/input"
	"app/model/output"
	"app/repo"
	"app/util"
	"app/util/copier"
	"app/util/httputil"

	"app/log"

	"github.com/gofiber/fiber/v2"
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

func (o *userServ) Insert(c *fiber.Ctx, user *model.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return err
	}
	user.Password = util.EnPointer(string(password))
	return o.userRepo.Insert(c, user)
}

func (o *userServ) Delete(c *fiber.Ctx, id int) error {
	return o.userRepo.Delete(c, id)
}

func (o *userServ) Update(c *fiber.Ctx, user *model.User) error {
	if user.Password != nil && len(util.DePointer(user.Password)) > 0 {
		password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(httputil.GetRequestId(c), err)
			return err
		}
		user.Password = util.EnPointer(string(password))
	}

	return o.userRepo.Update(c, user)
}

func (o *userServ) Select(c *fiber.Ctx, userFilter *input.UserFilter) ([]output.UserOutput, error) {
	user := &model.User{}
	err := copier.CopyProperties(userFilter, user)
	if err != nil {
		return nil, err
	}
	users, err := o.userRepo.Select(c, user)
	if err != nil {
		return nil, err
	}
	var userOutputs []output.UserOutput
	err = copier.TransferListType(users, &userOutputs)
	if err != nil {
		return nil, err
	}
	return userOutputs, err
}

func (o *userServ) SelectById(c *fiber.Ctx, id int) (*output.UserOutput, error) {
	user, err := o.userRepo.SelectById(c, id)
	if err != nil {
		return nil, err
	}
	var userOutputs output.UserOutput
	err = copier.CopyProperties(user, &userOutputs)
	if err != nil {
		return nil, err
	}
	return &userOutputs, err
}

func (o *userServ) SelectWithPagination(c *fiber.Ctx, p *model.Pagination) error {
	err := o.userRepo.SelectWithPagination(c, p)
	if err != nil || p.Data == nil {
		return err
	}
	users := p.Data.([]model.User)
	var userOutputs []output.UserOutput
	err = copier.TransferListType(users, &userOutputs)
	if err != nil {
		return err
	}
	p.Data = userOutputs
	return nil
}

func (o *userServ) Login(c *fiber.Ctx, userLogin *input.UserLogin) (string, error) {
	userDB, err := o.userRepo.SelectByUsername(c, *userLogin.Username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*userDB.Password), []byte(*userLogin.Password)); err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return "", err
	}
	token, err := middleware.GenerateJwt(*userLogin.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (o *userServ) Register(c *fiber.Ctx, userRegister *input.UserRegister) error {
	var user model.User
	err := copier.CopyProperties(userRegister, &user)
	if err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(*userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(httputil.GetRequestId(c), err)
		return err
	}
	user.Password = util.EnPointer(string(password))

	return o.userRepo.Insert(c, &user)
}
