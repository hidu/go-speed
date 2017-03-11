# go-speed
Go Speed Statistics


example:
```go
package main

import (
    "github.com/hidu/go-speed"
    "time"
    "sync"
)

func main() {
    sp := speed.NewSpeed("test", 5, nil)
    
    var wg sync.WaitGroup
    
    for n:=0;n<1000;n++{
        wg.Add(1)
        go func(){
            for i := 0; i < 1000; i++ {
                //call api success num
                sp.Success("call_api", 1)
                
                //call api fail num
                sp.Fail("call_api", i/2)
                
                //all key wise "_size" Suffix is file size speed
                sp.Success("write_size", i*i*100)
                
                time.Sleep(10*time.Millisecond)
            }
            wg.Done()
        }()
    }
    wg.Wait()
    sp.Stop()
}

```

outputLog:
```
2017/03/11 22:21:57 speed [id=test] [sec=5] \
[call_api_all:(total:250500000,22829520.6/s),(suc:1000000,91135.8/s,0.4%),(fail:249500000,22738384.8/s,99.6%)]; \
[write_size_all:(total:30.27T,2.76T/s),(suc:30.27T,2.76T/s,100.0%),(fail:0.0B,0.0B/s,0.0%)]; \
[call_api_1s:(total:40854520,42014789.2/s),(suc:85258,87679.3/s,0.2%),(fail:40769262,41927109.9/s,99.8%)]; \
[write_size_1s:(total:7.10T,7.31T/s),(suc:7.10T,7.31T/s,100.0%),(fail:0.0B,0.0B/s,0.0%)]

```
