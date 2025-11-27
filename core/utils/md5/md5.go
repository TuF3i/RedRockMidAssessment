package md5

import (
	"crypto/md5"
	"fmt"
)

func GenMD5(password string) string {
	h := md5.New()                       // 1. 新建 MD5 摘要器
	h.Write([]byte(password))            // 2. 写入数据
	return fmt.Sprintf("%x", h.Sum(nil)) // 3. 取 16 字节并转小写 16 进制
}
