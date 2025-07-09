package usecase

import (
	"api-sample-with-echo-ddd/domain/model"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]*model.User, error) {
	args := m.Called()
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserUsecase_Create(t *testing.T) {
	t.Run("成功: ユーザーを作成できる", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		now := time.Now()
		expectedUser := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(expectedUser, nil)

		result, err := usecase.Create("testuser", "test@example.com", "password123")

		assert.NoError(t, err)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, "test@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: 無効なユーザー名", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		result, err := usecase.Create("ab", "test@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "ユーザー名は3文字以上20文字以下で入力してください")
	})

	t.Run("失敗: 無効なメールアドレス", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		result, err := usecase.Create("testuser", "invalid-email", "password123")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "メールアドレスが不正です")
	})

	t.Run("失敗: 無効なパスワード", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		result, err := usecase.Create("testuser", "test@example.com", "short")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "パスワードは8文字以上で入力してください")
	})

	t.Run("失敗: リポジトリエラー", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return((*model.User)(nil), errors.New("database error"))

		result, err := usecase.Create("testuser", "test@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_FindByID(t *testing.T) {
	t.Run("成功: ユーザーを取得できる", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		expectedUser := &model.User{
			ID:       "test-id",
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("FindByID", "test-id").Return(expectedUser, nil)

		result, err := usecase.FindByID("test-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: ユーザーが見つからない", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("FindByID", "nonexistent-id").Return(nil, errors.New("user not found"))

		result, err := usecase.FindByID("nonexistent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_FindAll(t *testing.T) {
	t.Run("成功: 全ユーザーを取得できる", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		expectedUsers := []*model.User{
			{ID: "1", Username: "user1", Email: "user1@example.com"},
			{ID: "2", Username: "user2", Email: "user2@example.com"},
		}

		mockRepo.On("FindAll").Return(expectedUsers, nil)

		result, err := usecase.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: リポジトリエラー", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("FindAll").Return(([]*model.User)(nil), errors.New("database error"))

		result, err := usecase.FindAll()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Update(t *testing.T) {
	t.Run("成功: ユーザーを更新できる", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		existingUser := &model.User{
			ID:       "test-id",
			Username: "olduser",
			Email:    "old@example.com",
			Password: "oldpassword",
		}

		updatedUser := &model.User{
			ID:       "test-id",
			Username: "newuser",
			Email:    "new@example.com",
			Password: "newpassword",
		}

		mockRepo.On("FindByID", "test-id").Return(existingUser, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(updatedUser, nil)

		result, err := usecase.Update("test-id", "newuser", "new@example.com", "newpassword")

		assert.NoError(t, err)
		assert.Equal(t, "newuser", result.Username)
		assert.Equal(t, "new@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: ユーザーが見つからない", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("FindByID", "nonexistent-id").Return(nil, errors.New("user not found"))

		result, err := usecase.Update("nonexistent-id", "newuser", "new@example.com", "newpassword")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Delete(t *testing.T) {
	t.Run("成功: ユーザーを削除できる", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		existingUser := &model.User{
			ID:       "test-id",
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("FindByID", "test-id").Return(existingUser, nil)
		mockRepo.On("Delete", existingUser).Return(nil)

		err := usecase.Delete("test-id")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: ユーザーが見つからない", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("FindByID", "nonexistent-id").Return(nil, errors.New("user not found"))

		err := usecase.Delete("nonexistent-id")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("失敗: 削除エラー", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		existingUser := &model.User{
			ID:       "test-id",
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("FindByID", "test-id").Return(existingUser, nil)
		mockRepo.On("Delete", existingUser).Return(errors.New("delete error"))

		err := usecase.Delete("test-id")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
		mockRepo.AssertExpectations(t)
	})
}