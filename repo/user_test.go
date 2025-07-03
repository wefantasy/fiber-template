package repo

import (
	"app/conf"
	"app/db"
	"app/log"
	"app/model"
	"app/util"
	"encoding/json"
	"testing"
)

func InitDbEnv() {
	conf.Initialize()
	log.Initialize()
	db.Initialize()
}

func Test_Insert(t *testing.T) {
	InitDbEnv()

	repo := NewUserRepo()
	user := &model.User{
		Username: util.EnPointer("username321212"),
		Password: util.EnPointer("password2"),
	}
	repo.Insert(nil, user)
}

func Test_Delete(t *testing.T) {
	InitDbEnv()

	repo := NewUserRepo()
	repo.Delete(nil, 8)
}

func Test_Update(t *testing.T) {
	InitDbEnv()
	repo := NewUserRepo()
	user := &model.User{
		Id:       util.EnPointer(9),
		Username: util.EnPointer("username13"),
	}
	repo.Update(nil, user)
}

func Test_Select(t *testing.T) {
	InitDbEnv()

	repo := NewUserRepo()
	user := &model.User{}
	users, _ := repo.Select(nil, user)
	jsonStr, _ := json.Marshal(users)
	t.Log(string(jsonStr))
}

func Test_SelectById(t *testing.T) {
	InitDbEnv()

	repo := NewUserRepo()
	user, _ := repo.SelectById(nil, 12)
	jsonStr, _ := json.Marshal(user)
	t.Log(string(jsonStr))
}

func Test_SelectByUsername(t *testing.T) {
	InitDbEnv()
	repo := NewUserRepo()
	user, _ := repo.SelectByUsername(nil, "username13")
	jsonStr, _ := json.Marshal(user)
	t.Log(string(jsonStr))
}

func Test_SelectWithPagination(t *testing.T) {
	InitDbEnv()
	repo := NewUserRepo()

	p := &model.Pagination{
		Page: 0,
		Size: 2,
	}

	repo.SelectWithPagination(nil, p)
	jsonStr, _ := json.Marshal(p.Data)
	t.Log(string(jsonStr))
}
