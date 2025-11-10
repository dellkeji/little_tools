package service

import (
	"path/filepath"
	"testing"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
)

func TestConfigService_GetConfig(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建仓库和服务
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := NewConfigService(configRepo)

	// 获取配置（应该创建默认配置）
	config, err := configService.GetConfig()
	if err != nil {
		t.Fatalf("获取配置失败: %v", err)
	}

	// 验证默认配置
	if config.ReminderTime != "10:00" {
		t.Errorf("默认提醒时间不正确，期望: 10:00，实际: %s", config.ReminderTime)
	}
	if config.ReminderEnabled {
		t.Errorf("默认提醒应该是禁用的")
	}
	if config.DataPath != "./data/tasks" {
		t.Errorf("默认数据路径不正确，期望: ./data/tasks，实际: %s", config.DataPath)
	}
}

func TestConfigService_UpdateConfig(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建仓库和服务
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := NewConfigService(configRepo)

	// 更新配置
	newConfig := &model.Config{
		WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test123",
		ReminderTime:    "09:30",
		ReminderEnabled: true,
		DataPath:        "./custom/path",
	}

	if err := configService.UpdateConfig(newConfig); err != nil {
		t.Fatalf("更新配置失败: %v", err)
	}

	// 重新获取配置验证
	config, err := configService.GetConfig()
	if err != nil {
		t.Fatalf("获取配置失败: %v", err)
	}

	if config.WebhookURL != newConfig.WebhookURL {
		t.Errorf("Webhook URL 不匹配")
	}
	if config.ReminderTime != newConfig.ReminderTime {
		t.Errorf("提醒时间不匹配")
	}
	if config.ReminderEnabled != newConfig.ReminderEnabled {
		t.Errorf("提醒开关不匹配")
	}
}

func TestConfigService_ValidateWebhook(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建仓库和服务
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := NewConfigService(configRepo)

	tests := []struct {
		name        string
		webhookURL  string
		expectError bool
	}{
		{
			name:        "有效的企业微信 Webhook URL",
			webhookURL:  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test123",
			expectError: false,
		},
		{
			name:        "空 URL",
			webhookURL:  "",
			expectError: true,
		},
		{
			name:        "无效的协议",
			webhookURL:  "ftp://example.com/webhook",
			expectError: true,
		},
		{
			name:        "缺少主机名",
			webhookURL:  "https:///webhook",
			expectError: true,
		},
		{
			name:        "企业微信 URL 缺少路径",
			webhookURL:  "https://qyapi.weixin.qq.com?key=test123",
			expectError: true,
		},
		{
			name:        "企业微信 URL 缺少 key 参数",
			webhookURL:  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send",
			expectError: true,
		},
		{
			name:        "有效的通用 Webhook URL",
			webhookURL:  "https://example.com/webhook",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configService.ValidateWebhook(tt.webhookURL)
			if tt.expectError && err == nil {
				t.Errorf("期望验证失败，但成功了")
			}
			if !tt.expectError && err != nil {
				t.Errorf("期望验证成功，但失败了: %v", err)
			}
		})
	}
}

func TestConfigService_ValidateReminderTime(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建仓库和服务
	configRepo := repository.NewFileConfigRepository(configPath)
	configService := NewConfigService(configRepo)

	tests := []struct {
		name         string
		reminderTime string
		expectError  bool
	}{
		{
			name:         "有效时间 10:00",
			reminderTime: "10:00",
			expectError:  false,
		},
		{
			name:         "有效时间 09:30",
			reminderTime: "09:30",
			expectError:  false,
		},
		{
			name:         "有效时间 23:59",
			reminderTime: "23:59",
			expectError:  false,
		},
		{
			name:         "有效时间 00:00",
			reminderTime: "00:00",
			expectError:  false,
		},
		{
			name:         "无效时间格式 - 缺少前导零",
			reminderTime: "9:30",
			expectError:  true,
		},
		{
			name:         "无效时间格式 - 小时超出范围",
			reminderTime: "25:00",
			expectError:  true,
		},
		{
			name:         "无效时间格式 - 分钟超出范围",
			reminderTime: "10:60",
			expectError:  true,
		},
		{
			name:         "无效时间格式 - 缺少冒号",
			reminderTime: "1000",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &model.Config{
				WebhookURL:      "",
				ReminderTime:    tt.reminderTime,
				ReminderEnabled: false,
				DataPath:        "./data/tasks",
			}

			err := configService.UpdateConfig(config)
			if tt.expectError && err == nil {
				t.Errorf("期望验证失败，但成功了")
			}
			if !tt.expectError && err != nil {
				t.Errorf("期望验证成功，但失败了: %v", err)
			}
		})
	}
}
