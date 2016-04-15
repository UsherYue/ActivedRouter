package netservice

import (
	"net/smtp"
	"strings"
)

//发送邮件
func SendEmailTo(user, password, host, to, subject, body, mailtype string) error {
	hostInfo := strings.Split(host, ":")
	//AUTH
	auth := smtp.PlainAuth("", user, password, hostInfo[0])
	//MIME TYPE
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	//email body
	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
