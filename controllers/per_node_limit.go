package controllers

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PerNodeLimit(r *TestAppReconciler, PerNodeLimit int, stopCh <-chan struct{}) {
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
