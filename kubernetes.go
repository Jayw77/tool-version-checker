package main

import (
	"context"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesImage struct {
	Name                  string
	Cluster               string
	PodName               string
	Namespace             string
	Image                 string
	Custom                bool
	LatestVersionEndpoint EndpointConfig
	Version               Version
}

func GetKubernetesImageVersions() []KubernetesImage {
	var kubernetesImage KubernetesImage
	var kubernetesImages []KubernetesImage

	for _, cluster := range config.Kubernetes.Clusters {
		clientset := getClientSet(cluster)

		pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				sContainerImage := strings.Split(container.Image, ":")
				kubernetesImage.Cluster = cluster.Name
				kubernetesImage.PodName = pod.Name
				kubernetesImage.Namespace = pod.Namespace
				kubernetesImage.Image = sContainerImage[0]
				kubernetesImage.Version.Current = sContainerImage[1]

				customImage := GetCustomImage(kubernetesImage.Image)
				if customImage.LatestVersion.Endpoint != "" {
					kubernetesImage.Name = customImage.Name
					kubernetesImage.Custom = true
					kubernetesImage.LatestVersionEndpoint = customImage.LatestVersion
				} else {
					kubernetesImage.Name = ContainerImageLatestVersionEndpointNames[kubernetesImage.Image]
					kubernetesImage.LatestVersionEndpoint = LatestVersionEndpoints[ContainerImageLatestVersionEndpointNames[kubernetesImage.Image]]
				}

				kubernetesImages = append(kubernetesImages, kubernetesImage)
			}
		}
	}

	return kubernetesImages
}

func getClientSet(cluster *KubernetesCluster) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	if cluster.KubeConfigPath != "" {
		// use kube config file
		kubeConfigFile, err := os.ReadFile(cluster.KubeConfigPath)
		if err != nil {
			panic(err.Error())
		}

		// use the current context in kubeconfig
		config, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigFile)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// default to in-cluster kube config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	return clientset
}

func GetCustomImage(image string) CustomImage {
	for _, customImage := range config.Kubernetes.CustomImages {
		if customImage.Image == image {
			return *customImage
		}
	}
	return CustomImage{}
}
