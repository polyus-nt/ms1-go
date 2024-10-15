package presentation

import (
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/entity"
)

func PacketPing(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "pi", Load: []Load{V{}}}
}

func PacketPong(addr Address) Packet {

	return Packet{Mark: 0, Addr: addr, Code: "po", Load: []Load{V{}}}
}

func PacketNuke(i int64, m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "np", Load: []Load{N{Value: i, Len: 2}}}
}

func PacketJump(addr Address) Packet {

	return Packet{Mark: 0, Addr: addr, Code: "jp", Load: []Load{V{}}}
}

func PacketResetSelf(addr Address) Packet {

	return Packet{Mark: 0, Addr: addr, Code: "rs", Load: []Load{V{}}}
}

func PacketResetTarget(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "rt", Load: []Load{V{}}}
}

func PacketPingTarget(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "gp", Load: []Load{V{}}}
}

func PacketMakeJump(addr Address) Packet {

	return Packet{Mark: 0, Addr: addr, Code: "jp", Load: []Load{V{}}}
}

func PacketTargetRef(addr Address) Packet {

	return Packet{Mark: 10, Addr: addr, Code: "RF", Load: []Load{V{}}}
}

func PacketTargetFrame(m uint8, p, i int64, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "rf", Load: append([]Load{}, N{Value: p, Len: 2}, N{Value: i, Len: 2})}
}

func PacketMode(m uint8, mode Mode, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "st", Load: append([]Load{}, N{Value: int64(mode), Len: 2})}
}

func PacketSetId(m uint8, newID string, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "id", Load: append([]Load{}, V{V: newID})}
}

func PacketGetId(m uint8) Packet {

	return Packet{Mark: m, Addr: config.ZeroAddress, Code: "ig", Load: []Load{}}
}

func File2Frames2Packets(filePath string, startMark uint8, addr Address) (packets []Packet, err error) {

	frames, err := FileToFrames(filePath)
	if err != nil {
		return
	}

	for i, frame := range frames {
		packets = append(packets, Packet{Mark: startMark + uint8(i), Addr: addr, Code: "fr", Load: []entity.Load{F{Frame: frame}}})
	}

	return
}

func PacketGetMeta(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "m1", Load: []Load{}}
}

func PacketGetMetadata2Direct(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "md", Load: []Load{}}
}
