package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
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

	var res Reply
	var rawData string

	switch response {

	case "pi":
		rawData, _ = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		var ping Ping
		err := presentation.Decoder([]interface{}{&ping.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "pi", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "pi" + rawData
		res = ping
	case "po":
		rawData, _ = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		var pong Pong
		err := presentation.Decoder([]interface{}{&pong.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "po" + rawData
		res = pong
	case "gp":
		rawData, _ = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		var genePong GenePong
		err := presentation.Decoder([]interface{}{&genePong.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "gp" + rawData
		res = genePong
	case "gA":
		rawData, _ = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		var geneAck GeneAck
		err := presentation.Decoder([]interface{}{&geneAck.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "gA" + rawData
		res = geneAck
	case "RF":
		rawData, _ = transport.GetSerialBytes(port, 4+2+16-config.CRC_LENGTH)
		a := rawData[:4]
		m := rawData[4 : 4+2]
		r := rawData[4+2:]
		var ref Ref
		// TODO: fix it (fields and value not correct (size not equal))
		err := presentation.Decoder([]interface{}{&ref.Value}, []presentation.Field{{4, 2, "mark"}, {6, 16, "ref"}}, a+m+r)
		if err != nil {
			return Garbage{Comment: "RF", Garbage: fmt.Sprintf("%v { error: %v }\n", r, err)}
		}
		rawData = "RF" + rawData
		res = ref
	case "OK":
		rawData, _ = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		var ack Ack
		err := presentation.Decoder([]interface{}{&ack.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "OK", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "OK" + rawData
		res = ack
	case "fr":
		rawData, _ = transport.GetSerialBytes(port, 281-config.CRC_LENGTH)
		var frame2 Frame2
		err := presentation.Decoder([]interface{}{&frame2.Mark, &frame2.Page, &frame2.Index, &frame2.Blob}, []presentation.Field{{16, 2, "mark"}, {18, 2, "page"}, {20, 2, "index"}, {22, 256, "blob"}}, rawData)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "fr" + rawData
		res = frame2
	case "ig":
		rawData, _ = transport.GetSerialBytes(port, 36-config.CRC_LENGTH)
		var id ID
		err := presentation.Decoder([]interface{}{&id.Mark, &id.Nanoid}, []presentation.Field{{16, 2, "mark"}, {18, 16, "nanoid"}}, rawData)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "ig" + rawData
		res = id
	case "NO":
		rawData, _ = transport.GetSerialBytes(port, 8-config.CRC_LENGTH)
		var ack Ack
		err := presentation.Decoder([]interface{}{&ack.Value}, []presentation.Field{{4, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "NO", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "NO" + rawData
		res = ack
	case "ER":
		rawData, _ = transport.GetSerialBytes(port, 22-config.CRC_LENGTH)
		var error Error
		err := presentation.Decoder([]interface{}{&error.Mark, &error.Message}, []presentation.Field{{16, 2, "mark"}, {18, 4, "msg"}}, rawData)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "ER" + rawData
		res = error
	case "m1":
		rawData, _ = transport.GetSerialBytes(port, 149-config.CRC_LENGTH)
		var meta Meta
		err := presentation.Decoder(
			[]interface{}{&meta.Mark, &meta.Valid, &meta.RefBlHw, &meta.RefBlFw,
				&meta.RefBlProtocol, &meta.RefBlChip, &meta.RefBlUserCode, &meta.RefCgHw,
				&meta.RefCgFw, &meta.RefCgProtocol},
			[]presentation.Field{{16, 2, "mark"}, {18, 1, "valid"},
				{19 + 16*0, 16, "refBlHw"}, {19 + 16*1, 16, "refBlFw"},
				{19 + 16*2, 16, "refBlProtocol"}, {19 + 16*3, 16, "refBlChip"},
				{19 + 16*4, 16, "refBlUserCode"}, {19 + 16*5, 16, "refCgHw"},
				{19 + 16*6, 16, "refCgFw"}, {19 + 16*7, 16, "refCgProtocol"}},
			rawData)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}
		}
		rawData = "m1" + rawData
		res = meta
	default:
		return Garbage{Comment: "other", Garbage: response}
	}

	// check crc8
	rawData = "." + rawData
	crcGot, _ := transport.GetSerialBytes(port, config.CRC_LENGTH)
	crcCalc := presentation.ToHex(int64(presentation.CalcCRC8([]byte(rawData))), 2)
	if crcGot != crcCalc {
		return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: crc8 for received packet is not correct (got: %v; calc: %v) }\n", rawData, crcGot, crcCalc)}
	}

	return res
}