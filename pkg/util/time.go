package util

import (
	"time"
)

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
