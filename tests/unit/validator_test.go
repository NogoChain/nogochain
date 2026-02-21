package unit

import (
	"testing"

	"nogochain/core/validator"
)

// TestValidator_NewValidator 测试创建新验证器
func TestValidator_NewValidator(t *testing.T) {
	// 创建新验证器
	v := validator.NewValidator()

	// 验证验证器不为空
	if v == nil {
		t.Fatal("验证器为nil")
	}
}
