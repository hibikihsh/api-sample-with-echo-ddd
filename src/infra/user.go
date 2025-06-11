package infra

import (
	"api-sample-with-echo-ddd/src/domain/model"
	"api-sample-with-echo-ddd/src/domain/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	user := &model.User{ID: id}

	if err := r.db.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	users := []*model.User{}

	if err := r.db.Find(users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Delete(user *model.User) error {
	if err := r.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}
