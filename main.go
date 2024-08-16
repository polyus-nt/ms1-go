package main

import (
	"fmt"
	"ms1-tool-go/internal"
)

func main() {

	for _, v := range internal.TestData {
		internal.PrintOneChunk(v)
		fmt.Println("")
	}
	internal.Xxd(internal.TestData)
	s := "123"
	fmt.Println(s)
}