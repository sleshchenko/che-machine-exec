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
	"github.com/eclipse/che-machine-exec/activity"
	"github.com/eclipse/che-machine-exec/auth"
	restUtil "github.com/eclipse/che-machine-exec/common/rest"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleActivityTick(c *gin.Context, manager *activity.Manager) {
	if auth.IsEnabled() {
		_, err := auth.Authenticate(c)
		if err != nil {
			restUtil.WriteErrorResponse(c, err)
			return
		}
	}

	// at this point, it's just stub handler that does nothing
	// but a bit later ActivityManager will appear and register the latest activity
	// to post pone workspace stopping by idle timeout
	manager.Tick()
	c.Writer.WriteHeader(http.StatusNoContent)
	return
}
