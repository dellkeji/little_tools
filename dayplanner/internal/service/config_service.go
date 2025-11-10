package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
	"daily-report-tool/internal/util"
)

// ConfigService 定义配置管理服务接口
type ConfigService interface {
	// GetConfig 获取配置
	GetConfig() (*model.Config, error)

	// UpdateConfig 更新配置
	UpdateConfig(config *model.Config) error

	// ValidateWebhook 验证 Webhook 配置
	ValidateWebhook(webhookURL string) error
}

// ConfigServiceImpl 配置管理服务实现
type ConfigServiceImpl struct {
	configRepo repository.ConfigRepository
}

// NewConfigService 创建新的配置管理服务
func NewConfigService(configRepo repository.ConfigRepository) *ConfigServiceImpl {
	return &ConfigServiceImpl{
		configRepo: configRepo,
	}
}

// GetConfig 获取配置
func (s *ConfigServiceImpl) GetConfig() (*model.Config, error) {
	util.Debug("获取配置")
	config, err := s.configRepo.Load()
	if err != nil {
		util.Error("获取配置失败: %v", err)
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	util.Info("配置加载成功")
	return config, nil
}

// UpdateConfig 更新配置
func (s *ConfigServiceImpl) UpdateConfig(config *model.Config) error {
	util.Info("更新配置: Webhook=%s, 提醒时间=%s, 启用=%v", 
		config.WebhookURL, config.ReminderTime, config.ReminderEnabled)

	// 验证 Webhook URL（如果已启用提醒）
	if config.ReminderEnabled && config.WebhookURL != "" {
		if err := s.ValidateWebhook(config.WebhookURL); err != nil {
			util.Warn("Webhook URL 验证失败: %v", err)
			return err
		}
		util.Debug("Webhook URL 验证通过")
	}

	// 验证提醒时间格式
	if err := s.validateReminderTime(config.ReminderTime); err != nil {
		util.Warn("提醒时间格式验证失败: %v", err)
		return err
	}
	util.Debug("提醒时间格式验证通过")

	if err := s.configRepo.Save(config); err != nil {
		util.Error("保存配置失败: %v", err)
		return fmt.Errorf("更新配置失败: %w", err)
	}

	util.Info("配置更新成功")
	return nil
}

// ValidateWebhook 验证 Webhook URL 格式的有效性
func (s *ConfigServiceImpl) ValidateWebhook(webhookURL string) error {
	util.Debug("验证 Webhook URL: %s", webhookURL)

	if webhookURL == "" {
		return fmt.Errorf("Webhook URL 不能为空")
	}

	// 解析 URL
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		util.Warn("Webhook URL 解析失败: %v", err)
		return fmt.Errorf("无效的 Webhook URL 格式: %w", err)
	}

	// 检查协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		util.Warn("Webhook URL 协议无效: %s", parsedURL.Scheme)
		return fmt.Errorf("Webhook URL 必须使用 http 或 https 协议")
	}

	// 检查主机名
	if parsedURL.Host == "" {
		util.Warn("Webhook URL 缺少主机名")
		return fmt.Errorf("Webhook URL 缺少主机名")
	}

	// 检查是否为企业微信 Webhook URL
	if strings.Contains(parsedURL.Host, "qyapi.weixin.qq.com") {
		// 验证企业微信 Webhook 路径格式
		if !strings.Contains(parsedURL.Path, "/cgi-bin/webhook/send") {
			util.Warn("企业微信 Webhook URL 路径无效: %s", parsedURL.Path)
			return fmt.Errorf("无效的企业微信 Webhook URL 路径")
		}

		// 检查是否包含 key 参数
		if parsedURL.Query().Get("key") == "" {
			util.Warn("企业微信 Webhook URL 缺少 key 参数")
			return fmt.Errorf("企业微信 Webhook URL 缺少 key 参数")
		}
	}

	util.Debug("Webhook URL 验证通过")
	return nil
}

// validateReminderTime 验证提醒时间格式
func (s *ConfigServiceImpl) validateReminderTime(reminderTime string) error {
	// 验证时间格式 HH:MM
	matched, err := regexp.MatchString(`^([01]\d|2[0-3]):([0-5]\d)$`, reminderTime)
	if err != nil {
		return fmt.Errorf("验证时间格式失败: %w", err)
	}

	if !matched {
		return fmt.Errorf("无效的时间格式，应为 HH:MM (例如: 10:00)")
	}

	return nil
}
