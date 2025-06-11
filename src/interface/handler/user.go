package handler

import (
	"api-sample-with-echo-ddd/src/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type UserHandler interface {
	Post(c echo.Context) error
	Get(c echo.Context) error
	GetAll(c echo.Context) error
	Put(c echo.Context) error
	Delete(c echo.Context) error
}

type userHandler struct {
	userUsecase usecase.UserUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase) UserHandler {
	return &userHandler{userUsecase: userUsecase}
}

type reqUser struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type resUser struct {
	ID        string `json:"id"`
	Name      string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *userHandler) Post(c echo.Context) error {
	var reqUser reqUser
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userUsecase.Create(reqUser.Name, reqUser.Email, reqUser.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resUser := resUser{
		ID:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusCreated, resUser)
}

func (h *userHandler) Get(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userUsecase.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resUser := resUser{
		ID:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, resUser)
}

func (h *userHandler) GetAll(c echo.Context) error {
	users, err := h.userUsecase.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resUsers := make([]resUser, len(users))
	for i, user := range users {
		resUsers[i] = resUser{
			ID:        user.ID,
			Name:      user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, resUsers)
}

func (h *userHandler) Put(c echo.Context) error {
	id := c.Param("id")

	var reqUser reqUser
	if err := c.Bind(&reqUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userUsecase.Update(id, reqUser.Name, reqUser.Email, reqUser.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resUser := resUser{
		ID:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, resUser)
}

func (h *userHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := h.userUsecase.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
