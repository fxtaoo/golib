package file

import "github.com/BurntSushi/toml"

// 从 toml 文件初始化 v 值
func TomlInitValue(filePath string, v interface{}) {
	confFilePath := AbsPath(filePath)
	if _, err := toml.DecodeFile(confFilePath, v); err != nil {
		panic(err)
	}
}
