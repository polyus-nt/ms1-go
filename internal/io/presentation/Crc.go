package presentation

import (
	"github.com/sigurn/crc8"
)

var table = crc8.MakeTable(crc8.CRC8_CDMA2000)

func CalcCRC8(data []byte) uint8 {

	return crc8.Checksum(data, table)
}