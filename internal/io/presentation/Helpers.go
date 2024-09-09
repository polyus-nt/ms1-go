package presentation

import (
	"encoding/hex"
	"fmt"
	"ms1-tool-go/internal/io/entity"
	"ms1-tool-go/internal/xxd"
	"os"
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
	fmt.Println(len(data))
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
func PrettyFrame(frame entity.Frame) {
	fmt.Printf("Frame %v.%v:\n", frame.Page, frame.Part)
	xxd.PrintOneChunk(xxd.Bin(frame.Blob)) // change it in the future...
}

// GetPart Изымает нужный кусочек (поле) из строки данных
func GetPart(field Field, s string) string {

	remains := min(len(s)-field.Start, field.Len)
	return s[field.Start:][:remains]
}

// GetHex иъятие беззнакового шестнацеричного числа из поля
func GetHex(field Field, s string) (int64, error) {

	blob := GetPart(field, s)
	num, err := strconv.ParseInt(blob, 16, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getHex [%v] failed! {error: %v}", field.Descr, err)
	}

	return num, err
}

// GetSignedHex иъятие знакового шестнацеричного числа из поля
func GetSignedHex(field Field, s string) (int64, error) {

	blob := GetPart(Field{Start: field.Start + 1, Len: field.Len - 1}, s)

	num, err := strconv.ParseInt(blob, 16, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getSignedHex [%v] failed! {error: %v}", field.Descr, err)
	}

	if GetPart(Field{Len: 1}, s)[0] == '-' {
		num = -num
	}

	return num, err
}