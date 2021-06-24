package metrics

import "time"

// Sink interface is used to transmit metrics information to an external system,
// like prometheus or tencent attr, atta, habo, etc.
type Sink interface {
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
	// SetGauge a Gauge should retain the last value it is set to
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
