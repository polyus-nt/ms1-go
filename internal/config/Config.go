package config

import (
	"github.com/polyus-nt/ms1-go/internal/io/entity"
	"time"
)

//goland:noinspection ALL
const (
	SIZE_PAGE  = 2048
	SIZE_FRAME = 128

	SERIAL_DEADLINE_DEFAULT     = 888 * time.Millisecond
	SERIAL_SEND_WAITING_DEFAULT = 10 * time.Millisecond
	SERIAL_READ_WAITING_DEFAULT = 25 * time.Millisecond

	ATTEMPTS_QUANTITY = 3

	CRC_LENGTH = 2

	DEBUG__ = true
)

var SERIAL_DEADLINE = 888 * time.Millisecond
var SERIAL_SEND_WAITING = 10 * time.Millisecond
var SERIAL_READ_WAITING = 25 * time.Millisecond
var DELTA_WAITING = 444 * time.Millisecond

var ZeroAddress = entity.Address{Val: "0000000000000000"}