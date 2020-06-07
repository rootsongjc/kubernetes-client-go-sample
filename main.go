package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	deploymentName := flag.String("deployment", "", "deployment name")
	imageName := flag.String("image", "", "new image name")
	appName := flag.String("app", "app", "application name")

	flag.Parse()
	if *deploymentName == "" {
		fmt.Println("You must specify the deployment name.")
		os.Exit(0)
	}
	if *imageName == "" {
		fmt.Println("You must specify the new image name.")
		os.Exit(0)
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deployment, err := clientset.AppsV1beta1().Deployments("default").Get(*deploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("Deployment not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting deployment%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment\n")
		name := deployment.GetName()
		fmt.Println("name ->", name)
		containers := &deployment.Spec.Template.Spec.Containers
		found := false
		for i := range *containers {
			c := *containers
			if c[i].Name == *appName {
				found = true
				fmt.Println("Old image ->", c[i].Image)
				fmt.Println("New image ->", *imageName)
				c[i].Image = *imageName
			}
		}
		if found == false {
			fmt.Println("The application container not exist in the deployment pods.")
			os.Exit(0)
		}
		_, err := clientset.AppsV1beta1().Deployments("default").Update(deployment)
		if err != nil {
			panic(err.Error())
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
