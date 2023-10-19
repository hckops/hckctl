package util

import (
	"os"
	"os/signal"
	"syscall"
)

func InterruptHandler(callback func()) {
	// captures CTRL+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		callback()
	}()
}
