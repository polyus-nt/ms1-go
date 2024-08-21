package main

import (
	"fmt"
	"go.bug.st/serial"
	"ms1-tool-go/internal"
)

var start1 = []internal.Packet{
	internal.PacketNuke(0, 30),
	internal.PacketNuke(1, 31),
	internal.PacketNuke(2, 32),
	internal.PacketNuke(3, 33),
	internal.PacketNuke(4, 34),
	internal.PacketNuke(5, 35),
}

var jump = []internal.Packet{
	internal.PacketPing(99),
	internal.PacketJump(),
}

func main() {

	for _, v := range internal.TestData {
		internal.PrintOneChunk(v)
		fmt.Println()
	}
	internal.Xxd(internal.TestData)

	fmt.Println("Start serial")

	ports, _ := serial.GetPortsList()
	fmt.Println(ports)
	serial := internal.MkSerial()
	defer serial.Close()

	frames := internal.FileToFrames("data/fast_blink_main.bin")
	frames = frames[:min(len(frames), 6*16)]

	var packets []internal.Packet
	for i, frame := range frames {
		packets = append(packets, internal.Packet{Mark: uint8(i), Addr: internal.TestAddress, Code: "fr", Load: []internal.Load{internal.F{Frame: frame}}})
	}

	packs := []internal.Packet{
		internal.PacketPing(12),
		internal.PacketMode(internal.ModeConf),
		internal.PacketGetId(10),
		internal.PacketMode(internal.ModeProg),
		internal.PacketMode(internal.ModeProg),
		internal.PacketResetTarget(13),
		internal.PacketPingTarget(14),
	}
	packs[3].Addr = internal.ZeroAddress

	packs = append(packs, start1...)
	packs = append(packs, packets...)

	var peekFrames []internal.Packet
	for i := 0; i < 16; i++ {
		peekFrames = append(peekFrames, internal.PacketTargetFrame(uint8(50+i), 0, int64(i)))
	}
	packs = append(packs, peekFrames...)
	packs = append(packs, internal.PacketMode(internal.ModeRun))

	internal.Worker(serial, packs)
	internal.WorkerFin(serial, []internal.Packet{internal.PacketResetSelf()})
}