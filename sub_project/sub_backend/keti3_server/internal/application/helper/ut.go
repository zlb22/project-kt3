package helper

import (
	"fmt"
	"os"
	"strings"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

// InitUnitTest 初始化单元测试
func InitUnitTest() {
	current, _ := os.Getwd()
	fmt.Println(current)

	dirPath := strings.Split(current, "keti3")
	if len(dirPath) == 0 {
		panic("app name not bus-mall-admin")
	}

	pontos.AppRootPath = dirPath[0] + "keti3"
	if err := pontos.InitApp(); err != nil {
		panic(err)
	}
}
