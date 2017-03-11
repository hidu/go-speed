package speed

import (
	"fmt"
)

func utilPer(num uint64, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(num) * 100 / float64(total)
}

func utilQps(num uint64, used float64) float64 {
	if num == 0 || used == 0.0 {
		return 0
	}
	return float64(num) / used
}

var _util_msize = float64(1024 * 1024)
var _util_gsize = _util_msize * 1024
var _util_tsize = _util_gsize * 1024

func utilSizeHumanFormat(fsize float64) string {
	if fsize > _util_tsize {
		return fmt.Sprintf("%.2fT", fsize/_util_tsize)
	} else if fsize > _util_gsize {
		return fmt.Sprintf("%.2fG", fsize/_util_gsize)
	} else if fsize > _util_msize {
		return fmt.Sprintf("%.2fM", fsize/_util_msize)
	} else if fsize > 1024.0 {
		return fmt.Sprintf("%.2fK", fsize/1024.0)
	} else {
		return fmt.Sprintf("%.1fB", fsize)
	}
}
