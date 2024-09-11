package presentation

import (
	"bytes"
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"os"
	"slices"
	"strconv"
)

// EncodeFrameLoad формирует на основе фрейма данные для отправки (байтовая строка)
func EncodeFrameLoad(frame Frame) string {

	var buf bytes.Buffer

	buf.WriteString(frame.Blob)
	slices.Reverse(buf.Bytes())

	encoded := enc(buf.Bytes())

	buf.Reset()
	buf.WriteString(encoded)

	fmt.Println(buf.String())

	fmt.Fprintf(&buf, "%s", enc([]byte{frame.Page}))
	fmt.Fprintf(&buf, "%s", enc([]byte{frame.Part}))

	fmt.Println(buf.String())

	return buf.String()
}

// EncodeFrame TODO разобраться с ее телом [deprecated --crc]
// формирует на основе фрейма данные для отправки (байтовая строка)
func EncodeFrame(frame Frame, addr Address, mark uint8) string {
	var buf bytes.Buffer

	buf.WriteString(frame.Blob)
	buf.WriteString(string(frame.Page))
	buf.WriteString(string(frame.Part))
	buf.WriteString(addr.Val)
	buf.WriteString(string(mark))
	encoded := enc(buf.Bytes())

	buf.Reset()
	buf.WriteString("--------") // CRC32
	buf.WriteString(encoded)
	buf.WriteString("fr:")

	return buf.String()
}

// fileToPages формирует страницы на основе данных из файла
func fileToPages(filePath string) []string {

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from fileToPages with filePath='%v'; Error desc -> %v", filePath, err)
	}

	res := chopBs(config.SIZE_PAGE, string(content))

	return res
}

// pageToFrames преобразует страницы в фреймы
func pageToFrames(pageIndex uint8, page string) []Frame {

	packets := chopBs(config.SIZE_FRAME, page)
	var frames []Frame

	for i, packet := range packets {
		frames = append(frames, Frame{Page: pageIndex, Part: uint8(i), Blob: packet})
	}

	return frames
}

// FileToFrames преобразует файл сразу в фреймы (файл -> страницы -> фреймы)
func FileToFrames(filePath string) []Frame {

	pages := fileToPages(filePath)
	var frames []Frame

	for i, page := range pages {

		frames = append(frames, pageToFrames(uint8(i), page)...)
	}

	return frames
}

// codeLoad кодирует Load данные в зависимости от типа их представления
func codeLoad(load Load) string {

	switch l := load.(type) {
	case V:
		return l.V
	case N:
		hex := strconv.FormatInt(l.Value, 16)
		str := make([]byte, l.Len)
		hexBegin := len(str) - len(hex)

		for i := range str {
			if i < hexBegin {
				str[i] = '0'
			} else {
				str[i] = hex[i-hexBegin]
			}
		}
		return string(str)
	case F:
		return EncodeFrameLoad(l.Frame)
	}

	panic(fmt.Errorf("no matches for the load argument (Given type: %T. Expected type N, V or F)", load))
}

// CodePacket формирует байтовую строку, готовую для отправки
func CodePacket(packet Packet) string {
	var data []byte

	data = append(data, "--"...)
	for _, l := range packet.Load {
		data = append(data, codeLoad(l)...)
	}
	fmt.Println(string(data))
	data = append(data, codeLoad(V{packet.Addr.Val})...)
	fmt.Println(string(data))
	data = append(data, codeLoad(N{int64(packet.Mark), 2})...)
	fmt.Println(string(data))
	data = append(data, packet.Code...)
	fmt.Println(string(data))
	data = append(data, ':')
	fmt.Println(string(data))

	return string(data)
}