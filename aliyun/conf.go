package aliyun

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type Conf struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

// 返回 OSSClient
func (c *Conf) OSSClient() (*oss.Client, error) {
	client, err := oss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}
