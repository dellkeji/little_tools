package ui

import (
	"daily-report-tool/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// PreviewView Markdown 预览视图组件
type PreviewView struct {
	container   *fyne.Container
	titleLabel  *widget.Label
	richText    *widget.RichText
	scrollContainer *container.Scroll
}

// NewPreviewView 创建新的预览视图
func NewPreviewView() *PreviewView {
	pv := &PreviewView{}

	// 创建标题标签
	pv.titleLabel = widget.NewLabel("预览")
	pv.titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 创建 RichText 组件用于显示渲染后的内容
	pv.richText = widget.NewRichText()
	pv.richText.Wrapping = fyne.TextWrapWord

	// 创建滚动容器
	pv.scrollContainer = container.NewScroll(pv.richText)

	// 创建容器布局
	pv.container = container.NewBorder(
		pv.titleLabel,      // top
		nil,                // bottom
		nil,                // left
		nil,                // right
		pv.scrollContainer, // center
	)

	return pv
}

// GetContainer 获取容器
func (pv *PreviewView) GetContainer() *fyne.Container {
	return pv.container
}

// UpdatePreview 更新预览内容
func (pv *PreviewView) UpdatePreview(markdown string) {
	if markdown == "" {
		pv.richText.ParseMarkdown("")
		return
	}

	// 使用 goldmark 将 Markdown 转换为 HTML
	_, err := util.MarkdownToHTML(markdown)
	if err != nil {
		// 如果转换失败，显示原始文本
		pv.richText.ParseMarkdown(markdown)
		return
	}

	// Fyne 的 RichText 支持 Markdown，但不直接支持 HTML
	// 我们使用 ParseMarkdown 来显示内容
	// 注意：这里直接使用 Markdown 而不是 HTML，因为 Fyne 的 RichText 更好地支持 Markdown
	pv.richText.ParseMarkdown(markdown)
}

// Clear 清空预览
func (pv *PreviewView) Clear() {
	pv.richText.ParseMarkdown("")
}
