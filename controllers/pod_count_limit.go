package controllers

import (
	corev1 "k8s.io/api/core/v1"
)

func PodLimit(r *TestAppReconciler, PerNodeLimit int, pods []corev1.Pod) bool {
	if len(pods) > PerNodeLimit {
		return true
	} else {
		return false
	}
}
