# project-kt3 系统架构说明

本档描述当前仓库的整体架构、技术栈、目录、运行方式、数据模型、存储设计、关键环境配置以及已知问题，供工程师开发与维护参考。

更新时间：2025-09-15 14:20

---

## 1. 总览

- 根目录：`/home/ubuntu/zlb/project-kt3/`
- 前端一（主站 React）：`frontend/`
- 前端二（子项目 Topic-Three Vue）：`sub_project/sub_frontend/`
- 后端：`backend/`
- 证书：`certs/`

整体为前后端分离方案：
- 后端提供认证、测评数据、Keti3 适配、MinIO 存储接口
- 前端一为综合教育测评工作台（React）
- 前端二为“在线实验 Topic-Three”子项目（Vue3 + Vite）

---

## 2. 技术栈

- 后端（`backend/`）
  - 框架：FastAPI（`backend/main.py`）
  - ORM/数据库：SQLAlchemy + MySQL（`backend/db_mysql.py`）
  - 存储：MinIO（`backend/keti3_storage.py`）
  - 认证：JWT（PyJWT）
  - 中间件：Keti3 请求签名校验（`backend/keti3_middleware.py`）
  - 配置：dotenv（`backend/.env`）
  - 运行：Uvicorn

- 前端一（`frontend/`）
  - React 18 + TypeScript + CRA（create-react-app）
  - UI：MUI
  - 路由：React Router
  - HTTP：Axios

- 前端二（`sub_project/sub_frontend/`）
  - Vue3 + Vite + TypeScript
  - 组件库：Element Plus
  - 状态：Pinia
  - 画布：fabric.js
  - 压缩：JSZip
  - HTTP：Axios

---

## 3. 目录结构要点

- `backend/`
  - `main.py`：FastAPI 入口和业务路由集中管理
  - `db_mysql.py`：SQLAlchemy Engine/Session 与模型定义、建表
  - `keti3_middleware.py`：Keti3 签名校验、统一响应格式与错误码
  - `keti3_storage.py`：MinIO 客户端、确保桶存在、预签名上传
  - `.env`：后端环境变量
  - `pyproject.toml`：后端依赖

- `frontend/`
  - `package.json`、`.env.development`、`src/`、`public/`、`tsconfig.json`

- `sub_project/sub_frontend/`
  - `package.json`、`.env.*`、`vite.config.ts`、`src/`、`index.html`、`项目使用文档.md`

- `certs/`
  - `server.crt`、`server.key`

---

## 4. 后端架构（`backend/`）

- 应用初始化
  - `load_dotenv()` 读取 `.env`
  - CORS 允许来源：取自 `CORS_ORIGINS`（逗号分隔）
  - MySQL 表：启动即执行 `create_tables()`（`db_mysql.py`）

- 认证与用户
  - 密码哈希：`SHA256`（`make_hashed_password()`）
  - 登录：`POST /api/auth/login` 返回 JWT，`sub=username`，`exp` 过期
  - 注册：`POST /api/auth/register`（写入 MySQL）
  - 当前用户：`GET /api/auth/me`
  - 改密：`POST /api/auth/change-password`

- 测试/测评
  - `GET /api/tests/available`、`GET /api/tests/{test_id}`（占位）
  - 24 点提交：`POST /api/tests/24point/submit`（将 payload 写入 MySQL `twentyfour_records`，并根据 `test_session_id` 绑定相关 `video_assets` 元数据）
  - 24 点大 JSON：
    - `POST /api/24pt/record/save`（Keti3 签名）
    - `GET /api/24pt/record/list`（Keti3 签名）

- 存储/上传（MinIO）
  - 直传授权（图片/音频）：`POST /web/keti3/oss/auth`（签名校验，返回 `PUT` 预签名 URL）
  - 代理上传（图片/音频）：`POST /web/keti3/oss/upload`（签名校验，服务端代传 `img|audio`，已改为流式转存）
  - 视频直传授权：`POST /api/videos/presign`（需 JWT，返回 `put_url`、`object_name` 等）
  - 视频直传登记：`POST /api/videos/commit`（直传成功后登记 `video_assets` 元数据）
  - 视频代理上传：`POST /api/upload/video`（需 JWT；对象路径：`{camera|screen}/{uid}/{test_session_id}/{filename}`；已改为流式转存）

- Keti3 适配
  - `GET /web/keti3/entry`：按 `username` upsert `Student` 并签发 JWT，通过 302 重定向至前端子项目（`KETI3_FRONTEND_URL?token=...`）
  - `POST /web/keti3/student/login`、`/web/keti3/config/list`、`/web/keti3/log/save`：均走签名校验

