package main

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/pkg/ms1"
	"log"
	"os"
)

func main() {

	buffer := make([]byte, 10)
	base := buffer
	fmt.Printf("base; len: %v, cap: %v\n", len(base), cap(base))
	fmt.Printf("buffer; len: %v, cap: %v\n", len(buffer), cap(buffer))

	buffer = buffer[2:]
	fmt.Printf("base; len: %v, cap: %v\n", len(base), cap(base))
	fmt.Printf("buffer; len: %v, cap: %v\n", len(buffer), cap(buffer))

	buffer = buffer[0:]
	fmt.Printf("base; len: %v, cap: %v\n", len(base), cap(base))
	fmt.Printf("buffer; len: %v, cap: %v\n", len(buffer), cap(buffer))

	fmt.Println("Start serial")

	ports := ms1.PortList()

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

	port, err := ms1.MkSerial(ports[usrInput-1])
	if err != nil {
		log.Fatalln(err)
	}
	defer port.Close()

	device := ms1.NewDevice(port)
	fmt.Printf("Device created: %v\n", device)
	fmt.Printf("Device address: %v\n", device.GetAddress())

	ping, err := device.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ping)

	id, err, updated := device.GetId(true, true)
	if err != nil || updated == false {
		log.Fatalf("Error get id { error: %v, isIDUpdated: %v}\n", err, updated)
	}
	fmt.Printf("Device id updated -> %v\n", id)
	fmt.Println(device)
	fmt.Printf("Device address: %v\n", device.GetAddress())

	// Можно создавать новый объект девайса при каждом обращении к устройству (начале сессии)
	deviceClone := ms1.NewDevice(port)
	err = deviceClone.SetAddress(device.GetAddress())
	if err != nil {
		log.Fatalln(fmt.Errorf("error setting address: %v", err))
	}

	ping, err = device.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ping)

	// Процесс прошивки платы
	fileName := "data\\mtrx.bin"
	fmt.Printf("Started process write firmware to board from file { %v }\n", fileName)
	firmware, err := deviceClone.WriteFirmware(fileName, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(firmware)

	// check reset functions
	fmt.Print("Enter for reset device")
	_, _ = os.Stdin.Read(make([]byte, 1))
	device.Reset(true)

	fmt.Print("Enter for reset target in device")
	_, _ = os.Stdin.Read(make([]byte, 1))
	resetTarget, err := device.ResetTarget()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resetTarget)

	fmt.Println("Finished!")
}