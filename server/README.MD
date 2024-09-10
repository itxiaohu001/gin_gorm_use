# 用户认证和文件管理系统

## 项目概述

这是一个基于 Go 语言和 Gin 框架开发的用户认证和文件管理系统。该系统提供了用户注册、登录、角色管理以及文件上传、下载、删除和列表功能。

## 主要特性

1. 用户认证
   - 用户注册
   - 用户登录
   - JWT token 认证

2. 用户管理
   - 删除用户（仅管理员）
   - 更新用户角色（仅管理员）

3. 文件管理
   - 文件上传
   - 文件下载
   - 文件删除
   - 文件列表

4. 存储支持
   - 本地存储
   - MinIO 对象存储

5. 日志系统
   - 使用 zap 进行日志记录
   - 支持日志轮转

6. 追踪系统
   - 为每个请求生成唯一的 traceID
   - 在响应头中返回 traceID

## 技术栈

- Go
- Gin Web 框架
- GORM (MySQL)
- JWT 认证
- Zap 日志库
- MinIO 客户端
- Viper 配置管理

## 项目结构

- `config/`: 配置相关
- `database/`: 数据库连接
- `handlers/`: HTTP 请求处理
- `middleware/`: 中间件（认证、追踪等）
- `models/`: 数据模型
- `services/`: 业务逻辑
- `storage/`: 存储相关（本地存储和 MinIO）
- `utils/`: 工具函数
- `main.go`: 应用入口

## 安装和运行

1. 克隆仓库：
   ```
   git clone <repository-url>
   cd <project-directory>
   ```

2. 安装依赖：
   ```
   go mod tidy
   ```

3. 配置：
   复制 `config.yaml.example` 到 `config.yaml` 并根据您的环境进行必要的修改。

4. 运行应用：
   ```
   go run main.go
   ```

## API 文档

项目集成了 Swagger，可通过 `/swagger/*any` 路径访问 API 文档。

## 配置

使用 `config.yaml` 文件进行配置，包括数据库连接、Redis、JWT、服务器端口、MinIO 和日志等设置。