---

## 5. 数据模型（`backend/db_mysql.py`）

- `Student`
  - `id, username(unique), school, grade, password(hash), is_active, created_at, updated_at`
- `OperationLog`
  - `id, uid, action, details(Text JSON 字符串), timestamp`
- `TwentyFourRecord`
  - `id(BigInt), uid, payload(Text), created_at`
- `VideoAsset`
  - `id(BigInt), uid, username, test_session_id, video_type(camera|screen), bucket, object_name, content_type, file_size, upload_time`

---

## 6. MinIO 存储与对象命名

- 连接参数：`MINIO_URL`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`，支持自签名证书（`MINIO_CA_CERT` 或 `MINIO_INSECURE_TLS`）
- 视频对象命名：
  - `video_type` 取 `camera|screen`（含常见拼写纠正）
  - 对象路径：`{video_type}/{uid}/{test_session_id}/{filename}`
  - 文件名：`{username}_{test_session_id}_{video_type}_{timestamp}.webm`
- 预签名上传：`/web/keti3/oss/auth` 返回 1 小时有效的 `PUT` URL
- 代理上传：`/web/keti3/oss/upload` 同时支持 `img`、`audio` 字段，自动解析 `uid`

### 直传与代理流程

- React 前端（24点视频）推荐直传流程：
  1) `POST /api/videos/presign` 获取 `put_url`、`bucket`、`object_name`、`public_url`
  2) 浏览器对 MinIO 执行 `PUT`（带 `Content-Type`）
  3) `POST /api/videos/commit` 将 `bucket/object_name/test_session_id/...` 登记到 `video_assets`
  4) 若直传失败，回退到 `POST /api/upload/video`（服务器代理，已流式）

- Vue 子前端（图片/音频）推荐直传流程：
  1) `POST /web/keti3/oss/auth` 获取各自 `PUT` 预签名 URL
  2) 浏览器对 MinIO 执行 `PUT`（带 `Content-Type`）
  3) 前端使用返回的 `public_url` 继续业务流程
  4) 若直传失败，回退到 `POST /web/keti3/oss/upload`（服务器代理，已流式）

---

## 7. 环境配置

- 后端 `.env`（`backend/.env`）核心项
  - JWT：`SECRET_KEY`, `ALGORITHM`, `ACCESS_TOKEN_EXPIRE_MINUTES`
  - CORS：`CORS_ORIGINS`（逗号分隔的来源列表）
  - 服务：`HOST`, `PORT`
  - MinIO：`MINIO_URL`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_BUCKET_KETI3`, `MINIO_BASE_PATH`, `MINIO_BASE_PATH_KETI3`, `MINIO_ALLOWED_BUCKETS`, `MINIO_CA_CERT`, `MINIO_INSECURE_TLS`
  - Keti3：`APP_ID`, `SKEY`, `KETI3_FRONTEND_URL`
  - MySQL：`MYSQL_DSN`

- 前端一 `.env.development`（`frontend/.env.development`）
  - `HOST`, `PORT`, `DANGEROUSLY_DISABLE_HOST_CHECK`, `WDS_SOCKET_HOST`, `WDS_SOCKET_PORT`
  - `package.json` 中 `proxy` 当前为 `https://localhost:8443`（建议改为后端真实地址，见“已知问题”）

- 前端二 `.env.development`（`sub_project/sub_frontend/.env.development`）
  - `VITE_SKEY`, `VITE_APP_ID`, `VITE_LOGIN_URL`（注意密钥泄露风险，见“已知问题”）

---

## 8. 接口清单（主要）

- 认证
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `GET /api/auth/me`
  - `POST /api/auth/change-password`

- 健康检查
  - `GET /api/health`

- 测试/测评
  - `GET /api/tests/available`
  - `GET /api/tests/{test_id}`
  - `POST /api/tests/24point/submit`
  - `POST /api/24pt/record/save`（签名）
  - `GET /api/24pt/record/list`（签名）

- 存储/上传
  - `POST /api/upload/video`（JWT）
  - `POST /api/videos/presign`（JWT）
  - `POST /api/videos/commit`（JWT）
  - `POST /web/keti3/oss/auth`（签名）
  - `POST /web/keti3/oss/upload`（签名）

- Keti3 适配
  - `GET /web/keti3/entry`
  - `POST /web/keti3/student/login`（签名）
  - `GET|POST /web/keti3/config/list`（签名）
  - `POST /web/keti3/log/save`（签名）

---

## 9. 运行指南（开发）

