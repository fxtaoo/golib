// 阿里云相关
package aliyun

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type OSSConf struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type OSSFile struct {
	BucketName     string `json:"bucketName"`
	BucketFilePath string `json:"bucketFilePath"`
	LoadFilePath   string `json:"loadFilePath"`
}

// 返回 OSSClient
func (c *OSSConf) OSSClient() (*oss.Client, error) {
	client, err := oss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 上传文件（文件存在覆盖）
func (f *OSSFile) UploadFile(conf *OSSConf) error {
	client, err := conf.OSSClient()
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(f.BucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObjectFromFile(f.BucketFilePath, f.LoadFilePath)
	if err != nil {
		return err
	}

	return nil
}
