package ms1

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
	"github.com/polyus-nt/ms1-go/internal/io/entity"
	"github.com/polyus-nt/ms1-go/internal/io/presentation"
	"io"
)

type Device struct {
	port io.ReadWriter

	addr entity.Address
	mark uint8
	id   map[uint8]string

	logger chan BackTrackMsg
}

func NewDevice(port io.ReadWriter) *Device {
	return &Device{port, config.ZeroAddress, 0, nil, nil}
}

// Stringer
func (d *Device) String() string {
	return fmt.Sprintf("Device { addr: %v, port: %v }", d.addr, PortName(d.port))
}

func (d *Device) ActivateLog() <-chan BackTrackMsg {

	/*  Если канал != nil, то он точно не закрыт. Внутренний инвариант структуры Device. (исключения close closed не будет, пользователь received channel закрыть не может)*/
	if d.logger != nil {
		close(d.logger)
	}

	// create new channel
	d.logger = make(chan BackTrackMsg, 8)

	return d.logger
}

// SetAddress Обновляет поле адреса (только у объекта, не затрагивая само устройство)
func (d *Device) SetAddress(addr string) (err error) {

	if len(addr) != 16 {
		return fmt.Errorf("address length must be 16")
	}

	d.addr.Val = addr
	return
}

func (d *Device) GetAddress() string {
	return d.addr.Val
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

	packs := []presentation.Packet{presentation.PacketGetId(d.getMark())}

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
			d.addr = entity.Address{Val: id.Nanoid}
			updated = true
		}
	}

	return
}

// SetId Присвоить устройству новый id, при этом id обновляется и у самой платы (отправляется соответсвующий пакет)
func (d *Device) SetId(id string) (res []Reply, err error) {

	if len(id) != 16 {
		return nil, fmt.Errorf("ID not correct! (len(ID) != 16), expected len: %v; id: %v", len(id), id)
	}

	res, err = worker(d.port, []presentation.Packet{presentation.PacketSetId(d.getMark(), id, d.addr)})
	if err == nil {
		d.addr = entity.Address{Val: id}
	}

	return
}

// Разрешить/запретить микроконтроллеру бутлоадера связываться с пользовательской шиной
func (d *Device) Allow(isBlock bool) (res Reply, err error) {

	// create packs
	packs := []presentation.Packet{presentation.PacketAllow(d.getMark(), d.addr, isBlock)}

	// exec
	resT, err := worker(d.port, packs)
	res = resT[0]

	return
}

func (d *Device) WriteFirmware(fileName string, checkFlashFirmware bool) (res []Reply, err error) {

	// ping device
	d.log(BackTrackMsg{UploadStage: PING, CurPack: 1, TotalPacks: 1})
	ping, err := d.Ping()
	res = append(res, ping)
	if err != nil {
		return
	}

	// Открытие прошивки и формирование пакетов
	d.log(BackTrackMsg{UploadStage: PREPARE_FIRMWARE, NoPacks: true})
	packs, err := presentation.File2Frames2Packets(fileName, d.mark, d.addr)
	if err != nil {
		return
	}
	d.mark += uint8(len(packs)) // Shift to mark len(Packets)

	// Перевод в режим программирования
	d.log(BackTrackMsg{UploadStage: CHANGE_MODE_TO_PROG, CurPack: 1, TotalPacks: 1})
	mode, err := d.changeMode(entity.ModeProg)
	res = append(res, mode...)
	if err != nil {
		return
	}

	// ping device
	d.log(BackTrackMsg{UploadStage: PING, CurPack: 1, TotalPacks: 1})
	ping, err = d.Ping()
	res = append(res, ping)
	if err != nil {
		return
	}

	// Очистка страниц
	d.log(BackTrackMsg{UploadStage: ERASE_OLD_FIRMWARE, NoPacks: true})
	pages, err := d.erasePages(len(packs)) // len(packs) must be equal to len(frames)
	res = append(res, pages...)
	if err != nil {
		return
	}

	// Загрузка прошивки
	d.log(BackTrackMsg{UploadStage: PUSH_FIRMWARE, NoPacks: true})
	replies, err := workerBackTrack(d.port, packs, d.log, BackTrackMsg{UploadStage: PUSH_FIRMWARE})
	res = append(res, replies...)
	if err != nil {
		return
	}

	// Проверка целостности загруженной прошивки (опционально)
	if checkFlashFirmware {
		d.log(BackTrackMsg{UploadStage: PULL_FIRMWARE, NoPacks: true})
		replies, err = d.getFrames(len(packs)) // Подтянули записанный код прошивки
		res = append(res, replies...)

		if err != nil {
			err = fmt.Errorf("device::WriteFirmware warning: failed loading frames from flash memory mk (%v)", err)
		} else {
			d.log(BackTrackMsg{UploadStage: VERIFY_FIRMWARE, NoPacks: true})
			var ok bool
			ok, err = d.verifyFirmware(packs, replies)
			if err != nil || !ok {
				err = fmt.Errorf("device::WriteFirmware warning: the firmware is loaded incorrectly (%v)", err)
			}
		}
	}

	// Перевод в режим Run
	d.log(BackTrackMsg{UploadStage: CHANGE_MODE_TO_RUN, CurPack: 1, TotalPacks: 1})
	mode, err2 := d.changeMode(entity.ModeRun)
	res = append(res, mode...)
	if err == nil {
		err = err2
	}

	d.deactivateLogger()

	return
}

