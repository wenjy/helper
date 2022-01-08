package helper

import "os"

// 判断文件是否存在
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 可读的大小
func HumanizeSize(s uint64) string {
	return IBytes(s)
}
