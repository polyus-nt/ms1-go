package ms1

import (
	"fmt"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"io"
	"log"
	"time"
)

func MkSerial(portName string) io.ReadWriteCloser {

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatalln(err)
	}
	err = port.SetReadTimeout(75 * time.Millisecond)
	if err != nil {
		log.Fatalln(err)
	}

	return port
}

func PortList() (res []string) {

	ports, _ := enumerator.GetDetailedPortsList()

	for _, port := range ports {

		res = append(res, port.Name)
	}
	return
}

func PortName(port interface{}) (res string) {

	if port, ok := port.(serial.Port); ok {
		res = fmt.Sprintf("%v", port)
	} else {
		res = "unknown"
	}

	return
}