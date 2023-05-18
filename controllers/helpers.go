package controllers

import (
	"context"
	"log"
	"strings"
	"fmt"
	"math"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListNodes(r *TestAppReconciler) ([]corev1.Node, error) {
	nodeList := &corev1.NodeList{}
	err := r.Client.List(context.TODO(), nodeList)
	if err != nil {
		return nil, err
	}

	return nodeList.Items, nil
}

func ListPodsOnNode(r *TestAppReconciler, nodeName string) ([]corev1.Pod, error) {
	podList := &corev1.PodList{}
	listOptions := &client.ListOptions{}
	err := r.Client.List(context.TODO(), podList, listOptions)
	if err != nil {
		return nil, err
	}

	filteredPods := make([]corev1.Pod, 0)
	for _, pod := range podList.Items {
		if pod.Spec.NodeName == nodeName {
			filteredPods = append(filteredPods, pod)
		}
	}

	return filteredPods, nil
}

func CordonNode(r *TestAppReconciler, nodeName string) error {
	node := &corev1.Node{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{Name: nodeName}, node)
	if err != nil {
		return err
	}

	if node.Spec.Unschedulable {
		log.Printf("Node '%s' is already cordoned", nodeName)
		return nil
	}

	node.Spec.Unschedulable = true
	err = r.Client.Update(context.TODO(), node)
	if err != nil {
		return err
	}

	log.Printf("Node '%s' cordoned successfully", nodeName)
	return nil
}

func UncordonNode(r *TestAppReconciler, nodeName string) error {
	node := &corev1.Node{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{Name: nodeName}, node)
	if err != nil {
		return err
	}

	if !node.Spec.Unschedulable {
		log.Printf("Node '%s' is already uncordoned", nodeName)
		return nil
	}

	node.Spec.Unschedulable = false
	err = r.Client.Update(context.TODO(), node)
	if err != nil {
		return err
	}

	log.Printf("Node '%s' uncordoned successfully", nodeName)
	return nil
}

func CalculateMaxIPs(cidr string) int {
	parts := strings.Split(cidr, "/")
	ip := parts[0]
	cidrBits := parts[1]

	ipParts := strings.Split(ip, ".")
	var ipBytes [4]int
	for i, part := range ipParts {
		ipBytes[i] = parseInt(part)
	}

	maskBits := parseInt(cidrBits)

	maxIPs := int(math.Pow(2, float64(32-maskBits)))

	return maxIPs
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
