package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	inFlightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})
	// 带有 "code", "method" 标签的计数器
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)
	// 带标签的 duration.如果除 `method,code` 外有其它标签,需要在包装器中调用 `CurryWith()` 或 `MustCurryWith()` 传入标签的值.
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)
	// 不带标签的 responseSize
	responseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)

	// 创建将被中间件包装的 Handlers
	pushHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Push"))
	})
	pullHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pull"))
	})

	// 在默认的注册表中注册所有的数据指标
	prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)

	// 按照以上定义的数据指标将 Handler 分组,并通过 `ObserverVec` 接口的 `MustCurryWith()` 方法传入 "handler" 标签
	pushChain := promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": "push"}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, pushHandler),
			),
		),
	)
	pullChain := promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": "pull"}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, pullHandler),
			),
		),
	)

	http.Handle("/metrics", promhttp.Handler())
	// 对不同的 HTTP 入口端点请求应用到带有不同标签的 Handler 中间件包装器
	http.Handle("/pull", pullChain)
	http.Handle("/push", pushChain)

	if err := http.ListenAndServe(":8087", nil); err != nil {
		log.Fatal(err)
	}
}
