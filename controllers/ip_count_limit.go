package controllers

import (
	"log"
	"time"
)


func IPLimit(r *TestAppReconciler, IpLimit int, stopCh <-chan struct{}) {
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
					current_usage_pecentage := 100*len(pods)/CalculateMaxIPs(node.Spec.PodCIDR) 
					if  current_usage_pecentage > IpLimit {
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