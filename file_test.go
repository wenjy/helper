package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试文件是否存在
func TestFileExists(t *testing.T) {
	assert.True(t, FileExists("file_test.go"))
	assert.False(t, FileExists("file_test.go.noexists"))
}

func TestHumanizeSize(t *testing.T) {
	assert.Equal(t, "1023 B", HumanizeSize(1023))
	assert.Equal(t, "1.0 KiB", HumanizeSize(1024))
	assert.Equal(t, "1.0 MiB", HumanizeSize(1024*1024))
	assert.Equal(t, "1.0 GiB", HumanizeSize(1024*1024*1024))
}
