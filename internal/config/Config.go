package config

import (
	"github.com/polyus-nt/ms1-go/internal/io/entity"
	"time"
)

//goland:noinspection ALL
const (
	SIZE_PAGE  = 2048
	SIZE_FRAME = 128

	SERIAL_DEADLINE     = 888 * time.Millisecond
	SERIAL_SEND_WAITING = 10 * time.Millisecond
	SERIAL_READ_WAITING = 25 * time.Millisecond

	DEBUG__ = false
)

var ZeroAddress = entity.Address{Val: "0000000000000000"}
