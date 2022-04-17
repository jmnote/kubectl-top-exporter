package client

type NodeMetrics struct {
	Name            string
	CPUMillicores   int64
	CPUPercent      float64
	MemoryMibibytes int64
	MemoryPercent   float64
}

type PodMetrics struct {
	Namespace       string
	Name            string
	CPUMillicores   int64
	MemoryMibibytes int64
}

type ContainerMetrics struct {
	Namespace       string
	Pod             string
	Name            string
	CPUMillicores   int64
	MemoryMibibytes int64
}

type PodAndContainerMetricsList struct {
	PodMetricsList       []PodMetrics
	ContainerMetricsList []ContainerMetrics
}
