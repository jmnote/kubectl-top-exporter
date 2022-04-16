package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type kubectlTopCollector struct {
	kubectlTopNodeMetric *prometheus.Desc
	kubectlTopPodMetric  *prometheus.Desc
}

func newKubectlTopCollector() *kubectlTopCollector {
	return &kubectlTopCollector{
		kubectlTopNodeMetric: prometheus.NewDesc("kubectl_top_node", "kubectl top node", nil, nil),
		kubectlTopPodMetric:  prometheus.NewDesc("kubectl_top_pod", "kubectl top pod -A", nil, nil),
	}
}

func (c *kubectlTopCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.kubectlTopNodeMetric
	ch <- c.kubectlTopPodMetric
}

func (c *kubectlTopCollector) Collect(ch chan<- prometheus.Metric) {
	m1 := prometheus.MustNewConstMetric(c.kubectlTopNodeMetric, prometheus.GaugeValue, 3.14)
	m2 := prometheus.MustNewConstMetric(c.kubectlTopPodMetric, prometheus.GaugeValue, 1.414)
	// m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
	// m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	ch <- m1
	ch <- m2
}

func realmain() {
	collector := newKubectlTopCollector()
	prometheus.Register(collector)
	http.Handle("/metrics", promhttp.Handler())
	log.Print("Starting http server")
	log.Fatal(http.ListenAndServe(":9100", nil))
}

func main() {
	main2()
}
