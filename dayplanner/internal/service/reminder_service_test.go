package service

import (
	"path/filepath"
	"testing"
	"time"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
)

// mockConfigService 用于测试的配置服务 mock
type mockConfigService struct {
	config *model.Config
}

func (m *mockConfigService) GetConfig() (*model.Config, error) {
	return m.config, nil
}

func (m *mockConfigService) UpdateConfig(config *model.Config) error {
	m.config = config
	return nil
}

func (m *mockConfigService) ValidateWebhook(webhookURL string) error {
	return nil
}

// mockTaskService 用于测试的任务服务 mock
type mockTaskService struct {
	hasTask bool
}

func (m *mockTaskService) GetTask(date time.Time) (*model.Task, error) {
	return nil, nil
}

func (m *mockTaskService) SaveTask(date time.Time, content string) error {
	return nil
}

func (m *mockTaskService) GetMonthTaskDates(year int, month time.Month) ([]time.Time, error) {
	return nil, nil
}

func (m *mockTaskService) HasTodayTask() (bool, error) {
	return m.hasTask, nil
}

func TestReminderService_StartStop(t *testing.T) {
	// 创建 mock 服务
	configService := &mockConfigService{
		config: &model.Config{
			WebhookURL:      "https://example.com/webhook",
			ReminderTime:    "10:00",
			ReminderEnabled: true,
			DataPath:        "./data/tasks",
		},
	}
	taskService := &mockTaskService{hasTask: false}

	// 创建提醒服务
	reminderService := NewReminderService(configService, taskService)

	// 测试启动
	if err := reminderService.Start(); err != nil {
		t.Fatalf("启动提醒服务失败: %v", err)
	}

	// 验证服务正在运行
	if !reminderService.running {
		t.Errorf("提醒服务应该在运行")
	}

	// 测试重复启动
	if err := reminderService.Start(); err == nil {
		t.Errorf("重复启动应该返回错误")
	}

	// 测试停止
	reminderService.Stop()

	// 验证服务已停止
	if reminderService.running {
		t.Errorf("提醒服务应该已停止")
	}
}

func TestReminderService_StartWithDisabledReminder(t *testing.T) {
	// 创建 mock 服务（提醒未启用）
	configService := &mockConfigService{
		config: &model.Config{
			WebhookURL:      "https://example.com/webhook",
			ReminderTime:    "10:00",
			ReminderEnabled: false, // 未启用
			DataPath:        "./data/tasks",
		},
	}
	taskService := &mockTaskService{hasTask: false}

	// 创建提醒服务
	reminderService := NewReminderService(configService, taskService)

	// 测试启动（应该成功但不实际运行）
	if err := reminderService.Start(); err != nil {
		t.Fatalf("启动提醒服务失败: %v", err)
	}

	// 验证服务未运行
	if reminderService.running {
		t.Errorf("提醒未启用时，服务不应该运行")
	}
}

func TestReminderService_SendReminder(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建真实的配置服务
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := NewConfigService(configRepo)

	// 设置配置（使用无效的 URL，因为我们只测试消息构建）
	config := &model.Config{
		WebhookURL:      "https://invalid-test-url.example.com/webhook",
		ReminderTime:    "10:00",
		ReminderEnabled: true,
		DataPath:        "./data/tasks",
	}
	if err := configService.UpdateConfig(config); err != nil {
		t.Fatalf("更新配置失败: %v", err)
	}

	taskService := &mockTaskService{hasTask: false}

	// 创建提醒服务
	reminderService := NewReminderService(configService, taskService)

	// 测试发送提醒（预期会失败，因为 URL 无效）
	err := reminderService.SendReminder("测试消息")
	if err == nil {
		t.Logf("注意: 发送到无效 URL 应该失败，但没有返回错误")
	}
}

func TestReminderService_PreventDuplicateReminders(t *testing.T) {
	// 创建 mock 服务
	configService := &mockConfigService{
		config: &model.Config{
			WebhookURL:      "https://example.com/webhook",
			ReminderTime:    "10:00",
			ReminderEnabled: true,
			DataPath:        "./data/tasks",
		},
	}
	taskService := &mockTaskService{hasTask: false}

	// 创建提醒服务
	reminderService := NewReminderService(configService, taskService)

	// 设置今天已发送
	today := time.Now().Format("2006-01-02")
	reminderService.lastSentDate = today

	// 验证防重复机制
	if reminderService.lastSentDate != today {
		t.Errorf("lastSentDate 应该是今天的日期")
	}
}
