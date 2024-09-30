package transport

import (
	"fmt"
	"github.com/polyus-nt/ms1-go/internal/config"
)

func Log__(msg string, args ...interface{}) {

	if config.DEBUG__ {

		fmt.Printf("\033[0;33m"+msg+"\033[0m", args...)
		//https://stackoverflow.com/questions/5762491/how-to-print-color-in-console-using-system-out-println
	}
}
