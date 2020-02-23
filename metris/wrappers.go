package metrics

import (
	"sync"
	"time"
)

var (
	// allow emit same metrics information to multi external system at same time
	metricsSinks    = map[string]MetricsSink{}
	muxMetricsSinks = sync.RWMutex{}

	counters   = map[string]MetricsCounter{}
	gauges     = map[string]MetricsGauge{}
	timers     = map[string]MetricsTimer{}
	histograms = map[string]MetricsHistogram{}
	lock       = sync.RWMutex{}

	meta = map[string]interface{}{}
)

// RegisterMetricsSink register one MetricsSink
func RegisterMetricsSink(sink MetricsSink) {
	muxMetricsSinks.Lock()
	metricsSinks[sink.Name()] = sink
	muxMetricsSinks.Unlock()
}

// Counter create a counter named `Name`
func Counter(name string) MetricsCounter {
	lock.Lock()
	defer lock.Unlock()

	c, ok := counters[name]
	if ok && c != nil {
		return c
	}

	return &counter{name: name}
}

// Gauge create a gauge named `Name`
func Gauge(name string) MetricsGauge {
	lock.Lock()
	defer lock.Unlock()

	c, ok := gauges[name]
	if ok && c != nil {
		return c
	}

	return &gauge{name: name}
}

// Timer create a timer named `Name`
func Timer(name string) MetricsTimer {
	lock.Lock()
	defer lock.Unlock()

	c, ok := timers[name]
	if ok && c != nil {
		return c
	}

	t := &timer{name: name, start: time.Now()}
	timers[name] = t
	return t
}

// Histogram create a histogram named `Name` with `Buckets`
func Histogram(name string, buckets BucketBounds, meta ...HistogramMeta) MetricsHistogram {
	lock.Lock()
	defer lock.Unlock()

	c, ok := histograms[name]
	if ok && c != nil {
		return c
	}

	h := newHistogram(name, buckets)
	histograms[name] = h

	return h
}

// IncrCounter increment counter `key` by `value`, Counters should accumulate values
func IncrCounter(key string, value float64) {
	muxMetricsSinks.RLock()
	sinks := metricsSinks
	muxMetricsSinks.RUnlock()

	for _, sink := range sinks {
		sink.IncrCounter(key, value)
	}
}

// SetGauge set gauge `key` with `value`, a MetricsGauge should retain the last value it is set to
func SetGauge(key string, value float64) {
	muxMetricsSinks.RLock()
	sinks := metricsSinks
	muxMetricsSinks.RUnlock()

	for _, sink := range sinks {
		sink.SetGauge(key, value)
	}
}

// RecordTimer record timer named `key` with `duration`
func RecordTimer(key string, duration time.Duration) {
	muxMetricsSinks.RLock()
	sinks := metricsSinks
	muxMetricsSinks.RUnlock()

	for _, sink := range sinks {
		sink.RecordTimer(key, duration)
	}
}

// AddSample add one sample `key` with `value`, Samples are added to histogram for timing or value distribution case
// `key` is the name of histogram
func AddSample(key string, value float64) {
	muxMetricsSinks.RLock()
	sinks := metricsSinks
	muxMetricsSinks.RUnlock()

	for _, sink := range sinks {
		sink.AddSample(key, value)
	}
}

// GetHistogram return the histogram by `key`, `key` is the name of histogram
func GetHistogram(key string) (v histogram, ok bool) {
	lock.RLock()
	h := histograms[key]
	lock.RUnlock()

	hist, ok := h.(*histogram)
	if !ok {
		return
	}

	v = *hist
	ok = true

	return
}
