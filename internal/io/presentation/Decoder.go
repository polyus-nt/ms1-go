package presentation

import (
	"fmt"
)

// Decoder декодирует поля в сообщении (байтовая строка)
func Decoder(fields []Field, s string) (res []int64, err error) {

	for _, field := range fields {

		num, err := GetHex(field, s)
		if err != nil {
			return res, fmt.Errorf("{ decoder error: %v}", err)
		}
		res = append(res, num)
	}

	return
}