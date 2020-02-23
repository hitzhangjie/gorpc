// Package metrics 定义了常见粒度的监控指标，如Counter、MetricsGauge、MetricsTimer、MetricsHistogram，
// 并在此基础上定义了与具体的外部监控系统对接的接口MetricsSink，对接具体的监控如公司
// Monitor或者外部开源的Prometheus等，只需是吸纳对应的MetricsSink接口即可.
//
// 为了使用方便，提供了两套常用方法：
// 1. counter
// - reqNumCounter := metrics.Counter("req.num")
//   reqNumCounter.Incr()
// - metrics.IncrCounter("req.num", 1)
//
// 2. gauge
// - cpuAvgLoad := metrics.Gauge("cpu.avgload")
//   cpuAvgLoad.Set(0.75)
// - metrics.SetGauge("cpu.avgload", 0.75)
//
// 3. timer
// - timeCostTimer := metrics.Timer("req.proc.timecost")
//   timeCostTimer.Record()
// - timeCostDuration := time.Millisecond * 2000
//   metrics.RecordTimer("req.proc.timecost", timeCostDuration)
//
// 4. histogram
// - Buckets := metrics.NewDurationBounds(time.Second, time.Second*2, time.Second*5),
//   timeCostDist := metrics.Histogram("timecost.distribution", Buckets)
//   timeCostDist.AddSample(float64(time.Second*3))
// - metrics.AddSample("timecost.distribution", float64(time.Second*3))
package metrics
