package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试实例化时间轮
func TestNewTimeWheel(t *testing.T) {
	_, err := NewTimeWheel(time.Second/10-1, 1, false)
	assert.Equal(t, "invalid params, must tick >= 100 ms", err.Error())

	_, err = NewTimeWheel(time.Second, -1, false)
	assert.Equal(t, "invalid params, must bucketsNum > 0", err.Error())

	_, err = NewTimeWheel(time.Second, 1, false)
	assert.Nil(t, err)
}

// 测试时间轮的功能
func TestTimeWheel(t *testing.T) {
	tw, _ := NewTimeWheel(time.Second, 1, false)

	t1 := tw.Add(time.Second, func() {})
	assert.NotEmpty(t, t1)
	err := tw.Remove(t1)
	assert.Nil(t, err)

	t2 := tw.AddCron(time.Second, func() {})
	assert.NotEmpty(t, t2)
	err = tw.Remove(t2)
	assert.Nil(t, err)

	tw.Start()
	tw.Start()
	tw.After(time.Millisecond)
	tw.Sleep(time.Millisecond)
	tw.Stop()
}
