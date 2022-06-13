// 文件读写
package gofile

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// 从 toml 文件读配置
// str 为配置文件名，从执行文件当前目录读取
// str 为绝对路径，从路径读取
func TomlFileRead(str string, v interface{}) {
	if !filepath.IsAbs(str) {
		str = filepath.Join(filepath.Dir(os.Args[0]), str)
	}
	if _, err := toml.DecodeFile(str, v); err != nil {
		panic(err)
	}
}
