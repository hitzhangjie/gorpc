package metrics

import (
	"time"
)

// timer 计时器定义，内部调用注册的插件MetricsSink传递gauge信息到外部系统
type timer struct {
	name  string
	start time.Time
}

// Record 记录定时器耗时
func (t *timer) Record() {
	duration := time.Since(t.start)
	for _, sink := range sinks {
		sink.RecordTimer(t.name, duration)
	}
}
