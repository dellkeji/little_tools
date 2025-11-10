# 图标转换快速指南

## 方法 1: 在线转换（最简单）

1. 访问 https://convertio.co/zh/png-ico/
2. 上传项目中的 `icon.png` 文件
3. 点击"转换"
4. 下载生成的 `icon.ico` 文件
5. 将 `icon.ico` 保存到项目根目录

## 方法 2: 使用 PowerShell（Windows 10+）

如果你有 PowerShell，可以使用以下脚本：

```powershell
# 需要先安装 ImageMagick
# 下载: https://imagemagick.org/script/download.php#windows

magick convert icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico
```

## 方法 3: 使用 Python（如果已安装）

```python
# 安装 Pillow
# pip install Pillow

from PIL import Image

img = Image.open('icon.png')
img.save('icon.ico', format='ICO', sizes=[(16,16), (32,32), (48,48), (256,256)])
```

## 完成后

运行构建脚本：

```bash
# Windows
.\build-with-icon.bat

# Linux/Mac
./build-with-icon.sh
```

## 推荐的图标尺寸

原始 PNG 图标建议尺寸：
- 最小: 256x256 像素
- 推荐: 512x512 像素
- 最佳: 1024x1024 像素（用于高 DPI 显示器）

ICO 文件应包含多个尺寸：
- 16x16 - 小图标
- 32x32 - 标准图标
- 48x48 - 大图标
- 256x256 - 高清图标
