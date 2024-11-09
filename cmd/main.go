package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	shplugexporter "github.com/ruupert/shplug_exporter"
)

type Shelly struct {
	Current     float64
	Ampere      float64
	Freq        float64
	Temperature float64
	Voltage     float64
	Running     float64
}

func (x *Shelly) Refresh() {
	d := shplugexporter.Plug{
		Hostname: "10.9.0.15",
		Device:   "",
	}
	c := shplugexporter.NewClient(d)
	stats, err := c.SwitchGetStatus()
	if err != nil {
		fmt.Println(err)
	}
	x.set("pwr", stats.Result.Current)
	x.set("amp", stats.Result.Apower)
	x.set("frq", stats.Result.Freq)
	x.set("vlt", stats.Result.Voltage)
	x.set("tmp", stats.Result.Temperature.TC)
	if stats.Result.Output {
		x.set("run", 1.0)
	} else {
		x.set("run", 0.0)
	}
}

func (x *Shelly) set(key string, value float64) error {
	switch key {
	case "pwr":
		x.Current = value
	case "amp":
		x.Ampere = value
	case "frq":
		x.Freq = value
	case "tmp":
		x.Temperature = value
	case "vlt":
		x.Voltage = value
	case "run":
		x.Running = value
	}
	return nil
}

func (x *Shelly) get(key string) float64 {
	switch key {
	case "pwr":
		return x.Current
	case "amp":
		return x.Ampere
	case "frq":
		return x.Freq
	case "tmp":
		return x.Temperature
	case "vlt":
		return x.Voltage
	}
	return 0.0
}

func recordMetrics(x *Shelly) {
	go func() {
		for {
			x.Refresh()
			pwrGauge.Set(x.get("pwr"))
			ampGauge.Set(x.get("amp"))
			frqGauge.Set(x.get("frq"))
			tmpGauge.Set(x.get("tmp"))
			vltGauge.Set(x.get("vlt"))
			runGauge.Set(x.get("run"))
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	pwrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_current",
		Help: "current current thingy gauge",
	})
	ampGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_amperes",
		Help: "ampere thingy gauge",
	})
	frqGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_frequency",
		Help: "frequency thingy gauge",
	})
	tmpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_temperature",
		Help: "temperature thingy gauge",
	})
	vltGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_voltage",
		Help: "voltage thingy gauge",
	})
	runGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_running",
		Help: "switch enabled thingy gauge",
	})
)

func main() {
	recordMetrics(&Shelly{})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
