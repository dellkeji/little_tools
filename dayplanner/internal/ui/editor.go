package ui

import (
	"time"

	"daily-report-tool/internal/service"
	"daily-report-tool/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// EditorView 编辑器视图组件
type EditorView struct {
	container       *fyne.Container
	titleLabel      *widget.Label
	editor          *widget.Entry
	currentDate     time.Time
	onContentChange func(content string)
	saveTimer       *time.Timer
	taskService     service.TaskService
	onSaveComplete  func() // 保存完成后的回调，用于刷新日历
	parentWindow    fyne.Window // 用于显示错误对话框
}

// NewEditorView 创建新的编辑器视图
func NewEditorView(taskService service.TaskService) *EditorView {
	ev := &EditorView{
		currentDate: time.Now(),
		taskService: taskService,
	}

	// 创建标题标签
	ev.titleLabel = widget.NewLabel("选择日期以开始编辑")
	ev.titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 创建多行文本编辑器
	ev.editor = widget.NewMultiLineEntry()
	ev.editor.SetPlaceHolder("在此输入您的日报内容，支持 Markdown 格式...")
	ev.editor.Wrapping = fyne.TextWrapWord

	// 监听内容变更事件
	ev.editor.OnChanged = func(content string) {
		// 触发自动保存
		ev.triggerAutoSave(content)
		
		// 调用外部回调（用于更新预览）
		if ev.onContentChange != nil {
			ev.onContentChange(content)
		}
	}

	// 创建容器布局
	ev.container = container.NewBorder(
		ev.titleLabel, // top
		nil,           // bottom
		nil,           // left
		nil,           // right
		ev.editor,     // center
	)

	return ev
}

// GetContainer 获取容器
func (ev *EditorView) GetContainer() *fyne.Container {
	return ev.container
}

// SetContent 设置编辑器内容
func (ev *EditorView) SetContent(content string) {
	ev.editor.SetText(content)
}

// GetContent 获取编辑器内容
func (ev *EditorView) GetContent() string {
	return ev.editor.Text
}

// SetDate 设置当前编辑的日期
func (ev *EditorView) SetDate(date time.Time) {
	ev.currentDate = date
	ev.titleLabel.SetText(date.Format("2006年01月02日 星期一"))
}

// GetDate 获取当前编辑的日期
func (ev *EditorView) GetDate() time.Time {
	return ev.currentDate
}

// SetOnContentChange 设置内容变更回调
func (ev *EditorView) SetOnContentChange(callback func(content string)) {
	ev.onContentChange = callback
}

// Clear 清空编辑器
func (ev *EditorView) Clear() {
	ev.editor.SetText("")
	ev.titleLabel.SetText("选择日期以开始编辑")
}

// SetOnSaveComplete 设置保存完成回调
func (ev *EditorView) SetOnSaveComplete(callback func()) {
	ev.onSaveComplete = callback
}

// triggerAutoSave 触发自动保存（延迟 2 秒）
func (ev *EditorView) triggerAutoSave(content string) {
	// 取消之前的定时器
	if ev.saveTimer != nil {
		ev.saveTimer.Stop()
	}

	// 创建新的定时器，2 秒后保存
	ev.saveTimer = time.AfterFunc(2*time.Second, func() {
		ev.saveContent(content)
	})
}

// saveContent 保存内容到任务服务
func (ev *EditorView) saveContent(content string) {
	if ev.taskService == nil {
		util.Warn("任务服务未初始化")
		return
	}

	util.Debug("自动保存任务: %s, 内容长度: %d", 
		ev.currentDate.Format("2006-01-02"), len(content))

	// 调用任务服务保存内容
	err := ev.taskService.SaveTask(ev.currentDate, content)
	if err != nil {
		util.Error("保存任务失败: %v", err)
		// 如果有父窗口，显示错误对话框
		if ev.parentWindow != nil {
			util.ShowErrorDialogWithMessage("保存失败", 
				"无法保存任务内容，请检查文件系统权限", err, ev.parentWindow)
		}
		return
	}

	util.Info("任务保存成功: %s", ev.currentDate.Format("2006-01-02"))

	// 保存成功后调用回调，刷新日历视图
	if ev.onSaveComplete != nil {
		ev.onSaveComplete()
	}
}

// SetParentWindow 设置父窗口（用于显示错误对话框）
func (ev *EditorView) SetParentWindow(window fyne.Window) {
	ev.parentWindow = window
}

// CancelAutoSave 取消自动保存定时器
func (ev *EditorView) CancelAutoSave() {
	if ev.saveTimer != nil {
		ev.saveTimer.Stop()
		ev.saveTimer = nil
	}
}
