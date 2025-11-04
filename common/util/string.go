// Package util @Author:冯铁城 [17615007230@163.com] 2025-11-04 11:40:27
package util

import "strings"

// RemoveSpacesAndNewlines 移除字符串中的空格、换行符
func RemoveSpacesAndNewlines(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == ' ' || r == '　' || r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
