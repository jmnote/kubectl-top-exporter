package client

import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (c *Client) GetPodAndContainerMetricsList() (*PodAndContainerMetricsList, error) {
	metrics, err := c.getPodMetricsFromMetricsAPI()
	if err != nil {
		return nil, err
	}
	return &PodAndContainerMetricsList{
		PodMetricsList:       c.getPodMetrics(metrics.Items),
		ContainerMetricsList: c.getContainerMetrics(metrics.Items),
	}, nil
}

func (c *Client) getPodMetricsFromMetricsAPI() (*metricsapi.PodMetricsList, error) {
	versionedMetrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	metrics := &metricsapi.PodMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_PodMetricsList_To_metrics_PodMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (c *Client) getPodMetrics(metrics []metricsapi.PodMetrics) []PodMetrics {
	var PodMetricsList []PodMetrics
	var podCPUQuantity int64
	var podMemoryQuantity int64
	for _, m := range metrics {
		podCPUQuantity = 0
		podMemoryQuantity = 0
		for _, container := range m.Containers {
			containerCPUQuantity := container.Usage[v1.ResourceCPU]
			containerMemoryQuantity := container.Usage[v1.ResourceMemory]
			podCPUQuantity += containerCPUQuantity.Value()
			podMemoryQuantity += containerMemoryQuantity.Value()
		}
		PodMetricsList = append(PodMetricsList, PodMetrics{
			Namespace:       m.ObjectMeta.Namespace,
			Name:            m.ObjectMeta.Name,
			CPUMillicores:   podCPUQuantity,
			MemoryMibibytes: podMemoryQuantity,
		})
	}
	return PodMetricsList
}

func (c *Client) getContainerMetrics(metrics []metricsapi.PodMetrics) []ContainerMetrics {
	var ContainerMetricsList []ContainerMetrics
	for _, m := range metrics {
		for _, container := range m.Containers {
			cpuQuantity := container.Usage[v1.ResourceCPU]
			memoryQuantity := container.Usage[v1.ResourceMemory]
			ContainerMetrics := ContainerMetrics{
				Namespace:       m.ObjectMeta.Namespace,
				Pod:             m.ObjectMeta.Name,
				Name:            container.Name,
				CPUMillicores:   cpuQuantity.MilliValue(),
				MemoryMibibytes: memoryQuantity.Value() / (1024 * 1024),
			}
			ContainerMetricsList = append(ContainerMetricsList, ContainerMetrics)
		}
	}
	return ContainerMetricsList
}
