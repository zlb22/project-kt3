# 互动学习场景下青少年创新能力智能化测评工具 V2.0

一个现代化的前后端分离教育测评系统，从 Streamlit 重构为 FastAPI + React 架构。

## 项目概述

本项目是一个综合性的教育测评工具，包含以下主要功能：
- 用户认证和权限管理
- 问题解决能力测评
- 互动讨论场景测评
- AUT（自闭症）测试功能
- 情绪评估测试
- 在线课程功能
- 结果评估和报告生成

## 技术栈

### 后端
- **FastAPI**: 现代、快速的 Web 框架
- **JWT**: 身份验证和授权
- **Pydantic**: 数据验证
- **Python 3.8+**: 编程语言

### 前端
- **React 18**: 用户界面库
- **TypeScript**: 类型安全的 JavaScript
- **Material-UI**: UI 组件库
- **React Router**: 路由管理
- **Axios**: HTTP 客户端

## 项目结构

```
project-kt3/
├── backend/                 # FastAPI 后端
│   ├── main.py             # 主应用文件
│   ├── requirements.txt    # Python 依赖
│   └── users.json         # 用户数据（临时存储）
├── frontend/               # React 前端
│   ├── public/            # 静态文件
│   ├── src/               # 源代码
│   │   ├── components/    # 可复用组件
│   │   ├── contexts/      # React Context
│   │   ├── pages/         # 页面组件
│   │   ├── App.tsx        # 主应用组件
│   │   └── index.tsx      # 入口文件
│   └── package.json       # Node.js 依赖
└── README.md              # 项目文档
```

## 环境要求

### 后端环境
- Python 3.8 或更高版本
- pip（Python 包管理器）

### 前端环境
- Node.js 16.0 或更高版本
- npm 或 yarn 包管理器

## 安装和配置

### 1. 克隆项目

```bash
git clone <repository-url>
cd project-kt3
```

### 2. 后端设置

```bash
# 进入后端目录
cd backend

# 创建虚拟环境（推荐）
python -m venv venv

# 激活虚拟环境
# Linux/Mac:
source venv/bin/activate
# Windows:
# venv\Scripts\activate

# 安装依赖
pip install -r requirements.txt
```

### 3. 前端设置

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install
# 或使用 yarn
# yarn install
```

## 运行应用

### 启动后端服务

```bash
# 在 backend 目录下
cd backend

# 激活虚拟环境（如果使用）
source venv/bin/activate

# 启动 FastAPI 服务器
python main.py

# 或使用 uvicorn 直接启动
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

后端服务将在 `http://localhost:8000` 启动
- API 文档: `http://localhost:8000/docs`
- 健康检查: `http://localhost:8000/api/health`

### 启动前端服务

```bash
# 在 frontend 目录下
cd frontend

# 启动开发服务器
npm start
# 或使用 yarn
# yarn start
```

前端应用将在 `http://localhost:3000` 启动

## 默认用户账户

- **用户名**: admin
- **密码**: 1

## API 端点

### 认证相关
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息
- `POST /api/auth/change-password` - 修改密码

### 测试相关
- `GET /api/tests/available` - 获取可用测试列表
- `GET /api/tests/{test_id}` - 获取特定测试详情

### 系统相关
- `GET /api/health` - 健康检查

## 开发指南

### 后端开发
1. 所有 API 端点都应该包含适当的错误处理
2. 使用 Pydantic 模型进行数据验证
3. 遵循 RESTful API 设计原则
4. 添加适当的日志记录

### 前端开发
1. 使用 TypeScript 确保类型安全
2. 遵循 React Hooks 最佳实践
3. 使用 Material-UI 组件保持一致的设计
4. 实现响应式设计

## 部署

### 后端部署
```bash
# 生产环境运行
uvicorn main:app --host 0.0.0.0 --port 8000
```

### 前端部署
```bash
# 构建生产版本
npm run build

# 构建文件将在 build/ 目录中
```

## 故障排除

### 常见问题

1. **端口冲突**: 确保 8000 和 3000 端口未被占用
2. **CORS 错误**: 检查后端 CORS 配置是否正确
3. **依赖安装失败**: 确保 Python 和 Node.js 版本符合要求
4. **认证失败**: 检查 JWT 密钥配置

### 日志查看
- 后端日志会在控制台输出
- 前端错误可在浏览器开发者工具中查看

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请联系项目维护者。
