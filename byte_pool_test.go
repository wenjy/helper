package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	buf12 := DefaultBufferPool.Alloc(12)
	assert.Equal(t, 12, len(buf12))
	assert.Equal(t, 16, cap(buf12))

	buf32 := DefaultBufferPool.Alloc(32)
	assert.Equal(t, 32, len(buf32))
	assert.Equal(t, 32, cap(buf32))
}

func TestAllocGtMax(t *testing.T) {
	bufMax := DefaultBufferPool.Alloc(32769)
	assert.Equal(t, 0, len(bufMax))
	assert.Equal(t, 32769, cap(bufMax))
}

func TestFree(t *testing.T) {
	buf12 := make([]byte, 12)

	// 无法放入
	DefaultBufferPool.Free(buf12)

	buf12 = DefaultBufferPool.Alloc(12)

	assert.Equal(t, 12, len(buf12))
	assert.Equal(t, 16, cap(buf12))

	buf12[0] = 1
	buf12[1] = 2
	buf12[2] = 3

	DefaultBufferPool.Free(buf12)

	buf12 = DefaultBufferPool.Alloc(12)
	assert.Equal(t, uint8(1), buf12[0])
	assert.Equal(t, uint8(2), buf12[1])
	assert.Equal(t, uint8(3), buf12[2])
}
