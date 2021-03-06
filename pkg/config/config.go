package config

import (
	"fmt"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitK8S() *rest.Config {

	config, err := clientcmd.BuildConfigFromFlags("", "./resources/kubeconfig")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config
}
