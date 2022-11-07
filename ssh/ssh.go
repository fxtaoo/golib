// ssh
package ssh

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/ssh"
)

// 连接 Session
func (gs *SSH) Connect() (*ssh.Session, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", gs.Host, gs.Port), gs.clientConfig)
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
func (gs *SSH) RunCmd(cmd string) (string, error) {
	session, err := gs.Connect()
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
func (gs *SSH) RunCmdsSequential(cmds []string) ([]string, error) {
	var outs []string

	for _, cmd := range cmds {
		str, err := gs.RunCmd(cmd)
		if err != nil {
			return nil, err
		}

		outs = append(outs, str)
	}

	return outs, nil
}

// 运行一组命令，并行执行，结果顺序与命令顺序一致
func (gs *SSH) RunCmdsParallel(cmds []string) ([]string, error) {
	var wg sync.WaitGroup
	cmdNum := len(cmds)
	outs := make([]string, cmdNum)

	for i := 0; i < cmdNum; i++ {
		wg.Add(1)
		go func(gs *SSH, cmd string, outs []string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			out, _ := gs.RunCmd(cmd)
			outs[i] = out
		}(gs, cmds[i], outs, i, &wg)
	}
	wg.Wait()

	return outs, nil
}
