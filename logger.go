package speed

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type logPrinter struct {
	ticker        *time.Ticker
	second        int
	current       *speedDetail
	last          *speedDetail
	printFn       func(string)
	lastPrintTime time.Time
}

func newLogPrinter(second int, printFn func(string), current *speedDetail) *logPrinter {
	l := &logPrinter{
		current:       current,
		second:        second,
		ticker:        time.NewTicker(time.Duration(second) * time.Second),
		lastPrintTime: time.Now(),
		printFn:       printFn,
		last:          newSpeedDetail(),
	}
	return l
}

func (l *logPrinter) print(now time.Time) {
	now_used := now.Sub(l.lastPrintTime).Seconds()
	l.lastPrintTime = now
	if now_used == 0 {
		return
	}
	ls := l.current.qpsInfo(l.last)
	l.copyDetail()
	l.printFn(fmt.Sprintf("[sec=%d] %s", l.second, strings.Join(ls, ";")))
}

func (l *logPrinter) copyDetail() {
	l.last.keysNum = l.current.keysNum
	l.last.keysSize = l.current.keysSize
	l.last.time = time.Now()

	l.current.rw.RLock()
	defer l.current.rw.RUnlock()

	l.last.rw.Lock()
	defer l.last.rw.Unlock()

	for k, v := range l.current.detail {
		l.last.detail[k] = v
	}

}

func (l *logPrinter) Start() {
	go func() {
		for now := range l.ticker.C {
			l.print(now)
		}
	}()
}
func (l *logPrinter) Stop(canPrint bool) {
	l.ticker.Stop()
	if canPrint || time.Now().Sub(l.lastPrintTime).Seconds() > float64(l.second) {
		l.print(time.Now())
	}
}

func DefaultPrintFn(msg string) {
	log.Println("speed", msg)
}
