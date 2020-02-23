package metrics

// counter 计数器定义, 内部调用注册的插件MetricsSink传递counter信息到外部系统
type counter struct {
	name string
}

// Inc 计数器值+1
func (c *counter) Incr() {
	for _, sink := range sinks {
		sink.IncrCounter(c.name, 1)
	}
}

// Inc 计数器值增加v
func (c *counter) IncrBy(v float64) {
	for _, sink := range sinks {
		sink.IncrCounter(c.name, v)
	}
}
