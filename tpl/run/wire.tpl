//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/gin-gonic/gin"
    {{ range .ImportList -}}
    "{{.}}"
    {{ end -}}
	"github.com/google/wire"
	"github.com/spf13/viper"
)

{{ range .SetList }}
var {{.Name -}}Set = wire.NewSet(
    {{range .News}}
	{{.Name}}.{{.Func}},
	{{ end }}
)
{{ end }}

func NewWire(*viper.Viper, *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
	{{ range .SetList }}
	    {{.SetName -}}Set,
	{{ end }}
	))
}
