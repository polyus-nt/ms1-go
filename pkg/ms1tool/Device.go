package ms1tool

import (
	"fmt"
	"io"
	"ms1-tool-go/internal/config"
	"ms1-tool-go/internal/io/entity"
	"ms1-tool-go/internal/io/presentation"
	"strconv"
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

// Здесь далее формируется API для работы с устройством
// (очистка страниц флеш-памяти, залить прошивку, осуществить прыжок и т.д.)

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

		packs := []presentation.Packet{presentation.PacketMode(d.getMark(), entity.ModeRun, d.addr)}

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

func (d *Device) SetId(id string) (res []Reply, err error) {

	if len(id) < 16 {
		return nil, fmt.Errorf("ID not correct! (len(ID) != 16), expected len: %v; id: %v", len(id), id)
	}

	idInt, err := strconv.ParseInt(id, 16, 64)
	if err != nil {
		return
	}

	res, err = worker(d.port, []presentation.Packet{presentation.PacketSetId(d.getMark(), idInt, d.addr)})
	if err == nil {
		d.addr = entity.Address{Val: id}
	}

	return
}

func (d *Device) WriteFirmware(fileName string) (res []Reply, err error) {

	// ping device
	ping, err := d.Ping()
	res = append(res, ping)
	if err != nil {
		return
	}

	// Открытие прошивки и формирование пакетов
	packs := presentation.File2Frames2Packets(fileName, d.mark, d.addr)
	d.mark += uint8(len(packs)) // Shift to mark len(Packets)

	//temp
	worker(d.port, []presentation.Packet{presentation.PacketMode(d.getMark(), entity.ModeConf, d.addr)})

	// temp
	worker(d.port, []presentation.Packet{presentation.PacketMode(d.getMark(), entity.ModeProg, config.ZeroAddress)})

	// Перевод в режим программирования
	mode, err := d.changeMode(entity.ModeProg) // TODO: maybe additional modeConf or modeProg for target
	res = append(res, mode...)
	if err != nil {
		return
	}

	//temp
	worker(d.port, []presentation.Packet{presentation.PacketResetTarget(d.getMark(), d.addr)})
	worker(d.port, []presentation.Packet{presentation.PacketPingTarget(d.getMark(), d.addr)})

	// ping device
	ping, err = d.Ping()
	fmt.Println(ping)
	res = append(res, ping)
	if err != nil {
		return
	}

	// Очистка страниц
	pages, err := d.erasePages(len(packs)) // TODO: Совпадает ли количество фреймов с количеством пакетов?
	res = append(res, pages...)
	if err != nil {
		return
	}

	// Загрузка прошивки
	replies, err := worker(d.port, packs)
	res = append(res, replies...)
	if err != nil {
		return
	}

	// Проверка целостности загруженной прошивки
	suspectFrames, err := d.getFrame(len(packs)) // Подтянули записанный код прошивки
	if err != nil {
		return
	}
	ok, err := d.verifyFirmware(packs, suspectFrames)
	if err != nil || !ok {
		_ = fmt.Errorf("the firmware is loaded incorrectly (%v)", err)
		return
	}

	// Перевод в режим Run
	mode, err = d.changeMode(entity.ModeRun)
	res = append(res, mode...)

	return
}

// Далее служебные функции

// changeMode Посылает пакет для переключения режима на кибергене
func (d *Device) changeMode(mode entity.Mode) (res []Reply, err error) {

	packs := []presentation.Packet{presentation.PacketMode(d.getMark(), mode, d.addr)}

	res, err = worker(d.port, packs)

	return
}

func (d *Device) Reset() {

}

// erasePages - очищает нужное количество страниц flash памяти для будущей прошивки
func (d *Device) erasePages(lenFrames int) (res []Reply, err error) {

	qPages := (lenFrames*config.SIZE_FRAME + config.SIZE_PAGE - 1) / config.SIZE_PAGE

	var packs []presentation.Packet
	for i := 0; i < qPages; i++ {
		packs = append(packs, presentation.PacketNuke(int64(i), d.getMark(), d.addr))
	}

	res, err = worker(d.port, packs)

	return
}

func (d *Device) getFrame(lenFrames int) (res []Reply, err error) {

	IPage := func(i int) (res int64) {
		res = int64(i) * config.SIZE_FRAME / config.SIZE_PAGE
		return
	}
	IFrame := func(i int) (res int64) {
		wholePart := i / config.SIZE_PAGE
		i -= config.SIZE_PAGE * wholePart
		res = int64(i) % config.SIZE_FRAME
		return
	}

	var packs []presentation.Packet

	for i := 0; i < lenFrames; i++ {
		packs = append(packs, presentation.PacketTargetFrame(d.getMark(),
			IPage(i), IFrame(i), d.addr))
	}

	res, err = worker(d.port, packs)

	return
}

func (d *Device) verifyFirmware(expected []presentation.Packet, suspected []Reply) (ok bool, err error) {

	if len(expected) != len(suspected) {
		return false, fmt.Errorf("error in 'verifyFirmware': len(expected) != len(suspected) { len(expected) = %v; len(suspected) = %v }", len(expected), len(suspected))
	}

	for i := 0; i < len(expected); i++ {
		expectedFrame, ok := expected[i].Load[0].(entity.F)
		if !ok {
			return false, fmt.Errorf("error read expectedFrame[%v] -> type mismath (expected: F; received: %T)", i, expected[i].Load[0])
		}
		suspectedFrame, ok := suspected[i].(Frame2)
		if !ok {
			return false, fmt.Errorf("error read suspectFrame[%v] -> type mismath (expected: Frame2; received: %T)", i, suspected[i])
		}
		if expectedFrame.Frame.Blob != suspectedFrame.Blob {
			return false, fmt.Errorf("frame[%v] data not match", i)
		}
	}

	return true, nil
}

func (d *Device) jump() (res bool, err error) {

	return
}