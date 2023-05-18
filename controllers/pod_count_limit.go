package controllers

import (
	//"context"
	"log"
	"time"

	//corev1 "k8s.io/api/core/v1"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodLimit(r *TestAppReconciler, PerNodeLimit int, stopCh <-chan struct{}) {
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
					pods, err := ListPodsOnNode(r, node.Name)
					if err != nil {
						log.Printf("Error listing pods on node '%s': %v", node.Name, err)
					}
					if len(pods) > PerNodeLimit {
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
			time.Sleep(1 * time.Second)
		}
	}
}