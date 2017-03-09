package speed

import "fmt"
import "time"
import "sync"
import "log"
import "strings"

type detailType struct {
	time   time.Time
	used   time.Duration
	detail map[string]uint64
	rw     sync.RWMutex
	keysNum    []string
	keysSize   []string
}


func (d *detailType) getKey(key string, event string) string{
	return fmt.Sprintf("%s_%s", key, event)
}

func (d *detailType) Add(key string, keyType string, num uint64) {
	d.rw.Lock()
	defer d.rw.Unlock()
	
	keyTotal := d.getKey(key, "total")
	keyData := d.getKey(key, keyType)
	
	if _, has := d.detail[keyTotal]; !has {
		d.detail[keyTotal] = 0
		if strings.HasSuffix(key, "_size"){
			d.keysSize=append(d.keysSize, key)
		}else{
			d.keysNum=append(d.keysNum, key)
		}
	}
	
	d.detail[keyTotal]+=num
	
	if _,has:=d.detail[keyData];!has{
		d.detail[keyData] = 0
	}
	
	d.detail[keyData]+=num
}

func (d *detailType) Get(key string) (total uint64, success uint64, fail uint64) {
	d.rw.RLock()
	defer d.rw.RUnlock()
	
	total, _= d.detail[d.getKey(key, "total")]
	success, _ = d.detail[d.getKey(key, "success")]
	fail, _ = d.detail[d.getKey(key, "fail")]
	return
}

func (d *detailType) compare(d2 *detailType) *detailType {
	d3 := newDetailType()
	for k, v := range d.detail {
		v2, _ := d2.detail[k]
		d3.detail[k] = v - v2
	}
	d3.used = d2.time.Sub(d.time)
	return d3
}

func (d *detailType) qpsInfo(d2 *detailType) []string {
	f1 := "{%s:<total:%d(qps:%.1f),success:%d(qps:%.1f),fail:%d(qps:%.1f%%)>}"
	var lines []string
	diff := d.compare(d2)
	used := diff.used.Seconds()
	for _, k := range d.keysNum {
		total, success, fail := diff.Get(k)
		line := fmt.Sprintf(
			f1, k,
			total, float64(total)/used,
			success, float64(success)/used,
			fail, float64(fail)/used,
		)
		lines = append(lines, line)
	}
	return lines
}

func newDetailType() *detailType {
	d := &detailType{
		time:   time.Now(),
		detail: make(map[string]uint64),
	}
	return d
}

type logPrinter struct {
	ticker        *time.Ticker
	current        *detailType
	last        *detailType
	printFn       func(string)
	lastPrintTime time.Time
}

func newLogPrinter(second int64, printFn func(string),current *detailType) *logPrinter {
	l := &logPrinter{
		current:current,
		ticker:        time.NewTicker(time.Duration(second) * time.Second),
		lastPrintTime: time.Now(),
		printFn:       printFn,
		last: newDetailType(),
	}
	return l
}

func (l *logPrinter) print(now time.Time) {
	now_used := now.Sub(l.lastPrintTime).Seconds()
	l.lastPrintTime = now
	if now_used == 0 {
		return
	}
	ls:=l.current.qpsInfo(l.last)
	l.copyDetail()
//	fmt.Println(ls)
//	msg:=fmt.Sprintf("hello:%x", l.current)
	l.printFn(strings.Join(ls, ","))
}

func (l *logPrinter)copyDetail(){
	l.last.keysNum=l.current.keysNum
	l.last.keysSize=l.current.keysSize
	l.last.detail=l.current.detail
	l.last.time=l.current.time
}

func (l *logPrinter) Start() {
	go func() {
		for now := range l.ticker.C {
			l.print(now)
		}
	}()
}
func (l *logPrinter) Stop() {
	l.ticker.Stop()
}

type Speed struct {
	detail     *detailType
	logPrinter map[int64]*logPrinter
}

func NewSpeed(id string, sec int64, printFn func(string)) *Speed {
	sp := new(Speed)

	if printFn == nil {
		printFn = DefaultPrintFn
	}

	fn := func(str string) {
		printFn(fmt.Sprintf("[id=%s] %s", id, str))
	}

	sp.detail = newDetailType()

	sp.logPrinter = make(map[int64]*logPrinter)

	for _,second := range []int64{sec, 60, 600, 1800, 3600} {
		sp.logPrinter[second] = newLogPrinter(second, fn,sp.detail)
	}
	for _, lp := range sp.logPrinter {
		lp.Start()
	}

	return sp
}

func DefaultPrintFn(msg string) {
	log.Println("speed", msg)
}

func (sp *Speed) Success(key string, num uint64) {
	sp.detail.Add(key, "success", num)
}

func (sp *Speed) Fail(key string, num uint64) {
	sp.detail.Add(key, "fail", num)
}

func (sp *Speed) Stop() {
	for _, lp := range sp.logPrinter {
		lp.Stop()
	}
}

var msize = float64(1024 * 1024)
var gsize = msize * 1024

func sizeHumanFormat(size uint64) string {
	fsize := float64(size)
	if fsize > gsize {
		return fmt.Sprintf("%.2fG", fsize/gsize)
	} else if fsize > msize {
		return fmt.Sprintf("%.2fM", fsize/msize)
	} else if fsize > 1024.0 {
		return fmt.Sprintf("%.2fK", fsize/1024.0)
	} else {
		return fmt.Sprintf("%.1fB", fsize)
	}
}
