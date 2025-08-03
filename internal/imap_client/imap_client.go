package imap_client

import (
	"crypto/tls"
	"log"
	
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ImapClient struct {
	Client *client.Client
}

func NewImapClient(server, email, password string) (*ImapClient, error) {
	// 连接到服务器
	c, err := client.DialTLS(server, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	
	// 登录
	if err := c.Login(email, password); err != nil {
		return nil, err
	}
	
	return &ImapClient{Client: c}, nil
}

func (ic *ImapClient) FetchUnreadMessages() (chan *imap.Message, error) {
	// 选择收件箱
	mbox, err := ic.Client.Select("INBOX", false)
	if err != nil {
		return nil, err
	}
	
	// 获取未读邮件
	seqset := new(imap.SeqSet)
	if mbox.Messages == 0 {
		return nil, nil
	}
	
	seqset.AddRange(1, mbox.Messages)
	
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- ic.Client.Fetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			imap.FetchBody,
			imap.FetchFlags,
			imap.FetchInternalDate,
		}, messages)
	}()
	
	return messages, nil
}

func (ic *ImapClient) Close() {
	if err := ic.Client.Logout(); err != nil {
		log.Printf("Error logging out: %v", err)
	}
}