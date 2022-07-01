// 阿里云相关
package goaliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	Endpoint        string // oss 地域节点
	AccessKeyId     string // 访问 key
	AccessKeySecret string // 访问 secret
}

type File struct {
	BucketName     string // oss bucket 名称
	BucketFilePath string // oss 文件路径
	LoadFilePath   string // 本地文件路径
}

// 返回 OSSClient
func OSSClient(a *OSS) (*oss.Client, error) {
	client, err := oss.New(a.Endpoint, a.AccessKeyId, a.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 上传文件
func UploadFile(client *oss.Client, file *File) error {
	bucket, err := client.Bucket(file.BucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(file.BucketFilePath, file.LoadFilePath)
	if err != nil {
		return err
	}

	return nil
}
