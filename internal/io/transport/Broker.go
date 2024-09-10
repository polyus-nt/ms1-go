package transport

import (
	"fmt"
	"io"
	"log"
	"ms1-tool-go/internal/config"
	"ms1-tool-go/internal/io/presentation"
	"time"
)

// PutMessage отправляет байтовую строку (сообщение) в порт
func PutMessage(port io.Writer, packet presentation.Packet) {

	var code = presentation.CodePacket(packet)

	fmt.Printf("Msg -> %v\n", code)
	write, err := port.Write([]byte(code))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Serial write %v bytes\n", write)
}

// GetSerialBytes считывает требуемое количество байт с порта
func GetSerialBytes(port io.Reader, count int) []byte {

	buffer := make([]byte, count)
	ready := 0
	bArr := buffer // base array of slice
	deadline := time.Now().Add(config.SERIAL_WAITING)

	for {
		qBytes, err := port.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		ready += qBytes
		if ready >= count {
			break
		}
		if time.Now().After(deadline) {
			break
		}
		buffer = buffer[qBytes:]
	}

	return bArr[:ready]
}