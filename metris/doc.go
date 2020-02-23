// Package metrics defines common metrics, like counter, gauge, timer, histogram.
// Also, it provides `Sinker` interface to integrate with other monitor systems,
// like Prometheus, etc.
//
// metrics  provides following operations:
//
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
