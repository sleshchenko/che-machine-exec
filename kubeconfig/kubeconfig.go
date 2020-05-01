//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//
package kubeconfig

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/eclipse/che-machine-exec/api/model"
	exec_info "github.com/eclipse/che-machine-exec/exec-info"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	APIVersion     string     `yaml:"apiVersion"`
	Clusters       []Clusters `yaml:"clusters"`
	Users          []Users    `yaml:"users"`
	Contexts       []Contexts `yaml:"contexts"`
	CurrentContext string     `yaml:"current-context"`
	Kind           string     `yaml:"kind"`
}

type Clusters struct {
	Cluster ClusterInfo `yaml:"cluster"`
	Name    string      `yaml:"name"`
}

type ClusterInfo struct {
	CertificateAuthority string `yaml:"certificate-authority"`
	Server               string `yaml:"server"`
}

type Users struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

type User struct {
	Token string `yaml:"token"`
}

type Contexts struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type Context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
}

func generateKubeConfig(token, server, namespace string) *KubeConfig {
	return &KubeConfig{
		APIVersion: "v1",
		Clusters: []Clusters{
			{
				Cluster: ClusterInfo{
					CertificateAuthority: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
					Server:               server,
				},
				Name: server,
			},
		},
		Users: []Users{
			{
				Name: "developer",
				User: User{
					Token: token,
				},
			},
		},
		Contexts: []Contexts{
			{
				Context: Context{
					Cluster:   server,
					Namespace: namespace,
					User:      "developer",
				},
				Name: "developer-context",
			},
		},
		CurrentContext: "developer-context",
		Kind:           "Config",
	}
}

func createKubeConfigText(token, namespace string) string {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if host == "" || port == "" {
		return ""
	}

	server := "https://" + net.JoinHostPort(host, port)
	kubeconfig := generateKubeConfig(token, server, namespace)

	bytes, err := yaml.Marshal(&kubeconfig)
	if err != nil {
		logrus.Error("error: %v", err)
		return ""
	}
	return string(bytes)
}

// CreateKubeConfig creates a kubeconfig located at $KUBECONFIG if set.
// If it is not set then fall back to $HOME/.kube/config
func CreateKubeConfig(cmdRslv exec_info.InfoExecCreator, namespace string, token string, containerInfo *model.ContainerInfo) error {
	config := createKubeConfigText(token, namespace)

	kubeconfigDir, err := findKubeConfigDir(cmdRslv, containerInfo)
	if err != nil {
		return err
	}

	err = createKubeConfigDir(cmdRslv, kubeconfigDir, containerInfo)
	if err != nil {
		return err
	}

	err = syncKubeConfig(cmdRslv, config, kubeconfigDir, containerInfo)
	if err != nil {
		return err
	}
	return nil
}

func findKubeConfigDir(cmdRslv exec_info.InfoExecCreator, containerInfo *model.ContainerInfo) (string, error) {
	infoExec := cmdRslv.CreateInfoExec([]string{"sh", "-c", "echo $KUBECONFIG"}, containerInfo)
	if err := infoExec.Start(); err != nil {
		logrus.Debugf("Error is not available in %s/%s. Error: %s", containerInfo.PodName, containerInfo.ContainerName, err.Error())
		return "", errors.New("Could not retrieve $KUBECONFIG")
	}
	kubeconfigDir := infoExec.GetOutput()
	if kubeconfigDir != "\n" {
		return strings.Replace(kubeconfigDir, "/config", "", 1), nil
	}

	infoExec = cmdRslv.CreateInfoExec([]string{"sh", "-c", "echo $HOME"}, containerInfo)
	if err := infoExec.Start(); err != nil {
		logrus.Debugf("Error is not available in %s/%s. Error: %s", containerInfo.PodName, containerInfo.ContainerName, err.Error())
		return "", errors.New("Could not retrieve $HOME")
	}
	kubeconfigDir = strings.TrimSuffix(infoExec.GetOutput(), "\n")
	return kubeconfigDir + "/.kube", nil
}

func createKubeConfigDir(cmdRslv exec_info.InfoExecCreator, kubeconfigDir string, containerInfo *model.ContainerInfo) error {
	infoExec := cmdRslv.CreateInfoExec([]string{"sh", "-c", "mkdir -p " + kubeconfigDir}, containerInfo)
	if err := infoExec.Start(); err != nil {
		logrus.Debugf("Error is not available in %s/%s. Error: %s", containerInfo.PodName, containerInfo.ContainerName, err.Error())
		return errors.New("Could not create directory: " + kubeconfigDir)
	}
	return nil
}

func syncKubeConfig(cmdRslv exec_info.InfoExecCreator, config string, kubeconfigLocation string, containerInfo *model.ContainerInfo) error {
	infoExec := cmdRslv.CreateInfoExec([]string{"sh", "-c", "echo \"" + config + "\" > " + kubeconfigLocation + "/config"}, containerInfo)
	if err := infoExec.Start(); err != nil {
		logrus.Debugf("Error is not available in %s/%s. Error: %s", containerInfo.PodName, containerInfo.ContainerName, err.Error())
		return errors.New("Could not write kubeconfig to: " + kubeconfigLocation)
	}
	return nil
}
