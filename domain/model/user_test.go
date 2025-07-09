package model

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	type TestCase struct {
		name          string
		username      string
		email         string
		password      string
		expectedError bool
		errorMessage  string
	}
	testCases := []TestCase{
		{"有効なユーザー", "testuser", "test@example.com", "password123", false, ""},
		{"ユーザー名が短すぎる", "test", "test@example.com", "pass123", true, "ユーザー名は3文字以上20文字以下で入力してください"},
		{"ユーザー名が長すぎる", "testuser12345678901234567890", "test@example.com", "password123", true, "ユーザー名は3文字以上20文字以下で入力してください"},
		{"メールアドレスが不正", "testuser", "test@example", "password123", true, "メールアドレスが不正です"},
		{"パスワードが短すぎる", "testuser", "test@example.com", "pass123", true, "パスワードは8文字以上で入力してください"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUser(tc.username, tc.email, tc.password)
			if tc.expectedError && err == nil {
				t.Errorf("Expected error, but got nil")
			}
		})
	}
}

func TestNewUserName(t *testing.T) {
	type TestCase struct {
		name          string
		input         string
		expectedError bool
		errorMessage  string
	}
	testCases := []TestCase{
		{"有効なユーザー名", "testuser", false, ""},
		{"最小長ユーザー名", "te", true, "ユーザー名は3文字以上20文字以下で入力してください"},
		{"最大長ユーザー名", "testuser12345678901234567890", true, "ユーザー名は3文字以上20文字以下で入力してください"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUserName(tc.input)
			if tc.expectedError && err == nil {
				t.Errorf("Expected error, but got nil")
			}
		})
	}
}

func TestNewUserEmail(t *testing.T) {
	type TestCase struct {
		name          string
		input         string
		expectedError bool
		errorMessage  string
	}
	testCases := []TestCase{
		{"有効なメールアドレス", "test@example.com", false, ""},
		{"ドット付きメールアドレス", "test.email@example.com", false, ""},
		{"アットマークなし", "testexample.com", true, "メールアドレスが不正です"},
		{"ドメインなし", "test@", true, "メールアドレスが不正です"},
		{"空文字列", "", true, "メールアドレスが不正です"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUserEmail(tc.input)
			if tc.expectedError && err == nil {
				t.Errorf("Expected error, but got nil")
			}
		})
	}
}

func TestNewPassword(t *testing.T) {
	type TestCase struct {
		name          string
		input         string
		expectedError bool
		errorMessage  string
	}
	testCases := []TestCase{
		{"有効なパスワード", "password123", false, ""},
		{"最小長パスワード", "pass123a", false, ""},
		{"短すぎるパスワード", "pass123", true, "パスワードは8文字以上で入力してください"},
		{"数字なし", "password", true, "パスワードは英数字を含む必要があります"},
		{"文字なし", "12345678", true, "パスワードは英数字を含む必要があります"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPassword(tc.input)
			if tc.expectedError && err == nil {
				t.Errorf("Expected error, but got nil")
			}
		})
	}
}
