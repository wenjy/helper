# helper
go helper

## 格式化

### 字节转换

`Bytes(82854982)` output `83 MB`

`IBytes(82854982)` output `79 MB`

### IP转换

`Ip2long`

`Long2ip`

## 文件

### 文件是否存在

`FileExists`

## 字符串

### 生成随机字符串

`Random`

### 字符串、版本比较

`Compare`

`VersionCompare`

### InArray

`InArrayString`

### MAC转换为纯小写字母

`NormalizeMac`

### MD5

`MD5Hash`

### 根据错误判断连接是否关闭

`ConnIsClose`

## 时间

### 中国时间

`CustomFormatTime` 格式时间格式为 `Y-m-d H:i:s`

`Strtotime` 字符串转换为时间格式

`NowTime` 当前时间

`DateTime` 当前时间格式化为 `Y-m-d H:i:s`

## Pool

### 字节切片池

`DefaultBufferPool.Alloc`

`DefaultBufferPool.Free`
