// 文件相关
package gofile

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

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

// 追加字符串到文件
func AppendFile(filepath, content string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	newLine := content
	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		return err
	}

	return nil
}

const tmpFileNamePrefix = "gofileTmpFile"

// 下载文件（缺省下载目录未临时文件夹）
func DownloadFile(url string, filePath string) (string, error) {
	var file *os.File
	var err error

	if filePath == "" {
		file, err = ioutil.TempFile("", tmpFileNamePrefix)
		if err != nil {
			return "", err
		}
	} else {
		file, err = os.Create(filePath)
		if err != nil {
			return "", err
		}
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

// 文件 MD5
func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 比较线上文件与本地文件是否一致
func OnlineLocalMD5Same(url, localFilePath string) (bool, error) {

	localMD5, err := FileMD5(localFilePath)
	if err != nil {
		return false, err
	}

	onlineFilePath, err := DownloadFile(url, "")
	if err != nil {
		return false, err
	}

	defer func(onlineFilePath string) {
		// 删除临时文件
		if strings.HasPrefix(path.Base(onlineFilePath), tmpFileNamePrefix) {
			os.Remove(onlineFilePath)
		}
	}(onlineFilePath)

	onlineMD5, err := FileMD5(onlineFilePath)
	if err != nil {
		return false, err
	}

	if localMD5 == onlineMD5 {
		return true, nil
	}

	return false, nil
}
