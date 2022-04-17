package collector

import (
	topclient "github.com/jmnote/kubectl-top-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type collector struct {
	topclient *topclient.Client

	// nodeCPUCoresDesc    *prometheus.Desc
	// nodeCPURatioDesc    *prometheus.Desc
	// nodeMemoryBytesDesc *prometheus.Desc
	// nodeMemoryRatioDesc *prometheus.Desc

	// podCPUCoresDesc          *prometheus.Desc
	// podMemoryBytesDesc       *prometheus.Desc
	// containerCPUCoresDesc    *prometheus.Desc
	// containerMemoryBytesDesc *prometheus.Desc

	nodeCPUMillicoresDesc   *prometheus.Desc
	nodeCPUPercentDesc      *prometheus.Desc
	nodeMemoryMibibytesDesc *prometheus.Desc
	nodeMemoryPercentDesc   *prometheus.Desc

	podCPUMillicoresDesc         *prometheus.Desc
	podMemoryMibibytesDesc       *prometheus.Desc
	containerCPUMillicoresDesc   *prometheus.Desc
	containerMemoryMibibytesDesc *prometheus.Desc
}

func NewCollector() (*collector, error) {
	topclient, err := topclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &collector{
		topclient: topclient,

		// nodeCPUCoresDesc:    prometheus.NewDesc("kubectl_top_node_cpu_cores", "kubectl top node; CPU(cores)", []string{"name"}, nil),
		// nodeCPURatioDesc:    prometheus.NewDesc("kubectl_top_node_cpu_ratio", "kubectl top node; CPU%", []string{"name"}, nil),
		// nodeMemoryBytesDesc: prometheus.NewDesc("kubectl_top_node_memory_bytes", "kubectl top node; MEMORY(bytes)", []string{"name"}, nil),
		// nodeMemoryRatioDesc: prometheus.NewDesc("kubectl_top_node_memory_ratio", "kubectl top node; MEMORY%", []string{"name"}, nil),

		// podCPUCoresDesc:          prometheus.NewDesc("kubectl_top_pod_cpu_cores", "kubectl top pod -A; CPU(cores)", []string{"namepace", "name"}, nil),
		// podMemoryBytesDesc:       prometheus.NewDesc("kubectl_top_pod_memory_bytes", "kubectl top pod; -A MEMORY(bytes)", []string{"namepace", "name"}, nil),
		// containerCPUCoresDesc:    prometheus.NewDesc("kubectl_top_pod_cpu_cores", "kubectl top pod -A --containers; CPU(cores)", []string{"namepace", "pod", "name"}, nil),
		// containerMemoryBytesDesc: prometheus.NewDesc("kubectl_top_pod_memory_bytes", "kubectl top pod -A --containers; MEMORY(bytes)", []string{"namepace", "pod", "name"}, nil),

		nodeCPUMillicoresDesc:   prometheus.NewDesc("kubectl_top_node_cpu_millicores", "kubectl top node; CPU(cores)", []string{"name"}, nil),
		nodeCPUPercentDesc:      prometheus.NewDesc("kubectl_top_node_cpu_percent", "kubectl top node; CPU%", []string{"name"}, nil),
		nodeMemoryMibibytesDesc: prometheus.NewDesc("kubectl_top_node_memory_mibibytes", "kubectl top node; MEMORY(bytes)", []string{"name"}, nil),
		nodeMemoryPercentDesc:   prometheus.NewDesc("kubectl_top_node_memory_percent", "kubectl top node; MEMORY%", []string{"name"}, nil),

		podCPUMillicoresDesc:         prometheus.NewDesc("kubectl_top_pod_cpu_millicores", "kubectl top pod -A; CPU(cores)", []string{"namepace", "name"}, nil),
		podMemoryMibibytesDesc:       prometheus.NewDesc("kubectl_top_pod_memory_mibibytes", "kubectl top pod -A; MEMORY(bytes)", []string{"namepace", "name"}, nil),
		containerCPUMillicoresDesc:   prometheus.NewDesc("kubectl_top_pod_container_cpu_millicores", "kubectl top pod -A --containers; CPU(cores)", []string{"namepace", "pod", "name"}, nil),
		containerMemoryMibibytesDesc: prometheus.NewDesc("kubectl_top_pod_container_memory_mibibytes", "kubectl top pod -A --containers; MEMORY(bytes)", []string{"namepace", "pod", "name"}, nil),
	}, nil
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.collectNodeMetrics(ch)
	c.collectPodAndContainerMetrics(ch)
}

func (c *collector) collectNodeMetrics(ch chan<- prometheus.Metric) {
	nodeMetricsList, err := c.topclient.GetNodeMetricsList()
	if err != nil {
		log.Println(err)
		return
	}
	for _, m := range nodeMetricsList {
		ch <- prometheus.MustNewConstMetric(c.nodeCPUMillicoresDesc, prometheus.GaugeValue, float64(m.CPUMillicores), []string{m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.nodeCPUPercentDesc, prometheus.GaugeValue, m.CPUPercent, []string{m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.nodeMemoryMibibytesDesc, prometheus.GaugeValue, float64(m.MemoryMibibytes), []string{m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.nodeMemoryPercentDesc, prometheus.GaugeValue, m.MemoryPercent, []string{m.Name}...)
	}
}

func (c *collector) collectPodAndContainerMetrics(ch chan<- prometheus.Metric) {
	podAndContainerMetricsList, err := c.topclient.GetPodAndContainerMetricsList()
	if err != nil {
		log.Println(err)
		return
	}
	for _, m := range podAndContainerMetricsList.PodMetricsList {
		ch <- prometheus.MustNewConstMetric(c.podCPUMillicoresDesc, prometheus.GaugeValue, float64(m.CPUMillicores), []string{m.Namespace, m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.podMemoryMibibytesDesc, prometheus.GaugeValue, float64(m.MemoryMibibytes), []string{m.Namespace, m.Name}...)
	}
	for _, m := range podAndContainerMetricsList.ContainerMetricsList {
		ch <- prometheus.MustNewConstMetric(c.containerCPUMillicoresDesc, prometheus.GaugeValue, float64(m.CPUMillicores), []string{m.Namespace, m.Pod, m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.containerMemoryMibibytesDesc, prometheus.GaugeValue, float64(m.MemoryMibibytes), []string{m.Namespace, m.Pod, m.Name}...)
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nodeCPUMillicoresDesc
	ch <- c.nodeCPUPercentDesc
	ch <- c.nodeMemoryMibibytesDesc
	ch <- c.nodeMemoryPercentDesc

	ch <- c.podCPUMillicoresDesc
	ch <- c.podMemoryMibibytesDesc
	ch <- c.containerCPUMillicoresDesc
	ch <- c.containerMemoryMibibytesDesc
}
