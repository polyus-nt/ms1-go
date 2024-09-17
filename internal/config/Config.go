package config

import (
	"github.com/polyus-nt/ms1-go/internal/io/entity"
	"time"
)

//goland:noinspection ALL
const (
	SIZE_PAGE          = 2048
	SIZE_FRAME         = 128
	SERIAL_WAITING     = 5 * time.Millisecond
	BOOTLOADER_WAITING = 10 * time.Millisecond
)

var ZeroAddress = entity.Address{Val: "0000000000000000"}