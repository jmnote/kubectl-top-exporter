package main

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	// "k8s.io/client-go/rest"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
	// "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

var MeasuredResources = []v1.ResourceName{
	v1.ResourceCPU,
	v1.ResourceMemory,
}

func main2() {
	// client
	// config, err := rest.InClusterConfig()
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// nodeMetrics
	nodeMetrics, err := getNodeMetrics(metricsClient)
	if err != nil {
		panic(err.Error())
	}
	for _, m := range nodeMetrics.Items {
		fmt.Println("m===", m.Name)
		fmt.Printf("%#v\n", m.Usage)
		nodeResources := getNodeResources(&m)
		fmt.Println(m.Name, nodeResources.Cpu().MilliValue(), nodeResources.Memory().Value()/1024/1024)
	}

	// podMetrics
	// podMetrics, err := getPodMetrics(metricsClient)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// for _, m := range podMetrics.Items {
	// 	podResources := getPodResources(&m)
	// 	fmt.Println(m.Namespace, m.Name, podResources.Cpu().MilliValue(), podResources.Memory().Value()/1024/1024)
	// }

}

func getNodeMetrics(metricsClient metricsclientset.Interface) (*metricsapi.NodeMetricsList, error) {
	versionedMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
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

func getPodMetrics(metricsClient metricsclientset.Interface) (*metricsapi.PodMetricsList, error) {
	versionedMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
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

func getPodResources(m *metricsapi.PodMetrics) v1.ResourceList {
	resources := make(v1.ResourceList)
	for _, res := range MeasuredResources {
		resources[res], _ = resource.ParseQuantity("0")
	}
	for _, c := range m.Containers {
		for _, res := range MeasuredResources {
			quantity := resources[res]
			quantity.Add(c.Usage[res])
			resources[res] = quantity
		}
	}
	return resources
}
