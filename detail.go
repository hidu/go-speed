package speed

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type speedDetail struct {
	time     time.Time
	used     time.Duration
	detail   map[string]uint64
	rw       sync.RWMutex
	keysNum  []string
	keysSize []string

	detailItemChan chan *detailItem
	workerWg       sync.WaitGroup
}

type detailItem struct {
	key     string
	keyType string
	num     uint64
}

func newDetailItem(key string, keyType string, num uint64) *detailItem {
	return &detailItem{
		key:     key,
		keyType: keyType,
		num:     num,
	}
}

func newSpeedDetail() *speedDetail {
	d := &speedDetail{
		time:   time.Now(),
		detail: make(map[string]uint64),
	}
	return d
}

func (d *speedDetail) getKey(key string, event string) string {
	return fmt.Sprintf("%s_%s", key, event)
}

func (d *speedDetail) initAddWorker() {

	d.detailItemChan = make(chan *detailItem, 1024)

	worker := func(item *detailItem) {
		d.rw.Lock()
		defer d.rw.Unlock()

		//	fmt.Println("Add",d.detail)
		keyTotal := d.getKey(item.key, "total")
		keyData := d.getKey(item.key, item.keyType)

		if _, has := d.detail[keyTotal]; !has {
			d.detail[keyTotal] = 0
			if strings.HasSuffix(item.key, "_size") {
				d.keysSize = append(d.keysSize, item.key)
			} else {
				d.keysNum = append(d.keysNum, item.key)
			}
		}

		d.detail[keyTotal] += item.num

		if _, has := d.detail[keyData]; !has {
			d.detail[keyData] = 0
		}

		d.detail[keyData] += item.num
	}

	d.workerWg.Add(1)
	go func() {
		for item := range d.detailItemChan {
			worker(item)
		}
		d.workerWg.Done()
	}()
}

func (d *speedDetail) compare(d2 *speedDetail) *speedDetail {
	d3 := newSpeedDetail()
	for k, v := range d.detail {
		if d2 != nil {
			v2, _ := d2.detail[k]
			d3.detail[k] = v - v2
		} else {
			d3.detail[k] = v
		}
	}
	if d2 != nil {
		d3.used = time.Now().Sub(d2.time)
	} else {
		d3.used = time.Now().Sub(d.time)
	}
	return d3
}

func (d *speedDetail) qpsInfo(d2 *speedDetail) []string {
	f1 := "[%s_%s:(total:%d,%.1f/s),(suc:%d,%.1f/s,%.1f%%),(fail:%d,%.1f/s,%.1f%%)]"
	f2 := "[%s_%s:(total:%s,%s/s),(suc:%s,%s/s,%.1f%%),(fail:%s,%s/s,%.1f%%)]"
	var lines []string

	d.rw.RLock()
	defer d.rw.RUnlock()

	var getQps = func(d2 *speedDetail) {
		diff := d.compare(d2)
		used := diff.used.Seconds()
		_type := "all"
		if d2 != nil {
			_type = fmt.Sprintf("%.0fs", used)
		}
		for _, k := range d.keysNum {
			total, success, fail := diff.Get(k)
			line := fmt.Sprintf(
				f1, k, _type,
				total, utilQps(total, used),
				success, utilQps(success, used), utilPer(success, total),
				fail, utilQps(fail, used), utilPer(fail, total),
			)
			lines = append(lines, line)
		}

		for _, k := range d.keysSize {
			total, success, fail := diff.Get(k)
			line := fmt.Sprintf(
				f2, k, _type,
				utilSizeHumanFormat(float64(total)), utilSizeHumanFormat(utilQps(total, used)),
				utilSizeHumanFormat(float64(success)), utilSizeHumanFormat(utilQps(success, used)), utilPer(success, total),
				utilSizeHumanFormat(float64(fail)), utilSizeHumanFormat(utilQps(fail, used)), utilPer(fail, total),
			)
			lines = append(lines, line)
		}
	}
	getQps(nil)
	getQps(d2)

	return lines
}

func (d *speedDetail) Add(key string, keyType string, num uint64) {
	item := newDetailItem(key, keyType, num)
	d.detailItemChan <- item
}

func (d *speedDetail) Get(key string) (total uint64, success uint64, fail uint64) {
	d.rw.RLock()
	defer d.rw.RUnlock()
	//	fmt.Println("GET",d.detail)
	total, _ = d.detail[d.getKey(key, "total")]
	success, _ = d.detail[d.getKey(key, "success")]
	fail, _ = d.detail[d.getKey(key, "fail")]
	return
}

func (d *speedDetail) Stop() {
	if d.detailItemChan != nil {
		close(d.detailItemChan)
	}
	d.workerWg.Wait()
}
