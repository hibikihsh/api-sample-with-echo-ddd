package handler

import (
	"api-sample-with-echo-ddd/domain/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase is a mock implementation of UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Create(username string, email string, password string) (*model.User, error) {
	args := m.Called(username, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUseCase) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUseCase) FindAll() ([]*model.User, error) {
	args := m.Called()
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserUseCase) Update(id string, username string, email string, password string) (*model.User, error) {
	args := m.Called(id, username, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUseCase) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserHandler_Post(t *testing.T) {
	t.Run("成功: ユーザーを作成できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockUseCase.On("Create", "testuser", "test@example.com", "password123").Return(user, nil)

		requestBody := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Post(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response resUser
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-id", response.ID)
		assert.Equal(t, "testuser", response.Name)
		assert.Equal(t, "test@example.com", response.Email)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("失敗: 無効なリクエストボディ", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Post(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("失敗: ユースケースエラー", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		mockUseCase.On("Create", "testuser", "test@example.com", "password123").Return(nil, errors.New("validation error"))

		requestBody := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(requestBody)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Post(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_Get(t *testing.T) {
	t.Run("成功: ユーザーを取得できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockUseCase.On("FindByID", "test-id").Return(user, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/users/test-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		err := handler.Get(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response resUser
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-id", response.ID)
		assert.Equal(t, "testuser", response.Name)
		assert.Equal(t, "test@example.com", response.Email)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("失敗: ユーザーが見つからない", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		mockUseCase.On("FindByID", "nonexistent-id").Return(nil, errors.New("user not found"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/users/nonexistent-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("nonexistent-id")

		err := handler.Get(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_GetAll(t *testing.T) {
	t.Run("成功: 全ユーザーを取得できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		now := time.Now()
		users := []*model.User{
			{
				ID:        "1",
				Username:  "user1",
				Email:     "user1@example.com",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "2",
				Username:  "user2",
				Email:     "user2@example.com",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mockUseCase.On("FindAll").Return(users, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []resUser
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "1", response[0].ID)
		assert.Equal(t, "user1", response[0].Name)
		assert.Equal(t, "2", response[1].ID)
		assert.Equal(t, "user2", response[1].Name)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("失敗: ユースケースエラー", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		mockUseCase.On("FindAll").Return(([]*model.User)(nil), errors.New("database error"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_Put(t *testing.T) {
	t.Run("成功: ユーザーを更新できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "updateduser",
			Email:     "updated@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockUseCase.On("Update", "test-id", "updateduser", "updated@example.com", "newpassword").Return(user, nil)

		requestBody := map[string]string{
			"username": "updateduser",
			"email":    "updated@example.com",
			"password": "newpassword",
		}
		jsonBody, _ := json.Marshal(requestBody)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/users/test-id", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		err := handler.Put(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response resUser
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test-id", response.ID)
		assert.Equal(t, "updateduser", response.Name)
		assert.Equal(t, "updated@example.com", response.Email)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("失敗: 無効なリクエストボディ", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/users/test-id", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		err := handler.Put(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("失敗: ユースケースエラー", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		mockUseCase.On("Update", "test-id", "updateduser", "updated@example.com", "newpassword").Return(nil, errors.New("user not found"))

		requestBody := map[string]string{
			"username": "updateduser",
			"email":    "updated@example.com",
			"password": "newpassword",
		}
		jsonBody, _ := json.Marshal(requestBody)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/users/test-id", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		err := handler.Put(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_Delete(t *testing.T) {
	t.Run("成功: ユーザーを削除できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		mockUseCase.On("Delete", "test-id").Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/users/test-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		err := handler.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deleted successfully", response["message"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("失敗: ユースケースエラー", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase)

		mockUseCase.On("Delete", "test-id").Return(errors.New("user not found"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/users/test-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		err := handler.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUseCase.AssertExpectations(t)
	})
}