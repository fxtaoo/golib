// 七牛云相关
package qiniu

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/cdn"
)

type Key struct {
	AccessKey string
	SecretKey string
}

func CdnManager(key Key) *cdn.CdnManager {
	mac := qbox.NewMac(key.AccessKey, key.SecretKey)
	return cdn.NewCdnManager(mac)
}

// 刷新文件，URL不做验证
func UrlsRefresh(cdn *cdn.CdnManager, urls []string) (cdn.RefreshResp, error) {
	ret, err := cdn.RefreshUrls(urls)
	return ret, err
}
