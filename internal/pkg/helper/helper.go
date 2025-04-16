package helper

import (
	"Skaarl/tpl"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

func GetProjectName(dir string) string {
	modFile, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}
func SplitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}
func FindMain(base, excludeDir string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}
	excludeDirArr := strings.Split(excludeDir, ",")
	cmdPath := make(map[string]string)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, s := range excludeDirArr {
			if strings.HasPrefix(strings.TrimPrefix(path, base), "/"+s) {
				return nil
			}
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if !strings.Contains(string(content), "package main") {
				return nil
			}
			re := regexp.MustCompile(`func\s+main\s*\(`)
			if re.Match(content) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				d, _ := filepath.Split(absPath)
				cmdPath[strings.TrimPrefix(absPath, wd)] = d

			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cmdPath, nil
}

// createFile 创建文件
// 在指定路径创建文件，如果文件已存在则返回nil
// 返回创建的文件指针
func CreateFile(dirPath string, filename string) *os.File {
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

// genFile 生成组件文件
// 根据模板生成对应的组件文件
func GenFile(filePath, CreateType, FileName, tplPath string, data any) {

	if filePath == "" {
		filePath = fmt.Sprintf("internal/%s/", CreateType)
	}
	f := CreateFile(filePath, strings.ToLower(FileName)+".go")
	if f == nil {
		log.Printf("warn: file %s%s %s", filePath, strings.ToLower(FileName)+".go", "already exists.")
		return
	}
	defer f.Close()
	var t *template.Template
	var err error
	switch tplPath {
	case "create":
		t, err = template.ParseFS(tpl.CreateTemplateFS, fmt.Sprintf("create/%s.tpl", CreateType))
	case "run":
		t, err = template.ParseFS(tpl.RunTemplateFS, fmt.Sprintf("run/%s.tpl", CreateType))
	default:
		t, err = template.ParseFiles(path.Join(tplPath, fmt.Sprintf("%s.tpl", CreateType)))
	}
	//if tplPath == "" {
	//	t, err = template.ParseFS(tpl.CreateTemplateFS, fmt.Sprintf("create/%s.tpl", CreateType))
	//} else {
	//	t, err = template.ParseFiles(path.Join(tplPath, fmt.Sprintf("%s.tpl", CreateType)))
	//}
	if err != nil {
		log.Fatalf("create %s error: %s", CreateType, err.Error())
	}
	err = t.Execute(f, data)
	if err != nil {
		log.Fatalf("create %s error: %s", CreateType, err.Error())
	}
	log.Printf("Created new %s: %s", CreateType, filePath+strings.ToLower(FileName)+".go")

}
func CapitalizeFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	// 将首字母转为大写，剩余部分保持原样
	return strings.ToUpper(string(str[0])) + str[1:]
}
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// 存在：返回存在状态，是否为目录，无错误
		return true
	}

	if os.IsNotExist(err) {
		// 明确的不存在：返回 false, 无需检查目录
		return false
	}

	// 其他错误（权限不足、路径非法等）
	return false
}
