package api

import (
	"fintech-banking-app/util"

	"github.com/go-playground/validator/v10"
)

// `validator.Func` 获取一个 `validator.FieldLevel` 接口作为入参并返回 true 如果验证成功
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// `fieldLevel.Field()` 获取字段值
	// `.Interface()` 获取值作为一个 `interface{}`
	// 接着尝试转换为字符串
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
