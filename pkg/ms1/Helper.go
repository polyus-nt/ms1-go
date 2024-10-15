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
		var res Ping
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "pi", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "po":
		raw, _ := transport.GetSerialBytes(port, 20)
		var res Pong
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "gp":
		raw, _ := transport.GetSerialBytes(port, 20)
		var res GenePong
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "gA":
		raw, _ := transport.GetSerialBytes(port, 20)
		var res GeneAck
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "RF":
		a, _ := transport.GetSerialBytes(port, 4)
		m, _ := transport.GetSerialBytes(port, 2)
		r, _ := transport.GetSerialBytes(port, 16)
		var res Ref
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{4, 2, "mark"}, {6, 16, "ref"}}, a+m+r)
		if err != nil {
			return Garbage{Comment: "RF", Garbage: fmt.Sprintf("%v { error: %v }\n", r, err)}
		}
		return res
	case "OK":
		raw, _ := transport.GetSerialBytes(port, 20)
		var res Ack
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{16, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "OK", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "fr":
		raw, _ := transport.GetSerialBytes(port, 281)
		var res Frame2
		err := presentation.Decoder([]interface{}{&res.Mark, &res.Page, &res.Index, &res.Blob}, []presentation.Field{{16, 2, "mark"}, {18, 2, "page"}, {20, 2, "index"}, {22, 256, "blob"}}, raw)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "ig":
		raw, _ := transport.GetSerialBytes(port, 36)
		var res ID
		err := presentation.Decoder([]interface{}{&res.Mark, &res.Nanoid}, []presentation.Field{{16, 2, "mark"}, {18, 16, "nanoid"}}, raw)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "NO":
		raw, _ := transport.GetSerialBytes(port, 8)
		var res Ack
		err := presentation.Decoder([]interface{}{&res.Value}, []presentation.Field{{4, 2, "mark"}}, raw)
		if err != nil {
			return Garbage{Comment: "NO", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "ER":
		raw, _ := transport.GetSerialBytes(port, 22)
		var res Error
		err := presentation.Decoder([]interface{}{&res.Mark, &res.Message}, []presentation.Field{{16, 2, "mark"}, {18, 4, "msg"}}, raw)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	case "m1":
		raw, _ := transport.GetSerialBytes(port, 150)
		var res Meta
		err := presentation.Decoder(
			[]interface{}{&res.Mark, &res.Valid, &res.RefBlHw, &res.RefBlFw,
				&res.RefBlProtocol, &res.RefBlChip, &res.RefBlUserCode, &res.RefCgHw,
				&res.RefCgFw, &res.RefCgProtocol},
			[]presentation.Field{{16, 2, "mark"}, {18, 1, "valid"},
				{19 + 16*0, 16, "refBlHw"}, {19 + 16*1, 16, "refBlFw"},
				{19 + 16*2, 16, "refBlProtocol"}, {19 + 16*3, 16, "refBlChip"},
				{19 + 16*4, 16, "refBlUserCode"}, {19 + 16*5, 16, "refCgHw"},
				{19 + 16*6, 16, "refCgFw"}, {19 + 16*7, 16, "refCgProtocol"}},
			raw)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", raw, err)}
		}
		return res
	default:
		return Garbage{Comment: "other", Garbage: response}
	}
}
