package main

import (
	"fmt"
	"os"

	"daily-report-tool/internal/repository"
	"daily-report-tool/internal/service"
	"daily-report-tool/internal/ui"
	"daily-report-tool/internal/util"

	"fyne.io/fyne/v2/app"
)

const (
	configPath = "./config/config.json"
	dataPath   = "./data/tasks"
	logPath    = "./logs/app.log"
)

func main() {
	// 初始化日志系统
	if err := util.InitLogger(logPath, util.INFO); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}
	defer util.GetLogger().Close()

	util.Info("应用程序启动")

	// 初始化 Fyne 应用程序
	fyneApp := app.NewWithID("com.dailyreport.tool")

	// 创建数据目录
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		util.Error("创建数据目录失败: %v", err)
		fmt.Printf("创建数据目录失败: %v\n", err)
		os.Exit(1)
	}
	util.Info("数据目录已创建: %s", dataPath)

	// 创建配置目录
	configDir := "./config"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		util.Error("创建配置目录失败: %v", err)
		fmt.Printf("创建配置目录失败: %v\n", err)
		os.Exit(1)
	}
	util.Info("配置目录已创建: %s", configDir)

	// 初始化仓库层
	taskRepo := repository.NewFileTaskRepository(dataPath)
	configRepo := repository.NewFileConfigRepository(configPath)

	// 初始化服务层
	taskService := service.NewTaskService(taskRepo, dataPath)
	configService := service.NewConfigService(configRepo)
	reminderService := service.NewReminderService(configService, taskService)

	// 加载配置文件（如果不存在会自动创建默认配置）
	config, err := configService.GetConfig()
	if err != nil {
		util.Error("加载配置失败: %v", err)
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}
	util.Info("配置文件已加载")

	// 启动提醒服务（如果配置启用）
	if config.ReminderEnabled {
		if err := reminderService.Start(); err != nil {
			util.Error("启动提醒服务失败: %v", err)
			fmt.Printf("启动提醒服务失败: %v\n", err)
			// 不退出程序，继续运行主界面
		} else {
			util.Info("提醒服务已启动")
			fmt.Println("提醒服务已启动")
		}
	}

	// 创建并显示主窗口
	mainWindow := ui.NewMainWindow(fyneApp, taskService, configService, reminderService)

	// 设置应用程序退出时的清理逻辑
	mainWindow.GetWindow().SetOnClosed(func() {
		// 停止提醒服务
		reminderService.Stop()
		util.Info("应用程序已退出")
		fmt.Println("应用程序已退出")
	})

	// 显示主窗口（阻塞直到窗口关闭）
	mainWindow.Show()
}
