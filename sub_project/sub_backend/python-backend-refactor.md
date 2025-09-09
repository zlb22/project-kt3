# project-kt3 Python 后端重构方案（最小改动）

## 目标与范围
- 保持前端无感知：接口路径、请求头签名、响应结构与 Go 版本一致（`/web/keti3/...`）。
- 存储改造：`/oss/auth` 改为返回 MinIO 预签名直传（保持我们已确认的“前端直传”方案）。
- 日志导出：继续“多 XLSX + ZIP + 临时下载链接”的现状，存储层换 MinIO 即可。
- 复用：MySQL/Redis 仍为主要状态存储；签名、权限规则保持一致。

## 技术栈
- 框架：FastAPI
- ORM：SQLAlchemy + Alembic（迁移）
- 缓存：redis-py（老师 token）
- 存储：minio（Python SDK）
- 配置：python-dotenv / pydantic-settings（读取 `.env`）
- 日志：structlog / logging

## 目录结构（建议）
```
backend-py/
  app/
    main.py                 # 入口，路由注册
    deps.py                 # 依赖注入（DB/Redis/MinIO 客户端）
    middleware/
      auth.py               # Authorization 校验
      sign.py               # X-Auth-AppId/TimeStamp/X-Sign 校验
    api/
      keti3.py              # 与 Go 对齐的 8 个接口
    services/
      storage_minio.py      # 预签名 GET/PUT、上传下载
      export_log.py         # 生成 XLSX、打包 ZIP、上传
    models/
      base.py               # SQLAlchemy Base
      keti3_student_info.py
      keti3_student_submit.py
      keti3_op_log.py
    schemas/
      req_resp.py           # Pydantic 请求/响应模型（与现有响应结构对齐）
    utils/
      response.py           # 统一响应：{errcode, errmsg, data, trace}
      security.py           # MD5 计算、时间戳窗口
      time.py
  alembic/                  # 迁移
  alembic.ini
  requirements.txt
  .env.example
  README.md
```

## 接口映射与最小实现
- 前缀：`/web/keti3`（Nginx/网关不变）
- 中间件：
  - `CheckSignature` 等价实现：校验 `X-Auth-AppId`、`X-Auth-TimeStamp`（可选时间窗口）、`X-Sign`（MD5：`secret + timestamp`）
  - `CheckUserToken` 等价实现（方案A：本地 JWT + Redis 双模）：
    - 先尝试将 `Authorization` 作为 JWT 验签（与 `backend/main.py` 相同 `SECRET_KEY`、`ALGORITHM`）。
    - 若非 JWT 或验签失败，则回退检查 Redis 键 `base-keti:teacher-login:{token}` 是否存在（老师口令登录）。
    - 两者均失败 → 返回 HTTP 200 + `errcode: 20004`。
    - 路径白名单（免登录）：`/web/keti3/student/login`、`/web/keti3/config/list`、`/web/keti3/oss/auth`。
- 路由与功能：
  - POST `/uid/get-by-increase`：用 Redis INCR 或自增表生成 uid
  - POST `/log/save`：写入 `keti3_op_log`
  - GET `/log/export`：按 UID 生成多份 XLSX → ZIP → MinIO 上传 → 返回预签名下载链接
  - POST `/oss/auth`：返回两个 PUT 预签名（img/audio），含 `key/method/url/content_type/public_url`，并保留旧字段（空）
  - POST `/config/list`：返回学校、年级等配置（与现有格式一致）
  - POST `/student/login`：查/建 `keti3_student_info`，返回 `uid`
  - POST `/teacher/login`：校验 `login_code`，生成 Redis token（1h）
  - POST `/log/list`：分页查询 `keti3_op_log`

## 数据库（MySQL）
保持与 Go `schema/` 一致（字段名/类型尽量对齐）：
- `keti3_student_info`：school_id, student_name, student_num, grade_id, class_name, created_at, updated_at
- `keti3_student_submit`：uid, school_id, student_name, student_num, grade_id, class_name, voice_url, oplog_url, screenshot_url, date, created_at, updated_at
- `keti3_op_log`：submit_id, uid, op_time, op_type, op_object, object_no, object_name, data_before(json), data_after(json), voice_url, screenshot_url, deleted_at, created_at, updated_at
- Alembic 迁移脚本提供建表/索引（含 json 列兼容 MySQL 8+）。

