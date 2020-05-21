package client

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func NewDynamicForToken(token string) (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.BearerTokenFile = ""
	config.BearerToken = token

	client, err := dynamic.NewForConfig(dynamic.ConfigFor(config))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewDynamicInCluster() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
