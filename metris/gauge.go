package metrics

// gauge 时刻量定义，内部调用注册的插件MetricsSink传递gauge信息到外部系统
type gauge struct {
	name string
}

// Update 更新时刻量
func (g *gauge) Set(v float64) {
	for _, sink := range sinks {
		sink.SetGauge(g.name, v)
	}
}
