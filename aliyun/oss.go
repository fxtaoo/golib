// 阿里云相关
package aliyun

type OSSFile struct {
	BucketName     string `json:"bucketName"`
	BucketFilePath string `json:"bucketFilePath"`
	LoadFilePath   string `json:"loadFilePath"`
}

// 上传文件（文件存在覆盖）
func (f *OSSFile) UploadFile(conf *Conf) error {
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
