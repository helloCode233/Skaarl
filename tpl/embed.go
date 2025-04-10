package tpl

import "embed"

//go:embed create/*.tpl
var CreateTemplateFS embed.FS

//go:embed run/*.tpl
var RunTemplateFS embed.FS
