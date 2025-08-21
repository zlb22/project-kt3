# 项目开发概要与结构说明（开发者参考）

本文件基于当前代码与 `README.md` 的梳理，总结了项目现状、目录结构、接口与页面逻辑关系，便于后续开发与扩展。

## 一、项目概览

* __技术栈__
  - 后端：`FastAPI`、`PyJWT`、`Pydantic`、`Uvicorn`
  - 前端：`React 18`、`TypeScript`、`React Router v6`、`Material-UI`、`Axios`
* __功能域__
  - 用户认证（登录、获取用户信息、修改密码）
  - 测评功能入口（问题解决、互动讨论、AUT 测试、情绪测评、在线课程）
  - 健康检查与（占位）测评接口

## 二、项目目录结构

```
project-kt3/
├── backend/                        # FastAPI 后端
│   ├── main.py                     # 主应用：认证、健康检查、测评占位接口
│   ├── requirements.txt            # Python 依赖
│   ├── users.json                  # 简易用户存储（默认：admin/1）
│   └── server.log / uvicorn.pid    # 运行产物（开发）
├── frontend/                       # React 前端
│   ├── public/
│   ├── src/
│   │   ├── App.tsx                 # 路由配置与主题
│   │   ├── App.css                 # 全局样式
│   │   ├── index.tsx               # 前端入口
│   │   ├── contexts/
│   │   │   └── AuthContext.tsx     # 认证上下文（token 管理、axios header 注入）
│   │   ├── components/
│   │   │   └── ProtectedRoute.tsx  # 受保护路由包装
│   │   └── pages/                  # 各页面（均已路由接入）
│   │       ├── Login.tsx
│   │       ├── Dashboard.tsx
│   │       ├── ChangePassword.tsx
│   │       ├── ProblemSolving.tsx
│   │       ├── InteractiveDiscussion.tsx
│   │       ├── AutTest.tsx
│   │       ├── EmotionTest.tsx
│   │       └── OnlineClass.tsx
│   ├── package.json                # 含 dev 代理到 http://localhost:8000
│   └── tsconfig.json
├── README.md
├── Makefile                        #（存在，未在本说明展开）
├── plan.md                         # 本文件
└── users.json                      # 项目根也有一份（注意与 backend/ 内文件区分）
```

## 三、后端接口与认证流程

文件：`backend/main.py`

* __JWT 配置__：`SECRET_KEY`、`HS256`、有效期默认 30 分钟。
* __CORS__：允许 `http://localhost:3000`（前端 dev）。
* __用户存储__：本地 `backend/users.json`，初次运行自动写入默认用户：admin/1（SHA256 存储）。

### 主要端点
* __认证__
  - `POST /api/auth/login`（`UserLogin`）→ 返回 `access_token`
  - `GET  /api/auth/me` → 返回当前用户（需 Bearer token）
  - `POST /api/auth/change-password`（`PasswordChange`）→ 修改密码（需 Bearer token）
* __系统__
  - `GET /api/health` → 健康检查
* __测评（占位）__
  - `GET /api/tests/available` → 返回可用测评列表（需 Bearer token）
  - `GET /api/tests/{test_id}` → 返回指定测评详情（需 Bearer token）

### 认证时序（前端视角）
1. `Login.tsx` 提交用户名/密码 → `AuthContext.login()` 调用 `POST /api/auth/login`。
2. 成功后将 `access_token` 写入 `localStorage`，并设置全局 axios `Authorization: Bearer <token>`。
3. 随后调用 `GET /api/auth/me` 拉取用户信息入 `AuthContext.user`。
4. `ProtectedRoute` 基于 `user` 与 `isLoading` 控制访问受保护页面，未登录跳转 `/login`。

## 四、前端路由与页面关系

文件：`frontend/src/App.tsx`

* __公共路由__
  - `/login` → `Login`
  - `/change-password` → `ChangePassword`（已登录与未登录均可进入，内部根据状态返回不同按钮/跳转）
* __受保护路由__（`<ProtectedRoute>` 包裹）
  - `/dashboard` → `Dashboard`
  - `/problem-solving` → `ProblemSolving`
  - `/interactive-discussion` → `InteractiveDiscussion`
  - `/aut-test` → `AutTest`
  - `/emotion-test` → `EmotionTest`
  - `/online-class` → `OnlineClass`
* __根路由__
  - `/` → 重定向到 `/dashboard`

### 页面跳转关系（核心路径）
* __Login__
  - 登录成功 → `/dashboard`
  - 底部链接 → `/change-password`
