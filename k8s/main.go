/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// creates the in-cluster config
	fmt.Println("start")
	k8sconfig := flag.String("config", "./config", "kubernetes config file path")
	flag.Parse()
	fmt.Println("flag.Parse()")
	//config, err := rest.InClusterConfig()
	config, err := clientcmd.BuildConfigFromFlags("", *k8sconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}
	{
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		//pods 遍历
		for _, pod := range pods.Items {
			fmt.Printf("NAMESPACE:%v \t NAME: %v \t STATUS: %v\n", pod.GetNamespace(), pod.GetName(), pod.Status.Phase)
		}
		//namespace
		namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})

		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d namespace in the cluster\n", len(namespaces.Items))
		//namespaces遍历
		for _, namespace := range namespaces.Items {

			fmt.Printf("NAMESPACE:%v \t NAME: %v \t STATUS: %v\n", namespace.GetName(), namespace.GetCreationTimestamp(), namespace.Status.Phase)
		}
		// 获取deloyments
		Deployments, err := clientset.AppsV1().Deployments("t5000").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d Deployments in the cluster\n", len(Deployments.Items))
		//deploys 遍历
		for _, d := range Deployments.Items {
			fmt.Printf("NAMESPACE:%v \t NAME: %v \t STATUS: %v\n", d.Namespace, d.GetName(), d.Status.AvailableReplicas)

		}

		//获取sevrice
		Services, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, Service := range Services.Items {
			fmt.Printf("NAMESPACE:%v ID:%v \t \t NAME: %v \t type: %v \t type: %v\n", Service.Namespace, Service.UID, Service.GetName(), Service.Spec.Type, Service.Spec.Ports)
		}

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get("alarm-save-657f44d8b-hphtp", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod\n")
		}

		//time.Sleep(10 * time.Second)

		//
		fmt.Printf("node 信息\n")
			//获取sevrice
		Ndoes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	
		if err != nil {
			panic(err.Error())
		}
		for _, Node := range Ndoes.Items {
			cpu , _ :=  Node.Status.Capacity.Memory().AsInt64() 
		//	cpu , _ :=  Node.
			fmt.Printf("node:%v cpu:%v \t \t  \n", 
			Node.Name, cpu / 1024 / 1024 /1024 )
		}


	}
}
