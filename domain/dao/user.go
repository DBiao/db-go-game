package dao

import (
	"db-go-game/domain/model"
	"db-go-game/pkg/common/dmysql"
	"db-go-game/pkg/entity"
)

type IUserDao interface {
	Create(user *model.User) error
	VerifyIdentity(w *entity.MysqlWhere) (*model.User, error)
}

type userDao struct {
}

func NewAuthRepository() IUserDao {
	return &userDao{}
}

func (*userDao) VerifyIdentity(w *entity.MysqlWhere) (*model.User, error) {
	user := new(model.User)
	db := dmysql.GetDB()
	err := db.Where(w.Query, w.Args...).First(user).Error
	return user, err
}

func (*userDao) Create(user *model.User) error {
	db := dmysql.GetDB()
	return db.Create(user).Error
}
