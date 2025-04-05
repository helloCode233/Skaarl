// root.go
// Nunu命令行工具根命令定义文件
// 包含主命令定义和子命令注册逻辑

package skaarl

import (
	"Skaarl/config"
	"Skaarl/internal/command/create"
	"Skaarl/internal/command/new"
	"Skaarl/internal/command/run"
	"Skaarl/internal/command/upgrade"
	"Skaarl/internal/command/wire"
	"fmt"
	"github.com/spf13/cobra"
)

// CmdRoot 根命令定义
// 包含Nunu CLI工具的:
// - 使用说明(Use)
// - 示例(Example)
// - 简短描述(Short)
// - 版本信息(Version)
var CmdRoot = &cobra.Command{
	Use:     "skaarl",
	Example: "skaarl new demo-api",
	Short:   fmt.Sprintf("   _____   _                             _ \n  / ____| | |                           | |\n | (___   | | __   __ _    __ _   _ __  | |\n  \\___ \\  | |/ /  / _` |  / _` | | '__| | |\n  ____) | |   <  | (_| | | (_| | | |    | |\n |_____/  |_|\\_\\  \\__,_|  \\__,_| |_|    |_|\n                                           \n " + "\x1B[38;2;66;211;146mA\x1B[39m \x1B[38;2;67;209;149mC\x1B[39m\x1B[38;2;68;206;152mL\x1B[39m\x1B[38;2;69;204;155mI\x1B[39m \x1B[38;2;70;201;158mt\x1B[39m\x1B[38;2;71;199;162mo\x1B[39m\x1B[38;2;72;196;165mo\x1B[39m\x1B[38;2;73;194;168ml\x1B[39m \x1B[38;2;74;192;171mf\x1B[39m\x1B[38;2;75;189;174mo\x1B[39m\x1B[38;2;76;187;177mr\x1B[39m \x1B[38;2;77;184;180mb\x1B[39m\x1B[38;2;78;182;183mu\x1B[39m\x1B[38;2;79;179;186mi\x1B[39m\x1B[38;2;80;177;190ml\x1B[39m\x1B[38;2;81;175;193md\x1B[39m\x1B[38;2;82;172;196mi\x1B[39m\x1B[38;2;83;170;199mn\x1B[39m\x1B[38;2;83;167;202mg\x1B[39m \x1B[38;2;84;165;205mg\x1B[39m\x1B[38;2;85;162;208mo\x1B[39m \x1B[38;2;86;160;211ma\x1B[39m\x1B[38;2;87;158;215mp\x1B[39m\x1B[38;2;88;155;218ml\x1B[39m\x1B[38;2;89;153;221mi\x1B[39m\x1B[38;2;90;150;224mc\x1B[39m\x1B[38;2;91;148;227ma\x1B[39m\x1B[38;2;92;145;230mt\x1B[39m\x1B[38;2;93;143;233mi\x1B[39m\x1B[38;2;94;141;236mo\x1B[39m\x1B[38;2;95;138;239mn\x1B[39m\x1B[38;2;96;136;243m.\x1B[39m"),
	Version: fmt.Sprintf(":\n   _____   _                             _ \n  / ____| | |                           | |\n | (___   | | __   __ _    __ _   _ __  | |\n  \\___ \\  | |/ /  / _` |  / _` | | '__| | |\n  ____) | |   <  | (_| | | (_| | | |    | |\n |_____/  |_|\\_\\  \\__,_|  \\__,_| |_|    |_|\n                                            \n  Skaarl %s - Copyright (c) 2025 Skaarl\n  Released under the MIT License.\n\n", config.Version),
}

// init 初始化函数
// 注册所有子命令到根命令:
// - new: 创建新项目
// - create: 创建项目组件
// - run: 运行项目
// - upgrade: 升级工具
// - wire: 依赖注入
func init() {
	CmdRoot.AddCommand(new.CmdNew)
	CmdRoot.AddCommand(create.CmdCreate)
	CmdRoot.AddCommand(run.CmdRun)

	CmdRoot.AddCommand(upgrade.CmdUpgrade)
	create.CmdCreate.AddCommand(create.CmdCreateHandler)
	create.CmdCreate.AddCommand(create.CmdCreateService)
	create.CmdCreate.AddCommand(create.CmdCreateRepository)
	create.CmdCreate.AddCommand(create.CmdCreateModel)
	create.CmdCreate.AddCommand(create.CmdCreateAll)

	CmdRoot.AddCommand(wire.CmdWire)
	wire.CmdWire.AddCommand(wire.CmdWireAll)
}

// Execute 执行根命令
// 启动Nunu命令行工具
// 返回错误时会传递到main函数处理
func Execute() error {
	return CmdRoot.Execute()
}
