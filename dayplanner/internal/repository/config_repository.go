package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"daily-report-tool/internal/model"
)

// ConfigRepository 定义配置数据访问接口
type ConfigRepository interface {
	// Load 加载配置
	Load() (*model.Config, error)

	// Save 保存配置
	Save(config *model.Config) error
}

// FileConfigRepository 基于文件系统的配置仓库实现
type FileConfigRepository struct {
	configPath string
}

// NewFileConfigRepository 创建新的文件配置仓库
func NewFileConfigRepository(configPath string) *FileConfigRepository {
	return &FileConfigRepository{
		configPath: configPath,
	}
}

// Load 加载配置
func (r *FileConfigRepository) Load() (*model.Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		defaultConfig := r.createDefaultConfig()
		if err := r.Save(defaultConfig); err != nil {
			return nil, fmt.Errorf("创建默认配置失败: %w", err)
		}
		return defaultConfig, nil
	}

	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config model.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置数据失败: %w", err)
	}

	return &config, nil
}

// Save 保存配置
func (r *FileConfigRepository) Save(config *model.Config) error {
	// 确保配置目录存在
	configDir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %w", err)
	}

	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// createDefaultConfig 创建默认配置
func (r *FileConfigRepository) createDefaultConfig() *model.Config {
	return &model.Config{
		WebhookURL:      "",
		ReminderTime:    "10:00",
		ReminderEnabled: false,
		DataPath:        "./data/tasks",
	}
}
