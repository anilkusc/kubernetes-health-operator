package controllers

import (
	//"context"
	appsv1 "github.com/anilkusc/ty-case/api/v1"
	"log"
	"time"
	//corev1 "k8s.io/api/core/v1"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

func LimitController(r *TestAppReconciler, testApp *appsv1.TestApp, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		default:
			nodes, err := ListNodes(r)
			if err != nil {
				log.Printf("Error on listing nodes: %v", err)
			} else {
				for _, node := range nodes {
					cordon := false
					pods, err := ListPodsOnNode(r, node.Name)
					if err != nil {
						log.Printf("Error listing pods on node '%s': %v", node.Name, err)
					}
					if PodLimit(r, testApp.Spec.PerNodePodLimit, pods) || IPLimit(r, testApp.Spec.PerNodeIpLimitPercentage, pods, node) || DiskPercentageLimit(r, testApp.Spec.DiskLimitPercentage, node) {
						cordon = true
					}
					if cordon {
						err = CordonNode(r, node.Name)
						if err != nil {
							log.Printf("Cannot cordon node '%s': %v", node.Name, err)
						}
					} else {
						err = UncordonNode(r, node.Name)
						if err != nil {
							log.Printf("Cannot uncordon node '%s': %v", node.Name, err)
						}
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}
}
