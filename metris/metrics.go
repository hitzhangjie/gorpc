package metrics

import "time"

// MetricsCounter is the interface for emitting counter type metrics.
type MetricsCounter interface {
	// Incr increments the counter by one.
	Incr()

	// IncrBy increments the counter by a delta.
	IncrBy(delta float64)
}

// MetricsGauge is the interface for emitting gauge metrics.
type MetricsGauge interface {
	// Update sets the gauges absolute Value.
	Set(value float64)
}

// MetricsHistogram is the interface for emitting histogram metrics
type MetricsHistogram interface {
	// AddSample records a sample into histogram
	AddSample(value float64)
}

// MetricsTimer is the interface for emitting timer metrics.
type MetricsTimer interface {
	// Record a specific duration directly.
	Record()
}

// MetricsSink interface is used to transmit metrics information to an external system,
// like prometheus or tencent attr, atta, habo, etc.
type MetricsSink interface {
	Name() string
	CounterSink
	GaugeSink
	TimerSink
	HistorgramSink
}

// CounterSink sink accumulate values
type CounterSink interface {
	// IncrCounter Counters should accumulate values
	IncrCounter(key string, value float64)
}

// GaugeSink sink gauge values
type GaugeSink interface {
	// SetGauge a MetricsGauge should retain the last value it is set to
	SetGauge(key string, value float64)
}

// TimerSink sink timing values
type TimerSink interface {
	// Record a timer should retain the duration since timer start
	RecordTimer(key string, duration time.Duration)
}

// HistorgramSink sink samples in histogram
type HistorgramSink interface {
	// AddSample Samples are added to histogram for timing or value distribution case
	AddSample(key string, value float64)
}
