package config

import (
	"ms1-tool-go/internal/io/entity"
	"time"
)

const (
	SIZE_PAGE          = 2048
	SIZE_FRAME         = 128
	SERIAL_WAITING     = 5 * time.Millisecond
	BOOTLOADER_WAITING = 10 * time.Millisecond
)

var ZeroAddress = entity.Address{Val: "0000000000000000"}