package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"github.com/polyus-nt/ms1-go/internal/io/transport"
	"io"
)

// GetReply считывает байты с порта и формирует сообщение
func getReply(port io.Reader) Reply {

	for {
		reply, err := transport.GetSerialBytes(port, 1)

		if len(reply) > 0 {
			if reply[0] == '.' {
				break
			}
			// else we are receive garbage and continue for next msg
			// transport.Wait()
		} else {
			return Error{0, fmt.Sprintf("Don't reply from serial { %v}", err)} // not reply...
		}
	}

	response, err := transport.GetSerialBytes(port, 2)
	if len(response) != 2 {
		return Error{0, fmt.Sprintf("Don't receive packet code from serial { %v}", err)}
	}

	switch response {

	case "pi":
		raw, _ := transport.GetSerialBytes(port, 20)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "pi", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Ping{Value: int(data[0])}
	case "po":
		raw, _ := transport.GetSerialBytes(port, 20)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Pong{Value: int(data[0])}
	case "gp":
		raw, _ := transport.GetSerialBytes(port, 20)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return GenePong{Value: int(data[0])}
	case "gA":
		raw, _ := transport.GetSerialBytes(port, 20)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return GeneAck{Value: int(data[0])}
	case "RF":
		a, _ := transport.GetSerialBytes(port, 4)
		m, _ := transport.GetSerialBytes(port, 2)
		r, _ := transport.GetSerialBytes(port, 16)
		data, err := presentation.Decoder([]presentation.Field{{4, 2, "mark"}, {6, 16, "ref"}}, a+m+r)
		if err != nil {
			return Garbage{Comment: "RF", Garbage: fmt.Sprintf("%v { error: %v }\n", r, err)}
		}
		return Ref{Value: data[1]}
	case "OK":
		raw, _ := transport.GetSerialBytes(port, 20)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "OK", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Ack{Value: int(data[0])}
	case "fr":
		raw, _ := transport.GetSerialBytes(port, 400)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}, {18, 2, "page"}, {20, 2, "index"}}, raw)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Frame2{Page: int(data[1]), Index: int(data[2]), Mark: int(data[0]), Blob: raw[22:][:256]}
	case "ig":
		raw, _ := transport.GetSerialBytes(port, 36)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return ID{Mark: int(data[0]), Nanoid: raw[18 : 18+16]}
	case "NO":
		raw, _ := transport.GetSerialBytes(port, 8)
		data, err := presentation.Decoder([]presentation.Field{{4, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "NO", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Ack{Value: int(data[0])}
	case "ER":
		raw, _ := transport.GetSerialBytes(port, 22)
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Error{Mark: int(data[0]), Message: presentation.GetPart(presentation.Field{Start: 18, Len: 4, Descr: "msg"}, raw)}
	default:
		return Garbage{Comment: "other", Garbage: response}
	}
}