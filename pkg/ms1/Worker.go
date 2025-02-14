package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"github.com/polyus-nt/ms1-go/internal/io/transport"
	"io"
	"time"
)

var ResetBuffers func()

func worker(port io.ReadWriter, packets []presentation.Packet) (res []Reply, err error) {

	return workerBackTrack(port, packets, nil, BackTrackMsg{})
}

func workerBackTrack(port io.ReadWriter, packets []presentation.Packet, logger func(msg BackTrackMsg), msg BackTrackMsg) (res []Reply, err error) {

	if logger != nil {
		msg.CurPack = 0
		msg.TotalPacks = uint16(len(packets))
	}

	for _, packet := range packets {

		if logger != nil {
			msg.CurPack++
			logger(msg)
		}

		var reply Reply
		transport.Log__("Transmit start\n")
		for i := 0; i < config.ATTEMPTS_QUANTITY; i++ {

			transport.Log__("Transmit: attempt %d\n", i+1)

			transport.PutMessage(port, packet)

			time.Sleep(777 * time.Millisecond)

			reply, err = getReply(port)
			if err == nil {
				break
			}
			changeTiming()
			ResetBuffers()
		}
		restoreTiming()

		if err != nil {
			transport.Log__("Transmit timeout error!\n")
			return []Reply{Error{0, fmt.Sprintf("%s", err)}}, err
		}
		transport.Log__("Transmit finished\n")

		res = append(res, reply)

		switch r := reply.(type) {
		case Error, Garbage:
			err = fmt.Errorf("worker interrupted: get bad packet: %v\n", r.String())
			return []Reply{r}, err
		}
	}
	return
}

func workerNoReply(port io.ReadWriter, packets []presentation.Packet) {

	for _, packet := range packets {

		transport.PutMessage(port, packet)
	}
	return
}

func changeTiming() {

	config.SERIAL_DEADLINE += config.DELTA_WAITING
}

func restoreTiming() {

	config.SERIAL_DEADLINE = config.SERIAL_DEADLINE_DEFAULT
}