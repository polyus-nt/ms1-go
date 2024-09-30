package config

import (
	"github.com/polyus-nt/ms1-go/internal/io/entity"
	"time"
)

//goland:noinspection ALL
const (
	SIZE_PAGE  = 2048
	SIZE_FRAME = 128

	SERIAL_DEADLINE     = 50 * time.Millisecond
	SERIAL_SEND_WAITING = 6 * time.Millisecond
	SERIAL_READ_WAITING = 8 * time.Millisecond
)

var ZeroAddress = entity.Address{Val: "0000000000000000"}