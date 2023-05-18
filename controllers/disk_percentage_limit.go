package controllers

import (
	"log"
	"time"
	"fmt"
	"strings"
	"strconv"
	corev1 "k8s.io/api/core/v1"
)


func DiskPercentageLimit(r *TestAppReconciler, DiskPercentageLimit int, stopCh <-chan struct{}) {
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
					allocatable := node.Status.Allocatable
					ephemeralStorage, _ := allocatable[corev1.ResourceEphemeralStorage]
					total_capacity_bytes , _ := ConvertToBytes(ephemeralStorage.String())				
					total_image_size_bytes := TotalNodeImageSize(r,node)
					current_usage_pecentage := 100*total_image_size_bytes/total_capacity_bytes
					fmt.Println(current_usage_pecentage)
					fmt.Println(DiskPercentageLimit)
					if current_usage_pecentage > int64(DiskPercentageLimit) {
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

func ConvertToBytes(value string) (int64, error) {
	unitMap := map[string]int64{
		"Ki": 1024,
		"Mi": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024,
		"Ti": 1024 * 1024 * 1024 * 1024,
		"Pi": 1024 * 1024 * 1024 * 1024 * 1024,
	}

	value = strings.TrimSpace(value)
	unit := value[len(value)-2:]
	numberStr := value[:len(value)-2]

	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return 0, err
	}

	if multiplier, ok := unitMap[unit]; ok {
		return number * multiplier, nil
	}

	return 0, fmt.Errorf("Invalid unit: %s", unit)
}