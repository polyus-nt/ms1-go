package presentation

import (
	"bytes"
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"os"
	"slices"
)

// EncodeFrameLoad формирует на основе фрейма данные для отправки (байтовая строка)
func EncodeFrameLoad(frame Frame) string {

	var buf bytes.Buffer

	buf.WriteString(frame.Blob)
	slices.Reverse(buf.Bytes())

	encoded := enc(buf.Bytes())

	buf.Reset()
	buf.WriteString(encoded)

	_, _ = fmt.Fprintf(&buf, "%s", enc([]byte{frame.Page}))
	_, _ = fmt.Fprintf(&buf, "%s", enc([]byte{frame.Part}))

	return buf.String()
}

// EncodeFrame TODO разобраться с ее телом [deprecated --crc32]
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
func fileToPages(filePath string) (res []string, err error) {

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	res = chopBs(config.SIZE_PAGE, string(content))

	return
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
func FileToFrames(filePath string) (frames []Frame, err error) {

	pages, err := fileToPages(filePath)
	if err != nil {
		return
	}

	for i, page := range pages {

		frames = append(frames, pageToFrames(uint8(i), page)...)
	}

	return
}

// codeLoad кодирует Load данные в зависимости от типа их представления
func codeLoad(load Load) string {

	switch l := load.(type) {
	case V:
		return l.V
	case N:
		return ToHex(l.Value, l.Len)
	case F:
		return EncodeFrameLoad(l.Frame)
	}

	panic(fmt.Errorf("no matches for the load argument (Given type: %T. Expected type N, V or F)", load))
}

// CodePacket формирует байтовую строку, готовую для отправки
func CodePacket(packet Packet) string {

	var data []byte

	data = append(data, "--"...) // crc8 (заполняем ниже, когда записали данные)
	for _, l := range packet.Load {
		data = append(data, codeLoad(l)...)
	}
	data = append(data, codeLoad(V{V: packet.Addr.Val})...)
	data = append(data, codeLoad(N{Value: int64(packet.Mark), Len: 2})...)
	data = append(data, packet.Code...)
	data = append(data, ':')

	copy(data, codeLoad(N{Value: int64(CalcCRC8(data[2:])), Len: 2}))

	return string(data)
}