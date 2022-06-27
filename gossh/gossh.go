// ssh
package gossh

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type GoSSH struct {
	Host          string            // 主机
	Port          int               // 端口
	User          string            // 用户
	Password      string            // 密码
	KeyPath       string            // 密钥路径
	KeyPathPasswd string            // 密钥密码
	KeyStrByte    []byte            // 字节密钥
	clientConfig  *ssh.ClientConfig // ssh 客户端配置
}

// 初始化 ssh 客户端配置
// 密钥存在，优先使用密钥
func (gossh *GoSSH) Init() error {
	gossh.clientConfig = &ssh.ClientConfig{}
	gossh.clientConfig.User = gossh.User
	gossh.clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	gossh.clientConfig.Timeout = time.Second * 5

	// KeyPath KeyStrByte 同时存在，优先使用 gossh.KeyStrByte
	// 否则使用 KeyPath
	if gossh.KeyStrByte == nil && gossh.KeyPath != "" {
		gossh.KeyStrByte, _ = ioutil.ReadFile(gossh.KeyPath)
	}

	//	密钥方式
	if gossh.KeyStrByte != nil {
		// 密钥存在密码
		if gossh.KeyPathPasswd != "" {
			signer, err := ssh.ParsePrivateKeyWithPassphrase(gossh.KeyStrByte, []byte(gossh.KeyPathPasswd))
			if err == nil {
				gossh.clientConfig.Auth = append(gossh.clientConfig.Auth, ssh.PublicKeys(signer))
			}
		} else {

			signer, err := ssh.ParsePrivateKey(gossh.KeyStrByte)
			if err == nil {
				gossh.clientConfig.Auth = append(gossh.clientConfig.Auth, ssh.PublicKeys(signer))
			}
		}
	}

	if gossh.Password != "" {
		if gossh.Password != "" {
			gossh.clientConfig.Auth = append(gossh.clientConfig.Auth, ssh.Password(gossh.Password))
		}
	}

	if len(gossh.clientConfig.Auth) == 0 {
		return fmt.Errorf("no auth method")
	}
	return nil
}

// 连接 Session
func (gossh *GoSSH) Connect() (*ssh.Session, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", gossh.Host, gossh.Port), gossh.clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil

}

// 运行命令
func (gossh *GoSSH) RunCmd(cmd string) (string, error) {
	session, err := gossh.Connect()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// 运行一组命令，顺序执行
func (gossh *GoSSH) RunCmdsSequential(cmds []string) ([]string, error) {
	var outs []string

	for _, cmd := range cmds {
		str, err := gossh.RunCmd(cmd)
		if err != nil {
			return nil, err
		}

		outs = append(outs, str)
	}

	return outs, nil
}

// 运行一组命令，并行执行，结果顺序与命令顺序一致
func (gossh *GoSSH) RunCmdsParallel(cmds []string) ([]string, error) {
	var wg sync.WaitGroup
	cmdNum := len(cmds)
	outs := make([]string, cmdNum)

	for i := 0; i < cmdNum; i++ {
		wg.Add(1)
		go func(gossh *GoSSH, cmd string, outs []string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			out, _ := gossh.RunCmd(cmd)
			outs[i] = out
		}(gossh, cmds[i], outs, i, &wg)
	}
	wg.Wait()

	return outs, nil
}
