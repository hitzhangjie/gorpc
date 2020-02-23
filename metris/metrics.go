package metrics

// Counter is the interface for emitting counter type metrics.
type Counter interface {
	// Incr increments the counter by one.
	Incr()

	// IncrBy increments the counter by a delta.
	IncrBy(delta float64)
}

// Gauge is the interface for emitting gauge metrics.
type Gauge interface {
	// Update sets the gauges absolute Value.
	Set(value float64)
}

// Histogram is the interface for emitting histogram metrics
type Histogram interface {
	// AddSample records a sample into histogram
	AddSample(value float64)
}

// Timer is the interface for emitting timer metrics.
type Timer interface {
	// Record a specific duration directly.
	Record()
}


