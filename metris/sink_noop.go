package metrics

import (
	"time"
)

type NoopSink struct {
}

func (n *NoopSink) Name() string {
	//panic("implement me")
	return "noop"
}

func (n *NoopSink) IncrCounter(key string, value float64) {
	//panic("implement me")
}

func (n *NoopSink) SetGauge(key string, value float64) {
	//panic("implement me")
}

func (n *NoopSink) RecordTimer(key string, duration time.Duration) {
	//panic("implement me")
}

func (n *NoopSink) AddSample(key string, value float64) {
	//panic("implement me")
}
