//
// Copyright (c) 2012-2019 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package activity

import (
	"errors"
	"github.com/eclipse/che-machine-exec/client"
	"github.com/eclipse/che-machine-exec/exec"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"time"
)

const (
	//TODO Make configurable and optional
	idleTimeout     = 1 * time.Minute
	stopRetryPeriod = 5 * time.Minute
)

type Manager struct {
	workspaceName string
	namespace     string
	activityC     chan bool
}

func New() (*Manager, error) {
	namespace := exec.GetNamespace()
	if namespace == "" {
		return nil, errors.New("unable to evaluate the current namespace that is needed for activity manager works correctly")
	}

	workspaceName, isFound := os.LookupEnv("CHE_WORKSPACE_NAME")
	if !isFound {
		return nil, errors.New("CHE_WORKSPACE_NAME env must be set for activity manager works correctly")
	}

	return &Manager{
		namespace:     namespace,
		workspaceName: workspaceName,
	}, nil
}

var (
	WorkspaceGroupVersionResource = schema.GroupVersionResource{
		Group:    "workspace.che.eclipse.org",
		Version:  "v1alpha1",
		Resource: "workspaces",
	}
)

func (m *Manager) Tick() {
	m.activityC <- true
}

func (m *Manager) Start() {
	if m.activityC != nil {
		//it's already started
		return
	}
	logrus.Infof("Activity tracker is run and workspace will be stopped in %t if there is no activity", idleTimeout)
	m.activityC = make(chan bool)
	// TODO Review if using timer is the best way
	// TODO Probably ticker could be faster but than we'll get correlation, so workspace will be started in idleTimeout +- tick period
	timer := time.NewTimer(idleTimeout)
	select {
	case <-timer.C:
		if err := m.stopWorkspace(); err != nil {
			timer.Reset(stopRetryPeriod)
			logrus.Errorf("Failed to stop workspace. Will retry in %t", stopRetryPeriod, err)
		} else {
			logrus.Info("Workspace is successfully stopped by inactivity. Bye")
			return
		}
	case <-m.activityC:
		logrus.Debug("Activity is reported. Resetting timer")
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(idleTimeout)
	}
}

func (m *Manager) stopWorkspace() error {
	c, err := client.NewDynamicInCluster()
	if err != nil {
		return err
	}

	stopWorkspacePath := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": map[string]interface{}{
				"started": false,
			},
		},
	}
	jsonPath, err := stopWorkspacePath.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = c.Resource(WorkspaceGroupVersionResource).Namespace(m.namespace).Patch(m.workspaceName, types.MergePatchType, jsonPath, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}
