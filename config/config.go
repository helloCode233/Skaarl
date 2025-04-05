// config.go
// Nunu项目配置文件
// 包含项目版本、依赖命令和仓库布局等配置信息

package config

var (
	// Version 项目版本号
	Version = "1.1.1"
	// WireCmd wire依赖注入工具命令
	WireCmd = "github.com/google/wire/cmd/wire@latest"
	// NunuCmd nunu命令行工具命令
	NunuCmd = "github.com/go-skaarl/skaarl@latest"
	// RepoBase 基础项目模板仓库
	RepoBase = "https://github.com/helloCode233/nunu-layout-basic.git"
	// RepoAdvanced 高级项目模板仓库
	RepoAdvanced = "https://github.com/go-nunu/nunu-layout-advanced.git"
	// RepoAdmin 后台管理项目模板仓库
	RepoAdmin = "https://github.com/go-nunu/nunu-layout-admin.git"
	// RepoChat 聊天项目模板仓库
	RepoChat = "https://github.com/go-nunu/nunu-layout-chat.git"
	// RunExcludeDir 运行时要排除的目录
	RunExcludeDir = ".git,.idea,tmp,vendor"
	// RunIncludeExt 运行时要包含的文件扩展名
	RunIncludeExt = "go,html,yaml,yml,toml,ini,json,xml,tpl,tmpl"
)

type Configuration struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	//Storage  Storage  `mapstructure:"storage" json:"storage" yaml:"storage"`
}
