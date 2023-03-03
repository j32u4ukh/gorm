package stmt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// 對傳入字串進行前處理，避免 SQL injection
////////////////////////////////////////////////////////////////////////////////////////////////////
var symbols []rune = []rune{'\\', '\''}

// 傳入參數可能是(字串形式的)數字 或 真的就是字串
func AntiInjection(compoment string) (error, string) {
	hasPrefix := strings.HasPrefix(compoment, "'")
	hasSuffix := strings.HasSuffix(compoment, "'")

	// 前後都沒有單引號，預期傳入的是數字
	if !hasPrefix && !hasSuffix {
		if containSymbols(compoment) {
			return errors.New(fmt.Sprintf("預期傳入數字，但卻包含特殊符號，此為不合法參數(%s)", compoment)), ""
		}
		return nil, compoment
	}

	// 數字不會有單引號，而字串要求前後都要有單引號
	if (hasPrefix && !hasSuffix) || (!hasPrefix && hasSuffix) {
		return errors.New(fmt.Sprintf("只有前端或後端有單引號(')，此為不合法參數(%s)", compoment)), ""
	}

	return nil, fmt.Sprintf("'%s'", AntiInjectionString(strings.Trim(compoment, "'")))
}

// compoment: 不包含前後單引號的字串
func AntiInjectionString(compoment string) string {
	var buffer bytes.Buffer
	var c rune
	for _, c = range compoment {
		if isSymbols(c) {
			buffer.WriteRune('\\')
		}
		buffer.WriteRune(c)
	}
	return buffer.String()
}

func containSymbols(str string) bool {
	var r rune
	for _, r = range str {
		if isSymbols(r) {
			return true
		}
	}
	return false
}

func isSymbols(r rune) bool {
	for _, s := range symbols {
		if s == r {
			return true
		}
	}
	return false
}
