package metrics_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/hitzhangjie/go-rpc/metris"
)

func Test_counter_Incr(t *testing.T) {
	metrics.RegisterSink(&metrics.ConsoleSink{})
	type fields struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"counter-1", fields{"counter-req.total"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := metrics.GetCounter(tt.fields.name)
			c.Incr()
			c.IncrBy(10)
		})
	}
}

func Test_gauge_Set(t *testing.T) {
	metrics.RegisterSink(&metrics.ConsoleSink{})
	type fields struct {
		name string
	}
	type args struct {
		v float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"gauge-cpu.avgload", fields{"cpu.avgload"}, args{0.75}},
		{"gauge-mem.avgload", fields{"mem.avgload"}, args{0.80}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := metrics.GetGauge(tt.fields.name)
			g.Set(tt.args.v)
		})
	}
}

func Test_histogram_AddSample(t *testing.T) {

	buckets := metrics.NewDurationBounds(time.Second, time.Second*2, time.Second*5)
	h := metrics.GetHistogram("req.timecost", buckets)

	metrics.RegisterSink(&metrics.ConsoleSink{})

	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"histogram-sample-0.5s", args{float64(time.Millisecond * 500)}},
		{"histogram-sample-1.5s", args{float64(time.Millisecond * 1500)}},
		{"histogram-sample-2.5s", args{float64(time.Millisecond * 2500)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.AddSample(tt.args.value)
		})
	}
}

func TestWithMeta(t *testing.T) {
	want := metrics.Options{}
	monitors := map[string]interface{}{"req.total": 10001, "req.fail": 10002, "req.succ": 10003}
	opt := metrics.WithMeta(monitors)
	opt(&want)

	type args struct {
		meta map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want metrics.Options
	}{
		{"monitor", args{meta: monitors}, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := metrics.Options{}
			opt := metrics.WithMeta(tt.args.meta)
			opt(&got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMeta() = %v, comp %v", got, tt.want)
			}
		})
	}
}
