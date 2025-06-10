package repo

import (
	"app/code"
	"app/model"
	"app/util"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (o *userRepo) Migrate() code.Error {
	return code.Nil
}

func (o *userRepo) Insert(c *fiber.Ctx, user *model.User) code.Error {
	sql := fmt.Sprintf("INSERT INTO user(%s) VALUES (%s)", util.ExtractDBNotZeroColumnStr(user), util.ExtractDBNotZeroColumnStrWithPrefix(user, ":"))
	user.CreatedAt = util.EnPointer(time.Now())
	result, err := util.DB.NamedExec(sql, user)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.DatabaseError
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.DatabaseError
	}
	user.Id = util.EnPointer(int(id))
	util.RDB.SetStruct(user.CacheKey(), user)
	return code.Nil
}

func (o *userRepo) Delete(c *fiber.Ctx, user *model.User) code.Error {
	sql := "DELETE FROM user WHERE id = :id"
	_, err := util.DB.NamedExec(sql, user)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.DatabaseError
	}
	util.RDB.Delete(user.CacheKey())
	return code.Nil
}

func (o *userRepo) Update(c *fiber.Ctx, user *model.User) code.Error {
	sql := fmt.Sprintf("UPDATE user SET %s WHERE id = :id", util.ExtractDBNotZeroColumnSet(user))
	user.UpdatedAt = util.EnPointer(time.Now())
	_, err := util.DB.NamedExec(sql, user)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.DatabaseError
	}
	util.RDB.Delete(user.CacheKey())
	return code.Nil
}

func (o *userRepo) Select(c *fiber.Ctx, user *model.User) ([]model.User, code.Error) {
	sql := "SELECT * FROM user"
	sets := util.ExtractDBNotZeroColumnSetWithPrefix(user, "user.")
	if sets != "" {
		sql = sql + " WHERE " + sets + " ORDER BY user.created_at desc"
	} else {
		sql = sql + " ORDER BY user.created_at desc"
	}
	var users []model.User
	stmt, err := util.DB.PrepareNamed(sql)
	if err != nil {
		log.Error(err)
		return nil, code.DatabaseError
	}
	err = stmt.Select(&users, user)
	if err != nil {
		log.Error(err)
		return nil, code.DatabaseError
	}
	return users, code.Nil
}

func (o *userRepo) SelectById(c *fiber.Ctx, user *model.User) (*model.User, code.Error) {
	if util.RDB.GetStruct(user.CacheKey(), user) == code.Nil {
		return user, code.Nil
	}

	sql := "SELECT * FROM user"
	sql = sql + " WHERE user.id=?"

	result := model.User{}
	err := util.DB.Get(&result, sql, user.Id)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return nil, code.DatabaseError
	}
	return &result, code.Nil
}

func (o *userRepo) SelectByUsername(c *fiber.Ctx, user *model.User) (*model.User, code.Error) {
	sql := "SELECT * FROM user"
	sql = sql + " WHERE user.username=?"

	result := model.User{}
	err := util.DB.Get(&result, sql, user.Username)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return nil, code.DatabaseError
	}
	return &result, code.Nil
}

func (o *userRepo) SelectWithPagination(c *fiber.Ctx, p *model.Pagination) code.Error {
	if p.Total == 0 {
		total, errCode := o.SelectTotalCount(c)
		if errCode.IsNotNil() {
			return errCode
		} else if total == 0 {
			p.Data = nil
			return code.Nil
		}
		p.Total = total
	}
	p.Offset = p.Page * p.Size
	p.Pages = (p.Total + p.Size - 1) / p.Size
	sql := "SELECT * FROM user"
	sql = sql + " LIMIT :offset,:size"
	stmt, err := util.DB.PrepareNamed(sql)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.DatabaseError
	}
	var users []model.User
	err = stmt.Select(&users, p)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return code.DatabaseError
	}
	p.Data = users
	return code.Nil
}

func (o *userRepo) SelectTotalCount(c *fiber.Ctx) (int, code.Error) {
	sql := "SELECT COUNT(id) AS total FROM user"
	var total int
	err := util.DB.Get(&total, sql)
	if err != nil {
		log.Error(util.GetRequestId(c), err)
		return 0, code.DatabaseError
	}
	return total, code.Nil
}
