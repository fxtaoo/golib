package ssh

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Host          string            `json:"host"`
	Port          int               `json:"port"`
	User          string            `json:"user"`
	Password      string            `json:"password"`
	KeyPath       string            `json:"keyPath"`
	KeyPathPasswd string            `json:"keyPathPasswd"`
	KeyStrByte    []byte            // 字节密钥
	clientConfig  *ssh.ClientConfig // ssh 客户端配置
}

// 初始化 ssh 客户端配置
// 密钥存在，优先使用密钥
func (gs *SSH) Init() error {
	gs.clientConfig = &ssh.ClientConfig{
		User:            gs.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 3,
	}

	// KeyPath KeyStrByte 同时存在，优先使用 gs.KeyStrByte
	// 否则使用 KeyPath
	if gs.KeyStrByte == nil && gs.KeyPath != "" {
		gs.KeyStrByte, _ = os.ReadFile(gs.KeyPath)
	}

	//	密钥方式
	if gs.KeyStrByte != nil {
		// 密钥存在密码
		if gs.KeyPathPasswd != "" {
			signer, err := ssh.ParsePrivateKeyWithPassphrase(gs.KeyStrByte, []byte(gs.KeyPathPasswd))
			if err == nil {
				gs.clientConfig.Auth = append(gs.clientConfig.Auth, ssh.PublicKeys(signer))
			}
		} else {

			signer, err := ssh.ParsePrivateKey(gs.KeyStrByte)
			if err == nil {
				gs.clientConfig.Auth = append(gs.clientConfig.Auth, ssh.PublicKeys(signer))
			}
		}
	}

	if gs.Password != "" {
		if gs.Password != "" {
			gs.clientConfig.Auth = append(gs.clientConfig.Auth, ssh.Password(gs.Password))
		}
	}

	if len(gs.clientConfig.Auth) == 0 {
		return fmt.Errorf("no auth method")
	}
	return nil
}
