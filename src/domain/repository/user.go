package repository

import "api-sample-with-echo-ddd/src/domain/model"

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindByID(id string) (*model.User, error)
	FindAll() ([]*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(user *model.User) error
}
