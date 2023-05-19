package controllers

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

func DiskPercentageLimit(r *TestAppReconciler, DiskPercentageLimit int, node corev1.Node) bool {
	allocatable := node.Status.Allocatable
	ephemeralStorage, _ := allocatable[corev1.ResourceEphemeralStorage]
	total_capacity_bytes, _ := ConvertToBytes(ephemeralStorage.String())
	total_image_size_bytes := TotalNodeImageSize(r, node)
	current_usage_pecentage := 100 * total_image_size_bytes / total_capacity_bytes
	if current_usage_pecentage > int64(DiskPercentageLimit) {
		return true
	} else {
		return false
	}
}

func ConvertToBytes(value string) (int64, error) {
	unitMap := map[string]int64{
		"": 1,
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
		return 1, err
	}

	if multiplier, ok := unitMap[unit]; ok {
		return number * multiplier, nil
	}

	return 1, fmt.Errorf("Invalid unit: %s", unit)
}
