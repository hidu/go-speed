package main

import (
	"github.com/hidu/go-speed"
	"time"
)

func main() {
	sp := speed.NewSpeed("test", 2, nil)
	for i := 0; i < 1000; i++ {
		sp.Success("deal", 1)
		time.Sleep(100 * time.Millisecond)
	}
	sp.Stop()

}
