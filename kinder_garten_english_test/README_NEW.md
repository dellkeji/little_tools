# 幼儿园学习APP 🎈

一个适合幼儿园儿童的学习应用，包含英语、数学和语文三个科目的学习与测试功能。

## ✨ 核心特性

### 📚 丰富的学习内容库
- **🔤 英语学习**：60+ 单词，涵盖水果、动物、颜色、食物、交通、天气、衣服、运动、乐器等 12 个分类
- **🔢 数学学习**：25+ 知识点，包括数字认知、加减法运算、形状识别、比较大小
- **📖 语文学习**：60+ 汉字，包括数字、方位词、形容词、动物、颜色、动词等 7 个分类

**总计：145+ 学习内容，125+ 测试题目**

### 🔄 智能随机刷新
- **学习模式**：每次随机显示 15 个内容，点击"🔄 换一批"即可刷新
- **测试模式**：每次随机抽取 10 道题，避免重复，保持新鲜感
- 支持无限刷新，每次都能看到不同的内容

### 🔊 真人语音发音
- 所有英语单词支持标准发音
- 学习模式：每个单词旁边有发音按钮
- 测试模式：问题和选项都可以发音
- 语速适中，适合儿童学习

### 📊 即时测试反馈
- 答对显示绿色 ✓，答错显示红色 ✗
- 自动统计得分和正确率
- 根据成绩给予不同鼓励

### 🎨 友好的界面设计
- 色彩鲜艳，大量使用表情符号
- 内容分类显示，清晰易懂
- 大字体，适合儿童阅读
- 流畅的动画效果

## 🚀 快速开始

### 方式一：本地浏览器直接打开（最简单）

直接双击 `index.html` 文件，或在浏览器中打开即可使用。

### 方式二：启动 Web 服务器（推荐用于外网访问）

#### 在 Linux 服务器上部署

**一键部署：**
```bash
chmod +x deploy.sh
sudo ./deploy.sh
```

**手动部署：**
```bash
# 1. 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 2. 启动应用
node server.js

# 3. 使用 PM2（推荐）
npm install -g pm2
pm2 start server.js --name kids-learning-app
pm2 save
pm2 startup
```

**访问应用：**
- 本地：`http://localhost:3000`
- 外网：`http://你的服务器IP:3000`

#### 在 Windows 上运行

```bash
node server.js
```

## 📖 使用说明

### 学习流程
1. 选择科目（英语/数学/语文）
2. 点击"📚 开始学习"浏览内容
3. 点击"🔄 换一批"查看更多内容
4. 点击"✏️ 开始测试"进行练习
5. 查看成绩，点击"🔄 再测一次"继续

### 英语学习技巧
- 点击 🔊 按钮听发音
- 跟读单词，练习发音
- 利用"换一批"功能学习更多单词
- 测试时可以点击选项的发音按钮

### 数学学习技巧
- 先学习数字和基本概念
- 理解加减法的图形表示
- 多次测试巩固计算能力

### 语文学习技巧
- 注意拼音的声调
- 理解汉字的意思
- 按分类学习，便于记忆

## 📁 项目结构

```
kids-learning-app/
├── index.html          # 主页面
├── app.js             # 应用逻辑（含随机刷新功能）
├── learningData.js    # 学习内容数据库（145+ 内容）
├── styles.css         # 样式文件
├── server.js          # Node.js 服务器
├── deploy.sh          # 一键部署脚本
├── nginx.conf         # Nginx 配置
├── package.json       # 项目配置
├── README.md          # 说明文档
└── FEATURES.md        # 详细功能说明
```

## 🎯 内容统计

| 科目 | 学习内容 | 测试题目 | 分类数 |
|------|---------|---------|--------|
| 英语 | 60+ 单词 | 40+ 题 | 12 类 |
| 数学 | 25+ 知识点 | 45+ 题 | 多种类型 |
| 语文 | 60+ 汉字 | 40+ 题 | 7 类 |

## 🔧 自定义内容

### 添加新的学习内容

编辑 `learningData.js` 文件：

```javascript
// 添加英语单词
learningData.english.lessons.push({
    word: 'NewWord',
    translation: '新单词',
    emoji: '🎯',
    category: '分类名'
});

// 添加测试题
learningData.english.tests.push({
    question: '问题？',
    options: ['选项1', '选项2', '选项3', '选项4'],
    answer: 0  // 正确答案的索引（0-3）
});
```

### 调整显示数量

在 `app.js` 中修改：
```javascript
const lessonsPerPage = 15;  // 学习模式每次显示数量
const questionsPerTest = 10; // 测试模式每次题目数量
```

## 🌐 部署到云服务器

### 安全组配置
开放端口：3000（应用）、80（HTTP）、443（HTTPS）

### 使用 Nginx 反向代理
```bash
sudo cp nginx.conf /etc/nginx/sites-available/kids-learning-app
sudo ln -s /etc/nginx/sites-available/kids-learning-app /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 配置 HTTPS（推荐）
```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## 💻 技术栈

- **前端**：HTML5 + CSS3 + 原生 JavaScript
- **语音**：Web Speech API
- **后端**：Node.js (可选)
- **部署**：Nginx + PM2

## 📱 浏览器兼容性

- ✅ Chrome/Edge (推荐)
- ✅ Firefox
- ✅ Safari
- ✅ 其他现代浏览器

## 🎓 教育价值

- **英语启蒙**：通过图像和发音建立单词认知
- **数学思维**：培养基础计算和逻辑能力
- **汉字学习**：掌握常用汉字和拼音
- **自主学习**：随机刷新保持学习兴趣
- **即时反馈**：帮助孩子及时纠正错误

## 📝 更新日志

### v2.0 (最新)
- ✨ 新增随机刷新功能
- 📚 扩充学习内容至 145+ 项
- 🔄 学习模式支持"换一批"
- 🎲 测试模式随机抽题
- 📊 优化统计显示

### v1.0
- 🎉 初始版本发布
- 📚 基础学习内容
- 🔊 语音发音功能
- ✏️ 测试功能

## 🤝 贡献

欢迎提交问题和建议！可以通过以下方式扩展：
- 添加更多单词和汉字
- 增加新的学习分类
- 优化界面设计
- 添加学习进度记录
- 支持多语言

## 📄 许可证

MIT License

## 📞 支持

如有问题，请查看 `FEATURES.md` 了解详细功能说明。

---

**祝小朋友们学习愉快！🎉**
