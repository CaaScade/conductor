package app

import (
	"os/user"
	"strings"

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
		usr, err := user.Current()
		if err != nil {
			revel.AppLog.Fatalf("Error getting current user %v", err)
		}
		kubeConfig := strings.Replace(revel.Config.StringDefault("kube.config.path", ""), "~", usr.HomeDir, -1)
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
