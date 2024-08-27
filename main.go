package main

import (
	"fmt"
	"go.bug.st/serial/enumerator"
	"ms1-tool-go/internal"
	"os"
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

	//internal.Test1()

	fmt.Println("Start serial")

	ports, _ := enumerator.GetDetailedPortsList()
	fmt.Println(ports)
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
			fmt.Printf("   USB product %s\n", port.Product)
		}
	}

	port := internal.MkSerial()

	setup := true
	if setup {
		internal.Worker(port, []internal.Packet{
			internal.PacketGetId(7),
			//internal.PacketSetId(5, 13286090),
		})
		internal.TestAddress = internal.LastAdrress
	}
	(*port).Close()

	fmt.Printf("Address -> %v\n", internal.TestAddress)
	fmt.Println("Ready. Reset board and press enter for continue...")
	b := make([]byte, 1)
	os.Stdin.Read(b)

	port = internal.MkSerial()
	defer (*port).Close()

	frames := internal.FileToFrames("data/usercode-mtrx.bin")
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
	for i := 0; i <= 16; i++ {
		peekFrames = append(peekFrames, internal.PacketTargetFrame(uint8(50+i), 0, int64(i)))
	}
	//packs = append(packs, peekFrames...)
	packs = append(packs, internal.PacketMode(internal.ModeRun))

	internal.Worker(port, packs)
	internal.WorkerFin(port, []internal.Packet{internal.PacketResetSelf()})
}