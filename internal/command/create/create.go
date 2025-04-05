// create.go
// Nunu命令行工具create子命令实现
// 负责创建项目组件(handler/service/repository/model)的逻辑

package create

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/go-nunu/nunu/tpl"
	"github.com/spf13/cobra"
)

// Create 创建组件结构体
// 包含创建项目组件所需的信息
type Create struct {
	ProjectName          string // 项目名称
	CreateType           string // 创建类型(handler/service/repository/model)
	FilePath             string // 文件路径
	FileName             string // 文件名
	StructName           string // 结构体名称(首字母大写)
	StructNameLowerFirst string // 结构体名称(首字母小写)
	StructNameFirstChar  string // 结构体名称首字符
	StructNameSnakeCase  string // 结构体名称(snake_case)
	IsFull               bool   // 是否创建完整组件(all)
}

// NewCreate 创建Create实例
// 返回一个初始化的Create指针
func NewCreate() *Create {
	return &Create{}
}

// CmdCreate create子命令定义
// 用于创建项目组件
var CmdCreate = &cobra.Command{
	Use:     "create [type] [handler-name]",                  // 命令格式
	Short:   "Create a new handler/service/repository/model", // 简短描述
	Example: "nunu create handler user",                      // 使用示例
	Args:    cobra.ExactArgs(2),                              // 参数数量
	Run: func(cmd *cobra.Command, args []string) {
		// 主逻辑由子命令实现
	},
}
var (
	tplPath string
)

func init() {
	CmdCreateHandler.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateService.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateRepository.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateModel.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateAll.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")

}

// CmdCreateHandler 创建handler子命令
// 用于创建handler组件
var CmdCreateHandler = &cobra.Command{
	Use:     "handler",                  // 命令名称
	Short:   "Create a new handler",     // 简短描述
	Example: "nunu create handler user", // 使用示例
	Args:    cobra.ExactArgs(1),         // 参数数量
	Run:     runCreate,                  // 执行函数
}
var CmdCreateService = &cobra.Command{
	Use:     "service",
	Short:   "Create a new service",
	Example: "nunu create service user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateRepository = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new repository",
	Example: "nunu create repository user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateModel = &cobra.Command{
	Use:     "model",
	Short:   "Create a new model",
	Example: "nunu create model user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateAll = &cobra.Command{
	Use:     "all",
	Short:   "Create a new handler & service & repository & model",
	Example: "nunu create all user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

// runCreate 创建组件主逻辑
// 根据类型创建对应的项目组件
func runCreate(cmd *cobra.Command, args []string) {
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	c.CreateType = cmd.Use
	c.FilePath, c.StructName = filepath.Split(args[0])
	c.FileName = strings.ReplaceAll(c.StructName, ".go", "")
	c.StructName = strutil.UpperFirst(strutil.CamelCase(c.FileName))
	c.StructNameLowerFirst = strutil.LowerFirst(c.StructName)
	c.StructNameFirstChar = string(c.StructNameLowerFirst[0])
	c.StructNameSnakeCase = strutil.SnakeCase(c.StructName)

	switch c.CreateType {
	case "handler", "service", "repository", "model":
		c.genFile()
	case "all":

		c.CreateType = "handler"
		c.genFile()

		c.CreateType = "service"
		c.genFile()

		c.CreateType = "repository"
		c.genFile()

		c.CreateType = "model"
		c.genFile()
	default:
		log.Fatalf("Invalid handler type: %s", c.CreateType)
	}

}

// genFile 生成组件文件
// 根据模板生成对应的组件文件
func (c *Create) genFile() {
	filePath := c.FilePath
	if filePath == "" {
		filePath = fmt.Sprintf("internal/%s/", c.CreateType)
	}
	f := createFile(filePath, strings.ToLower(c.FileName)+".go")
	if f == nil {
		log.Printf("warn: file %s%s %s", filePath, strings.ToLower(c.FileName)+".go", "already exists.")
		return
	}
	defer f.Close()
	var t *template.Template
	var err error
	if tplPath == "" {
		t, err = template.ParseFS(tpl.CreateTemplateFS, fmt.Sprintf("create/%s.tpl", c.CreateType))
	} else {
		t, err = template.ParseFiles(path.Join(tplPath, fmt.Sprintf("%s.tpl", c.CreateType)))
	}
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	err = t.Execute(f, c)
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	log.Printf("Created new %s: %s", c.CreateType, filePath+strings.ToLower(c.FileName)+".go")

}

// createFile 创建文件
// 在指定路径创建文件，如果文件已存在则返回nil
// 返回创建的文件指针
func createFile(dirPath string, filename string) *os.File {
	filePath := filepath.Join(dirPath, filename)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
	stat, _ := os.Stat(filePath)
	if stat != nil {
		return nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}

	return file
}
