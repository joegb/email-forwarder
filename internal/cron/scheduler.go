package cron

import (
	"time"
	"log"
	
	"github.com/robfig/cron/v3"
	"github.com/joegb/email-forwarder/internal/services"
	"github.com/joegb/email-forwarder/internal/config"

)

func StartEmailCron() {
	cfg := config.GetConfig()
	
	c := cron.New()
	
	// 使用配置中的定时计划
	_, err := c.AddFunc(cfg.CronSchedule, func() {
		start := time.Now()
		services.ProcessEmails()
		log.Printf("Email processing took %v", time.Since(start))
	})
	
	if err != nil {
		log.Fatalf("Failed to start cron job: %v", err)
	}
	
	c.Start()
	log.Println("Email cron job started")
}