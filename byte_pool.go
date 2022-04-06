package helper

import (
	"sort"
	"sync"
)

// DefaultBufferPool 默认全局buffer池.
var DefaultBufferPool = NewBytePool(32768, []int{8, 16, 32, 64, 128, 256, 512, 1024, 1032, 2048, 4096, 8192, 16384, 32768})

// BytePool 字节缓冲池结构体
type BytePool struct {
	classes     []sync.Pool //各大小池子
	classesSize []int       //池子对应大小
	maxSize     int         //最大字节数
}

// NewBytePool 创建字节缓冲池
func NewBytePool(maxSize int, sizeList []int) *BytePool {
	sort.Ints(sizeList)

	pool := &BytePool{
		classes:     make([]sync.Pool, len(sizeList)),
		classesSize: make([]int, len(sizeList)),
		maxSize:     maxSize,
	}
	n := 0
	for _, chunkSize := range sizeList {
		pool.classesSize[n] = chunkSize
		pool.classes[n].New = func(size int) func() interface{} {
			return func() interface{} {
				buf := make([]byte, size)
				return &buf
			}
		}(chunkSize)
		n++
	}
	return pool
}

func (pool *BytePool) Alloc(size int) []byte {
	if size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] >= size {
				mem := pool.classes[i].Get().(*[]byte)
				return (*mem)[:size]
			}
		}
	}
	return make([]byte, 0, size)
}

func (pool *BytePool) Free(mem []byte) {
	if size := cap(mem); size <= pool.maxSize {
		for i := 0; i < len(pool.classesSize); i++ {
			if pool.classesSize[i] == size {
				pool.classes[i].Put(&mem)
				return
			}
		}
	}
}