## 存储（MinIO）
- `.env`：MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET, MINIO_USE_SSL=false, MINIO_BASE_PATH=24game
- 预签名：
  - PUT：前端直传；强制 `Content-Type`（image/*、audio/*）
  - GET：导出 ZIP、语音/截图列写临时私链（有效期可配置）
- 目录：`{base_path}/{uid}/{date}/...`（与现有前端约定保持一致）

## 配置与环境
- `.env` 示例：
  - APP_PORT=8096
  - APP_ID=keti-3
  - SKEY=xxx（签名密钥，需与前端 `VITE_SKEY` 对齐）
  - MYSQL_DSN=mysql://user:pass@host:3306/base_keti?charset=utf8mb4
  - REDIS_URL=redis://host:6379/0
  - MINIO_*
  - TEACHER_LOGIN_CODE=xxxx

## 统一响应
- 成功：`{"errcode":0,"errmsg":"ok","data":{...},"trace":""}`
- 失败码：与前端 `axios.ts` 已处理的 20004/20005/60000 对齐

## 登录态集成（方案A：本地 JWT + Redis 双模）
- 前端保持不变：每次请求自动携带 `Authorization`、`X-Auth-AppId`、`X-Auth-TimeStamp`、`X-Sign`。
- 服务端校验流程：
  1) `sign.py`：验证 `X-Sign == md5(secret + timestamp)`，并可校验时间戳漂移。
  2) `auth.py`：
     - 将 `Authorization` 作为 JWT 验签（`SECRET_KEY`、`ALGORITHM`）。成功则通过，并在 `request.state.user` 注入用户信息。
     - 否则回退用 Redis 检查 `base-keti:teacher-login:{token}` 是否存在。存在则通过，并注入 `{role: 'teacher'}`。
     - 两者均失败 → 返回 HTTP 200 + `errcode: 20004`（前端会按现有逻辑处理）。
- 白名单接口（免登录）：`/student/login`、`/config/list`、`/oss/auth`。
- 环境变量：需与 `backend/main.py` 一致提供 `SECRET_KEY`、`ALGORITHM`；Redis 配置保持不变以兼容口令登录 token。

## 导出逻辑（维持现状）
- 每 UID 生成一份 `.xlsx`（列与 Go 版一致），目录打包为 `.zip`
- ZIP 上传至 MinIO，返回 1 小时预签名下载链接
- Excel 内“语音/截图”列使用长时效临时链接或 `bucket/key`（默认继续临时链接）

## 开发步骤（建议）
1. 初始化骨架与依赖、统一响应、签名/鉴权中间件
2. 接上 MySQL/Redis/MinIO 客户端
3. 实现 `/student/login`、`/config/list` 以完成基本联通
4. 实现 `/oss/auth` → 前端改 PUT 直传
5. 实现 `/teacher/login`、`/uid/get-by-increase`
6. 实现 `/log/save`、`/log/list`、`/log/export`
7. 编写 Alembic 迁移与 README（启动、环境变量、样例请求）

## 里程碑与工期（估算）
- M1（0.5 天）：项目骨架、配置、签名/鉴权
- M2（0.5 天）：DB/Redis/MinIO 接入
- M3（0.5–1 天）：登录/配置/UID/授权
- M4 (0.5–1 天)：日志保存/列表/导出（XLSX+ZIP）
- M5（0.5 天）：联调与文档
- 合计：约 2–3 天

## 风险与回滚
- 风险：签名规则细节偏差；Excel 列/内容不一致；时钟偏差导致签名失败
- 缓解：编写对照测试、抓包比对；先在测试环境并行运行 Go 与 Python 版本
- 回滚：Nginx 路由切回 Go 服务；配置 `storage.type` 切回 AliOSS（如仍保留）

## 与 Go 的取舍
- 研发速度：Python > Go（本项目业务简单，生态与开发效率更佳）
- 运行效率：Go 更强；但本项目 QPS 不高，Python 足够
- 生态/维护：你已在主项目使用 Python + MinIO，更易统一维护