func (d *Device) GetFirmware(w io.Writer, qFrames int) (err error) {

	// ping device
	d.log(BackTrackMsg{UploadStage: PING, CurPack: 1, TotalPacks: 1})
	_, err = d.Ping()
	if err != nil {
		return
	}

	// Перевод в режим программирования
	d.log(BackTrackMsg{UploadStage: CHANGE_MODE_TO_PROG, CurPack: 1, TotalPacks: 1})
	_, err = d.changeMode(entity.ModeProg)
	if err != nil {
		return
	}

	// ping device
	d.log(BackTrackMsg{UploadStage: PING, CurPack: 1, TotalPacks: 1})
	_, err = d.Ping()
	if err != nil {
		return
	}

	// pull firmware
	d.log(BackTrackMsg{UploadStage: PULL_FIRMWARE, NoPacks: true})
	frames, err := d.getFrames(qFrames)
	if err != nil {
		return fmt.Errorf("device::getFirmware error: failed loading frames from flash memory mk (%w)", err)
	}

	// convert from frame to bin data and put in w
	d.log(BackTrackMsg{UploadStage: GET_FIRMWARE, NoPacks: true})
	for i, frame := range frames {

		d.log(BackTrackMsg{UploadStage: GET_FIRMWARE, CurPack: uint16(i), TotalPacks: uint16(qFrames)})

		if frame2, ok := frame.(Frame2); ok {

			data, err := presentation.Frame2Bin(frame2.Blob)
			if err != nil {
				return err
			}

			_, err = w.Write(data)
			if err != nil {
				return fmt.Errorf("device::getFirmware error: failed writing frame to writer (%w)", err)
			}
		} else {
			return fmt.Errorf("device::getFirmware error: invalid frame type")
		}
	}

	// Перевод в режим Run
	d.log(BackTrackMsg{UploadStage: CHANGE_MODE_TO_RUN, CurPack: 1, TotalPacks: 1})
	_, err = d.changeMode(entity.ModeRun)
	if err == nil {
		return
	}

	d.deactivateLogger()

	return nil
}

// TRY AND DELETE !? [deprecated]
func (d *Device) GetMetadata2Direct() (res []Reply, err error) {

	res, err = worker(d.port, []entity.Packet{presentation.PacketGetMetadata2Direct(d.mark, d.addr)})

	return
}

func (d *Device) GetMeta() (res Meta, err error) {

	reply, err := worker(d.port, []entity.Packet{presentation.PacketGetMeta(d.mark, d.addr)})

	if err != nil {
		return
	}

	if meta, ok := reply[0].(Meta); ok {
		res = meta
	} else {
		err = fmt.Errorf("device::getMeta warning: the meta is incorrect (%T)", reply[0])
	}

	return
}

