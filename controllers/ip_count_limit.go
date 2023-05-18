package controllers

import (
	corev1 "k8s.io/api/core/v1"
)

func IPLimit(r *TestAppReconciler, IpLimit int, pods []corev1.Pod, node corev1.Node) bool {
	current_usage_pecentage := 100 * len(pods) / CalculateMaxIPs(node.Spec.PodCIDR)
	if current_usage_pecentage > IpLimit {
		return true
	} else {
		return false
	}
}
