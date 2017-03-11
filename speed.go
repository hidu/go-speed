package speed

import (
	"fmt"
	"sort"
)

type Speed struct {
	detail     *speedDetail
	logPrinter map[int]*logPrinter
}

func NewSpeed(id string, sec int, printFn func(string)) *Speed {
	sp := new(Speed)

	if printFn == nil {
		printFn = DefaultPrintFn
	}

	fn := func(str string) {
		printFn(fmt.Sprintf("[id=%s] %s", id, str))
	}

	sp.detail = newSpeedDetail()
	sp.detail.initAddWorker()

	sp.logPrinter = make(map[int]*logPrinter)

	for _, second := range []int{sec, 60, 600, 1800, 3600} {
		if _, has := sp.logPrinter[second]; !has {
			lp := newLogPrinter(second, fn, sp.detail)
			sp.logPrinter[second] = lp
			lp.Start()
		}
	}

	return sp
}

func (sp *Speed) Success(key string, num int) {
	sp.detail.Add(key, "success", uint64(num))
}

func (sp *Speed) Fail(key string, num int) {
	sp.detail.Add(key, "fail", uint64(num))
}

func (sp *Speed) Stop() {
	sp.detail.Stop()
	var secs []int
	for sec := range sp.logPrinter {
		secs = append(secs, sec)
	}
	if len(secs) < 1 {
		return
	}
	sort.Ints(secs)

	for sec, lp := range sp.logPrinter {
		lp.Stop(sec == secs[0])
	}
}
