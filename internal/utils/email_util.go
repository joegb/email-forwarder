package utils

import (
	"crypto/tls"
	"log"
	"net/smtp"
	"os"
)

func ForwardEmail(from, to, subject string, body []byte) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	username := os.Getenv("GMAIL_EMAIL")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	
	// 设置认证
	auth := smtp.PlainAuth("", username, password, smtpHost)
	
	// 构建邮件
	msg := []byte(
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		string(body) + "\r\n")
	
	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
	
	// 连接服务器
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		log.Printf("TLS connection error: %v", err)
		return
	}
	defer conn.Close()
	
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("SMTP client error: %v", err)
		return
	}
	defer client.Close()
	
	// 认证
	if err = client.Auth(auth); err != nil {
		log.Printf("SMTP auth error: %v", err)
		return
	}
	
	// 设置发件人
	if err = client.Mail(from); err != nil {
		log.Printf("Mail from error: %v", err)
		return
	}
	
	// 设置收件人
	if err = client.Rcpt(to); err != nil {
		log.Printf("Rcpt to error: %v", err)
		return
	}
	
	// 发送数据
	w, err := client.Data()
	if err != nil {
		log.Printf("Data command error: %v", err)
		return
	}
	
	_, err = w.Write(msg)
	if err != nil {
		log.Printf("Write error: %v", err)
		return
	}
	
	err = w.Close()
	if err != nil {
		log.Printf("Close writer error: %v", err)
		return
	}
	
	// 退出
	client.Quit()
}