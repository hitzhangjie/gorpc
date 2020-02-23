package metrics

import (
	"sync"
	"time"
)

var (
	// allow emit same metrics information to multi external system at same time
	sinks    = map[string]Sink{}
	sinksLck = sync.RWMutex{}

	counters    = map[string]Counter{}
	countersLck = sync.RWMutex{}

	gauges    = map[string]Gauge{}
	gaugesLck = sync.RWMutex{}

	timers    = map[string]Timer{}
	timersLck = sync.RWMutex{}

	histograms    = map[string]Histogram{}
	histogramsLck = sync.RWMutex{}

	meta = map[string]interface{}{}
)

// RegisterSink register one Sink
func RegisterSink(sink Sink) {
	sinksLck.Lock()
	sinks[sink.Name()] = sink
	sinksLck.Unlock()
}

// GetCounter create a counter named `name`
func GetCounter(name string) Counter {
	countersLck.RLock()
	c, ok := counters[name]
	countersLck.RUnlock()

	if ok && c != nil {
		return c
	}

	c = &counter{name: name}

	countersLck.Lock()
	counters[name] = c
	countersLck.Unlock()

	return c
}

// GetGauge create a gauge named `Name`
func GetGauge(name string) Gauge {
	histogramsLck.Lock()
	defer histogramsLck.Unlock()

	c, ok := gauges[name]
	if ok && c != nil {
		return c
	}

	return &gauge{name: name}
}

// GetTimer create a timer named `Name`
func GetTimer(name string) Timer {
	histogramsLck.Lock()
	defer histogramsLck.Unlock()

	c, ok := timers[name]
	if ok && c != nil {
		return c
	}

	t := &timer{name: name, start: time.Now()}
	timers[name] = t
	return t
}

// GetHistogram create a histogram named `Name` with `Buckets`
func GetHistogram(name string, buckets BucketBounds, meta ...HistogramMeta) Histogram {
	histogramsLck.RLock()
	c, ok := histograms[name]
	histogramsLck.RUnlock()

	if ok && c != nil {
		return c
	}

	h := newHistogram(name, buckets)
	histogramsLck.Lock()
	histograms[name] = h
	histogramsLck.Unlock()

	return h
}

// IncrCounter increment counter `key` by `value`, Counters should accumulate values
func IncrCounter(key string, value float64) {
	sinksLck.RLock()
	sinks := sinks
	sinksLck.RUnlock()

	for _, sink := range sinks {
		sink.IncrCounter(key, value)
	}
}

// SetGauge set gauge `key` with `value`, a Gauge should retain the last value it is set to
func SetGauge(key string, value float64) {
	sinksLck.RLock()
	sinks := sinks
	sinksLck.RUnlock()

	for _, sink := range sinks {
		sink.SetGauge(key, value)
	}
}

// RecordTimer record timer named `key` with `duration`
func RecordTimer(key string, duration time.Duration) {
	sinksLck.RLock()
	sinks := sinks
	sinksLck.RUnlock()

	for _, sink := range sinks {
		sink.RecordTimer(key, duration)
	}
}

// AddSample add one sample `key` with `value`, Samples are added to histogram for timing or value distribution case
// `key` is the name of histogram
func AddSample(key string, value float64) {
	sinksLck.RLock()
	sinks := sinks
	sinksLck.RUnlock()

	for _, sink := range sinks {
		sink.AddSample(key, value)
	}
}