// ChangeModeToConf - Переключить устройство в режим конфигурации
func (d *Device) ChangeModeToConf() (res []Reply, err error) {

	res, err = worker(d.port, []entity.Packet{presentation.PacketMode(d.getMark(), entity.ModeConf, d.addr)})

	return
}

// ChangeModeToRun - Переключить устройство в режим Run
func (d *Device) ChangeModeToRun() (res []Reply, err error) {

	res, err = worker(d.port, []entity.Packet{presentation.PacketMode(d.getMark(), entity.ModeRun, d.addr)})

	return
}

// ChangeModeToProg - Переключить устройство в режим программирования
func (d *Device) ChangeModeToProg() (res []Reply, err error) {

	res, err = worker(d.port, []entity.Packet{presentation.PacketMode(d.getMark(), entity.ModeProg, d.addr)})

	return
}

// Reset - сброс кибергена. Это действие влечет за собой сброс bootloader-a
func (d *Device) Reset(resetMark bool) {

	packs := []presentation.Packet{presentation.PacketResetSelf(d.addr)}

	workerNoReply(d.port, packs)

	if resetMark {
		d.mark = 0
	}
}

// ResetTarget - сброс bootloader-a (пользовательский микроконтроллер)
func (d *Device) ResetTarget() (res []Reply, err error) {

	packs := []presentation.Packet{presentation.PacketResetTarget(d.getMark(), d.addr)}

	return worker(d.port, packs)
}

// Далее служебные функции

// changeMode Посылает пакет для переключения режима на кибергене
func (d *Device) changeMode(mode entity.Mode) (res []Reply, err error) {

	packs := []presentation.Packet{presentation.PacketMode(d.getMark(), mode, d.addr)}

	res, err = worker(d.port, packs)

	return
}

// erasePages - очищает нужное количество страниц flash памяти для будущей прошивки
func (d *Device) erasePages(lenFrames int) (res []Reply, err error) {

	qPages := (lenFrames*config.SIZE_FRAME + config.SIZE_PAGE - 1) / config.SIZE_PAGE

	var packs []presentation.Packet
	for i := 0; i < qPages; i++ {
		packs = append(packs, presentation.PacketNuke(int64(i), d.getMark(), d.addr))
	}

	res, err = workerBackTrack(d.port, packs, d.log, BackTrackMsg{UploadStage: ERASE_OLD_FIRMWARE})

	return
}

func (d *Device) getFrames(lenFrames int) (res []Reply, err error) {

	IPage := func(i int) (res int64) {
		res = int64(i) * config.SIZE_FRAME / config.SIZE_PAGE
		return
	}
	IFrame := func(i int) (res int64) {
		qPages := int(IPage(i))
		res = int64((i*config.SIZE_FRAME - qPages*config.SIZE_PAGE) / config.SIZE_FRAME)
		return
	}

	var packs []presentation.Packet

	for i := 0; i < lenFrames; i++ {
		packs = append(packs, presentation.PacketTargetFrame(d.getMark(),
			IPage(i), IFrame(i), d.addr))
	}

	res, err = workerBackTrack(d.port, packs, d.log, BackTrackMsg{UploadStage: PULL_FIRMWARE})

	return
}

func (d *Device) verifyFirmware(expected []presentation.Packet, suspected []Reply) (ok bool, err error) {

	if len(expected) != len(suspected) {
		return false, fmt.Errorf("error in 'Device::verifyFirmware': len(expected) != len(suspected) { len(expected) = %v; len(suspected) = %v }", len(expected), len(suspected))
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

		encoded := presentation.EncodeFrameLoad(expectedFrame.Frame)
		encoded = encoded[:len(encoded)-4]

		if encoded != suspectedFrame.Blob {
			return false, fmt.Errorf("frame[%v] data not match", i)
		}
	}

	return true, nil
}

func (d *Device) loggerIsActive() bool {

	return d.logger != nil
}

func (d *Device) log(msg BackTrackMsg) {

	if d.loggerIsActive() {
		// write only if channel has place for data else skip
		select {
		case d.logger <- msg:
		default:
		}
	}
}

func (d *Device) deactivateLogger() {

	if d.loggerIsActive() {
		close(d.logger)
	}

	d.logger = nil
}

func (d *Device) jump() (res bool, err error) {

	return
}