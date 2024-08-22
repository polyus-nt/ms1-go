package internal

import (
	"fmt"
	"go.bug.st/serial"
	"log"
	"os"
	"time"
)

const (
	_Com string = "COM6"
)

// general type (for enum impl)
type Reply interface{}

// derived types for Reply
type Ping struct{ value int }
type GenePong struct{ value int }
type GeneAck struct{ value int }
type Pong struct{ value int }
type Ack struct{ value int }
type Nack struct{ value int }
type Ref struct{ value int64 }
type ID struct {
	mark   int
	nanoid int64
}
type Frame2 struct {
	page  int
	index int
	mark  int
	blob  string
}
type Garbage struct {
	comment string
	garbage string
}
type Error struct {
	mark    int
	message string
}

func wait() {
	time.Sleep(3 * time.Millisecond)
}

func MkSerial() *serial.Port {

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open("COM6", mode)
	if err != nil {
		log.Fatalln(err)
	}
	err = port.SetReadTimeout(time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	return &port
}

func PutMessage(port *serial.Port, packet Packet) {

	var code = CodePacket(packet)

	fmt.Printf("Msg -> %v\n", code)
	write, err := (*port).Write([]byte(code))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Serial write %v bytes\n", write)
}

// Считывает требуемое количество байт с порта
func getSerialBytes(port *serial.Port, count int) []byte {

	buffer := make([]byte, count)
	ready := 0
	bArr := buffer

	for {
		qBytes, err := (*port).Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		ready += qBytes
		if ready >= count {
			break
		}
		buffer = buffer[qBytes:]
	}

	return bArr
}

func GetReply(port *serial.Port) Reply {

	for getSerialBytes(port, 1)[0] != '.' {
		wait()
	}

	response := string(getSerialBytes(port, 2))
	switch response {
	case "pi":
		raw := string(getSerialBytes(port, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"pi", raw}
		}
		return Ping{value: int(data[0])}
	case "po":
		raw := string(getSerialBytes(port, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"po", raw}
		}
		return Pong{value: int(data[0])}
	case "gp":
		raw := string(getSerialBytes(port, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"po", raw}
		}
		return GenePong{value: int(data[0])}
	case "gA":
		raw := string(getSerialBytes(port, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"po", raw}
		}
		return GeneAck{value: int(data[0])}
	case "RF":
		a := string(getSerialBytes(port, 4))
		m := string(getSerialBytes(port, 2))
		r := string(getSerialBytes(port, 16))
		data, err := Decoder(append([]Field{}, Field{4, 2, "mark"}, Field{6, 16, "ref"}), a+m+r)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"RF", r}
		}
		fmt.Println(r)
		return Ref{value: data[1]}
	case "OK":
		raw := string(getSerialBytes(port, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			return Garbage{"OK", raw}
		}
		return Ack{value: int(data[0])}
	case "fr":
		raw := string(getSerialBytes(port, 400))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}, Field{18, 2, "page"}, Field{20, 2, "index"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"fr", raw}
		}
		fmt.Println(raw)
		return Frame2{page: int(data[1]), index: int(data[2]), mark: int(data[0]), blob: raw[22:][:256]}
	case "ig":
		raw := string(getSerialBytes(port, 36))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}, Field{18, 16, "id"}), raw)
		if err != nil {
			return Garbage{"fr", raw}
		}
		fmt.Println(raw)
		return ID{mark: int(data[0]), nanoid: data[1]}
	case "NO":
		raw := string(getSerialBytes(port, 8))
		data, err := Decoder(append([]Field{}, Field{4, 2, "mark"}), raw)
		if err != nil {
			return Garbage{"NO", raw}
		}
		return Ack{value: int(data[0])}
	case "ER":
		raw := string(getSerialBytes(port, 22))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			return Garbage{"ERR", raw}
		}
		return Error{mark: int(data[0]), message: GetPart(Field{18, 4, "msg"}, raw)}
	default:
		return Garbage{"other", response}
	}
}

func _worker(ans bool, port *serial.Port, packets []Packet) {

	for _, packet := range packets {

		PutMessage(port, packet)

		if ans {

			reply := GetReply(port)

			switch r := reply.(type) {
			case Pong:
				fmt.Printf("Got pong [%#x]\n", r.value)
				wait()
			case Ping:
				fmt.Printf("Got ping [%#x]\n", r.value)
			case GenePong:
				fmt.Printf("Got pong from gene [%#x]\n", r.value)
			case GeneAck:
				fmt.Printf("Got Ack from gene [%#x]\n", r.value)
			case Garbage:
				fmt.Printf("Got Garbege %v : %v\n", r.comment, r.garbage)
				wait()
			case Ack:
				fmt.Printf("Got Ack [%#x]\n", r.value)
				wait()
			case Frame2:
				fmt.Printf("Frame [ %#x ] %v.%v\n", r.mark, r.page, r.index)
				var bin []Bin
				for i := 0; i < len(r.blob); i += 16 {
					rI := min(i+16, len(r.blob))
					bin = append(bin, Bin(r.blob[i:rI]))
				}
				Xxd(bin)
				wait()
			case Ref:
				fmt.Printf("Got ref [%#x]\n", r.value)
				wait()
			case Nack:
				fmt.Printf("Got Ack [%#x]\n", r.value)
				wait()
			case ID:
				fmt.Printf("Got id [%#x]: %#x\n", r.mark, r.nanoid)
				wait()
			case Error:
				fmt.Printf("Got error [%#x] \"%v\"\n", r.mark, r.message)
				wait()
			}
		} else {
			fmt.Println("Skipping answer!")
		}
	}

	fmt.Println("Work done!")
}

func measured(action func()) {
	start := time.Now()
	action()
	end := time.Since(start)
	fmt.Printf("Worker executed %#v ", end)
}

func Worker(port *serial.Port, packets []Packet) {
	measured(func() {
		_worker(true, port, packets)
	})
}

func WorkerFin(port *serial.Port, packets []Packet) {
	measured(func() {
		_worker(false, port, packets)
	})
}