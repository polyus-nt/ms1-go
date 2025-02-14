package transport

import (
	"errors"
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"go.bug.st/serial"
	"io"
	"os"
	"time"
)

var TimeoutError = errors.New("timeout expired")

// PutMessage отправляет байтовую строку (сообщение) в порт
func PutMessage(port io.Writer, packet presentation.Packet) {

	var code = presentation.CodePacket(packet)

	Log__("Msg -> %v; size = %v\n", code, len(code))

	bytesSent, err := port.Write([]byte(code))

	if port, ok := port.(serial.Port); ok {
		Log__("Waiting sent bytes")
		err = port.Drain()
		Log__("Wait finished; err -> %v\n", err)
	}

	Log__("Bytes sent -> %v; err = %v\n", bytesSent, err)

	time.Sleep(config.SERIAL_SEND_WAITING) // time.Sleep for signal to OS scheduler

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Put message error: %v\n", err)
		return
	}
}

// GetSerialBytes считывает требуемое количество байт с порта
func GetSerialBytes(port io.Reader, count int) (string, error) {

	Log__("started get bytes...\n")
	buffer := make([]byte, count)
	ready := 0
	bArr := buffer // base array of slice
	deadline := time.Now().Add(config.SERIAL_DEADLINE)
	Log__("Buffer before loop: buffer=%v; len = %v\n", buffer, len(buffer))
	isTimeout := false

	for {
		Log__("GetSerialBytes: ready=%d deadline=%v buffer= %v\n", ready, time.Since(deadline), buffer)
		qBytes, err := port.Read(buffer)
		if err != nil {
			return "", err
		}
		ready += qBytes
		if ready >= count {
			break
		}
		if time.Now().After(deadline) {
			isTimeout = true
			break
		}
		buffer = buffer[qBytes:]
		time.Sleep(config.SERIAL_READ_WAITING) // time.Sleep for signal to OS scheduler
	}
	Log__("GetSerialBytes Final: ready=%d deadline=%v buffer= %v\n", ready, time.Since(deadline), buffer)

	if isTimeout {
		return string(bArr[:ready]), fmt.Errorf("%w: expected: %d, received: %d", TimeoutError, count, ready)
	}

	return string(bArr[:ready]), nil
}