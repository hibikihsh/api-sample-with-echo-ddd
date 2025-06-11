package model

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username string, email string, password string) (User, error) {
	userID := NewUserID()
	userName, err := NewUserName(username)
	if err != nil {
		return User{}, err
	}
	userEmail, err := NewUserEmail(email)
	if err != nil {
		return User{}, err
	}
	userPassword, err := NewPassword(password)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        userID.value,
		Username:  userName.value,
		Email:     userEmail.value,
		Password:  userPassword.hashedValue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

type UserID struct {
	value string
}

func NewUserID() UserID {
	return UserID{
		value: "",
	}
}

type UserName struct {
	value string
}

func NewUserName(username string) (UserName, error) {
	if len(username) < 3 || len(username) > 20 {
		return UserName{}, errors.New("ユーザー名は3文字以上20文字以下で入力してください")
	}

	return UserName{
		value: username,
	}, nil
}

type Email struct {
	value string
}

func NewUserEmail(email string) (Email, error) {
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return Email{}, errors.New("メールアドレスが不正です")
	}

	return Email{
		value: email,
	}, nil
}

type Password struct {
	hashedValue string
}

func NewPassword(password string) (Password, error) {
	if len(password) < 8 {
		return Password{}, errors.New("パスワードは8文字以上で入力してください")
	}

	hasLetter := false
	hasNumber := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return Password{}, errors.New("パスワードは英数字を含む必要があります")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, errors.New("パスワードのハッシュ化に失敗しました")
	}

	return Password{
		hashedValue: string(hashedPassword),
	}, nil
}
