package utils

import (
	"fmt"
	"time"
)

func Log(message string) {
	var now = time.Now()
	var nowFormatted = now.Format("2006/01/02 15:04:05")
	fmt.Println("["+nowFormatted+"] " + message)
}