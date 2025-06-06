// main.go
// Nunu项目主入口文件

package main

//Skaarl
import (
	"Skaarl/cmd/skaarl"
	"fmt"
)

// main 项目主函数
// 执行nunu命令行工具的主逻辑
// 返回错误时会打印错误信息
func main() {
	err := skaarl.Execute()
	if err != nil {
		fmt.Println("execute error: ", err.Error())
	}
}
