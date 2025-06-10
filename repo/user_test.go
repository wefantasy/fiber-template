package repo

import (
	"app/conf"
	"app/db"
	"app/logger"
	"app/model"
	"app/util"
	"encoding/json"
	"testing"
)

func init() {
	conf.Initialize()
	logger.Initialize()
	db.Initialize()
}

func Test_Insert(t *testing.T) {
	repo := NewUserRepo()

	user := &model.User{
		Username: util.EnPointer("username1"),
		Password: util.EnPointer("password2"),
	}
	repo.Insert(nil, user)
}

func Test_Delete(t *testing.T) {
	repo := NewUserRepo()

	user := &model.User{
		Id: util.EnPointer(13),
	}

	repo.Delete(nil, user)
}

func Test_Update(t *testing.T) {
	repo := NewUserRepo()

	user := &model.User{
		Id:       util.EnPointer(13),
		Username: util.EnPointer("username13"),
	}
	repo.Update(nil, user)
}

func Test_Select(t *testing.T) {
	repo := NewUserRepo()

	user := &model.User{
		Username: util.EnPointer("username1"),
	}
	users, _ := repo.Select(nil, user)
	jsonStr, _ := json.Marshal(users)
	t.Log(string(jsonStr))
}

func Test_SelectById(t *testing.T) {
	repo := NewUserRepo()

	user := &model.User{
		Id: util.EnPointer(10),
	}
	user, _ = repo.SelectById(nil, user)
	jsonStr, _ := json.Marshal(user)
	t.Log(string(jsonStr))
}

func Test_SelectByUsername(t *testing.T) {
	repo := NewUserRepo()

	user := &model.User{
		Username: util.EnPointer("username13"),
	}
	user, _ = repo.SelectByUsername(nil, user)
	jsonStr, _ := json.Marshal(user)
	t.Log(string(jsonStr))
}

func Test_SelectWithPagination(t *testing.T) {
	repo := NewUserRepo()

	p := &model.Pagination{
		Page: 1,
		Size: 2,
	}

	repo.SelectWithPagination(nil, p)
	jsonStr, _ := json.Marshal(p.Data)
	t.Log(string(jsonStr))
}
