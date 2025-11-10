package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"daily-report-tool/internal/service"
	"daily-report-tool/internal/util"
)

// CalendarView 日历视图组件
type CalendarView struct {
	widget.BaseWidget
	
	// 状态
	currentYear  int
	currentMonth time.Month
	selectedDate time.Time
	taskDates    map[string]bool // 存储有任务的日期，格式: "2006-01-02"
	
	// 服务
	taskService service.TaskService
	
	// UI 组件
	monthLabel   *widget.Label
	prevButton   *widget.Button
	nextButton   *widget.Button
	dateButtons  [][]*widget.Button // 7x6 网格
	container    *fyne.Container
	
	// 回调函数
	onDateSelected func(date time.Time)
}

// NewCalendarView 创建新的日历视图
func NewCalendarView(taskService service.TaskService) *CalendarView {
	cv := &CalendarView{
		taskService: taskService,
		taskDates:   make(map[string]bool),
	}
	
	// 初始化为当前年月
	now := time.Now()
	cv.currentYear = now.Year()
	cv.currentMonth = now.Month()
	cv.selectedDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	
	cv.ExtendBaseWidget(cv)
	cv.buildUI()
	
	return cv
}

// buildUI 构建用户界面
func (cv *CalendarView) buildUI() {
	// 创建月份标签
	cv.monthLabel = widget.NewLabel(cv.getMonthYearString())
	
	// 创建导航按钮
	cv.prevButton = widget.NewButton("上一月", cv.previousMonth)
	cv.nextButton = widget.NewButton("下一月", cv.nextMonth)
	
	// 创建顶部导航栏
	header := container.NewBorder(
		nil, nil,
		cv.prevButton,
		cv.nextButton,
		container.NewCenter(cv.monthLabel),
	)
	
	// 创建星期标题行
	weekdayLabels := cv.createWeekdayLabels()
	
	// 创建日期网格（将在 renderCalendar 中填充）
	cv.dateButtons = make([][]*widget.Button, 6)
	for i := 0; i < 6; i++ {
		cv.dateButtons[i] = make([]*widget.Button, 7)
	}
	
	// 创建日历网格容器
	calendarGrid := cv.createCalendarGrid()
	
	// 组合所有组件
	cv.container = container.NewVBox(
		header,
		weekdayLabels,
		calendarGrid,
	)
	
	// 初始渲染
	cv.refresh()
}

// createWeekdayLabels 创建星期标题行
func (cv *CalendarView) createWeekdayLabels() *fyne.Container {
	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	labels := make([]fyne.CanvasObject, 7)
	
	for i, day := range weekdays {
		label := widget.NewLabel(day)
		labels[i] = container.NewCenter(label)
	}
	
	return container.NewGridWithColumns(7, labels...)
}

// createCalendarGrid 创建日历网格
func (cv *CalendarView) createCalendarGrid() *fyne.Container {
	cells := make([]fyne.CanvasObject, 42) // 6行 x 7列
	
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			index := row*7 + col
			btn := widget.NewButton("", nil)
			cv.dateButtons[row][col] = btn
			cells[index] = btn
		}
	}
	
	return container.NewGridWithColumns(7, cells...)
}

// previousMonth 切换到上一月
func (cv *CalendarView) previousMonth() {
	cv.currentMonth--
	if cv.currentMonth < time.January {
		cv.currentMonth = time.December
		cv.currentYear--
	}
	cv.refresh()
}

// nextMonth 切换到下一月
func (cv *CalendarView) nextMonth() {
	cv.currentMonth++
	if cv.currentMonth > time.December {
		cv.currentMonth = time.January
		cv.currentYear++
	}
	cv.refresh()
}

// getMonthYearString 获取月份年份字符串
func (cv *CalendarView) getMonthYearString() string {
	return fmt.Sprintf("%d年%d月", cv.currentYear, cv.currentMonth)
}

// SetOnDateSelected 设置日期选择回调函数
func (cv *CalendarView) SetOnDateSelected(callback func(date time.Time)) {
	cv.onDateSelected = callback
}

