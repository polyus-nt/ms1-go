package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"github.com/polyus-nt/ms1-go/internal/io/transport"
	"io"
)

// GetReply считывает байты с порта и формирует сообщение
func getReply(port io.Reader) (Reply, error) {

	for {
		reply, err := transport.GetSerialBytes(port, 1)
		if err != nil { // not reply...
			return nil, err
		}

		if reply[0] == '.' {
			break
		}
	}

	response, err := transport.GetSerialBytes(port, 2)
	if err != nil {
		return nil, err
	}

	var res Reply
	var rawData string

	switch response {

	case "pi":
		rawData, err = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var ping Ping
		err := presentation.Decoder([]interface{}{&ping.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "pi", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "pi" + rawData
		res = ping
	case "po":
		rawData, err = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var pong Pong
		err := presentation.Decoder([]interface{}{&pong.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "po" + rawData
		res = pong
	case "gp":
		rawData, err = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var genePong GenePong
		err := presentation.Decoder([]interface{}{&genePong.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "gp" + rawData
		res = genePong
	case "gA":
		rawData, err = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var geneAck GeneAck
		err := presentation.Decoder([]interface{}{&geneAck.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "po", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "gA" + rawData
		res = geneAck
	case "OK":
		rawData, err = transport.GetSerialBytes(port, 20-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var ack Ack
		err := presentation.Decoder([]interface{}{&ack.Value}, []presentation.Field{{16, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "OK", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "OK" + rawData
		res = ack
	case "fr":
		rawData, err = transport.GetSerialBytes(port, 281-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var frame2 Frame2
		err := presentation.Decoder([]interface{}{&frame2.Mark, &frame2.Page, &frame2.Index, &frame2.Blob}, []presentation.Field{{16, 2, "mark"}, {18, 2, "page"}, {20, 2, "index"}, {22, 256, "blob"}}, rawData)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "fr" + rawData
		res = frame2
	case "ig":
		rawData, err = transport.GetSerialBytes(port, 36-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var id ID
		err := presentation.Decoder([]interface{}{&id.Mark, &id.Nanoid}, []presentation.Field{{16, 2, "mark"}, {18, 16, "nanoid"}}, rawData)
		if err != nil {
			return Garbage{Comment: "fr", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "ig" + rawData
		res = id
	case "NO":
		rawData, err = transport.GetSerialBytes(port, 8-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var ack Ack
		err := presentation.Decoder([]interface{}{&ack.Value}, []presentation.Field{{4, 2, "mark"}}, rawData)
		if err != nil {
			return Garbage{Comment: "NO", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "NO" + rawData
		res = ack
	case "ER":
		rawData, err = transport.GetSerialBytes(port, 22-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
		var error1 Error
		err := presentation.Decoder([]interface{}{&error1.Mark, &error1.Message}, []presentation.Field{{16, 2, "mark"}, {18, 4, "msg"}}, rawData)
		if err != nil {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "ER" + rawData
		res = error1
	case "m1":
		rawData, err = transport.GetSerialBytes(port, 149-config.CRC_LENGTH)
		if err != nil {
			return nil, err
		}
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
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: %v }\n", rawData, err)}, nil
		}
		rawData = "m1" + rawData
		res = meta
	default:
		return Garbage{Comment: "other", Garbage: response}, nil
	}

	// check crc8
	rawData = "." + rawData
	crcGot, err := transport.GetSerialBytes(port, config.CRC_LENGTH)
	if err != nil {
		return nil, err
	}
	crcCalc := presentation.ToHex(int64(presentation.CalcCRC8([]byte(rawData))), 2)
	if crcGot != crcCalc {
		// for packet 'fr' crc8 calc is not correct (skip this packet)
		if rawData[1:3] != "fr" {
			return Garbage{Comment: "ER", Garbage: fmt.Sprintf("%v { error: crc8 for received packet is not correct (got: %v; calc: %v) }\n", rawData, crcGot, crcCalc)}, nil
		}
	}

	return res, nil
}