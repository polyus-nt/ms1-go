package main

import (
	"fmt"
	"log"
	"ms1-tool-go/pkg/ms1tool"
	"os"
	"sort"
)

func main() {

	fmt.Println("Start serial")

	ports := ms1tool.PortList()
	sort.Strings(ports)

	fmt.Println("Available ports:")
	for i, port := range ports {
		fmt.Printf("%v. %v\n", i+1, port)
	}

	fmt.Print("Choose port (1, 2, 3...): ")
	var usrInput int
	_, err := fmt.Scanf("%d", &usrInput)
	if err != nil {
		_ = fmt.Errorf("Error input for port: %v\n", err)
	}

	port := ms1tool.MkSerial(ports[usrInput-1])
	defer port.Close()

	device := ms1tool.NewDevice(port)
	fmt.Printf("Device created: %v\n", device)

	ping, err := device.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ping)

	id, err, b := device.GetId(true, true)
	if err != nil || b == false {
		log.Fatalf("Error get id { error: %v, isUpdateID: %v}\n", err, b)
	}
	fmt.Printf("Device id updated -> %v\n", id)
	fmt.Println(device)

	ping, err = device.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ping)

	// Процесс прошивки платы
	fileName := "data/usercode-mtrx.bin"
	fmt.Printf("Started process write firmware to board from file { %v }\n", fileName)
	firmware, err := device.WriteFirmware(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(firmware)

	// check reset functions
	os.Stdin.Read(make([]byte, 1))
	device.Reset(true)

	os.Stdin.Read(make([]byte, 1))
	resetTarget, err := device.ResetTarget()
	if err != nil {
		return
	}
	fmt.Println(resetTarget)

	fmt.Println("Finished!")
}