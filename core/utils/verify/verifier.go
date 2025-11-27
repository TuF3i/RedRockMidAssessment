package verify

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func VerifyUserName(userName string) bool {
	var (
		miniLength = 2 //长度控制，最短2个，最长15个
		maxLength  = 15
		reserved   = map[string]struct{}{"admin": {}, "root": {}, "user": {}, "api": {}, "bob": {}}                        // 0. 系统保留字（bob是测试账号）
		userRe     = regexp.MustCompile(`^[a-zA-Z0-9\u4e00-\u9fa5]([a-zA-Z0-9_\u4e00-\u9fa5]*[a-zA-Z0-9\u4e00-\u9fa5])?$`) // 1. 合法字符：字母数字中文下划线；禁止首尾下划线 *AI生成的正则表达式*
	)

	// 校验用户名
	if strings.TrimSpace(userName) == "" {
		// 用户名不为空
		return false
	}
	if length := utf8.RuneCountInString(userName); length < miniLength || length > maxLength {
		// 用户名在2~15个字符
		return false
	}
	if _, ok := reserved[strings.ToLower(userName)]; ok {
		// 不可用保留字
		return false
	}
	if s := userName; strings.HasPrefix(s, "_") || strings.HasSuffix(s, "_") {
		// 不能用 _ 开头或结尾
		return false
	}
	if !userRe.MatchString(userName) {
		// 只能包含字母、数字、中文、_
		return false
	}
	allDigit := true
	for _, r := range []rune(userName) {
		if !unicode.IsDigit(r) {
			allDigit = false
			break
		}
	}
	if allDigit {
		// 不能为纯数字
		return false
	}

	return true
}

func VerifyUserID(stuID uint) bool {
	// 10位用户ID
	if len(strconv.Itoa(int(stuID))) == 10 {
		return true
	}
	return false
}

func VerifyStudentClass(stuClass string) bool {
	// 班级字符串在3~15个字符间
	length := utf8.RuneCountInString(stuClass)
	if !(length < 3 || length > 10) {
		return true
	}
	return false
}

func VerifyPassword(passwd string) bool {
	var (
		miniLength       = 6 //长度控制，最短6个，最长20个
		maxLength        = 20
		miniNumberNum    = 1                                // 最少含有一个数字
		miniLowerCharNum = 1                                // 至少一个小写字母
		miniUpperCharNum = 1                                // 至少一个大写字母
		miniSpecialsNum  = 1                                // 至少一个特殊字符
		specials         = "!@#$%^&*()_+-=[]{}|;':\",./<>?" // 特殊字符集

	)

	if length := utf8.RuneCountInString(passwd); length < miniLength || length > maxLength {
		// 检测密码长度
		return false
	}

	countNum := 0
	countLower := 0
	countUpper := 0
	countSpecial := 0

	for _, char := range passwd {
		if unicode.IsDigit(char) { //检测数字个数
			countNum++
		}
		if unicode.IsLower(char) { //检测小写字母个数
			countLower++
		}
		if unicode.IsUpper(char) { //检测大写字母个数
			countUpper++
		}
		if strings.ContainsRune(specials, char) {
			countSpecial++
		}
	}

	if countLower < miniLowerCharNum || countUpper < miniUpperCharNum || countNum < miniNumberNum || countSpecial < miniSpecialsNum {
		return false
	}

	return true
}

func VerifySexSetting(sex uint) bool {
	if sex > 2 {
		return false
	}
	return true
}

func VerifyGrade(grade uint) bool {
	if grade > 0 && grade < 5 {
		return true
	}
	return false
}

func VerifyAge(age uint) bool {
	if age > 10 && age < 60 {
		return true
	}
	return false
}
