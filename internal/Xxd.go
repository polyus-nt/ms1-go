package internal

import (
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

func PrintOneChunk(bin Bin) {
	i := 0
	for i < len(bin) {
		remains := min(i+4, len(bin))
		if remains-i > 2 {
			fmt.Printf("%v ", bin[i:remains])
		}
		i = remains
	}
}

func Xxd(bin []Bin) {
	num := 0
	for i := 1; i < len(bin); i += 2 {
		fmt.Printf("%08x: ", num*16)
		PrintOneChunk(bin[i-1])
		PrintOneChunk(bin[i])
		num++
		fmt.Println()
	}
}