// CreateRenderer 实现 fyne.Widget 接口
func (cv *CalendarView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(cv.container)
}

// GetSelectedDate 获取当前选中的日期
func (cv *CalendarView) GetSelectedDate() time.Time {
	return cv.selectedDate
}

// refresh 刷新日历显示
func (cv *CalendarView) refresh() {
	// 更新月份标签
	cv.monthLabel.SetText(cv.getMonthYearString())
	
	// 加载当月有任务的日期
	cv.loadTaskDates()
	
	// 渲染日历网格
	cv.renderCalendar()
}

// loadTaskDates 加载当月有任务的日期
func (cv *CalendarView) loadTaskDates() {
	// 清空现有的任务日期
	cv.taskDates = make(map[string]bool)
	
	util.Debug("加载任务日期: %d年%d月", cv.currentYear, cv.currentMonth)
	
	// 从服务获取当月有任务的日期
	dates, err := cv.taskService.GetMonthTaskDates(cv.currentYear, cv.currentMonth)
	if err != nil {
		// 如果出错，记录日志但不阻塞UI
		util.Error("加载任务日期失败: %v", err)
		fmt.Printf("加载任务日期失败: %v\n", err)
		return
	}
	
	util.Debug("找到 %d 个有任务的日期", len(dates))
	
	// 将日期存储到 map 中
	for _, date := range dates {
		dateStr := date.Format("2006-01-02")
		cv.taskDates[dateStr] = true
	}
}

// renderCalendar 渲染日历网格
func (cv *CalendarView) renderCalendar() {
	// 获取当月第一天
	firstDay := time.Date(cv.currentYear, cv.currentMonth, 1, 0, 0, 0, 0, time.Local)
	
	// 获取第一天是星期几 (0=Sunday, 1=Monday, ..., 6=Saturday)
	firstWeekday := int(firstDay.Weekday())
	
	// 获取当月天数
	lastDay := firstDay.AddDate(0, 1, 0).Add(-24 * time.Hour)
	daysInMonth := lastDay.Day()
	
	// 当前日期（用于高亮今天）
	today := time.Now()
	
	// 填充日历网格
	day := 1
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			btn := cv.dateButtons[row][col]
			
			// 计算当前单元格应该显示的日期
			cellIndex := row*7 + col
			
			if cellIndex < firstWeekday || day > daysInMonth {
				// 空白单元格
				btn.SetText("")
				btn.OnTapped = nil
				btn.Importance = widget.MediumImportance
			} else {
				// 有效日期
				currentDay := day
				btn.SetText(fmt.Sprintf("%d", currentDay))
				
				// 创建日期对象
				date := time.Date(cv.currentYear, cv.currentMonth, currentDay, 0, 0, 0, 0, time.Local)
				dateStr := date.Format("2006-01-02")
				
				// 设置点击事件
				btn.OnTapped = func() {
					cv.selectDate(date)
				}
				
				// 设置按钮样式
				if cv.taskDates[dateStr] {
					// 有任务的日期 - 使用高重要性样式
					btn.Importance = widget.HighImportance
				} else if date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day() {
					// 今天 - 使用警告样式
					btn.Importance = widget.WarningImportance
				} else if date.Year() == cv.selectedDate.Year() && date.Month() == cv.selectedDate.Month() && date.Day() == cv.selectedDate.Day() {
					// 选中的日期 - 使用成功样式
					btn.Importance = widget.SuccessImportance
				} else {
					// 普通日期
					btn.Importance = widget.MediumImportance
				}
				
				day++
			}
			
			btn.Refresh()
		}
	}
}

// selectDate 选择日期
func (cv *CalendarView) selectDate(date time.Time) {
	cv.selectedDate = date
	cv.renderCalendar() // 重新渲染以更新选中状态
	
	// 触发回调
	if cv.onDateSelected != nil {
		cv.onDateSelected(date)
	}
}

// Refresh 刷新日历（公开方法，供外部调用）
func (cv *CalendarView) Refresh() {
	cv.refresh()
}
