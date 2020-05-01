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

package rest

import (
	"net/http"

	"github.com/eclipse/che-machine-exec/exec"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	execManager = exec.GetExecManager()
)

func HandleKubeConfig(request *http.Request, response gin.ResponseWriter) {
	token := request.Header.Get("X-Forwarded-Access-Token")
	if token == "" {
		response.WriteHeader(http.StatusUnauthorized)
		_, err := response.Write([]byte("Authorization token must not be empty"))
		if err != nil {
			logrus.Error("Failed to write error response", err)
		}
	}

	err := execManager.CreateKubeConfig(token)

	if err != nil {
		logrus.Errorf("Unable to create kubeconfig. Cause: %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, err := response.Write([]byte(err.Error()))
		if err != nil {
			logrus.Error("Failed to write error response", err)
		}
		return
	}
}
