# Windows 右键打开命令行工具

这是一个用 Go 编写的工具，可以在 Windows 资源管理器中添加右键菜单项，快速在当前目录打开命令提示符。

## 编译

在 Windows 上编译：
```bash
go build -o opencmd.exe
```

在 macOS/Linux 上交叉编译为 Windows 可执行文件：
```bash
GOOS=windows GOARCH=amd64 go build -o opencmd.exe
```

## 使用方法

### 1. 安装右键菜单（需要管理员权限）

以管理员身份运行命令提示符，然后执行：
```cmd
opencmd.exe install
```

### 2. 使用右键菜单

安装后，在文件夹中右键点击空白处，选择"在此处打开命令行"即可。

### 3. 卸载右键菜单

以管理员身份运行：
```cmd
opencmd.exe uninstall
```

## 功能特点

- 在文件夹内右键打开命令行
- 在文件夹上右键打开命令行
- 使用 cmd.exe 图标
- 简单的安装和卸载

## 注意事项

- 安装和卸载需要管理员权限
- 仅支持 Windows 系统
- 修改注册表，请谨慎使用
