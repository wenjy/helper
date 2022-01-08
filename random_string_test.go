package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试随机字符串
func TestRandom(t *testing.T) {
	assert.Equal(t, 3, len(Random(3)))

	str1 := Random(32)
	str2 := Random(32)
	assert.Equal(t, 32, len(str1))
	assert.NotEqual(t, str1, str2)
}

// 基准测试随机字符串
func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Random(32)
	}
}
