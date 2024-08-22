package internal

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
)

type Address struct {
	val string
}

//goland:noinspection SpellCheckingInspection
var TestAddress = Address{"facedeadbeef0000"}
var ZeroAddress = Address{"0000000000000000"}

type Frame struct {
	page uint8
	part uint8
	blob string
}

func enc(data []byte) string {
	return hex.EncodeToString(data)
}

func EncodeFrameLoad(frame Frame) string {

	var buf bytes.Buffer

	buf.WriteString(frame.blob)
	slices.Reverse(buf.Bytes())

	encoded := enc(buf.Bytes())

	buf.Reset()
	buf.WriteString(encoded)

	fmt.Fprintf(&buf, "%d", frame.page)
	fmt.Fprintf(&buf, "%d", frame.part)

	return buf.String()
}

func EncodeFrame(frame Frame, addr Address, mark uint8) string {
	var buf bytes.Buffer

	buf.WriteString(frame.blob)
	buf.WriteString(string(frame.page))
	buf.WriteString(string(frame.part))
	buf.WriteString(addr.val)
	buf.WriteString(string(mark))
	encoded := enc(buf.Bytes())

	buf.Reset()
	buf.WriteString("--------") // CRC32
	buf.WriteString(encoded)
	buf.WriteString("fr:")

	return buf.String()
}

var sizePage int = 2048
var sizeFrame int = 128

func chopBs(size int, data string) []string {

	res := make([]string, 0, (len(data)+size)/size)

	for i := 0; i < len(data); i += size {
		if i+size < len(data) {
			res = append(res, data[i:i+size])
		} else {
			str := make([]byte, 0, size)
			remains := copy(str, data[i:])
			for i := remains; i < len(str); i++ {
				str[i] = '\x00'
			}
			res = append(res, string(str))
		}
	}

	return res
}

func fileToPages(filePath string) []string {

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from fileToPages with filePath='%v'; Error desc -> %v", filePath, err)
	}

	res := chopBs(sizePage, string(content))

	return res
}

func pageToFrames(pageIndex uint8, page string) []Frame {

	packets := chopBs(sizeFrame, page)
	var frames []Frame

	for i, packet := range packets {
		frames = append(frames, Frame{page: pageIndex, part: uint8(i), blob: packet})
	}

	return frames
}

func FileToFrames(filePath string) []Frame {

	pages := fileToPages(filePath)
	var frames []Frame

	for i, page := range pages {

		frames = append(frames, pageToFrames(uint8(i), page)...)
	}

	return frames
}

func PrettyFrame(frame Frame) {
	fmt.Printf("Frame %v.%v:\n", frame.page, frame.part)
	PrintOneChunk(Bin(frame.blob)) // change it in the future...
}

func Test1() {
	frames := FileToFrames("data/test1.bin")

	for _, frame := range frames {
		fmt.Printf("Frame %v.%v:\n", frame.page, frame.part)
		EncodeFrame(frame, TestAddress, 0xee)
	}
}

func Test2() {
	content, err := os.ReadFile("data/test1.bin")

	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Printf("First 20 bytes from file: %v\n", content[:20])
	fmt.Print("Check encoding: ")
	fmt.Println(hex.DecodeString(hex.EncodeToString(content[:20])))
}

func Test3() {

	file, err := os.Create("data/test3.bin")
	defer file.Close()

	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	data := make([]byte, 256*4)
	for i := 0; i < 256*4; i++ {
		data[i] = byte(i / 4)
	}

	for i := 0; i < 1000; i++ {
		file.Write(data)
	}
}