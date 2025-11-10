package service

import (
	"fmt"
	"os"
	"time"

	"daily-report-tool/internal/model"
	"daily-report-tool/internal/repository"
)

// TaskService 定义任务管理服务接口
type TaskService interface {
	// GetTask 获取指定日期的任务
	GetTask(date time.Time) (*model.Task, error)

	// SaveTask 保存或更新任务
	SaveTask(date time.Time, content string) error

	// GetMonthTaskDates 获取月份内有任务的日期
	GetMonthTaskDates(year int, month time.Month) ([]time.Time, error)

	// HasTodayTask 检查今天是否有任务
	HasTodayTask() (bool, error)
}

// TaskServiceImpl 任务管理服务实现
type TaskServiceImpl struct {
	taskRepo repository.TaskRepository
	dataPath string
}

// NewTaskService 创建新的任务管理服务
func NewTaskService(taskRepo repository.TaskRepository, dataPath string) *TaskServiceImpl {
	return &TaskServiceImpl{
		taskRepo: taskRepo,
		dataPath: dataPath,
	}
}

// GetTask 获取指定日期的任务
func (s *TaskServiceImpl) GetTask(date time.Time) (*model.Task, error) {
	// 确保数据目录存在
	if err := s.ensureDataDirectory(); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.GetByDate(date)
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	return task, nil
}

// SaveTask 保存或更新任务
func (s *TaskServiceImpl) SaveTask(date time.Time, content string) error {
	// 确保数据目录存在
	if err := s.ensureDataDirectory(); err != nil {
		return err
	}

	// 获取现有任务
	existingTask, err := s.taskRepo.GetByDate(date)
	if err != nil {
		return fmt.Errorf("获取现有任务失败: %w", err)
	}

	now := time.Now()
	var task *model.Task

	if existingTask != nil {
		// 更新现有任务
		task = existingTask
		task.Content = content
		task.UpdatedAt = now
	} else {
		// 创建新任务
		task = &model.Task{
			Date:      date,
			Content:   content,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	if err := s.taskRepo.Save(task); err != nil {
		return fmt.Errorf("保存任务失败: %w", err)
	}

	return nil
}

// GetMonthTaskDates 获取月份内有任务的日期
func (s *TaskServiceImpl) GetMonthTaskDates(year int, month time.Month) ([]time.Time, error) {
	// 确保数据目录存在
	if err := s.ensureDataDirectory(); err != nil {
		return nil, err
	}

	// 计算月份的开始和结束日期
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second) // 下个月第一天减一秒

	dates, err := s.taskRepo.GetTaskDates(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("获取月份任务日期失败: %w", err)
	}

	return dates, nil
}

// HasTodayTask 检查今天是否有任务
func (s *TaskServiceImpl) HasTodayTask() (bool, error) {
	// 确保数据目录存在
	if err := s.ensureDataDirectory(); err != nil {
		return false, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	hasTask, err := s.taskRepo.HasTask(today)
	if err != nil {
		return false, fmt.Errorf("检查今天任务失败: %w", err)
	}

	return hasTask, nil
}

// ensureDataDirectory 确保数据目录存在
func (s *TaskServiceImpl) ensureDataDirectory() error {
	if err := os.MkdirAll(s.dataPath, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}
	return nil
}
