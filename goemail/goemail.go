// 发送邮件
package goemail

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Host         string
	Port         int
	User, UserPW string
}

type Mail struct {
	To         string // 接收邮箱
	Subject    string // 主题
	Body       string // 内容
	AttachPath string // 附件路径
}

// 发送单封邮件
func SendEmail(smtp *Smtp, mail *Mail) error {
	// 收件人不能为空
	if mail.To == "" {
		return fmt.Errorf("%#v can not empty", mail.To)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", smtp.User)
	m.SetHeader("To", mail.To)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)

	if mail.AttachPath != "" {
		m.Attach(mail.AttachPath)
	}

	e := gomail.NewDialer(smtp.Host, smtp.Port, smtp.User, smtp.UserPW)
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

// 发送单封邮件给多人
func SendEmailMP(smtp *Smtp, mail *Mail, mailList []string) []error {
	var errList []error
	for _, to := range mailList {
		mail.To = to
		if err := SendEmail(smtp, mail); err != nil {
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
		if err := SendEmail(smtp, &m); err != nil {
			errList = append(errList, err)
		}

		// 间隔 0.5 秒
		time.Sleep(500 * time.Millisecond)
	}
	return errList
}

// 发送单封邮件，相关信息从参数读取
// 参数顺序固定
// 依次为：SMTP：Host、Port、User、UserPW，邮件：接收邮箱、主题、内容,以上为 7 项必填
// 可选参数：邮件：、附件路径
func SendEmailSmtpMail(info ...string) error {
	if len(info) < 7 {
		return fmt.Errorf("%#v Missing parameter", info)
	}

	smtp := &Smtp{
		Host:   info[0],
		User:   info[2],
		UserPW: info[3],
	}
	smtp.Port, _ = strconv.Atoi(info[1])

	mail := &Mail{
		To:      info[4],
		Subject: info[5],
		Body:    info[6],
	}
	if len(info) == 8 {
		mail.AttachPath = info[7]
	}

	if err := SendEmail(smtp, mail); err != nil {
		return err
	}

	return nil
}
