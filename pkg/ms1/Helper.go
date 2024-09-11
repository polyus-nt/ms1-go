package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"github.com/polyus-nt/ms1-go/internal/io/transport"
	"io"
	"time"
)

// GetReply считывает байты с порта и формирует сообщение
func getReply(port io.Reader) Reply {

	timer := time.Now()
	for transport.GetSerialBytes(port, 1)[0] != '.' {
		transport.Wait()
	}
	fmt.Printf("Elapsed1: %v\n", time.Since(timer))

	response := string(transport.GetSerialBytes(port, 2))
	fmt.Printf("Elapsed2: %v\n", time.Since(timer))

	switch response {

	case "pi":
		raw := string(transport.GetSerialBytes(port, 20))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "pi", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Ping{Value: int(data[0])}
	case "po":
		raw := string(transport.GetSerialBytes(port, 20))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Pong{Value: int(data[0])}
	case "gp":
		raw := string(transport.GetSerialBytes(port, 20))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return GenePong{Value: int(data[0])}
	case "gA":
		raw := string(transport.GetSerialBytes(port, 20))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return GeneAck{Value: int(data[0])}
	case "RF":
		a := string(transport.GetSerialBytes(port, 4))
		m := string(transport.GetSerialBytes(port, 2))
		r := string(transport.GetSerialBytes(port, 16))
		data, err := presentation.Decoder([]presentation.Field{{4, 2, "mark"}, {6, 16, "ref"}}, a+m+r)
		if err != nil {
			return Garbage{Comment: "RF", Garbage: fmt.Sprintf("%v { error: %v }\n", r, err)}
		}
		return Ref{Value: data[1]}
	case "OK":
		raw := string(transport.GetSerialBytes(port, 20))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "OK", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Ack{Value: int(data[0])}
	case "fr":
		timer = time.Now()
		raw := string(transport.GetSerialBytes(port, 400))
		fmt.Printf("Elapsed3: %v\n", time.Since(timer))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}, {18, 2, "page"}, {20, 2, "index"}}, raw)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Frame2{Page: int(data[1]), Index: int(data[2]), Mark: int(data[0]), Blob: raw[22:][:256]}
	case "ig":
		raw := string(transport.GetSerialBytes(port, 36))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}, {18, 16, "id"}}, raw)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return ID{Mark: int(data[0]), Nanoid: data[1]}
	case "NO":
		raw := string(transport.GetSerialBytes(port, 8))
		data, err := presentation.Decoder([]presentation.Field{{4, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "NO", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Ack{Value: int(data[0])}
	case "ER":
		raw := string(transport.GetSerialBytes(port, 22))
		data, err := presentation.Decoder([]presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return Error{Mark: int(data[0]), Message: presentation.GetPart(presentation.Field{Start: 18, Len: 4, Descr: "msg"}, raw)}
	default:
		return Garbage{Comment: "other", Garbage: response}
	}
}

func intToID(id int64) string {

	res := make([]byte, 16)

	idStr := fmt.Sprintf("%x", id)

	for i, _ := range res {

		I := len(res) - len(idStr)

		if i < I {
			res[i] = '0'
		} else {
			res[i] = idStr[i-I]
		}
	}

	return string(res)
}