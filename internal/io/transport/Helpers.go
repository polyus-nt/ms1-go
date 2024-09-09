package transport

import (
	"ms1-tool-go/internal/config"
	"time"
)

func Wait() {
	time.Sleep(config.BOOTLOADER_WAITING)
}