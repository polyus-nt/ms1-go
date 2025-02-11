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
	//fileName := "C:\\Users\\mrxmr\\Boss\\gitFolders\\Polyus_group\\ms-tuc\\LapkiIdePlatformEdition\\stm32G030\\build\\main_btn_lmp_main.bin"
	//fileName := "C:\\Users\\mrxmr\\Boss\\gitFolders\\Polyus_group\\ms-tuc\\LapkiIdePlatformEdition\\stm32G030\\build\\main_mtrx.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\dump_firmware3492523334"

	fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\firmware test\\ms1-fw\\tests\\tjc-ms1-btn-a3\\allowCheck\\blinkOnMsgSerial\\build\\main.bin"
	//fileName := "C:\\Users\\mrxmr\\Downloads\\Work\\firmware test\\ms1-fw\\tests\\tjc-ms1-btn-a3\\allowCheck\\writeSerial\\build\\main.bin"

	fmt.Printf("Started process write firmware to board from file { %v }\n", fileName)
	deviceClone.ActivateLog()
	backTrack := deviceClone.ActivateLog()
	fmt.Print("Start writing firmware:")
	go printer(backTrack)

	firmware, err := deviceClone.WriteFirmware(fileName, false)
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

	// Извлечение прошивки из памяти мк
	if false {

		fmt.Println("Temp dir: ", os.TempDir())
		dumpFirmware, err := os.CreateTemp("", "dump_firmware")
		defer func() {
			_ = dumpFirmware.Close()
			_ = os.Remove(dumpFirmware.Name())
		}()
		fmt.Println("Temp file with dump firmware: ", dumpFirmware.Name())

		backTrack = deviceClone.ActivateLog()
		go printer(backTrack)

		err = deviceClone.GetFirmware(dumpFirmware, 180)
		_ = dumpFirmware.Close()
		if err != nil {
			log.Fatalln(fmt.Errorf("error getting firmware: %v", err))
		}
	}

	// TODO: turn off this branch
	allow, err := device.Allow(true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(allow)

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

func printer(backTrack <-chan ms1.BackTrackMsg) {

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
		case ms1.GET_FIRMWARE:
			msg += "GET_FIRMWARE"
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
}