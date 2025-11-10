package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"daily-report-tool/internal/util"
)

// ReminderService 定义提醒服务接口
type ReminderService interface {
	// Start 启动提醒服务
	Start() error

	// Stop 停止提醒服务
	Stop()

	// SendReminder 发送提醒消息
	SendReminder(message string) error
}

// ReminderServiceImpl 提醒服务实现
type ReminderServiceImpl struct {
	configService ConfigService
	taskService   TaskService
	ticker        *time.Ticker
	stopChan      chan bool
	lastSentDate  string // 记录上次发送提醒的日期，防止重复发送
	mu            sync.Mutex
	running       bool
}

// NewReminderService 创建新的提醒服务
func NewReminderService(configService ConfigService, taskService TaskService) *ReminderServiceImpl {
	return &ReminderServiceImpl{
		configService: configService,
		taskService:   taskService,
		stopChan:      make(chan bool),
		lastSentDate:  "",
	}
}

// Start 启动提醒服务
func (s *ReminderServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		util.Warn("提醒服务已经在运行")
		return fmt.Errorf("提醒服务已经在运行")
	}

	// 获取配置
	config, err := s.configService.GetConfig()
	if err != nil {
		util.Error("启动提醒服务失败，获取配置失败: %v", err)
		return fmt.Errorf("获取配置失败: %w", err)
	}

	// 检查是否启用提醒
	if !config.ReminderEnabled {
		util.Info("提醒服务未启用")
		return nil // 未启用提醒，直接返回
	}

	// 创建定时器，每分钟检查一次
	s.ticker = time.NewTicker(1 * time.Minute)
	s.running = true

	util.Info("提醒服务已启动，提醒时间: %s", config.ReminderTime)

	// 启动后台 goroutine 进行定时检查
	go s.checkReminder()

	return nil
}

// Stop 停止提醒服务
func (s *ReminderServiceImpl) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.stopChan <- true
	s.running = false
	util.Info("提醒服务已停止")
}

// checkReminder 定时检查是否需要发送提醒
func (s *ReminderServiceImpl) checkReminder() {
	for {
		select {
		case <-s.ticker.C:
			s.performReminderCheck()
		case <-s.stopChan:
			return
		}
	}
}

// performReminderCheck 执行提醒检查逻辑
func (s *ReminderServiceImpl) performReminderCheck() {
	// 获取配置
	config, err := s.configService.GetConfig()
	if err != nil {
		util.Error("获取配置失败: %v", err)
		fmt.Printf("获取配置失败: %v\n", err)
		return
	}

	// 检查是否启用提醒
	if !config.ReminderEnabled {
		return
	}

	// 检查当前时间是否匹配提醒时间
	now := time.Now()
	currentTime := now.Format("15:04")

	if currentTime != config.ReminderTime {
		return // 不是提醒时间
	}

	util.Debug("到达提醒时间: %s", currentTime)

	// 检查今天是否已经发送过提醒
	today := now.Format("2006-01-02")
	s.mu.Lock()
	if s.lastSentDate == today {
		s.mu.Unlock()
		util.Debug("今天已发送过提醒，跳过")
		return // 今天已经发送过提醒
	}
	s.mu.Unlock()

	// 检查今天是否有任务
	hasTask, err := s.taskService.HasTodayTask()
	if err != nil {
		util.Error("检查今天任务失败: %v", err)
		fmt.Printf("检查今天任务失败: %v\n", err)
		return
	}

	if hasTask {
		util.Debug("今天已有任务，不需要提醒")
		return // 今天已有任务，不需要提醒
	}

	// 发送提醒
	message := "提醒：您今天还没有填写日报，请及时记录工作内容。"
	if err := s.SendReminder(message); err != nil {
		util.Error("发送提醒失败: %v", err)
		fmt.Printf("发送提醒失败: %v\n", err)
		return
	}

	// 记录发送日期，防止重复发送
	s.mu.Lock()
	s.lastSentDate = today
	s.mu.Unlock()

	util.Info("提醒已发送: %s", today)
	fmt.Printf("提醒已发送: %s\n", today)
}

// SendReminder 发送提醒消息
func (s *ReminderServiceImpl) SendReminder(message string) error {
	util.Info("准备发送提醒消息: %s", message)

	// 获取配置
	config, err := s.configService.GetConfig()
	if err != nil {
		util.Error("获取配置失败: %v", err)
		return fmt.Errorf("获取配置失败: %w", err)
	}

	// 检查 Webhook URL
	if config.WebhookURL == "" {
		util.Warn("Webhook URL 未配置")
		return fmt.Errorf("Webhook URL 未配置")
	}

	// 构建企业微信消息体
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		util.Error("序列化消息失败: %v", err)
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	// 发送 HTTP POST 请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	util.Debug("发送 Webhook 请求到: %s", config.WebhookURL)
	resp, err := client.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		util.Error("发送 Webhook 请求失败: %v", err)
		return fmt.Errorf("发送 Webhook 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		util.Error("Webhook 请求失败，状态码: %d", resp.StatusCode)
		return fmt.Errorf("Webhook 请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		util.Error("解析响应失败: %v", err)
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查企业微信返回的错误码
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg := result["errmsg"].(string)
		util.Error("企业微信返回错误: %s (错误码: %.0f)", errmsg, errcode)
		return fmt.Errorf("企业微信返回错误: %s (错误码: %.0f)", errmsg, errcode)
	}

	util.Info("提醒消息发送成功")
	return nil
}
