// Package ranking 提供排行相关的参数校验。
package ranking

// ValidLimit 判断排行数量 limit 是否合法。
//
// limit 必须为正整数。0 或负数在截断逻辑中等价于“不限制”，
// 会返回完整排行而非 Top-N，调用方极易将其误认为“无限制”，
// 因此一律视为非法。
func ValidLimit(limit int) bool {
	return limit > 0
}
