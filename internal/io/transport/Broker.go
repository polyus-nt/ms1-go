package transport

import (
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"io"
	"time"
)

// PutMessage отправляет байтовую строку (сообщение) в порт
func PutMessage(port io.Writer, packet presentation.Packet) {

	var code = presentation.CodePacket(packet)

	//fmt.Printf("Msg -> %v\n", code)
	_, err := port.Write([]byte(code))
	if err != nil {
		return
	}
	//fmt.Printf("Serial write %v bytes\n", write)
}

// GetSerialBytes считывает требуемое количество байт с порта
func GetSerialBytes(port io.Reader, count int) (string, error) {

	buffer := make([]byte, count)
	ready := 0
	bArr := buffer // base array of slice
	deadline := time.Now().Add(config.SERIAL_WAITING)

	for {
		qBytes, err := port.Read(buffer)
		if err != nil {
			return "", err
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

	return string(bArr[:ready]), nil
}