package main

import (
	"github.com/hidu/go-speed"
	"sync"
	"time"
)

func main() {
	sp := speed.NewSpeed("test", 5, nil)

	var wg sync.WaitGroup

	for n := 0; n < 1000; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				//call api success num
				sp.Success("call_api", 1)

				//call api fail num
				sp.Fail("call_api", i/2)

				//all key wise "_size" Suffix is file size speed
				sp.Success("write_size", i*i*100)

				time.Sleep(10 * time.Millisecond)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	sp.Stop()
}
