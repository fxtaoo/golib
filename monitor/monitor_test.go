// 监控相关
package monitor

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(StartProcess("echo", "echo"))
}
