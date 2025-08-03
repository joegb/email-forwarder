package services

import (
	"bytes"
	"log"
	"strings"
	
	"github.com/joegb/email-forwarder/internal/database"
	"github.com/joegb/email-forwarder/internal/models"
	"github.com/joegb/email-forwarder/internal/utils"
	"github.com/joegb/email-forwarder/internal/logger"
	"github.com/joegb/email-forwarder/internal/config"
	
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

const (
	Keyword = "FORWARD"
)

//添加日志
func ProcessEmails() {
	log.Println("Starting email processing...")
	
	// 获取配置
	cfg := config.GetConfig()
	
	// 使用配置中的邮箱信息
	email := cfg.GmailEmail
	password := cfg.GmailAppPassword
	if email == "" || password == "" {
		log.Fatal("Email credentials not configured")
	}
	
	// 创建IMAP客户端
	imapClient, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Printf("Failed to create IMAP client: %v", err)
		return
	}
	defer imapClient.Logout()
	
	// 登录
	if err := imapClient.Login(email, password); err != nil {
		log.Printf("Failed to login: %v", err)
		return
	}
	
	// 选择收件箱
	_, err = imapClient.Select("INBOX", false)
	if err != nil {
		log.Printf("Failed to select inbox: %v", err)
		return
	}
	
	// 搜索未读邮件
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := imapClient.Search(criteria)
	if err != nil {
		log.Printf("Failed to search emails: %v", err)
		return
	}
	
	if len(ids) == 0 {
		log.Println("No unread messages to process")
		return
	}

	logger.Info("Found %d unread messages to process", len(ids))
	
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)
	
	// 获取邮件
	section := &imap.BodySectionName{}
	messages := make(chan *imap.Message, 10)
	go func() {
		if err := imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages); err != nil {
			log.Printf("Failed to fetch messages: %v", err)
		}
	}()
	
	// 处理邮件
	for msg := range messages {
		processSingleEmail(msg, email, section)
	}
	
	log.Println("Email processing completed")
}

//添加日志
func processSingleEmail(msg *imap.Message, fromEmail string, section *imap.BodySectionName) {
	subject := msg.Envelope.Subject
	logger.Info("Processing email: %s", subject)

	// 检查是否包含关键字
	if !strings.Contains(subject, Keyword) {
		return
	}
	
	// 解析主题
	parts := strings.SplitN(subject, "-", 2)
	if len(parts) < 2 {
		log.Printf("Invalid subject format: %s", subject)
		return
	}
	
	targetName := strings.TrimSpace(parts[1])
	
	// 查询目标邮箱
	var target models.ForwardTarget
	if err := database.DB.Where("name = ?", targetName).First(&target).Error; err != nil {
		log.Printf("Target not found: %s, error: %v", targetName, err)
		return
	}
	
	// 读取邮件正文
	body := msg.GetBody(section)
	if body == nil {
		log.Printf("Failed to get message body for subject: %s", subject)
		return
	}
	
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(body); err != nil {
		log.Printf("Failed to read message body: %v", err)
		return
	}
	
	// 转发邮件
	utils.ForwardEmail(fromEmail, target.Email, subject, buf.Bytes())
	log.Printf("Forwarded email to %s (%s)", target.Name, target.Email)
}