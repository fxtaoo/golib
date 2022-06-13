// 发送邮件
package goemail

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/fxtaoo/golib/gofile"
	"gopkg.in/gomail.v2"
)

type SendSmtp struct {
	Host, User, UserPW string
	Port               int
}

type Config struct {
	Smtp SendSmtp
}

func (c *Config) Read(str string) {
	gofile.TomlFileRead(str, c)
}

// 发送邮件
// 参数顺序固定
// info[0-3] 发送邮箱 toml 配置文件名或路径，接收邮箱 主题 内容 必选
// info[4] 附件路径 可选
func SendEmail(info ...string) error {

	if len(info) < 4 {
		return errors.New("发送邮箱 toml 配置文件名或路径，接收邮箱 主题 内容 必选")
	}
	var conf Config
	conf.Read(info[0])

	m := gomail.NewMessage()

	m.SetHeader("From", conf.Smtp.User)
	m.SetHeader("To", info[1])
	m.SetHeader("Subject", info[2])
	m.SetBody("text/html", info[3])

	if len(info) == 5 {
		m.Attach(info[4])
	}

	e := gomail.NewDialer(conf.Smtp.Host, conf.Smtp.Port, conf.Smtp.User, conf.Smtp.UserPW)
	e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := e.DialAndSend(m); err != nil {
		time.Sleep(3 * time.Second)
		if err := e.DialAndSend(m); err != nil {
			return err
		}
	}
	return nil
}

// func main() {
// 	flag.Parse()
// 	SendEmail(flag.Args()...)
// }
