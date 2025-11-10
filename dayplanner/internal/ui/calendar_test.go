package ui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
	"daily-report-tool/internal/service"
)

// mockTaskService 用于测试的模拟任务服务
type mockTaskService struct {
	taskDates []time.Time
}

func (m *mockTaskService) GetTask(date time.Time) (*model.Task, error) {
	return nil, nil
}

func (m *mockTaskService) SaveTask(date time.Time, content string) error {
	return nil
}

func (m *mockTaskService) GetMonthTaskDates(year int, month time.Month) ([]time.Time, error) {
	return m.taskDates, nil
}

func (m *mockTaskService) HasTodayTask() (bool, error) {
	return false, nil
}

func TestNewCalendarView(t *testing.T) {
	// 初始化测试应用
	test.NewApp()
	
	// 创建模拟服务
	mockService := &mockTaskService{
		taskDates: []time.Time{},
	}

	// 创建日历视图
	cv := NewCalendarView(mockService)

	// 验证初始状态
	if cv == nil {
		t.Fatal("CalendarView 不应该为 nil")
	}

	now := time.Now()
	if cv.currentYear != now.Year() {
		t.Errorf("期望年份 %d, 得到 %d", now.Year(), cv.currentYear)
	}

	if cv.currentMonth != now.Month() {
		t.Errorf("期望月份 %d, 得到 %d", now.Month(), cv.currentMonth)
	}

	if cv.taskService == nil {
		t.Error("taskService 不应该为 nil")
	}

	if cv.taskDates == nil {
		t.Error("taskDates 不应该为 nil")
	}
}

func TestCalendarNavigation(t *testing.T) {
	// 初始化测试应用
	test.NewApp()
	
	mockService := &mockTaskService{
		taskDates: []time.Time{},
	}

	cv := NewCalendarView(mockService)

	// 测试下一月
	originalMonth := cv.currentMonth
	originalYear := cv.currentYear

	cv.nextMonth()

	if originalMonth == time.December {
		if cv.currentMonth != time.January {
			t.Errorf("12月的下一月应该是1月, 得到 %d", cv.currentMonth)
		}
		if cv.currentYear != originalYear+1 {
			t.Errorf("跨年后年份应该增加, 期望 %d, 得到 %d", originalYear+1, cv.currentYear)
		}
	} else {
		if cv.currentMonth != originalMonth+1 {
			t.Errorf("期望月份 %d, 得到 %d", originalMonth+1, cv.currentMonth)
		}
	}

	// 测试上一月
	cv.previousMonth()

	if cv.currentMonth != originalMonth {
		t.Errorf("返回后应该是原始月份 %d, 得到 %d", originalMonth, cv.currentMonth)
	}
	if cv.currentYear != originalYear {
		t.Errorf("返回后应该是原始年份 %d, 得到 %d", originalYear, cv.currentYear)
	}
}

func TestGetMonthYearString(t *testing.T) {
	// 初始化测试应用
	test.NewApp()
	
	mockService := &mockTaskService{
		taskDates: []time.Time{},
	}

	cv := NewCalendarView(mockService)
	cv.currentYear = 2025
	cv.currentMonth = time.November

	expected := "2025年11月"
	result := cv.getMonthYearString()

	if result != expected {
		t.Errorf("期望 %s, 得到 %s", expected, result)
	}
}

func TestLoadTaskDates(t *testing.T) {
	// 初始化测试应用
	test.NewApp()
	
	// 创建测试日期
	testDates := []time.Time{
		time.Date(2025, time.November, 10, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.November, 15, 0, 0, 0, 0, time.Local),
		time.Date(2025, time.November, 20, 0, 0, 0, 0, time.Local),
	}

	mockService := &mockTaskService{
		taskDates: testDates,
	}

	cv := NewCalendarView(mockService)
	cv.currentYear = 2025
	cv.currentMonth = time.November

	cv.loadTaskDates()

	// 验证任务日期已加载
	if len(cv.taskDates) != 3 {
		t.Errorf("期望加载 3 个任务日期, 得到 %d", len(cv.taskDates))
	}

	// 验证特定日期存在
	if !cv.taskDates["2025-11-10"] {
		t.Error("2025-11-10 应该在任务日期中")
	}
	if !cv.taskDates["2025-11-15"] {
		t.Error("2025-11-15 应该在任务日期中")
	}
	if !cv.taskDates["2025-11-20"] {
		t.Error("2025-11-20 应该在任务日期中")
	}
}

func TestSelectDate(t *testing.T) {
	// 初始化测试应用
	test.NewApp()
	
	mockService := &mockTaskService{
		taskDates: []time.Time{},
	}

	cv := NewCalendarView(mockService)

	// 测试日期选择
	testDate := time.Date(2025, time.November, 15, 0, 0, 0, 0, time.Local)
	
	callbackCalled := false
	var callbackDate time.Time
	
	cv.SetOnDateSelected(func(date time.Time) {
		callbackCalled = true
		callbackDate = date
	})

	cv.selectDate(testDate)

	// 验证选中日期已更新
	if cv.selectedDate != testDate {
		t.Errorf("期望选中日期 %v, 得到 %v", testDate, cv.selectedDate)
	}

	// 验证回调被调用
	if !callbackCalled {
		t.Error("日期选择回调应该被调用")
	}

	if callbackDate != testDate {
		t.Errorf("回调日期期望 %v, 得到 %v", testDate, callbackDate)
	}
}

func TestCalendarWithRealService(t *testing.T) {
	// 初始化测试应用
	test.NewApp()
	
	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 创建真实的服务
	taskRepo := repository.NewFileTaskRepository(tempDir)
	taskService := service.NewTaskService(taskRepo, tempDir)

	// 保存一些测试任务
	testDate := time.Date(2025, time.November, 10, 0, 0, 0, 0, time.Local)
	err := taskService.SaveTask(testDate, "测试任务内容")
	if err != nil {
		t.Fatalf("保存测试任务失败: %v", err)
	}

	// 创建日历视图
	cv := NewCalendarView(taskService)
	cv.currentYear = 2025
	cv.currentMonth = time.November

	// 加载任务日期
	cv.loadTaskDates()

	// 验证任务日期已加载
	if !cv.taskDates["2025-11-10"] {
		t.Error("2025-11-10 应该在任务日期中")
	}
}
