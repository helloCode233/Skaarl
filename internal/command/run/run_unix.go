//go:build !plan9 && !windows
// +build !plan9,!windows

// run包实现了Unix平台下的文件监控和自动重启功能
package run

import (
	// 标准库导入
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"Skaarl/config"

	"github.com/AlecAivazis/survey/v2"

	"Skaarl/internal/pkg/helper"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// quit 用于接收系统信号的通知通道
var quit = make(chan os.Signal, 1)

// Run 结构体定义了运行命令的相关功能
type Run struct {
}

// excludeDir 定义要排除监控的目录列表，多个目录用逗号分隔
var excludeDir string

// includeExt 定义要监控的文件扩展名列表，多个扩展名用逗号分隔
var includeExt string

// init 函数初始化run命令的参数
func init() {
	// 设置排除目录参数
	CmdRun.Flags().StringVarP(&excludeDir, "excludeDir", "", excludeDir, `eg: skaarl run --excludeDir="tmp,vendor,.git,.idea"`)
	// 设置包含文件扩展名参数
	CmdRun.Flags().StringVarP(&includeExt, "includeExt", "", includeExt, `eg: skaarl run --includeExt="go,tpl,tmpl,html,yaml,yml,toml,ini,json"`)
	// 如果没有设置参数，则使用默认配置
	if excludeDir == "" {
		excludeDir = config.RunExcludeDir
	}
	if includeExt == "" {
		includeExt = config.RunIncludeExt
	}
}

// CmdRun 定义了skaarl run命令
var CmdRun = &cobra.Command{
	Use:     "run",                       // 命令名称
	Short:   "skaarl run [main.go path]", // 简短描述
	Long:    "skaarl run [main.go path]", // 详细描述
	Example: "skaarl run cmd/server",     // 使用示例
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs, programArgs := helper.SplitArgs(cmd, args)
		var dir string
		if len(cmdArgs) > 0 {
			dir = cmdArgs[0]
		}
		base, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if dir == "" {
			cmdPath, err := helper.FindMain(base, excludeDir)

			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
				return
			}
			switch len(cmdPath) {
			case 0:
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current directory")
				return
			case 1:
				for _, v := range cmdPath {
					dir = v
				}
			default:
				var cmdPaths []string
				for k := range cmdPath {
					cmdPaths = append(cmdPaths, k)
				}
				sort.Strings(cmdPaths)

				prompt := &survey.Select{
					Message:  "Which directory do you want to run?",
					Options:  cmdPaths,
					PageSize: 10,
				}
				e := survey.AskOne(prompt, &dir)
				if e != nil || dir == "" {
					return
				}
				dir = cmdPath[dir]
			}
		}
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		fmt.Printf("\033[35mNunu run %s.\033[0m\n", dir)
		fmt.Printf("\033[35mWatch excludeDir %s\033[0m\n", excludeDir)
		fmt.Printf("\033[35mWatch includeExt %s\033[0m\n", includeExt)
		watch(dir, programArgs)

	},
}

// watch 函数监控文件变化并自动重启程序
func watch(dir string, programArgs []string) {
	// 监控当前目录
	watchPath := "./"

	// 创建文件监控器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer watcher.Close()

	// 处理排除目录和包含扩展名配置
	excludeDirArr := strings.Split(excludeDir, ",")
	includeExtArr := strings.Split(includeExt, ",")
	includeExtMap := make(map[string]struct{})
	for _, s := range includeExtArr {
		includeExtMap[s] = struct{}{}
	}

	// 遍历目录添加监控文件
	err = filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过排除目录
		for _, s := range excludeDirArr {
			if s == "" {
				continue
			}
			if strings.HasPrefix(path, s) {
				return nil
			}
		}
		// 只监控指定扩展名的文件
		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			if _, ok := includeExtMap[strings.TrimPrefix(ext, ".")]; ok {
				err = watcher.Add(path)
				if err != nil {
					fmt.Println("Error:", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 启动初始程序
	cmd := start(dir, programArgs)
	flag_ := true
	// 主监控循环
	for {
		select {
		case <-quit:
			// 处理退出信号
			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
			if err != nil {
				fmt.Printf("\033[31mserver exiting...\033[0m\n")
				return
			}
			fmt.Printf("\033[31mserver exiting...\033[0m\n")
			os.Exit(0)

		case event := <-watcher.Events:
			// 文件创建、修改或删除时重启程序
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove {
				if flag_ {
					fmt.Printf("\033[36mfile modified: %s\033[0m\n", event.Name)
					flag_ = false
					syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
					cmd = start(dir, programArgs)
				}
			}
		case err := <-watcher.Errors:
			// 处理监控错误
			fmt.Println("Error:", err)
		}
	}
}

// start 启动新的go程序进程
func start(dir string, programArgs []string) *exec.Cmd {
	// 构造go run命令
	cmd := exec.Command("go", append([]string{"run", dir}, programArgs...)...)
	// 设置新的进程组以便终止时能杀死所有子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 重定向标准输出和错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动进程
	err := cmd.Start()
	if err != nil {
		log.Fatalf("\033[33;1mcmd run failed\u001B[0m")
	}
	// 等待1秒确保程序启动
	time.Sleep(time.Second)
	fmt.Printf("\033[32;1mrunning...\033[0m\n")
	return cmd
}
