// 发送邮件
package mail

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	To         []string // 接收邮箱
	Subject    string   // 主题
	Body       string   // 内容
	AttachPath string   // 附件路径
}

// 发送邮件
func (mail *Mail) Send(s *Smtp) error {
	// 收件人不能为空
	if len(mail.To) == 0 {
		return fmt.Errorf("%#v can not empty", mail.To)
	}

	m := gomail.NewMessage()

	if s.Name == "" {
		s.Name = s.Address
	}

	m.SetHeaders(map[string][]string{
		"From":      {m.FormatAddress(s.Address, s.Name)},
		"To":        mail.To,
		"Subject":   {mail.Subject},
		"text/html": {mail.Body},
	})

	if mail.AttachPath != "" {
		m.Attach(mail.AttachPath)
	}

	e := gomail.NewDialer(s.Host, s.Port, s.Address, s.PassWd)
	e.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := e.DialAndSend(m); err != nil {
		// 失败暂停 1s 重发
		time.Sleep(1 * time.Second)
		if err := e.DialAndSend(m); err != nil {
			return err
		}
	}
	return nil
}

// 单独发送邮件给每一个人
func (mail *Mail) SendAlone(s *Smtp) []error {
	var errList []error
	tmpMail := Mail{
		Subject:    mail.Subject,
		Body:       mail.Body,
		AttachPath: mail.AttachPath,
	}
	for _, to := range mail.To {
		tmpMail.To[0] = to
		if err := tmpMail.Send(s); err != nil {
			errList = append(errList, err)
		}
		// 间隔 0.5 秒
		time.Sleep(500 * time.Millisecond)
	}
	return errList
}

// 发送多封邮件
func SendEmailList(smtp *Smtp, mail []Mail) []error {
	var errList []error
	for _, m := range mail {
		if err := m.Send(smtp); err != nil {
			errList = append(errList, err)
		}
		// 间隔 0.5 秒
		time.Sleep(500 * time.Millisecond)
	}
	return errList
}

// 发送邮件，相关配置信息从参数读取
// 参数顺序固定
// 依次为：SMTP：Host、Port、Address、UserPW，Name 邮件：接收邮箱、主题、内容,以上为 8 项必填
// 可选参数：邮件：附件路径
func SendEmail(info ...string) error {
	if len(info) < 8 {
		return fmt.Errorf("%#v Missing parameter", info)
	}

	smtp := &Smtp{
		Host:    info[0],
		Address: info[2],
		PassWd:  info[3],
		Name:    info[4],
	}
	smtp.Port, _ = strconv.Atoi(info[1])

	mail := &Mail{
		To:      []string{info[5]},
		Subject: info[6],
		Body:    info[7],
	}
	if len(info) == 9 {
		mail.AttachPath = info[8]
	}

	if err := mail.Send(smtp); err != nil {
		return err
	}

	return nil
}
