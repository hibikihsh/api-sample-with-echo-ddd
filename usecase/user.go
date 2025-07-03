package usecase

import (
	"api-sample-with-echo-ddd/domain/model"
	"api-sample-with-echo-ddd/domain/repository"
)

type UserUseCase interface {
	Create(username string, email string, password string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	FindAll() ([]*model.User, error)
	Update(id string, username string, email string, password string) (*model.User, error)
	Delete(id string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUseCase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Create(username string, email string, password string) (*model.User, error) {
	user, err := model.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	if _, err := u.userRepo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userUsecase) FindByID(id string) (*model.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) FindAll() ([]*model.User, error) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) Update(id string, username string, email string, password string) (*model.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	user.Username = username
	user.Email = email
	user.Password = password
	if _, err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Delete(id string) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if err := u.userRepo.Delete(user); err != nil {
		return err
	}
	return nil
}
