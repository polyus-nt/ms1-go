package presentation

import (
	"fmt"
	"os"
)

// Decoder декодирует поля в сообщении (байтовая строка)
func Decoder(fields []Field, s string) ([]int64, error) {

	var res []int64

	for _, field := range fields {
		num, err := GetHex(field, s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decoder error: %v", err)
			return res, err
		}
		res = append(res, num)
	}

	return res, nil
}