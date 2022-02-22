package helper

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试规范化MAC
func TestNormalizeMac(t *testing.T) {
	mac := "00:0c:29:22:e0:AA"
	assert.Equal(t, "000c2922e0aa", NormalizeMac(mac), "they should be equal")
	assert.Equal(t, "000c2922e013", NormalizeMac("000c2922e013"), "they should be equal")
}

// 基准测试规范化MAC
func BenchmarkNormalizeMac(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NormalizeMac("00:0c:29:22:e0:AA")
	}
}

// 测试判断连接是否关闭
func TestConnIsClose(t *testing.T) {
	assert.True(t, ConnIsClose(errors.New("use of closed network connection")))
	assert.True(t, ConnIsClose(errors.New("the mux has closed")))
	assert.True(t, ConnIsClose(io.EOF))
}

func TestMin(t *testing.T) {
	min1 := Min(1, 2)
	assert.Equal(t, 1, min1)
	min2 := Min(3, 2)
	assert.Equal(t, 2, min2)
}

func TestMax(t *testing.T) {
	max1 := Max(1, 2)
	assert.Equal(t, 2, max1)
	max2 := Max(3, 2)
	assert.Equal(t, 3, max2)
}

func TestIp2long(t *testing.T) {
	assert.Equal(t, uint32(3071182317), Ip2long("183.14.133.237"))
	assert.Equal(t, uint32(2130706433), Ip2long("127.0.0.1"))
	assert.Equal(t, uint32(0), Ip2long("a.0.0.1"))
	assert.Equal(t, uint32(0), Ip2long("aaa"))
	assert.Equal(t, uint32(0), Ip2long("256.1.1.1"))
}

func TestLong2Ip(t *testing.T) {
	assert.Equal(t, "183.14.133.237", Long2ip(uint32(3071182317)))
	assert.Equal(t, "127.0.0.1", Long2ip(uint32(2130706433)))
}

func TestMD5Hash(t *testing.T) {
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", MD5Hash("test"))
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e", MD5Hash("123456"))
}

func TestInArrayString(t *testing.T) {
	strs := []string{"test1", "test2"}

	assert.True(t, InArrayString("test1", strs))
	assert.True(t, InArrayString("test2", strs))
	assert.False(t, InArrayString("test3", strs))
}
