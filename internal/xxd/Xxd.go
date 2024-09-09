package xxd

import (
	"bytes"
	"fmt"
)

type Bin string

var TestData = []Bin{
	"0800585920002000",
	"080058a9080058a9",
	"0000000000000000",
	"0000000000000000",
	"0000000000000000",
	"080058a900000000",
	"0000000000000000",
	"0800538d080058a9",
	"00000000080058a9",
	"080058a9080058a9",
	"080058a9080058a9",
	"080058a9080058a9",
	"080058a900000000",
	"080058a9080058a9",
	"080058a9080058a9",
	"00000000080058a9",
}

func PrintOneChunk(bin Bin) string {

	buf := bytes.Buffer{}

	i := 0
	for i < len(bin) {
		remains := min(i+4, len(bin))
		if remains-i > 2 {
			buf.WriteString(fmt.Sprintf("%v ", bin[i:remains]))
		}
		i = remains
	}

	return buf.String()
}

func Xxd(bin []Bin) string {

	buf := bytes.Buffer{}

	num := 0
	for i := 1; i < len(bin); i += 2 {
		buf.WriteString(fmt.Sprintf("%08x: ", num*16))
		buf.WriteString(PrintOneChunk(bin[i-1]))
		buf.WriteString(PrintOneChunk(bin[i]))
		num++
		buf.WriteByte('\n')
	}

	return buf.String()
}