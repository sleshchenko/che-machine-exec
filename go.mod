module github.com/eclipse/che-machine-exec

go 1.13

replace (
	github.com/eclipse/che-go-jsonrpc => github.com/eclipse/che-go-jsonrpc v0.0.0-20200317130110-931966b891fe
	github.com/gin-contrib/sse => github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v0.0.0-20180501062418-bd4f73af679e
	github.com/google/gofuzz => github.com/google/gofuzz v0.0.0-20161122191042-44d81051d367
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
	github.com/gorilla/websocket => github.com/gorilla/websocket v0.0.0-20151102191034-361d4c0ffd78
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.3
	github.com/pkg/errors => github.com/pkg/errors v0.0.0-20161002052512-839d9e913e06
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify => github.com/stretchr/testify v1.2.1
	github.com/ugorji/go => github.com/ugorji/go v1.1.1
	golang.org/x/sys => golang.org/x/sys v0.0.0-20171031081856-95c657629925
	gopkg.in/go-playground/validator.v8 => gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.0.0-20170721113624-670d4cfef054
)

require (
	github.com/eclipse/che-go-jsonrpc v0.0.0-00010101000000-000000000000
	github.com/elazarl/goproxy v0.0.0-20200426045556-49ad98f6dac1 // indirect
	github.com/gin-contrib/sse v0.0.0-00010101000000-000000000000 // indirect
	github.com/gin-gonic/gin v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v0.0.0-00010101000000-000000000000
	github.com/mattn/go-isatty v0.0.0-00010101000000-000000000000 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0
	github.com/ugorji/go v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
)
