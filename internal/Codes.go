package internal

import (
	"fmt"
	"os"
	"strconv"
)

// for enum
// general type
type Load interface{}

// derived types (3)
type V struct {
	v string
}

type N struct {
	value int64
	len   int
}

type F struct {
	Frame Frame
}

func codeLoad(load Load) string {

	switch l := load.(type) {
	case V:
		return l.v
	case N:
		hex := strconv.FormatInt(l.value, 16)
		str := make([]byte, l.len)
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

type Packet struct {
	Mark uint8
	Addr Address
	Code string
	Load []Load
}

type Mode int64

const (
	ModeBlank Mode = iota
	ModeRun
	ModeProg
	ModeConf
	ModeErr
)

type Field struct {
	start int
	len   int
	descr string
}

func CodePacket(packet Packet) string {
	var data []byte

	data = append(data, "--"...)
	for _, l := range packet.Load {
		data = append(data, codeLoad(l)...)
	}
	data = append(data, codeLoad(V{packet.Addr.val})...)
	data = append(data, codeLoad(N{int64(packet.Mark), 2})...)
	data = append(data, packet.Code...)
	data = append(data, ':')

	return string(data)
}

func PacketPing(m uint8) Packet {

	return Packet{Mark: m, Addr: TestAddress, Code: "pi", Load: append([]Load{}, V{""})}
}

func PacketPong() Packet {

	return Packet{Mark: 0, Addr: TestAddress, Code: "po", Load: append([]Load{}, V{""})}
}

func PacketNuke(i int64, m uint8) Packet {

	return Packet{Mark: m, Addr: TestAddress, Code: "np", Load: append([]Load{}, N{i, 2})}
}

func PacketJump() Packet {

	return Packet{Mark: 0, Addr: TestAddress, Code: "jp", Load: append([]Load{}, V{""})}
}

func PacketResetSelf() Packet {

	return Packet{Mark: 0, Addr: TestAddress, Code: "rs", Load: append([]Load{}, V{""})}
}

func PacketResetTarget(m uint8) Packet {

	return Packet{Mark: m, Addr: TestAddress, Code: "rt", Load: append([]Load{}, V{""})}
}

func PacketPingTarget(m uint8) Packet {

	return Packet{Mark: m, Addr: TestAddress, Code: "gp", Load: append([]Load{}, V{""})}
}

func PacketMakeJump() Packet {

	return Packet{Mark: 0, Addr: TestAddress, Code: "jp", Load: append([]Load{}, V{""})}
}

func PacketTargetRef() Packet {

	return Packet{Mark: 10, Addr: TestAddress, Code: "RF", Load: append([]Load{}, V{""})}
}

func PacketTargetFrame(m uint8, p, i int64) Packet {

	return Packet{Mark: m, Addr: TestAddress, Code: "rf", Load: append([]Load{}, N{p, 2}, N{i, 2})}
}

func PacketMode(mode Mode) Packet {

	return Packet{Mark: 0, Addr: TestAddress, Code: "st", Load: append([]Load{}, N{int64(mode), 2})}
}

func PacketSetId(m uint8, w int64) Packet {

	return Packet{Mark: m, Addr: ZeroAddress, Code: "id", Load: append([]Load{}, N{w, 16})}
}

func PacketGetId(m uint8) Packet {

	return Packet{Mark: m, Addr: ZeroAddress, Code: "ig", Load: []Load{}}
}

func GetPart(field Field, s string) string {

	remains := min(len(s)-field.start, field.len)
	return s[field.start:][:remains]
}

func GetHex(field Field, s string) (int64, error) {

	blob := GetPart(field, s)
	num, err := strconv.ParseInt(blob, 16, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getHex [%v] failed! {error: %v}", field.descr, err)
	}

	return num, err
}

func GetSignedHex(field Field, s string) (int64, error) {

	blob := GetPart(Field{start: field.start + 1, len: field.len - 1}, s)

	num, err := strconv.ParseInt(blob, 16, 64)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getSignedHex [%v] failed! {error: %v}", field.descr, err)
	}

	if GetPart(Field{len: 1}, s)[0] == '-' {
		num = -num
	}

	return num, err
}

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