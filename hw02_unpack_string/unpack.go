package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func UnpackEscape(letter *rune) {
	switch *letter {
	case 't':
		*letter = '\t'
	case 'n':
		*letter = '\n'
	case 'r':
		*letter = '\r'
	case 'b':
		*letter = ' '
	case 'v':
		*letter = '\v'
	}
}

func GetWindow(runestr []rune, i int) (prefix rune, letter rune, postfix rune, e error) {
	letter = runestr[i]
	e = nil
	if i == 0 && unicode.IsDigit(letter) {
		return 0, 0, 0, ErrInvalidString
	}
	if i == 0 {
		prefix = 0
	} else {
		prefix = runestr[i-1]
	}
	if i == len(runestr)-1 {
		postfix = 0
	} else {
		postfix = runestr[i+1]
	}
	if unicode.IsDigit(letter) && (prefix != '\\' && unicode.IsDigit(postfix) || unicode.IsDigit(prefix)) { //число
		return 0, 0, 0, ErrInvalidString
	}
	return
}
func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	var b strings.Builder
	var runestr = []rune(str)
	for i := 0; i < len(runestr); i++ {
		prefix, letter, postfix, err := GetWindow(runestr, i)
		if err != nil {
			return "", err
		}
		if prefix == '\\' { //эскейп символы
			UnpackEscape(&letter)
		}
		if letter == '\\' {
			if prefix != '\\' {
				continue //переходим к экранированному символу
			}
		}
		if unicode.IsDigit(postfix) {
			postfix -= 48                                               //приводим руну к числу
			b.WriteString(strings.Repeat(string(letter), int(postfix))) //повторяем любой символ
			i++
		} else {
			if letter == '\\' {
				if prefix == '\\' {
					b.WriteRune(letter) //пишим экранированный слеш
					i++
				}
				continue
			}
			b.WriteRune(letter)
		}
	}
	return b.String(), nil
}
