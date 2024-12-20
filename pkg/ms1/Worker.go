package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"github.com/polyus-nt/ms1-go/internal/io/transport"
	"io"
)

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

		transport.PutMessage(port, packet)

		reply := getReply(port)
		// TODO: impl recover (for interrupt serial problems, how it resolve and how update state device?)

		res = append(res, reply)

		switch r := reply.(type) {
		case Error, Garbage:
			err = fmt.Errorf("worker interrupted: %v\n", r.String())
			return
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