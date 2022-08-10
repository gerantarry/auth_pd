package router

import (
	"io"
	"strings"
	"unicode"
)

type Logger interface {
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})
}

//убрать символы переноса строки и т.д. из входящих данных
func formatReaderData(r io.Reader) (string, error) {
	buf := make([]byte, 2048)
	n, err := r.Read(buf)
	var str, resStr string

	for {
		for _, b := range buf[:n] {
			str = str + string(b)
		}
		if err == io.EOF {
			resStr = cleanString(str)
			return resStr, nil
		}
		if err != nil {
			return "", err
		}

		resStr = cleanString(str)
		return resStr, nil
	}

}

//убирает сначала управляющие символы затем '\'
func cleanString(str string) string {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, str)
	return clean
}
