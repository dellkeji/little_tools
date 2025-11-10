package ui

import (
	"fmt"
	"strings"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/service"
	"daily-report-tool/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SettingsView 设置界面组件
type SettingsView struct {
	window          fyne.Window
	configService   service.ConfigService
	webhookEntry    *widget.Entry
	hourSelect      *widget.Select
	minuteSelect    *widget.Select
	reminderCheck   *widget.Check
	saveButton      *widget.Button
	cancelButton    *widget.Button
	onConfigUpdated func() // 配置更新后的回调
}

// NewSettingsView 创建新的设置界面
func NewSettingsView(parent fyne.Window, configService service.ConfigService) *SettingsView {
	sv := &SettingsView{
		window:        parent,
		configService: configService,
	}

	sv.initializeComponents()
	return sv
}

// initializeComponents 初始化所有 UI 组件
func (sv *SettingsView) initializeComponents() {
	// 创建 Webhook URL 输入框
	sv.webhookEntry = widget.NewEntry()
	sv.webhookEntry.SetPlaceHolder("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxx")

	// 创建小时选择器
	hours := make([]string, 24)
	for i := 0; i < 24; i++ {
		hours[i] = fmt.Sprintf("%02d", i)
	}
	sv.hourSelect = widget.NewSelect(hours, nil)
	sv.hourSelect.SetSelected("10")

	// 创建分钟选择器
	minutes := make([]string, 60)
	for i := 0; i < 60; i++ {
		minutes[i] = fmt.Sprintf("%02d", i)
	}
	sv.minuteSelect = widget.NewSelect(minutes, nil)
	sv.minuteSelect.SetSelected("00")

	// 创建提醒开关
	sv.reminderCheck = widget.NewCheck("启用每日提醒", nil)

	// 创建保存按钮
	sv.saveButton = widget.NewButton("保存", sv.onSave)

	// 创建取消按钮
	sv.cancelButton = widget.NewButton("取消", sv.onCancel)
}

// Show 显示设置窗口
func (sv *SettingsView) Show() {
	// 加载当前配置
	sv.loadConfig()

	// 创建表单布局
	form := sv.createForm()

	// 创建按钮容器
	buttons := container.NewHBox(
		sv.saveButton,
		sv.cancelButton,
	)

	// 创建主容器
	content := container.NewBorder(
		nil,     // top
		buttons, // bottom
		nil,     // left
		nil,     // right
		form,    // center
	)

	// 创建并显示对话框
	settingsDialog := dialog.NewCustom("设置", "关闭", content, sv.window)
	settingsDialog.Resize(fyne.NewSize(500, 300))
	settingsDialog.Show()
}

// createForm 创建表单布局
func (sv *SettingsView) createForm() *fyne.Container {
	// Webhook URL 表单项
	webhookLabel := widget.NewLabel("企业微信 Webhook URL:")
	webhookForm := container.NewVBox(
		webhookLabel,
		sv.webhookEntry,
	)

	// 提醒时间表单项
	timeLabel := widget.NewLabel("提醒时间:")
	timeContainer := container.NewHBox(
		sv.hourSelect,
		widget.NewLabel(":"),
		sv.minuteSelect,
	)
	timeForm := container.NewVBox(
		timeLabel,
		timeContainer,
	)

	// 提醒开关表单项
	reminderForm := container.NewVBox(
		sv.reminderCheck,
	)

	// 组合所有表单项
	form := container.NewVBox(
		webhookForm,
		timeForm,
		reminderForm,
	)

	return form
}

// loadConfig 加载当前配置到 UI
func (sv *SettingsView) loadConfig() {
	config, err := sv.configService.GetConfig()
	if err != nil {
		util.Error("加载配置失败: %v", err)
		util.ShowWarningDialog("警告", "无法加载配置，将使用默认值", sv.window)
		return
	}

	// 设置 Webhook URL
	sv.webhookEntry.SetText(config.WebhookURL)

	// 解析并设置提醒时间
	if config.ReminderTime != "" {
		parts := strings.Split(config.ReminderTime, ":")
		if len(parts) == 2 {
			sv.hourSelect.SetSelected(parts[0])
			sv.minuteSelect.SetSelected(parts[1])
		}
	}

	// 设置提醒开关
	sv.reminderCheck.SetChecked(config.ReminderEnabled)
}

// SetOnConfigUpdated 设置配置更新回调
func (sv *SettingsView) SetOnConfigUpdated(callback func()) {
	sv.onConfigUpdated = callback
}

// onSave 保存按钮点击事件处理
func (sv *SettingsView) onSave() {
	// 输入验证
	webhookURL := strings.TrimSpace(sv.webhookEntry.Text)
	
	// 验证提醒时间选择
	if sv.hourSelect.Selected == "" || sv.minuteSelect.Selected == "" {
		util.Warn("提醒时间未选择")
		util.ShowWarningDialog("输入错误", "请选择提醒时间", sv.window)
		return
	}

	// 如果启用提醒，验证 Webhook URL
	if sv.reminderCheck.Checked && webhookURL == "" {
		util.Warn("启用提醒但未配置 Webhook URL")
		util.ShowWarningDialog("输入错误", "启用提醒功能需要配置企业微信 Webhook URL", sv.window)
		return
	}

	// 构建配置对象
	config := &model.Config{
		WebhookURL:      webhookURL,
		ReminderTime:    fmt.Sprintf("%s:%s", sv.hourSelect.Selected, sv.minuteSelect.Selected),
		ReminderEnabled: sv.reminderCheck.Checked,
		DataPath:        "./data/tasks", // 保持默认数据路径
	}

	util.Info("保存配置: Webhook=%s, 提醒时间=%s, 启用=%v", 
		config.WebhookURL, config.ReminderTime, config.ReminderEnabled)

	// 调用配置服务验证和保存
	err := sv.configService.UpdateConfig(config)
	if err != nil {
		// 显示错误提示
		util.ShowErrorDialogWithMessage("保存失败", "无法保存配置", err, sv.window)
		return
	}

	// 保存成功，显示成功提示
	util.ShowSuccessNotification("配置保存成功！", sv.window)

	// 调用配置更新回调
	if sv.onConfigUpdated != nil {
		sv.onConfigUpdated()
	}
}

// onCancel 取消按钮点击事件处理
func (sv *SettingsView) onCancel() {
	// 关闭设置窗口（通过对话框的关闭按钮实现）
	// 这里不需要额外操作，用户点击对话框的关闭按钮即可
}

// showError 显示错误提示对话框
func (sv *SettingsView) showError(message string) {
	dialog.ShowError(fmt.Errorf("%s", message), sv.window)
}

// showSuccess 显示成功提示对话框
func (sv *SettingsView) showSuccess(message string) {
	dialog.ShowInformation("成功", message, sv.window)
}
