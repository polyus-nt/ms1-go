package internal

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"time"
)

const (
	_Com string = "/dev/ttyACM0"
)

// general type (for enum impl)
type Reply interface{}

// derrived types for Reply
type Ping struct{ value int }
type GenePong struct{ value int }
type GeneAck struct{ value int }
type Pong struct{ value int }
type Ack struct{ value int }
type Nack struct{ value int }
type Ref struct{ value int }
type ID struct {
	mark   int
	nanoid int
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

	c := &serial.Config{Name: _Com, Baud: 115200, Parity: serial.ParityNone, StopBits: serial.Stop1, ReadTimeout: time.Second}

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalln(err)
	}

	return s
}

func PutMessage(serial *serial.Port, packet Packet) {

	var code = CodePacket(packet)

	fmt.Printf("Msg -> %v\n", code)
	serial.Write([]byte(code))
}

// Считывает требуемое количество байт с порта
func getSerialBytes(serial *serial.Port, count int) []byte {

	buffer := make([]byte, count)
	ready := 0
	bArr := buffer

	for {
		qBytes, err := serial.Read(buffer)
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

func GetReply(serial *serial.Port) Reply {

	for getSerialBytes(serial, 1)[0] != '.' {
		wait()
	}

	response := string(getSerialBytes(serial, 2))
	switch response {
	case "pi":
		raw := string(getSerialBytes(serial, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"pi", raw}
		}
		return Ping{value: data[0]}
	case "po":
		raw := string(getSerialBytes(serial, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"po", raw}
		}
		return Pong{value: data[0]}
	case "gp":
		raw := string(getSerialBytes(serial, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"po", raw}
		}
		return GenePong{value: data[0]}
	case "gA":
		raw := string(getSerialBytes(serial, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"po", raw}
		}
		return GeneAck{value: data[0]}
	case "RF":
		a := string(getSerialBytes(serial, 4))
		m := string(getSerialBytes(serial, 2))
		r := string(getSerialBytes(serial, 16))
		data, err := Decoder(append([]Field{}, Field{4, 2, "mark"}, Field{6, 16, "ref"}), a+m+r)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"RF", r}
		}
		fmt.Println(r)
		return Ref{value: data[1]}
	case "OK":
		raw := string(getSerialBytes(serial, 20))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			return Garbage{"OK", raw}
		}
		return Ack{value: data[0]}
	case "fr":
		raw := string(getSerialBytes(serial, 400))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}, Field{18, 2, "page"}, Field{20, 2, "index"}), raw)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return Garbage{"fr", raw}
		}
		fmt.Println(raw)
		return Frame2{page: data[1], index: data[2], mark: data[0], blob: raw[22:][:256]}
	case "ig":
		raw := string(getSerialBytes(serial, 36))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}, Field{18, 16, "id"}), raw)
		if err != nil {
			return Garbage{"fr", raw}
		}
		fmt.Println(raw)
		return ID{mark: data[0], nanoid: data[1]}
	case "NO":
		raw := string(getSerialBytes(serial, 8))
		data, err := Decoder(append([]Field{}, Field{4, 2, "mark"}), raw)
		if err != nil {
			return Garbage{"NO", raw}
		}
		return Ack{value: data[0]}
	case "ER":
		raw := string(getSerialBytes(serial, 22))
		data, err := Decoder(append([]Field{}, Field{16, 2, "mark"}), raw)
		if err != nil {
			return Garbage{"ERR", raw}
		}
		return Error{mark: data[0], message: GetPart(Field{18, 4, "msg"}, raw)}
	default:
		return Garbage{"other", response}
	}
}

func _worker(ans bool, serial *serial.Port, packets []Packet) {

	for _, packet := range packets {

		PutMessage(serial, packet)

		if ans {

			reply := GetReply(serial)

			switch r := reply.(type) {
			case Pong:
				fmt.Printf("Got pong [%#xd]\n", r.value)
				wait()
			case Ping:
				fmt.Printf("Got ping [%#xd]\n", r.value)
			case GenePong:
				fmt.Printf("Got pong from gene [%#xd]\n", r.value)
			case GeneAck:
				fmt.Printf("Got Ack from gene [%#xd]\n", r.value)
			case Garbage:
				fmt.Printf("Got Garbege %v : %v\n", r.comment, r.garbage)
				wait()
			case Ack:
				fmt.Printf("Got Ack [%#xd]\n", r.value)
				wait()
			case Frame2:
				fmt.Printf("Frame [ %#xd ] %v.%v\n", r.mark, r.page, r.index)
				var bin []Bin
				for i := 0; i < len(r.blob); i += 16 {
					rI := min(i+16, len(r.blob))
					bin = append(bin, Bin(r.blob[i:rI]))
				}
				Xxd(bin)
				wait()
			case Ref:
				fmt.Printf("Got ref [%#xd]\n", r.value)
				wait()
			case Nack:
				fmt.Printf("Got Ack [%#xd]\n", r.value)
				wait()
			case ID:
				fmt.Printf("Got id [%#xd]: %#xd\n", r.mark, r.nanoid)
				wait()
			case Error:
				fmt.Printf("Got error [%#xd] \"%v\"\n", r.mark, r.message)
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

func Worker(serial *serial.Port, packets []Packet) {
	measured(func() {
		_worker(true, serial, packets)
	})
}

func WorkerFin(serial *serial.Port, packets []Packet) {
	measured(func() {
		_worker(false, serial, packets)
	})
}