package client

import (
	"context"
	"errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (c *Client) GetNodeMetricsList() ([]NodeMetrics, error) {
	metrics, err := c.getNodeMetricsFromMetricsAPI()
	if err != nil {
		return nil, err
	}
	if len(metrics.Items) == 0 {
		return nil, errors.New("metrics not available yet")
	}
	var nodes []v1.Node
	nodeList, err := c.nodeClient.Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, nodeList.Items...)

	allocatable := make(map[string]v1.ResourceList)
	for _, n := range nodes {
		allocatable[n.Name] = n.Status.Allocatable
	}

	NodeMetricsList := getNodeMetrics(metrics.Items, allocatable)
	return NodeMetricsList, nil
}

func (c *Client) getNodeMetricsFromMetricsAPI() (*metricsapi.NodeMetricsList, error) {
	versionedMetrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	metrics := &metricsapi.NodeMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func getNodeMetrics(metrics []metricsapi.NodeMetrics, availableResources map[string]v1.ResourceList) []NodeMetrics {
	var NodeMetricsList []NodeMetrics
	for _, m := range metrics {
		available := availableResources[m.Name]
		cpuQuantity := m.Usage[v1.ResourceCPU]
		cpuAvailable := available[v1.ResourceCPU]
		memoryQuantity := m.Usage[v1.ResourceMemory]
		memoryAvailable := available[v1.ResourceMemory]
		NodeMetricsList = append(NodeMetricsList, NodeMetrics{
			Name:            m.Name,
			CPUMillicores:   cpuQuantity.MilliValue(),
			CPUPercent:      float64(cpuQuantity.MilliValue()) / float64(cpuAvailable.MilliValue()) * 100,
			MemoryMibibytes: memoryQuantity.Value() / (1024 * 1024),
			MemoryPercent:   float64(memoryQuantity.MilliValue()) / float64(memoryAvailable.MilliValue()) * 100,
		})
	}
	return NodeMetricsList
}
