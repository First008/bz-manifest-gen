package main

import (
	"context"
	"fmt"

	//v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	clierntset := kubernetes.NewForConfigOrDie(config)

	nodeList, err := clierntset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	for _, n := range nodeList.Items {
		fmt.Println(n.Name)
	}

	// newStatefulset := &v1.StatefulSet{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "geth",
	// 	},
	// 	Spec: v1.StatefulSetSpec{
	// 		Replicas: 1,
	// 		Spec: metav1.ObjectMeta{

	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{

			Containers: []corev1.Container{
				{Name: "busybox", Image: "busybox:latest", Command: []string{"sleep", "1000000"}},
			},
		},
	}

	pod, err := clierntset.CoreV1().Pods("default").Create(context.Background(), newPod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pod)

}
