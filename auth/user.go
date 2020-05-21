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

package auth

import (
	"errors"
	"github.com/eclipse/che-machine-exec/client"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	UserGroupVersionResource = schema.GroupVersionResource{
		Group:    "user.openshift.io",
		Version:  "v1",
		Resource: "users",
	}
)

func getCurrentUserID(token string) (string, error) {
	c, err := client.NewDynamicForToken(token)
	if err != nil {
		return "", err
	}

	userInfo, err := c.Resource(UserGroupVersionResource).Get("~", metav1.GetOptions{})
	if err != nil {
		return "", errors.New("Failed to retrieve the current user info. Cause: " + err.Error())
	}

	logrus.Debugf("Current user info %s, %s", userInfo.GetUID(), userInfo.GetName())
	return string(userInfo.GetUID()), nil
}