#!/bin/bash

set -e

echo "=========================================="
echo "  PromQL Insights - Build & Install"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}Checking dependencies...${NC}"
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
        echo "Visit: https://golang.org/dl/"
        exit 1
    fi
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓ Go installed: $GO_VERSION${NC}"
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${RED}Error: Node.js is not installed. Please install Node.js first.${NC}"
        echo "Visit: https://nodejs.org/"
        exit 1
    fi
    NODE_VERSION=$(node -v)
    echo -e "${GREEN}✓ Node.js installed: $NODE_VERSION${NC}"
    
    # 检查 npm
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}Error: npm is not installed. Please install npm first.${NC}"
        exit 1
    fi
    NPM_VERSION=$(npm -v)
    echo -e "${GREEN}✓ npm installed: $NPM_VERSION${NC}"
    
    # 检查 Docker（可选）
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker -v)
        echo -e "${GREEN}✓ Docker installed: $DOCKER_VERSION${NC}"
    else
        echo -e "${YELLOW}⚠ Docker is not installed. Optional for Prometheus deployment.${NC}"
    fi
    
    echo ""
}

# 构建后端
build_backend() {
    echo -e "${YELLOW}Building backend...${NC}"
    cd backend
    
    echo "Downloading Go dependencies..."
    go mod download
    
    echo "Building backend binary..."
    go build -o bin/api cmd/api/main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Backend built successfully: backend/bin/api${NC}"
    else
        echo -e "${RED}✗ Backend build failed${NC}"
        exit 1
    fi
    
    cd ..
    echo ""
}

# 构建前端
build_frontend() {
    echo -e "${YELLOW}Building frontend...${NC}"
    cd frontend
    
    echo "Installing npm dependencies..."
    npm install
    
    echo "Building frontend..."
    npm run build
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Frontend built successfully: frontend/dist${NC}"
    else
        echo -e "${RED}✗ Frontend build failed${NC}"
        exit 1
    fi
    
    cd ..
    echo ""
}

# 启动 Docker 服务（可选）
start_docker() {
    if command -v docker &> /dev/null; then
        echo -e "${YELLOW}Starting Docker services...${NC}"
        
        # 检查容器是否已存在
        if docker ps -a | grep -q "prometheus"; then
            echo -e "${GREEN}✓ Prometheus container already running${NC}"
        else
            echo "Starting Prometheus container..."
            docker run -d --name prometheus --network insights_monitoring -p 9090:9090 \
                -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
                prom/prometheus:v2.48.0 \
                --config.file=/etc/prometheus/prometheus.yml \
                --storage.tsdb.path=/prometheus \
                --storage.tsdb.retention.time=30d \
                --web.enable-lifecycle
            echo -e "${GREEN}✓ Prometheus started on port 9090${NC}"
        fi
        
        if docker ps -a | grep -q "pushgateway"; then
            echo -e "${GREEN}✓ Pushgateway container already running${NC}"
        else
            echo "Starting Pushgateway container..."
            docker run -d --name pushgateway --network insights_monitoring -p 9091:9091 \
                prom/pushgateway:v1.6.2
            echo -e "${GREEN}✓ Pushgateway started on port 9091${NC}"
        fi
        
        echo ""
    fi
}

# 显示使用说明
show_usage() {
    echo -e "${GREEN}=========================================="
    echo "  Build & Install Complete!"
    echo "==========================================${NC}"
    echo ""
    echo "Services:"
    echo "  - Backend API:  http://localhost:8888"
    echo "  - Frontend App: http://localhost:3000"
    echo "  - Prometheus:    http://localhost:9090"
    echo "  - Pushgateway:   http://localhost:9091"
    echo ""
    echo "To start services:"
    echo "  1. Start backend:"
    echo "     cd backend && ./bin/api"
    echo ""
    echo "  2. Start frontend (dev mode):"
    echo "     cd frontend && npm run dev"
    echo ""
    echo "  3. Start frontend (production):"
    echo "     cd frontend && npm run preview"
    echo ""
    echo "  4. Push metrics to Prometheus (optional):"
    echo "     python3 push_metrics.py"
    echo ""
    echo "To stop Docker services:"
    echo "     docker stop prometheus pushgateway"
    echo ""
    echo "To view logs:"
    echo "     docker logs prometheus"
    echo "     docker logs pushgateway"
    echo ""
}

# 主函数
main() {
    check_dependencies
    build_backend
    build_frontend
    start_docker
    show_usage
    
    echo -e "${GREEN}All done! 🎉${NC}"
}

# 运行主函数
main
