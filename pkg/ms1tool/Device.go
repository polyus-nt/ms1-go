package ms1tool

import (
	"fmt"
	"io"
	"ms1-tool-go/internal/config"
	"ms1-tool-go/internal/io/entity"
	"ms1-tool-go/internal/io/presentation"
)

type Device struct {
	port io.ReadWriter

	addr entity.Address
	mark uint8
	id   map[uint8]string
}

func NewDevice(port io.ReadWriter) *Device {
	return &Device{port, config.ZeroAddress, 0, nil}
}

// Stringer
func (d *Device) String() string {
	return fmt.Sprintf("Device { addr: %v, port: %v }", d.addr, PortName(d.port))
}

func (d *Device) getMark() uint8 {

	res := d.mark
	d.mark++

	return res
}

func (d *Device) ResetPort(port io.ReadWriter) {
	d.port = port
}

func (d *Device) Ping() (res Reply, err error) {

	// create packs
	packs := []presentation.Packet{presentation.PacketPing(d.getMark(), d.addr)}

	// exec
	resT, err := worker(d.port, packs)
	res = resT[0]

	return
}

func (d *Device) GetId(updateID, exitConfMode bool) (res []Reply, err error, updated bool) {

	packs := []presentation.Packet{presentation.PacketGetId(d.getMark(), d.addr)}

	resT, err := worker(d.port, packs)
	if err != nil {
		return
	}
	res = append(res, resT...)

	if exitConfMode {

		packs := []presentation.Packet{presentation.PacketMode(entity.ModeRun, d.addr)}

		resT, err = worker(d.port, packs)
		res = append(res, resT...)
	}

	if updateID {

		if id, ok := res[0].(ID); ok {
			d.addr = entity.Address{Val: intToID(id.Nanoid)}
			updated = true
		}
	}

	return
}

func (d *Device) SetId() (res bool, err error) {

	return
}

func (d *Device) Jump() (res bool, err error) {

	return
}

func (d *Device) ErasePages() (res bool, err error) {

	return
}

// TODO: здесь далее формируются пакеты для работы с устройством
// (очистка страниц флеш-памяти, залить прошивку, осуществить прыжок и т.д.)