package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//k8pod "k8s.io/kubernetes/pkg/api/v1/pod"
)
var ns, label, field, maxClaims string

func main() {

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.StringVar(&label, "l", "", "Label selector")
	flag.StringVar(&field, "f", "", "Field selector")
	flag.StringVar(&maxClaims, "max-claims", "200Gi", "Maximum total claims to watch")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()

	// bootstrap config
	fmt.Println()
	fmt.Println("Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset : clientset == *CoreV1Client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	api := clientset.CoreV1()

	// initial list
	listOptions := metav1.ListOptions{LabelSelector: label, FieldSelector: field}
	// pvcs, err := api.PersistentVolumeClaims(ns).List(listOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	namespaces, err := api.Namespaces().List(listOptions)

	// printPVCs(pvcs)
	printNamespaces(namespaces, api)
	fmt.Println()
}

// printPVCs prints a list of PersistentVolumeClaim on console
// func printPVCs(pvcs *v1.PersistentVolumeClaimList) {
// 	if len(pvcs.Items) == 0 {
// 		log.Println("No claims found")
// 		return
// 	}
// 	template := "%-32s%-8s%-8s\n"
// 	fmt.Println("--- PVCs ----")
// 	fmt.Printf(template, "NAME", "STATUS", "CAPACITY")
// 	var cap resource.Quantity
// 	for _, pvc := range pvcs.Items {
// 		quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]
// 		cap.Add(quant)
// 		fmt.Printf(template, pvc.Name, string(pvc.Status.Phase), quant.String())
// 	}

// 	fmt.Println("-----------------------------")
// 	fmt.Printf("Total capacity claimed: %s\n", cap.String())
// 	fmt.Println("-----------------------------")
// }

func printPods(ns string, api *v1.CoreV1Client) {	

	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}
	pods, err := api.Pods(ns).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}
	if len(pods.Items) == 0 {
		log.Println("No pods found")
		return
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Spec.RestartPolicy, pod.Labels, pod.Annotations, pod.Status.Phase, pod.Spec.Volumes)
	}
}

func printNamespaces(namespaces *v1.NamespaceList, api *v1.CoreV1Client) {
	if len(namespaces.Items) == 0 {
		log.Println("No namespaces found")
		return
	}
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
		printPods(ns.Name, api)
	}
}
