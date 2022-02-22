package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomFormatTime(t *testing.T) {
	time := NowTime()
	date := time.Format("2006-01-02 15:04:05")
	assert.Equal(t, date, CustomFormatTime(time))
}

func TestStrtotime(t *testing.T) {
	time, err := Strtotime("2021-12-30 12:12:12")
	assert.Nil(t, err)
	assert.Equal(t, int64(1640837532), time.Unix())
}

func TestNowTime(t *testing.T) {
	assert.NotNil(t, NowTime())
}

func TestDateTime(t *testing.T) {
	time, _ := Strtotime(DateTime())
	assert.GreaterOrEqual(t, NowTime().Unix(), time.Unix())
}
