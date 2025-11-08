package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  安装右键菜单: program.exe install")
		fmt.Println("  卸载右键菜单: program.exe uninstall")
		fmt.Println("  打开命令行: program.exe open <目录路径>")
		return
	}

	command := os.Args[1]

	switch command {
	case "install":
		installContextMenu()
	case "uninstall":
		uninstallContextMenu()
	case "open":
		if len(os.Args) < 3 {
			fmt.Println("错误: 需要指定目录路径")
			return
		}
		openCmd(os.Args[2])
	default:
		fmt.Printf("未知命令: %s\n", command)
	}
}

func installContextMenu() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取程序路径失败: %v\n", err)
		return
	}

	// 为文件夹背景添加右键菜单
	err = addRegistryKey(
		`Directory\Background\shell\OpenCmdHere`,
		`在此处打开命令行`,
		exePath,
	)
	if err != nil {
		fmt.Printf("安装失败: %v\n", err)
		return
	}

	// 为文件夹添加右键菜单
	err = addRegistryKey(
		`Directory\shell\OpenCmdHere`,
		`在此处打开命令行`,
		exePath,
	)
	if err != nil {
		fmt.Printf("安装失败: %v\n", err)
		return
	}

	fmt.Println("右键菜单安装成功！")
	fmt.Println("现在可以在文件夹中右键选择 '在此处打开命令行'")
}

func addRegistryKey(shellPath, menuText, exePath string) error {
	// 创建主菜单项
	key, _, err := registry.CreateKey(registry.CLASSES_ROOT, shellPath, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("创建注册表项失败: %v", err)
	}
	defer key.Close()

	err = key.SetStringValue("", menuText)
	if err != nil {
		return fmt.Errorf("设置菜单文本失败: %v", err)
	}

	// 设置图标（使用 cmd.exe 的图标）
	err = key.SetStringValue("Icon", "cmd.exe")
	if err != nil {
		return fmt.Errorf("设置图标失败: %v", err)
	}

	// 创建 command 子项
	cmdKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, shellPath+`\command`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("创建 command 项失败: %v", err)
	}
	defer cmdKey.Close()

	// 设置命令
	cmdValue := fmt.Sprintf(`"%s" open "%%V"`, exePath)
	err = cmdKey.SetStringValue("", cmdValue)
	if err != nil {
		return fmt.Errorf("设置命令失败: %v", err)
	}

	return nil
}

func uninstallContextMenu() {
	// 删除文件夹背景的右键菜单
	err := registry.DeleteKey(registry.CLASSES_ROOT, `Directory\Background\shell\OpenCmdHere\command`)
	if err == nil {
		registry.DeleteKey(registry.CLASSES_ROOT, `Directory\Background\shell\OpenCmdHere`)
	}

	// 删除文件夹的右键菜单
	err = registry.DeleteKey(registry.CLASSES_ROOT, `Directory\shell\OpenCmdHere\command`)
	if err == nil {
		registry.DeleteKey(registry.CLASSES_ROOT, `Directory\shell\OpenCmdHere`)
	}

	fmt.Println("右键菜单已卸载")
}

func openCmd(dir string) {
	// 确保目录存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("目录不存在: %s\n", dir)
		return
	}

	// 获取绝对路径
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Printf("获取绝对路径失败: %v\n", err)
		return
	}

	// 启动命令提示符
	cmd := fmt.Sprintf(`start cmd.exe /K "cd /d %s"`, absDir)
	err = runCommand(cmd)
	if err != nil {
		fmt.Printf("打开命令行失败: %v\n", err)
	}
}

func runCommand(command string) error {
	// 使用 cmd.exe 执行命令
	cmd := exec.Command("cmd.exe", "/C", command)
	return cmd.Run()
}
