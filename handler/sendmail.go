package handler
import (
	"fmt"
	"strings"
	"net/smtp"
)

var (
	user = "quanix@163.com"
	password = "123456"
	host = "smtp.163.com:25"
	to = "domac@qq.com"
	subject = "文件发生变更"
)

//发送电子邮件
func SendMail(content string) error {
	body := `
    <html>
    <body>
    <h3>
    文件列表：
    </h3>
	<br>
	<h5>
	%s
	</h5>
    </body>
    </html>
    `

	body = fmt.Sprintf(body, content)
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	content_type := "Content-Type: text/" + "html" + "; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	if err != nil {
		fmt.Println("send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("send mail success!")
	}
	return err
}
