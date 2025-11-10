#!/usr/bin/env python3
"""
图标转换脚本
将 icon.png 转换为 icon.ico（包含多个尺寸）
"""

import sys
import os

def convert_icon():
    """转换图标文件"""
    
    # 检查 Pillow 是否已安装
    try:
        from PIL import Image
    except ImportError:
        print("错误: 未安装 Pillow 库")
        print("")
        print("请运行以下命令安装:")
        print("  pip install Pillow")
        print("")
        return False
    
    # 检查输入文件
    if not os.path.exists('icon.png'):
        print("错误: 找不到 icon.png 文件")
        print("")
        print("请确保项目根目录有 icon.png 文件")
        print("")
        return False
    
    print("正在转换图标...")
    print("")
    
    try:
        # 打开 PNG 图像
        img = Image.open('icon.png')
        
        # 转换为 RGBA 模式（如果不是）
        if img.mode != 'RGBA':
            img = img.convert('RGBA')
        
        # 定义要生成的尺寸
        sizes = [(16, 16), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)]
        
        print(f"原始图像尺寸: {img.size[0]}x{img.size[1]}")
        print(f"生成 ICO 文件，包含以下尺寸:")
        for size in sizes:
            print(f"  - {size[0]}x{size[1]}")
        print("")
        
        # 保存为 ICO 格式
        img.save('icon.ico', format='ICO', sizes=sizes)
        
        # 检查生成的文件
        if os.path.exists('icon.ico'):
            file_size = os.path.getsize('icon.ico')
            print(f"✓ 转换成功！")
            print(f"  输出文件: icon.ico")
            print(f"  文件大小: {file_size:,} 字节")
            print("")
            print("下一步:")
            print("  运行 build-with-icon.bat (Windows) 或 build-with-icon.sh (Linux/Mac)")
            print("")
            return True
        else:
            print("错误: ICO 文件未生成")
            return False
            
    except Exception as e:
        print(f"错误: 转换失败")
        print(f"详细信息: {e}")
        print("")
        return False

if __name__ == '__main__':
    success = convert_icon()
    sys.exit(0 if success else 1)
