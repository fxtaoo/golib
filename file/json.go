package file

import (
	"encoding/json"
	"os"
)

// Json 文件初始化 v 值
func JsonInitValue(filePath string, v interface{}) error {
	confFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(confFile, &v)
	if err != nil {
		return err
	}
	return nil
}
