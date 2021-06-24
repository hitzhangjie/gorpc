package metrics

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

func NewConsoleSink() Sink {
	return &ConsoleSink{
		counters:   make(map[string]float64),
		gauges:     make(map[string]float64),
		timers:     make(map[string]timer),
		histograms: make(map[string]histogram),
		cm:         sync.RWMutex{},
		gm:         sync.RWMutex{},
		tm:         sync.RWMutex{},
		hm:         sync.RWMutex{},
	}
}

type ConsoleSink struct {
	counters   map[string]float64
	gauges     map[string]float64
	timers     map[string]timer
	histograms map[string]histogram

	cm sync.RWMutex
	gm sync.RWMutex
	tm sync.RWMutex
	hm sync.RWMutex
}

func (c *ConsoleSink) Name() string {
	return "console"
}

func (c *ConsoleSink) IncrCounter(key string, value float64) {
	if c.counters == nil {
		return
	}
	c.cm.Lock()
	c.counters[key]++
	c.cm.Unlock()
	fmt.Printf("metrics counter[key] = %s val = %v", key, value)
}

func (c *ConsoleSink) SetGauge(key string, value float64) {
	if c.gauges == nil {
		return
	}
	c.cm.Lock()
	c.gauges[key] = value
	c.cm.Unlock()
	fmt.Printf("metrics gauge[key] = %s val = %v\n", key, value)
}

func (c *ConsoleSink) RecordTimer(key string, duration time.Duration) {
}

func (c *ConsoleSink) AddSample(key string, value float64) {

	// fixme! using atomic instead of locks

	//if c.histograms == nil {
	//	return
	//}
	//c.hm.RLock()
	//hist, ok := c.histograms[key]
	//c.hm.RUnlock()
	//
	//if ok {
	//	idx := sort.SearchFloat64s(hist.LookupByValue, value)
	//	upperBound := hist.Buckets[idx].ValueUpperBound
	//	fmt.Printf("metrics histogram[%s.%v] = %v", hist.Name, upperBound, hist.Buckets[idx].samples)
	//	return
	//}

	histogramsLck.RLock()
	h := histograms[key]
	histogramsLck.RUnlock()

	v, ok := h.(*histogram)
	if !ok {
		return
	}
	//c.hm.Lock()
	hist := *v

	histogramsLck.Lock()
	idx := sort.SearchFloat64s(hist.LookupByValue, value)
	upperBound := hist.Buckets[idx].ValueUpperBound
	hist.Buckets[idx].samples += value
	histogramsLck.Unlock()

	//c.histograms[key] = hist
	//c.hm.Unlock()

	fmt.Printf("metrics histogram[%s.%v] = %v", hist.Name, upperBound, hist.Buckets[idx].samples)
}
