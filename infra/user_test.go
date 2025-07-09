package infra

import (
	"api-sample-with-echo-ddd/domain/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic("failed to migrate database")
	}
	
	return db
}

func TestUserRepository_Create(t *testing.T) {
	t.Run("成功: ユーザーを作成できる", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Act
		result, err := repo.Create(user)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, "test@example.com", result.Email)

		// データベースに保存されていることを確認
		var savedUser model.User
		err = db.First(&savedUser, "id = ?", "test-id").Error
		assert.NoError(t, err)
		assert.Equal(t, "testuser", savedUser.Username)
	})

	t.Run("失敗: 重複したIDでの作成", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user1 := &model.User{
			ID:        "duplicate-id",
			Username:  "user1",
			Email:     "user1@example.com",
			Password:  "password1",
			CreatedAt: now,
			UpdatedAt: now,
		}
		user2 := &model.User{
			ID:        "duplicate-id",
			Username:  "user2",
			Email:     "user2@example.com",
			Password:  "password2",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Act
		_, err1 := repo.Create(user1)
		_, err2 := repo.Create(user2)

		// Assert
		assert.NoError(t, err1)
		assert.Error(t, err2)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Run("成功: ユーザーを取得できる", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		db.Create(user)

		// Act
		result, err := repo.FindByID("test-id")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("失敗: 存在しないIDでの取得", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}

		// Act
		result, err := repo.FindByID("nonexistent-id")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "record not found")
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	t.Run("成功: 全ユーザーを取得できる", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		users := []*model.User{
			{
				ID:        "user1",
				Username:  "testuser1",
				Email:     "user1@example.com",
				Password:  "password1",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "user2",
				Username:  "testuser2",
				Email:     "user2@example.com",
				Password:  "password2",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		for _, user := range users {
			db.Create(user)
		}

		// Act
		result, err := repo.FindAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		
		// ユーザー名でソートして比較
		usernames := []string{result[0].Username, result[1].Username}
		assert.Contains(t, usernames, "testuser1")
		assert.Contains(t, usernames, "testuser2")
	})

	t.Run("成功: ユーザーが存在しない場合は空のスライスを返す", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}

		// Act
		result, err := repo.FindAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("成功: ユーザーを更新できる", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		db.Create(user)

		// ユーザー情報を更新
		user.Username = "updateduser"
		user.Email = "updated@example.com"
		user.UpdatedAt = time.Now()

		// Act
		result, err := repo.Update(user)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, "updateduser", result.Username)
		assert.Equal(t, "updated@example.com", result.Email)

		// データベースで更新されていることを確認
		var updatedUser model.User
		err = db.First(&updatedUser, "id = ?", "test-id").Error
		assert.NoError(t, err)
		assert.Equal(t, "updateduser", updatedUser.Username)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	t.Run("成功: 存在しないユーザーの更新（新規作成）", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "nonexistent-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Act
		result, err := repo.Update(user)

		// Assert - GORMのSaveは存在しないレコードに対して新規作成を行うため、エラーは発生しない
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// 新規作成されたことを確認
		var savedUser model.User
		err = db.First(&savedUser, "id = ?", "nonexistent-id").Error
		assert.NoError(t, err)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("成功: ユーザーを削除できる", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		db.Create(user)

		// Act
		err := repo.Delete(user)

		// Assert
		assert.NoError(t, err)

		// データベースから削除されていることを確認
		var deletedUser model.User
		err = db.First(&deletedUser, "id = ?", "test-id").Error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")
	})

	t.Run("成功: 存在しないユーザーの削除（エラーなし）", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "nonexistent-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Act
		err := repo.Delete(user)

		// Assert - GORMのDeleteは存在しないレコードに対してもエラーを返さない
		assert.NoError(t, err)
	})
}

func TestUserRepository_ConcurrentAccess(t *testing.T) {
	t.Run("成功: 並行アクセスでの整合性", func(t *testing.T) {
		// Arrange
		db := setupTestDB()
		repo := &UserRepository{db: db}
		
		now := time.Now()
		user := &model.User{
			ID:        "concurrent-test-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		_, err := repo.Create(user)
		assert.NoError(t, err)

		// Act - 複数のgoroutineで同時にアクセス
		done := make(chan bool, 2)
		
		go func() {
			_, err := repo.FindByID("concurrent-test-id")
			assert.NoError(t, err)
			done <- true
		}()

		go func() {
			user.Username = "updated-concurrent"
			_, err := repo.Update(user)
			assert.NoError(t, err)
			done <- true
		}()

		// 両方のgoroutineが完了するまで待機
		<-done
		<-done

		// Assert
		var finalUser model.User
		err = db.First(&finalUser, "id = ?", "concurrent-test-id").Error
		assert.NoError(t, err)
		assert.Equal(t, "updated-concurrent", finalUser.Username)
	})
}