package repo

import (
	"app/db"
	"app/model"
	"app/util"
	"app/util/dbutil"
	"fmt"
	"time"

	"app/log"

	"github.com/gofiber/fiber/v2"
)

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (o *userRepo) Insert(c *fiber.Ctx, user *model.User) error {
	sql := fmt.Sprintf("INSERT INTO user(%s) VALUES (%s)",
		dbutil.NewBuilder(user).OnlyNonZero().BuildColumns(", "),
		dbutil.NewBuilder(user).OnlyNonZero().WithPrefix(":").BuildNamedPlaceholders(", "))
	result, err := db.DB.NamedExec(sql, user)
	if err != nil {
		log.F(c).Info(err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.F(c).Error(err)
		return err
	}
	user.Id = util.EnPointer(int(id))
	db.RDB.SetStruct(user.CacheKey(), user)
	return nil
}

func (o *userRepo) Delete(c *fiber.Ctx, id int) error {
	sql := "DELETE FROM user WHERE id = ?"
	_, err := db.DB.Exec(sql, id)
	if err != nil {
		log.F(c).Error(err)
		return err
	}
	db.RDB.Delete(model.UserCacheKey(id))
	return nil
}

func (o *userRepo) Update(c *fiber.Ctx, user *model.User) error {
	sql := fmt.Sprintf("UPDATE user SET %s WHERE id = :id",
		dbutil.NewBuilder(user).ExcludePK().OnlyNonZero().BuildSetClauses(","))
	user.UpdatedAt = util.EnPointer(time.Now())
	_, err := db.DB.NamedExec(sql, user)
	if err != nil {
		log.F(c).Error(err)
		return err
	}
	db.RDB.Delete(user.CacheKey())
	return nil
}

func (o *userRepo) Select(c *fiber.Ctx, userFilter *model.User) ([]model.User, error) {
	sql := dbutil.NewBuilder(userFilter).
		OnlyNonZero().
		WithOrderBy("created_at desc").
		BuildSelectQuery("user")
	var users []model.User
	stmt, err := db.DB.PrepareNamed(sql)
	if err != nil {
		log.F(c).Error(err)
		return nil, err
	}
	err = stmt.Select(&users, userFilter)
	if err != nil {
		log.F(c).Error(err)
		return nil, err
	}
	return users, nil
}

func (o *userRepo) SelectById(c *fiber.Ctx, id int) (*model.User, error) {
	user := model.User{}
	if db.RDB.GetStruct(model.UserCacheKey(id), &user) != nil {
		return &user, nil
	}

	sql := dbutil.NewBuilder(user).
		OnlyNonZero().
		WithCustomWhere("id = ?").
		BuildSelectQuery("user")

	err := db.DB.Get(&user, sql, id)
	if err != nil {
		log.F(c).Error(err)
		return nil, err
	}
	return &user, nil
}

func (o *userRepo) SelectByUsername(c *fiber.Ctx, username string) (*model.User, error) {
	sql := dbutil.NewBuilder(&model.User{}).
		OnlyNonZero().
		WithCustomWhere("username = ?").
		BuildSelectQuery("user")

	result := model.User{}
	err := db.DB.Get(&result, sql, username)
	if err != nil {
		log.F(c).Error(err)
		return nil, err
	}
	return &result, nil
}

func (o *userRepo) SelectWithPagination(c *fiber.Ctx, p *model.Pagination) error {
	if p.Total == 0 {
		total, err := o.SelectTotalCount(c)
		if err != nil {
			return err
		} else if total == 0 {
			p.Data = nil
			return nil
		}
		p.Total = total
	}
	p.Format()
	sql := dbutil.NewBuilder(&model.User{}).
		OnlyNonZero().
		WithLimitOffset(p.Size, p.Offset).
		BuildSelectQuery("user")
	var users []model.User
	err := db.DB.Select(&users, sql)
	if err != nil {
		log.F(c).Error(err)
		return err
	}
	p.Data = users
	return nil
}

func (o *userRepo) SelectTotalCount(c *fiber.Ctx) (int, error) {
	sql := "SELECT COUNT(id) AS total FROM user"
	var total int
	err := db.DB.Get(&total, sql)
	if err != nil {
		log.F(c).Error(err)
		return 0, err
	}
	return total, nil
}
