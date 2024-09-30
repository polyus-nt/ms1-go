package transport

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"io"
	"os"
	"time"
)

// PutMessage отправляет байтовую строку (сообщение) в порт
func PutMessage(port io.Writer, packet presentation.Packet) {

	var code = presentation.CodePacket(packet)

	fmt.Printf("Msg -> %v; size = %v\n", code, len(code))
	bytesSended, err := port.Write([]byte(code))
	fmt.Printf("BytesSended -> %v; err = %v\n", bytesSended, err)
	time.Sleep(config.SERIAL_SEND_WAITING) // time.Sleep for signal to OS scheduler
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Put message error: %v\n", err)
		return
	}
	//fmt.Printf("Serial write %v bytes\n", write)
}

// GetSerialBytes считывает требуемое количество байт с порта
func GetSerialBytes(port io.Reader, count int) (string, error) {

	fmt.Println("started get bytes...")
	buffer := make([]byte, count)
	ready := 0
	bArr := buffer // base array of slice
	deadline := time.Now().Add(config.SERIAL_DEADLINE)
	fmt.Printf("Buffer before loop: buffer=%v; len = %v\n", buffer, len(buffer))
	//fmt.Println(buffer[0])

	for {
		fmt.Printf("GetSerialBytes: ready=%d deadline=%v buffer= %v\n", ready, deadline, buffer)
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
		time.Sleep(config.SERIAL_READ_WAITING) // time.Sleep for signal to OS scheduler
	}
	fmt.Printf("RETURN FROM GetSerialBytes: ready=%d deadline=%v buffer=%v\n", ready, deadline, buffer)

	return string(bArr[:ready]), nil
}