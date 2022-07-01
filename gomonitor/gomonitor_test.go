// 监控相关
package gomonitor

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(StartProcess("echo", "echo"))
}
