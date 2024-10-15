package presentation

import (
	"fmt"
)

// Decoder декодирует поля в сообщении (байтовая строка). В зависимости от типа, куда сохраняется результат, вызывается соответствующая функция
func Decoder(res []interface{}, fields []Field, s string) (err error) {

	if len(res) != len(fields) {
		return fmt.Errorf("decoder: field count mismatch")
	}

	for i, field := range fields {

		switch f := res[i].(type) {
		case *int:
			var num int64
			num, err = GetHex(field, s)
			*f = int(num)
		case *int8:
			var num int64
			num, err = GetHex(field, s)
			*f = int8(num)
		case *int16:
			var num int64
			num, err = GetHex(field, s)
			*f = int16(num)
		case *int32:
			var num int64
			num, err = GetHex(field, s)
			*f = int32(num)
		case *int64:
			var num int64
			num, err = GetHex(field, s)
			*f = num
		case *uint:
			var num int64
			num, err = GetHex(field, s)
			*f = uint(num)
		case *uint8:
			var num int64
			num, err = GetHex(field, s)
			*f = uint8(num)
		case *uint16:
			var num int64
			num, err = GetHex(field, s)
			*f = uint16(num)
		case *uint32:
			var num int64
			num, err = GetHex(field, s)
			*f = uint32(num)
		case *uint64:
			var num int64
			num, err = GetHex(field, s)
			*f = uint64(num)
		case *bool:
			*f, err = GetBool(field, s)
		case *string:
			*f, err = GetString(field, s)
		default:
			return fmt.Errorf("decoder: field type mismatch: %T", f)
		}

		if err != nil {
			return
		}
	}

	return
}
