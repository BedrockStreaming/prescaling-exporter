package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type Client struct {
	Clientset kubernetes.Interface
	Config    *rest.Config
}

func NewClient() (*Client, error) {
	var config *rest.Config
	var err error

	if os.Getenv("KUBECONFIG") != "" {
		pathOptions := clientcmd.NewDefaultPathOptions()
		pathOptions.LoadingRules.DoNotResolvePaths = false
		c, err := pathOptions.GetStartingConfig()
		if err != nil {
			return nil, err
		}

		configOverrides := clientcmd.ConfigOverrides{}
		clientConfig := clientcmd.NewDefaultClientConfig(*c, &configOverrides)
		config, err = clientConfig.ClientConfig()
		if err != nil {
			return nil, err
		}

	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{clientset, config}, nil
}
