package main

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/pkg/ms1"
	"log"
	"os"
)

func main() {

	fmt.Println("Start serial")

	//data := []byte(".dr668e739880610dc1320000000000000800409d0800409d0800409d0800409d0800409d0800409d000000000800409d0800409d0800409d0800409d0800409d0800409d000000000800409d08003b8d0800409d00000000000000000800409d000000000000000000000000000000000000000000000000000000000800409d0800409d0800404d20002000")
	//fmt.Println(crc8.Checksum(data, crc8.MakeTable(crc8.CRC8_CDMA2000)))

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

	_, err = device.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ping)

	// Процесс прошивки платы
	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\UART\\UARTModules\\DataBus\\build\\main1DataBus.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\repo\\stm32\\ms-tuc\\firmwares\\portingLapkiIDE\\UARTModules\\DataBus\\build\\main2DataBus.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\UART\\UARTModules\\SimpleBus\\build\\main1SimpleBus.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\UART\\UARTModules\\SimpleBus\\build\\main2SimpleBus.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\repo\\stm32\\ms-tuc\\buildFiles\\mainUART1ListenOnly.bin"
	//fileName := "C:\\Users\\mrxmr\\OneDrive\\Документы\\Cache\\TGCache\\sketch.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\repo\\stm32\\ms-tuc\\firmwares\\portingLapkiIDE\\mtrx\\build\\mtrx.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\sketch\\mtrx\\sketch\\build\\mtrxSketch.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\mtrx\\build\\mtrx.bin"
	fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\UARTModules_PlatformEdition_2\\build\\main1DataBus.bin"

	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\UARTModules\\UARTModules\\SimpleBus\\build\\main3_lmp_SimpleBus.bin"
	fmt.Printf("Started process write firmware to board from file { %v }\n", fileName)
	deviceClone.ActivateLog()
	backTrack := deviceClone.ActivateLog()
	go func() {

		fmt.Print("Start writing firmware:")
		var lastRecordType ms1.UploadStage = 255
		var noPackRecord *ms1.BackTrackMsg = nil
		var noPackMsg string

		for record := range backTrack {

			// parse type of msg
			var msg string
			if lastRecordType != record.UploadStage {
				msg += "\n"
				lastRecordType = record.UploadStage

				if noPackRecord != nil {
					fmt.Print("\r" + noPackMsg + ": done")
				}
			} else {
				msg += "\r"
			}

			switch record.UploadStage {
			case ms1.PING:
				msg += "PING"
			case ms1.CHANGE_MODE_TO_PROG:
				msg += "CHANGE_MODE_TO_PROG"
			case ms1.PREPARE_FIRMWARE:
				msg += "PREPARE_FIRMWARE"
			case ms1.ERASE_OLD_FIRMWARE:
				msg += "ERASE_OLD_FIRMWARE"
			case ms1.PUSH_FIRMWARE:
				msg += "PUSH_FIRMWARE"
			case ms1.PULL_FIRMWARE:
				msg += "PULL_FIRMWARE"
			case ms1.VERIFY_FIRMWARE:
				msg += "VERIFY_FIRMWARE"
			case ms1.CHANGE_MODE_TO_RUN:
				msg += "CHANGE_MODE_TO_RUN"
			default:
				msg += "SOME ACTION"
			}

			// fill progress
			// if packet has num of packets
			if !record.NoPacks {
				noPackRecord = nil
				msg += fmt.Sprintf(": %v/%v", record.CurPack, record.TotalPacks)
			} else {
				noPackRecord = &record
				noPackMsg = msg[1:]
				msg += ": ..."
			}

			fmt.Print(msg)
		}
		fmt.Println("\n--- FINISHED STATUS GOROUTINE ---")
	}()

	firmware, err := deviceClone.WriteFirmware(fileName, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(firmware)

	// Получение мета информации об устройстве
	meta, err := deviceClone.GetMeta()
	if err != nil {
		log.Fatalln(fmt.Errorf("error getting meta info: %v", err))
	}
	fmt.Println(meta)

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