package ui

import (
	"time"

	"daily-report-tool/internal/service"
	"daily-report-tool/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// MainWindow 主窗口
type MainWindow struct {
	app             fyne.App
	window          fyne.Window
	taskService     service.TaskService
	configService   service.ConfigService
	reminderService service.ReminderService

	// UI 组件
	calendarView *CalendarView
	editorView   *EditorView
	previewView  *PreviewView
	settingsView *SettingsView
}

// NewMainWindow 创建新的主窗口
func NewMainWindow(
	app fyne.App,
	taskService service.TaskService,
	configService service.ConfigService,
	reminderService service.ReminderService,
) *MainWindow {
	mw := &MainWindow{
		app:             app,
		taskService:     taskService,
		configService:   configService,
		reminderService: reminderService,
	}

	// 创建窗口
	mw.window = app.NewWindow("日报工具")

	// 初始化 UI 组件
	mw.initializeComponents()

	// 设置窗口布局
	mw.setupLayout()

	// 设置窗口大小
	mw.window.Resize(fyne.NewSize(1200, 800))

	return mw
}

// initializeComponents 初始化所有 UI 组件
func (mw *MainWindow) initializeComponents() {
	// 创建日历视图
	mw.calendarView = NewCalendarView(mw.taskService)

	// 创建编辑器视图
	mw.editorView = NewEditorView(mw.taskService)
	mw.editorView.SetParentWindow(mw.window)

	// 创建预览视图
	mw.previewView = NewPreviewView()

	// 创建设置视图
	mw.settingsView = NewSettingsView(mw.window, mw.configService)

	// 设置组件间交互
	mw.setupInteractions()
}

// setupInteractions 设置组件间交互
func (mw *MainWindow) setupInteractions() {
	// 1. 日历日期选择事件 - 加载对应日期的任务到编辑器
	mw.calendarView.SetOnDateSelected(func(date time.Time) {
		mw.onDateSelected(date)
	})

	// 2. 编辑器内容变更事件 - 更新预览
	mw.editorView.SetOnContentChange(func(content string) {
		mw.previewView.UpdatePreview(content)
	})

	// 3. 编辑器保存完成事件 - 刷新日历标记
	mw.editorView.SetOnSaveComplete(func() {
		mw.calendarView.Refresh()
	})

	// 4. 设置更新事件 - 重启提醒服务
	mw.settingsView.SetOnConfigUpdated(func() {
		mw.onConfigUpdated()
	})
}

// onDateSelected 处理日期选择事件
func (mw *MainWindow) onDateSelected(date time.Time) {
	util.Info("选择日期: %s", date.Format("2006-01-02"))

	// 加载选中日期的任务
	task, err := mw.taskService.GetTask(date)
	if err != nil {
		util.Error("加载任务失败: %v", err)
		util.ShowErrorDialogWithMessage("加载失败", 
			"无法加载选中日期的任务，请检查文件系统权限", err, mw.window)
		return
	}

	// 设置编辑器日期
	mw.editorView.SetDate(date)

	// 设置编辑器内容
	if task != nil && task.Content != "" {
		mw.editorView.SetContent(task.Content)
		// 同时更新预览
		mw.previewView.UpdatePreview(task.Content)
		util.Debug("加载任务内容成功，长度: %d", len(task.Content))
	} else {
		// 没有任务内容，清空编辑器和预览
		mw.editorView.SetContent("")
		mw.previewView.Clear()
		util.Debug("该日期无任务内容")
	}
}

// onConfigUpdated 处理配置更新事件
func (mw *MainWindow) onConfigUpdated() {
	util.Info("配置已更新，重启提醒服务")

	// 停止现有的提醒服务
	if mw.reminderService != nil {
		mw.reminderService.Stop()
	}

	// 重新启动提醒服务（如果配置启用）
	if err := mw.reminderService.Start(); err != nil {
		util.Error("重启提醒服务失败: %v", err)
		util.ShowWarningDialog("警告", 
			"提醒服务启动失败，请检查配置是否正确", mw.window)
	} else {
		util.Info("提醒服务已重启")
	}
}

// setupLayout 设置窗口布局
func (mw *MainWindow) setupLayout() {
	// 创建右侧分栏：编辑器和预览
	rightSplit := container.NewVSplit(
		mw.editorView.GetContainer(),
		mw.previewView.GetContainer(),
	)
	rightSplit.SetOffset(0.5) // 设置分割比例为 50:50

	// 创建主分栏：日历和右侧内容
	mainSplit := container.NewHSplit(
		mw.calendarView,
		rightSplit,
	)
	mainSplit.SetOffset(0.25) // 设置分割比例为 25:75

	// 创建菜单栏
	mainMenu := mw.createMenu()

	// 设置窗口内容
	mw.window.SetContent(mainSplit)
	mw.window.SetMainMenu(mainMenu)
}

// createMenu 创建菜单栏
func (mw *MainWindow) createMenu() *fyne.MainMenu {
	// 创建设置菜单项
	settingsItem := fyne.NewMenuItem("设置", func() {
		mw.settingsView.Show()
	})

	// 创建文件菜单
	fileMenu := fyne.NewMenu("文件", settingsItem)

	// 创建主菜单
	mainMenu := fyne.NewMainMenu(fileMenu)

	return mainMenu
}

// Show 显示主窗口
func (mw *MainWindow) Show() {
	mw.window.ShowAndRun()
}

// GetWindow 获取窗口对象
func (mw *MainWindow) GetWindow() fyne.Window {
	return mw.window
}
