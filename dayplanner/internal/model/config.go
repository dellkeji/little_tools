package model

// Config 表示应用程序的配置信息
type Config struct {
	WebhookURL      string `json:"webhook_url"`       // 企业微信 Webhook 地址
	ReminderTime    string `json:"reminder_time"`     // 提醒时间 (格式: "10:00")
	ReminderEnabled bool   `json:"reminder_enabled"`  // 是否启用提醒
	DataPath        string `json:"data_path"`         // 数据存储路径
}
