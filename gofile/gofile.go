// 文件读写
package gofile

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// 返回绝对路径
// str 为当前目录下文件名或文件绝对路径
func AbsPath(str string) string {
	if filepath.IsAbs(str) {
		return str
	}
	return filepath.Join(filepath.Dir(os.Args[0]), str)
}

// 从 toml 文件读配置
// str 为当前目录下 toml 文件名或 toml 文件绝对路径
func TomlFileRead(str string, v interface{}) {
	confFilePath := AbsPath(str)
	if _, err := toml.DecodeFile(confFilePath, v); err != nil {
		panic(err)
	}
}

// CSV 文件读数据
// str 为当前目录下 CSV 文件名或 CSV 文件绝对路径
func CSVFileRead(str string) ([][]string, error) {
	str = AbsPath(str)

	csvFile, err := os.Open(str)
	if err != nil {
		return nil, err
	}

	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.LazyQuotes = true

	var csvdata [][]string
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// 追加
		csvdata = append(csvdata, row)
	}
	return csvdata, nil
}
