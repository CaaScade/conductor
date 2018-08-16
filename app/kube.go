package app

import (
	"github.com/revel/revel"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var Client *kubernetes.Clientset

func init() {
	revel.OnAppStart(func() {
		var cfg *rest.Config
		var err error
		kubeConfig := revel.Config.StringDefault("kube.config.path", "")
		if kubeConfig != "" {
			cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		} else {
			cfg, err = rest.InClusterConfig()
		}
		if err != nil {
			revel.AppLog.Fatalf("Error building kubernetes config %v", err)
		}

		client, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			revel.AppLog.Fatalf("unable to create client %v", err)
		}
		Client = client
	})
}
