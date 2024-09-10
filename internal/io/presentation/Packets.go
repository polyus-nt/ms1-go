package presentation

import (
	"ms1-tool-go/internal/config"
	"ms1-tool-go/internal/io/entity"
)

func PacketPing(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "pi", Load: []Load{V{}}}
}

func PacketPong(addr Address) Packet {

	return Packet{Mark: 0, Addr: addr, Code: "po", Load: []Load{V{}}}
}

func PacketNuke(i int64, m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "np", Load: []Load{N{i, 2}}}
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

	return Packet{Mark: m, Addr: addr, Code: "rf", Load: append([]Load{}, N{p, 2}, N{i, 2})}
}

func PacketMode(m uint8, mode Mode, addr Address) Packet {

	return Packet{Mark: m, Addr: addr, Code: "st", Load: append([]Load{}, N{int64(mode), 2})}
}

func PacketSetId(m uint8, w int64, addr Address) Packet {

	return Packet{Mark: m, Addr: config.ZeroAddress, Code: "id", Load: append([]Load{}, N{w, 16})}
}

func PacketGetId(m uint8, addr Address) Packet {

	return Packet{Mark: m, Addr: config.ZeroAddress, Code: "ig", Load: []Load{}}
}

func File2Frames2Packets(filePath string, startMark uint8, addr Address) (packets []Packet) {

	frames := FileToFrames(filePath)

	for i, frame := range frames {
		packets = append(packets, Packet{Mark: startMark + uint8(i), Addr: addr, Code: "fr", Load: []entity.Load{F{Frame: frame}}})
	}

	return
}