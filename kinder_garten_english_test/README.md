# 幼儿园学习APP

一个适合幼儿园儿童的学习应用，包含英语、数学和语文三个科目的学习与测试功能。

## 功能特点

### 📚 三大学习模块
- **🔤 英语学习**：基础单词学习（Apple、Ball、Cat等）+ 🔊 真人发音
- **🔢 数学学习**：数字认知、简单加减法、形状识别
- **📖 语文学习**：汉字、拼音学习

### ✏️ 互动测试
- 每个科目都有配套的测试题
- 即时反馈答题结果
- 显示测试分数和正确率
- 英语测试支持单词发音，帮助孩子听音辨词

### 🔊 语音功能
- 使用浏览器内置的 Web Speech API
- 英语单词标准发音
- 点击 🔊 按钮即可听发音
- 学习和测试模式都支持发音

## 如何使用

### 方式一：本地浏览器直接打开

1. 双击打开 `index.html` 文件
2. 或者在命令行运行：
   ```bash
   # Windows
   start index.html
   
   # Linux/Mac
   open index.html
   ```

### 方式二：启动 Web 服务器（推荐）

#### 在 Linux 服务器上部署（外网访问）

**快速部署（一键脚本）：**
```bash
# 1. 赋予执行权限
chmod +x deploy.sh

# 2. 运行部署脚本
sudo ./deploy.sh
```

**手动部署：**
```bash
# 1. 安装 Node.js（如果未安装）
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 2. 安装 PM2 进程管理器
sudo npm install -g pm2

# 3. 启动应用
node server.js
# 或使用 PM2（推荐，支持自动重启）
pm2 start server.js --name kids-learning-app
pm2 save
pm2 startup

# 4. 开放防火墙端口
sudo ufw allow 3000/tcp
sudo ufw allow 80/tcp
```

**访问应用：**
- 本地访问：`http://localhost:3000`
- 外网访问：`http://你的服务器IP:3000`

#### 在 Windows 上运行服务器

```bash
# 安装 Node.js 后运行
node server.js
```

### 方式三：使用 Nginx 反向代理（生产环境推荐）

```bash
# 1. 安装 Nginx
sudo apt-get install nginx

# 2. 复制配置文件
sudo cp nginx.conf /etc/nginx/sites-available/kids-learning-app

# 3. 编辑配置文件，修改域名和路径
sudo nano /etc/nginx/sites-available/kids-learning-app

# 4. 启用站点
sudo ln -s /etc/nginx/sites-available/kids-learning-app /etc/nginx/sites-enabled/

# 5. 测试配置
sudo nginx -t

# 6. 重启 Nginx
sudo systemctl restart nginx
```

配置完成后，可以通过 `http://你的域名` 或 `http://你的IP` 访问（80端口）。

## 打包为独立应用

### 打包为 Windows 应用

可以使用以下工具将网页打包为 Windows 应用：

1. **使用 Electron**（需要安装 Node.js）：
   ```
   npm install -g electron
   ```

2. **使用 NW.js**：下载 NW.js 并将项目文件放入其中

3. **使用在线工具**：
   - ToDesktop (https://www.todesktop.com/)
   - Electron Forge (https://www.electronforge.io/)

### 打包为 Android APK

可以使用以下方法：

1. **使用 Apache Cordova**：
   ```
   npm install -g cordova
   cordova create KidsLearningApp
   # 将 HTML/CSS/JS 文件复制到 www 目录
   cordova platform add android
   cordova build android
   ```

2. **使用在线工具**：
   - AppGeyser (https://appgeyser.com/)
   - Appy Pie (https://www.appypie.com/)
   - WebIntoApp (https://webintoapp.com/)

## 项目结构

```
kids-learning-app/
├── index.html      # 主页面
├── styles.css      # 样式文件
├── app.js          # 应用逻辑
├── package.json    # 项目配置
└── README.md       # 说明文档
```

## 自定义内容

你可以在 `app.js` 文件中的 `learningData` 对象中修改或添加学习内容：

```javascript
const learningData = {
    english: { lessons: [...], tests: [...] },
    math: { lessons: [...], tests: [...] },
    chinese: { lessons: [...], tests: [...] }
};
```

## 技术栈

- HTML5
- CSS3（响应式设计）
- 原生 JavaScript（无需框架）

## 浏览器兼容性

- ✅ Chrome/Edge (推荐)
- ✅ Firefox
- ✅ Safari
- ✅ 其他现代浏览器

## 未来改进

- [ ] 添加语音朗读功能
- [ ] 添加更多学习内容
- [ ] 添加进度保存功能
- [ ] 添加奖励系统
- [ ] 添加家长监控面板

## 许可证

MIT License


## 云服务器部署注意事项

### 1. 安全组配置
确保云服务器安全组已开放以下端口：
- **3000**：Node.js 应用端口
- **80**：HTTP 端口（使用 Nginx 时）
- **443**：HTTPS 端口（使用 SSL 时）

### 2. 防火墙配置
```bash
# Ubuntu/Debian
sudo ufw allow 3000/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=3000/tcp
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --reload
```

### 3. 进程管理
使用 PM2 管理应用进程：
```bash
# 查看应用状态
pm2 status

# 查看日志
pm2 logs kids-learning-app

# 重启应用
pm2 restart kids-learning-app

# 停止应用
pm2 stop kids-learning-app

# 开机自启动
pm2 startup
pm2 save
```

### 4. 域名配置（可选）
如果有域名，可以配置 DNS 解析：
1. 在域名服务商添加 A 记录，指向服务器 IP
2. 修改 `nginx.conf` 中的 `server_name` 为你的域名
3. 重启 Nginx

### 5. HTTPS 配置（推荐）
使用 Let's Encrypt 免费 SSL 证书：
```bash
# 安装 Certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取证书并自动配置 Nginx
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

## 性能优化建议

1. **启用 Gzip 压缩**：已在 Nginx 配置中启用
2. **静态资源缓存**：已配置 1 年缓存
3. **CDN 加速**：可将静态文件上传到 CDN
4. **负载均衡**：高并发时可配置多个 Node.js 实例

## 监控和维护

```bash
# 查看服务器资源使用
htop

# 查看应用日志
pm2 logs kids-learning-app --lines 100

# 查看 Nginx 日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# 重启所有服务
pm2 restart all
sudo systemctl restart nginx
```
