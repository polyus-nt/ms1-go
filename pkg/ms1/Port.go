package ms1

import (
	"fmt"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"io"
	"slices"
	"time"
)

func MkSerial(portName string) (io.ReadWriteCloser, error) {

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	err = port.SetReadTimeout(88 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	return port, nil
}

func PortList() (res []string) {

	ports, _ := enumerator.GetDetailedPortsList()

	for _, port := range ports {

		res = append(res, port.Name)
	}
	slices.Sort(res)

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