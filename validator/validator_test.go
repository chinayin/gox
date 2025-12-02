package validator

import (
	"testing"
)

func TestValidate_OfficialRules(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name: "valid email",
			data: struct {
				Email string `validate:"required,email"`
			}{Email: "user@example.com"},
			wantErr: false,
		},
		{
			name: "invalid email",
			data: struct {
				Email string `validate:"required,email"`
			}{Email: "invalid-email"},
			wantErr: true,
		},
		{
			name: "valid port range",
			data: struct {
				Port int `validate:"required,min=1,max=65535"`
			}{Port: 8080},
			wantErr: false,
		},
		{
			name: "invalid port range",
			data: struct {
				Port int `validate:"required,min=1,max=65535"`
			}{Port: 70000},
			wantErr: true,
		},
		{
			name: "valid oneof",
			data: struct {
				Level string `validate:"oneof=debug info warn error"`
			}{Level: "info"},
			wantErr: false,
		},
		{
			name: "invalid oneof",
			data: struct {
				Level string `validate:"oneof=debug info warn error"`
			}{Level: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	v1 := New()
	v2 := New()

	// 每次调用 New() 应该返回不同的实例
	if v1 == v2 {
		t.Error("New() should return different instances")
	}
}

func TestGlobal(t *testing.T) {
	v1 := Global()
	v2 := Global()

	// Global() 应该返回相同的实例
	if v1 != v2 {
		t.Error("Global() should return the same instance")
	}
}

func TestValidate_SnowflakeID(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name: "valid snowflake id (int64)",
			data: struct {
				ID int64 `validate:"snowflake_id"`
			}{ID: 1234567890123456789},
			wantErr: false,
		},
		{
			name: "valid snowflake id (string)",
			data: struct {
				ID string `validate:"snowflake_id"`
			}{ID: "1234567890123456789"},
			wantErr: false,
		},
		{
			name: "invalid snowflake id (zero)",
			data: struct {
				ID int64 `validate:"snowflake_id"`
			}{ID: 0},
			wantErr: true,
		},
		{
			name: "invalid snowflake id (negative)",
			data: struct {
				ID int64 `validate:"snowflake_id"`
			}{ID: -1},
			wantErr: true,
		},
		{
			name: "invalid snowflake id (empty string)",
			data: struct {
				ID string `validate:"snowflake_id"`
			}{ID: ""},
			wantErr: true,
		},
		{
			name: "invalid snowflake id (non-numeric string)",
			data: struct {
				ID string `validate:"snowflake_id"`
			}{ID: "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetLocale(t *testing.T) {
	v := New()

	// 默认应该是英文
	if v.Locale() != "en" {
		t.Errorf("Default locale = %s, want en", v.Locale())
	}

	// 切换到中文
	if err := v.SetLocale("zh"); err != nil {
		t.Fatalf("SetLocale(zh) error = %v", err)
	}
	if v.Locale() != "zh" {
		t.Errorf("Locale = %s, want zh", v.Locale())
	}

	// 切换回英文
	if err := v.SetLocale("en"); err != nil {
		t.Fatalf("SetLocale(en) error = %v", err)
	}
	if v.Locale() != "en" {
		t.Errorf("Locale = %s, want en", v.Locale())
	}

	// 无效的 locale
	if err := v.SetLocale("invalid"); err == nil {
		t.Error("SetLocale(invalid) should return error")
	}
}

func TestWithLocale(t *testing.T) {
	v := New(WithLocale("zh"))
	if v.Locale() != "zh" {
		t.Errorf("Locale = %s, want zh", v.Locale())
	}
}

func TestTranslatedError(t *testing.T) {
	type TestStruct struct {
		Email string `validate:"required,email"`
	}

	tests := []struct {
		name   string
		locale string
		data   TestStruct
		want   string
	}{
		{
			name:   "english error",
			locale: "en",
			data:   TestStruct{Email: "invalid"},
			want:   "Email must be a valid email address",
		},
		{
			name:   "chinese error",
			locale: "zh",
			data:   TestStruct{Email: "invalid"},
			want:   "Email必须是一个有效的邮箱",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(WithLocale(tt.locale))
			err := v.Validate(tt.data)
			if err == nil {
				t.Fatal("Expected validation error")
			}

			// 检查错误消息是否包含预期的关键词
			errMsg := err.Error()
			if tt.locale == "en" && !contains(errMsg, "email") {
				t.Errorf("Error message = %s, should contain 'email'", errMsg)
			}
			if tt.locale == "zh" && !contains(errMsg, "邮箱") {
				t.Errorf("Error message = %s, should contain '邮箱'", errMsg)
			}
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && (s == substr || len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
