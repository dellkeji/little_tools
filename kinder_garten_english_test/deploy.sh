#!/bin/bash

# 幼儿园学习APP部署脚本

echo "🚀 开始部署幼儿园学习APP..."

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}请使用 sudo 运行此脚本${NC}"
    exit 1
fi

# 1. 更新系统包
echo -e "${YELLOW}步骤 1: 更新系统包...${NC}"
apt-get update

# 2. 安装 Node.js（如果未安装）
if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}步骤 2: 安装 Node.js...${NC}"
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
    apt-get install -y nodejs
else
    echo -e "${GREEN}Node.js 已安装: $(node -v)${NC}"
fi

# 3. 安装 Nginx（如果未安装）
if ! command -v nginx &> /dev/null; then
    echo -e "${YELLOW}步骤 3: 安装 Nginx...${NC}"
    apt-get install -y nginx
else
    echo -e "${GREEN}Nginx 已安装${NC}"
fi

# 4. 安装 PM2（进程管理器）
if ! command -v pm2 &> /dev/null; then
    echo -e "${YELLOW}步骤 4: 安装 PM2...${NC}"
    npm install -g pm2
else
    echo -e "${GREEN}PM2 已安装${NC}"
fi

# 5. 配置防火墙
echo -e "${YELLOW}步骤 5: 配置防火墙...${NC}"
if command -v ufw &> /dev/null; then
    ufw allow 80/tcp
    ufw allow 443/tcp
    ufw allow 3000/tcp
    echo -e "${GREEN}防火墙规则已添加${NC}"
fi

# 6. 获取当前目录
CURRENT_DIR=$(pwd)
echo -e "${GREEN}项目目录: $CURRENT_DIR${NC}"

# 7. 启动 Node.js 服务器（使用 PM2）
echo -e "${YELLOW}步骤 6: 启动应用服务器...${NC}"
pm2 stop kids-learning-app 2>/dev/null || true
pm2 delete kids-learning-app 2>/dev/null || true
pm2 start server.js --name kids-learning-app
pm2 save
pm2 startup

# 8. 配置 Nginx（可选）
echo -e "${YELLOW}步骤 7: 配置 Nginx（可选）...${NC}"
read -p "是否配置 Nginx 反向代理？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # 创建 Nginx 配置
    cat > /etc/nginx/sites-available/kids-learning-app << EOF
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF

    # 启用站点
    ln -sf /etc/nginx/sites-available/kids-learning-app /etc/nginx/sites-enabled/
    
    # 测试 Nginx 配置
    nginx -t
    
    # 重启 Nginx
    systemctl restart nginx
    
    echo -e "${GREEN}Nginx 配置完成！${NC}"
fi

# 9. 显示访问信息
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 部署完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "📍 访问方式："
echo -e "   本地: http://localhost:3000"
echo -e "   外网: http://$(curl -s ifconfig.me):3000"
echo ""
echo -e "🔧 管理命令："
echo -e "   查看日志: pm2 logs kids-learning-app"
echo -e "   重启应用: pm2 restart kids-learning-app"
echo -e "   停止应用: pm2 stop kids-learning-app"
echo -e "   查看状态: pm2 status"
echo ""
echo -e "${YELLOW}⚠️  注意：请确保云服务器安全组已开放相应端口！${NC}"
echo ""
