package speed

import "fmt"
import "time"
import "sync/atomic"

type Speed struct {
	Name           string
	Sec            int64
	LastTime       time.Time
	Total          int64 //总数量
	TotalSec      int64 //n秒内总数量
	TotalSuccess   int64 //成功总数量
	TotalSecSuccess      int64 //n秒内成功总数量
	TotalSecSize int64 //n秒内总尺寸
	TotalSize     int64 //总尺寸
	PrintFn        func(string,*Speed)
	ticker         *time.Ticker
}

func NewSpeed(name string, sec int, PrintFn func(string,*Speed)) *Speed {
	sp := new(Speed)
	sp.Name = name
	sp.Sec = int64(sec)
	sp.LastTime = time.Now()
	sp.PrintFn = PrintFn
	sp.ticker = time.NewTicker(time.Duration(sec) * time.Second)
	go func() {
		for now := range sp.ticker.C {
			sp.printLog(now)
		}
	}()
	return sp
}
func (sp *Speed) Inc(num, size,suc_num int) {
	
	atomic.AddInt64(&sp.Total, int64(num))
	atomic.AddInt64(&sp.TotalSec, int64(num))
	atomic.AddInt64(&sp.TotalSecSize, int64(size))
	atomic.AddInt64(&sp.TotalSize, int64(size))

	atomic.AddInt64(&sp.TotalSuccess, int64(suc_num))
	atomic.AddInt64(&sp.TotalSecSuccess, int64(suc_num))
}

func (sp *Speed) Stop() {
	sp.printLog(time.Now())
	sp.ticker.Stop()
}
var msize=int64(1024*1024)
var gsize=msize*1024;

func (sp *Speed) printLog(now time.Time) {
	now_used := now.Sub(sp.LastTime).Seconds()
	sp.LastTime = now
	if now_used == 0 || sp.TotalSec == 0 {
		return
	}
	log_format := "total=%d,qps=%d,total_%ds=%d,speed=%.2fMps,total_size=%s,total_suc=%d,"+fmt.Sprintf("total_%ds_suc",int64(now_used))+"=%d"
	size_speed := float64(sp.TotalSecSize) / now_used / (1024 * 1024)

	TotalSize_str := ""
	if sp.TotalSize > gsize {
		TotalSize_str = fmt.Sprintf("%.2fG", float64(sp.TotalSize)/float64(gsize))
	} else if sp.TotalSize > msize {
		TotalSize_str = fmt.Sprintf("%.2fM", float64(sp.TotalSize)/float64(msize))
	} else {
		TotalSize_str = fmt.Sprintf("%.2fK", float64(sp.TotalSize)/float64(1024))
	}

	logMsg := fmt.Sprintf(log_format, sp.Total, int64(float64(sp.TotalSec)/now_used), int64(now_used), sp.TotalSec, size_speed, TotalSize_str,sp.TotalSuccess,sp.TotalSecSuccess)
	sp.PrintFn(logMsg,sp)

	atomic.StoreInt64(&sp.TotalSec, 0)
	atomic.StoreInt64(&sp.TotalSecSize, 0)
	atomic.StoreInt64(&sp.TotalSecSuccess, 0)
}
