package md5

import "testing"

func TestMD5(t *testing.T) {
	md5 := GenMD5("root")
	t.Log(md5)
}
