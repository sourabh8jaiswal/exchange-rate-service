package cache

import (
	"fmt"
	"sync"
)

var rateCache = struct {
	sync.RWMutex
	data map[string]float64
}{data: make(map[string]float64)}

func GetRate(from, to, date string) float64 {
	key := fmt.Sprintf("%s_%s_%s", from, to, date)
	rateCache.RLock()
	defer rateCache.RUnlock()
	return rateCache.data[key]
}

func StoreRate(from, to, date string, rate float64) {
	key := fmt.Sprintf("%s_%s_%s", from, to, date)
	rateCache.Lock()
	defer rateCache.Unlock()
	rateCache.data[key] = rate
}
