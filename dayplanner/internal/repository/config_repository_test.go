package repository

import (
	"os"
	"path/filepath"
	"testing"

	"daily-report-tool/internal/model"
)

func TestFileConfigRepository_LoadDefaultConfig(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "config_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	repo := NewFileConfigRepository(configPath)

	// 加载配置（文件不存在，应该创建默认配置）
	config, err := repo.Load()
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if config == nil {
		t.Fatal("配置为 nil")
	}

	// 验证默认配置
	if config.ReminderTime != "10:00" {
		t.Errorf("默认提醒时间不正确: 期望 10:00, 实际 %s", config.ReminderTime)
	}

	if config.ReminderEnabled != false {
		t.Error("默认提醒开关应该为 false")
	}

	if config.DataPath != "./data/tasks" {
		t.Errorf("默认数据路径不正确: 期望 ./data/tasks, 实际 %s", config.DataPath)
	}

	// 验证配置文件已创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("配置文件应该已创建")
	}
}

func TestFileConfigRepository_SaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	repo := NewFileConfigRepository(configPath)

	// 创建测试配置
	testConfig := &model.Config{
		WebhookURL:      "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
		ReminderTime:    "09:30",
		ReminderEnabled: true,
		DataPath:        "./custom/path",
	}

	// 保存配置
	err = repo.Save(testConfig)
	if err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}

	// 加载配置
	loadedConfig, err := repo.Load()
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 验证配置内容
	if loadedConfig.WebhookURL != testConfig.WebhookURL {
		t.Errorf("Webhook URL 不匹配: 期望 %s, 实际 %s", testConfig.WebhookURL, loadedConfig.WebhookURL)
	}

	if loadedConfig.ReminderTime != testConfig.ReminderTime {
		t.Errorf("提醒时间不匹配: 期望 %s, 实际 %s", testConfig.ReminderTime, loadedConfig.ReminderTime)
	}

	if loadedConfig.ReminderEnabled != testConfig.ReminderEnabled {
		t.Errorf("提醒开关不匹配: 期望 %v, 实际 %v", testConfig.ReminderEnabled, loadedConfig.ReminderEnabled)
	}

	if loadedConfig.DataPath != testConfig.DataPath {
		t.Errorf("数据路径不匹配: 期望 %s, 实际 %s", testConfig.DataPath, loadedConfig.DataPath)
	}
}

func TestFileConfigRepository_SaveCreatesDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_repo_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 使用嵌套的不存在的目录
	configPath := filepath.Join(tempDir, "nested", "dir", "config.json")
	repo := NewFileConfigRepository(configPath)

	testConfig := &model.Config{
		WebhookURL:      "",
		ReminderTime:    "10:00",
		ReminderEnabled: false,
		DataPath:        "./data/tasks",
	}

	// 保存配置应该自动创建目录
	err = repo.Save(testConfig)
	if err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}

	// 验证文件已创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("配置文件应该已创建")
	}

	// 验证目录已创建
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("配置目录应该已创建")
	}
}
