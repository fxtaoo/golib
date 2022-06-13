// 文件读写
package gofile

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// 返回当前目录下文件名为 str 文件的绝对路径
func IsAbsStr(str string) string {
	if filepath.IsAbs(str) {
		return str
	}
	return filepath.Join(filepath.Dir(os.Args[0]), str)
}

// 从 toml 文件读配置
func TomlFileRead(str string, v interface{}) {
	str = IsAbsStr(str)
	if _, err := toml.DecodeFile(str, v); err != nil {
		panic(err)
	}
}

// 从当前目录文件名为 str CSV文件读数据
// 接收数据列表，改数据列表数据追加函数
func CSVFileRead(str string) [][]string {
	str = IsAbsStr(str)

	csvFile, err := os.Open(str)
	if err != nil {
		log.Println(str, err)
		panic(err)
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
			log.Println("csvFile", err)
			panic(err)
		}
		// 追加
		csvdata = append(csvdata, row)
	}
	return csvdata
}
