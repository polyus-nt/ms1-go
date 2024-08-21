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
	value int
	len   int
}

type F struct {
	frame Frame
}

func codeLoad(load Load) string {

	switch l := load.(type) {
	case V:
		return l.v
	case N:
		hex := strconv.FormatInt(int64(l.value), 16)
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
		return EncodeFrameLoad(l.frame)
	}

	panic(fmt.Errorf("no matches for the load argument (Given type: %T. Expected type N, V or F)", load))
}

type Packet struct {
	mark uint8
	addr Address
	code string
	load []Load
}

type Mode int

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
	data = append(data, codeLoad(packet.load)...)
	data = append(data, codeLoad(V{packet.addr.val})...)
	data = append(data, codeLoad(N{int(packet.mark), 2})...)
	data = append(data, packet.code...)

	return string(data)
}

func PacketPing(m uint8) Packet {

	return Packet{mark: m, addr: TestAddress, code: "pi", load: append([]Load{}, V{""})}
}

func PacketPong() Packet {

	return Packet{mark: 0, addr: TestAddress, code: "po", load: append([]Load{}, V{""})}
}

func PacketNuke(i int, m uint8) Packet {

	return Packet{mark: m, addr: TestAddress, code: "np", load: append([]Load{}, N{i, 2})}
}

func PacketJump() Packet {

	return Packet{mark: 0, addr: TestAddress, code: "jp", load: append([]Load{}, V{""})}
}

func PacketResetSelf() Packet {

	return Packet{mark: 0, addr: TestAddress, code: "rs", load: append([]Load{}, V{""})}
}

func PacketResetTarget(m uint8) Packet {

	return Packet{mark: m, addr: TestAddress, code: "rt", load: append([]Load{}, V{""})}
}

func PacketPingTarget(m uint8) Packet {

	return Packet{mark: m, addr: TestAddress, code: "gp", load: append([]Load{}, V{""})}
}

func PacketMakeJump() Packet {

	return Packet{mark: 0, addr: TestAddress, code: "jp", load: append([]Load{}, V{""})}
}

func PacketTargetRef() Packet {

	return Packet{mark: 10, addr: TestAddress, code: "RF", load: append([]Load{}, V{""})}
}

func PacketTargetFrame(m uint8, p, i int) Packet {

	return Packet{mark: m, addr: TestAddress, code: "rf", load: append([]Load{}, N{p, 2}, N{i, 2})}
}

func PacketMode(mode Mode) Packet {

	return Packet{mark: 0, addr: TestAddress, code: "st", load: append([]Load{}, N{int(mode), 2})}
}

func PacketSetId(m uint8, w int) Packet {

	return Packet{mark: m, addr: ZeroAddress, code: "id", load: append([]Load{}, N{w, 16})}
}

func PacketGetId(m uint8) Packet {

	return Packet{mark: m, addr: ZeroAddress, code: "ig", load: []Load{}}
}

func GetPart(field Field, s string) string {

	remains := min(len(s)-field.start, field.len)
	return s[field.start:][:remains]
}

func GetHex(field Field, s string) (int, error) {

	blob := GetPart(field, s)
	num, err := strconv.ParseInt(blob, 16, 32)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getHex [%v] failed! {error: %v}", field.descr, err)
	}

	return int(num), err
}

func GetSignedHex(field Field, s string) (int, error) {

	blob := GetPart(Field{start: field.start + 1, len: field.len - 1}, s)

	num, err := strconv.ParseInt(blob, 16, 32)

	if err != nil {
		fmt.Fprintf(os.Stderr, "getSignedHex [%v] failed! {error: %v}", field.descr, err)
	}

	if GetPart(Field{len: 1}, s)[0] == '-' {
		num = -num
	}

	return int(num), err
}

func Decoder(fields []Field, s string) ([]int, error) {

	var res []int

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