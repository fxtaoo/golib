// 文件相关
package file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 文件绝对路径
func FileAbsPath(filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	return filepath.Join(filepath.Dir(os.Args[0]), filePath)
}

// 文件末尾追加内容
func FileAppendContent(filePath, content string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

// 下载文件（缺省下载目录为临时文件夹）
func FileDownload(url string, tmpFilePath string) (string, error) {
	var file *os.File
	var err error

	if tmpFilePath == "" {
		file, err = os.CreateTemp("", TmpFileNamePrefix)
		if err != nil {
			return "", err
		}
	} else {
		file, err = os.Create(tmpFilePath)
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

// 删除临时文件
func FileTmpDel(filePath string) {
	// 检查是不是临时文件
	if strings.HasPrefix(path.Base(filePath), TmpFileNamePrefix) {
		os.Remove(filePath)
	}
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

// 文件本地线上比较
func FileLocalOnline(localFilePath, url string) (bool, error) {

	localMD5, err := FileMD5(localFilePath)
	if err != nil {
		return false, err
	}

	onlineFilePath, err := FileDownload(url, "")
	if err != nil {
		return false, err
	}

	defer FileTmpDel(onlineFilePath)

	onlineMD5, err := FileMD5(onlineFilePath)
	if err != nil {
		return false, err
	}

	if localMD5 == onlineMD5 {
		return true, nil
	}

	return false, nil
}