- 后端
  - 确保 MySQL 与 MinIO 可用，`backend/.env` 中 `MYSQL_DSN`、`MINIO_*` 配置正确
  - 运行：
    ```bash
    cd backend
    uvicorn main:app --host 0.0.0.0 --port 8000 --reload
    ```
  - 文档：`http://localhost:8000/docs`

- 前端一（React CRA）
  - 建议将 `frontend/package.json` 的 `"proxy"` 指向后端地址（例如 `http://localhost:8000`）
  - 运行：
    ```bash
    cd frontend
    npm install
    npm start
    ```

- 前端二（Vue Vite）
  - 运行：
    ```bash
    cd sub_project/sub_frontend
    npm install
    npm run dev
    ```

---

## 10. 已知/潜在问题与建议

1) 未定义常量导致 NameError
- 位置：`backend/main.py` 的 `keti3_oss_upload()` 使用了 `ERRCODE_PARAM_ERROR`
- 问题：`keti3_middleware.py` 未定义此常量，`main.py` 也未导入
- 影响：当请求没有 `img` 和 `audio` 时会抛 `NameError`
- 建议：改为已导入的 `ERRCODE_INVALID_PARAMS`，或在 `keti3_middleware.py` 中新增并导入

2) CRA 代理端口与后端不一致
- 位置：`frontend/package.json` 的 `proxy` 为 `https://localhost:8443`
- 问题：后端默认 `:8000`，8443 可能并非你的后端
- 影响：本地联调请求被代理到错误端口
- 建议：改为 `http://localhost:8000` 或你的真实网关地址

3) 启动即建表对可用性敏感
- 位置：`backend/main.py` 顶部 `create_tables()`
- 问题：数据库未就绪时会导致应用启动失败
- 建议：迁移到启动脚本/命令，或添加重试/超时与友好日志

4) 不再使用的 Mongo 依赖/配置残留
- 位置：`backend/pyproject.toml` 仍包含 `motor`、`pymongo`；`main.py` 中也保留 Mongo 配置变量
- 建议：若已完全迁移 MySQL，清理依赖与配置，更新 README

5) 前端暴露服务端签名密钥
- 位置：`sub_project/sub_frontend/.env.development` 含 `VITE_SKEY`
- 问题：打包构建后前端环境变量可被访问，泄露密钥
- 建议：签名逻辑放在服务端；若需要前端参与，改为短时令牌/临时密钥方案

6) TLS/证书与安全
- 位置：`backend/.env` 设置 `MINIO_INSECURE_TLS=true`
- 建议：生产关闭此项，使用可信 CA 证书（或配置 `MINIO_CA_CERT`）

7) 密码哈希
- 现状：单次 `SHA256`
- 建议：生产改为 `bcrypt`/`argon2`（带盐与成本因子）

8) 其余建议
- 清理 `users.json` 等遗留文件
- `frontend/.env.development` 中 `DANGEROUSLY_DISABLE_HOST_CHECK` 仅限开发，勿用于生产
- 子项目依赖 `ali-oss` 与后端 MinIO 直传并存，注意不要混淆两套上云路径与鉴权方式

---

## 11. 环境变量对齐建议

- `CORS_ORIGINS` 覆盖实际访问来源（协议+域名/IP+端口）
- CRA `proxy` 对齐后端真实地址
- `KETI3_FRONTEND_URL` 指向可达的子前端路由（示例：`https://<host>:5173/topic-three`）
- MinIO 相关参数在 `backend/.env`、`backend/main.py` 默认值、`backend/keti3_storage.py` 中均出现，建议以 `.env` 为准并统一

---

## 12. 生产化加固路线图（可选）

- 认证：切换到 `bcrypt/argon2`；引入刷新 Token 与 Token 失效策略
- 存储：MinIO 使用可信证书，关闭 `INSECURE_TLS`
- 数据库：引入 Alembic 做迁移；建表从运行时移除
- 网关：Nginx/Caddy 统一 HTTPS、反向代理后端与 MinIO，简化前端跨域
- 日志：结构化日志与追踪 ID，重点端点（上传/签名失败/持久化失败）加告警
- 依赖：移除未使用依赖，定期升级安全版本

---

## 13. 附录：关键文件引用

- 后端入口：`backend/main.py`
- 数据库与模型：`backend/db_mysql.py`
- 存储封装：`backend/keti3_storage.py`
- 中间件与错误码：`backend/keti3_middleware.py`
- 后端环境：`backend/.env`
- 前端一配置：`frontend/package.json`、`frontend/.env.development`
- 前端二配置：`sub_project/sub_frontend/package.json`、`sub_project/sub_frontend/.env.development`

---

如需我直接提交修复（错误码、CRA 代理、密钥治理等），请告知优先级。