* __Dashboard__
  - 右上角退出 → 清除登录态，跳转 `/login`
  - 功能卡片按钮：
    - 问题解决 → `/problem-solving`
    - 互动讨论 → `/interactive-discussion`
    - 在线实验 → `/online-class`
  - 底部按钮 → `/change-password`
* __ChangePassword__
  - 修改成功 2 秒后：若已登录 → `/dashboard`；未登录 → `/login`
* __各测评页面（ProblemSolving/InteractiveDiscussion/AutTest/EmotionTest/OnlineClass）__
  - 头部返回按钮 → `/dashboard`
  - 当前均为占位内容，描述未来功能点，方便后续接入具体测评流程与组件。

## 五、数据与网络层

* __Axios 全局设置__：在 `AuthContext` 中基于 token 写入/删除 `Authorization` header。
* __开发代理__：`frontend/package.json` 中 "proxy": `http://localhost:8000`，使前端在开发时可直接以 `/api/*` 调用后端。
* __状态管理__：`AuthContext` 提供 `user`、`token`、`isLoading`、`login/logout/changePassword` 能力；页面通过 `useAuth()` 消费。

## 六、运行与环境

* __后端__：`backend/requirements.txt`，`python main.py` 或 `uvicorn main:app --reload`，端口 `8000`
* __前端__：`npm start`，端口 `3000`（开发）
* __默认账号__：用户名 `admin`，密码 `1`

## 七、当前进度与待办建议

* __已完成__
  - 前后端分离架构与认证闭环（登录/鉴权/修改密码）
  - 前端导航与受保护路由框架
  - 测评功能页面框架（占位 UI），统一样式与返回导航
  - 基础 README、运行说明
* __待办建议__
  - 后端：将占位测评接口替换为真实业务模块，细化 `tests/*` 路由、数据模型与存储层；拆分 `main.py`（路由、服务、模型、配置、日志）。
  - 前端：为各测评页接入实际流程（表单/媒体/可视化），抽象服务层（`/src/services/*`）与类型定义（`/src/types/*`）。
  - 安全与配置：`SECRET_KEY` 环境化；完善 CORS 域名；token 续期策略；生产部署脚本。
  - 账号体系：注册/角色/权限模型；用户管理页面与接口。
  - 测试与质量：前后端单元测试、E2E，以及 ESLint/Prettier 配置。

## 八、演进规划（建议）

* __后端分层__：`app/routers`、`app/services`、`app/models`、`app/core`、`app/schemas`、`app/deps`，提高可维护性。
* __前端结构__：新增 `services/`、`types/`、`hooks/`、`layouts/`，将页面逻辑与数据访问解耦。
* __日志与监控__：完善 `server.log` 规范化输出，引入请求追踪与错误上报。
* __持久化__：替换本地 `users.json` 为数据库（如 SQLite/PostgreSQL），引入迁移工具。

—— 以上即当前开发成果与结构说明，如需我继续按上述待办拆分任务并实施，请告知优先级。

### 已完成功能

✅ **在线实验测试页面**
- 修改了 `OnlineClass.tsx`，将其转换为"在线实验测试"入口页面
- 添加了测试选择对话框，提供24点游戏和Fool Planning两个选项
- 使用Material-UI设计了现代化的卡片式选择界面

✅ **24点游戏完整实现** (`TwentyFourGame.tsx`)
- **微课学习阶段**：
  - 4段微课视频的模拟播放器（基础规则、运算技巧、高级策略、实战演练）
  - 视频进度条、播放/暂停控制、自动切换下一视频
  - 观看时长统计，必须观看完所有视频才能进入测试
  
- **答题测试阶段**：
  - 10道24点题目，使用预设和随机生成的数字组合
  - 30分钟倒计时，实时显示剩余时间
  - 支持提交答案和跳过题目功能
  - 每题记录答题次数、用时、答案内容
  - 简单的表达式验证逻辑（可后续完善）
  
- **结果统计阶段**：
  - 总得分、总答题用时、视频观看用时统计
  - 每道题目的详细答题情况展示
  - 数据保存到localStorage（可扩展为API调用）

### 下一步计划

1. **数据采集功能**：实现人像和屏幕数据采集，保存为mp4文件（需要摄像头和屏幕录制权限）
2. **后端API扩展**：添加24点测试数据存储接口，支持用户测试记录管理
3. **题库集成**：替换当前的模拟题目为真实题库数据
4. **Fool Planning测试**：实现第二个测试选项
5. **数据分析**：添加测试结果的详细分析和可视化展示
