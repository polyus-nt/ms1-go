package presentation

import (
	"encoding/hex"
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/io/entity"
	"github.com/polyus-nt/ms1-go/internal/xxd"
	"strconv"
)

/* aliases */

type Packet = entity.Packet
type Field = entity.Field
type Frame = entity.Frame
type Address = entity.Address
type Mode = entity.Mode
type Load = entity.Load
type V = entity.V
type N = entity.N
type F = entity.F

func enc(data []byte) string {
	return hex.EncodeToString(data)
}

// chopBs разбивает строку на равные кусочки размером size, заполняя пустоту в конце при помощи \x00
func chopBs(size int, data string) []string {

	res := make([]string, 0, (len(data)+size)/size)

	for i := 0; i < len(data); i += size {

		if i+size < len(data) {
			res = append(res, data[i:i+size])
		} else {
			str := make([]byte, size)
			remains := copy(str, data[i:])
			for i := remains; i < len(str); i++ {
				str[i] = '\x00'
			}
			res = append(res, string(str))
		}
	}

	return res
}

// PrettyFrame выводит данные фрейма
func PrettyFrame(frame entity.Frame) (res string) {
	res = fmt.Sprintf("Frame %v.%v:\n", frame.Page, frame.Part)
	res += xxd.PrintOneChunk(xxd.Bin(frame.Blob))
	return
}

// GetPart Изымает нужный кусочек (поле) из строки данных
func GetPart(field Field, s string) string {

	remains := min(len(s)-field.Start, field.Len)
	return s[field.Start:][:remains]
}

// GetHex изъятие без знакового шестнадцатеричного числа из поля
func GetHex(field Field, s string) (num int64, err error) {

	blob := GetPart(field, s)
	num, err = strconv.ParseInt(blob, 16, 64)

	if err != nil {
		return num, fmt.Errorf("getHex failed: { fieldDescr: %v; err: %v }", field.Descr, err)
	}
	return
}

// GetBool изъятие логического значения из поля
func GetBool(field Field, s string) (res bool, err error) {

	blob := GetPart(field, s)

	if blob[0] == 't' {
		res = true
	} else if blob[0] == 'f' {
		res = false
	} else {
		err = fmt.Errorf("getBool failed: { fieldDescr: %v; err: %v}", field.Descr, err)
	}

	return
}

// GetString изъятие строки из поля
func GetString(field Field, s string) (string, error) {

	return GetPart(field, s), nil
}

// GetSignedHex изъятие знакового шестнадцатеричного числа из поля
func GetSignedHex(field Field, s string) (num int64, err error) {

	blob := GetPart(Field{Start: field.Start + 1, Len: field.Len - 1}, s)

	num, err = strconv.ParseInt(blob, 16, 64)

	if err != nil {
		return num, fmt.Errorf("getSignedHex failed: { fieldDescr: %v; err: %v }", field.Descr, err)
	}

	if GetPart(Field{Len: 1}, s)[0] == '-' {
		num = -num
	}
	return
}
