package transport

import (
	"github.com/polyus-nt/ms1-go/internal/config"
	"time"
)

func Wait() {
	time.Sleep(config.BOOTLOADER_WAITING)
}