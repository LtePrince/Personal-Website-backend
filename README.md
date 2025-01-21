# 个人网站后端

这是一个使用 Go 语言编写的个人网站后端项目。

## 项目结构

```
wzy-website-backend/
├── main.go
├── handlers/
│   └── userHandler.go
├── models/
│   └── user.go
├── routers/
│   └── router.go
├── utils/
│   └── db.go
└── README.md
```

## 功能

- 用户注册和登录
- 用户信息管理
- 博客文章管理
- 评论功能

## 安装和运行

1. 克隆仓库

```bash
git clone https://github.com/yourusername/wzy-website-backend.git
cd wzy-website-backend
```

2. 安装依赖

```bash
go mod tidy
```

3. 运行项目

```bash
go run main.go
```

## API 文档

### 用户

- `POST /register` - 用户注册
- `POST /login` - 用户登录
- `GET /user/:id` - 获取用户信息

### 博客

- `POST /blog` - 创建博客文章
- `GET /blog/:id` - 获取博客文章
- `PUT /blog/:id` - 更新博客文章
- `DELETE /blog/:id` - 删除博客文章

### 评论

- `POST /comment` - 添加评论
- `GET /comment/:id` - 获取评论
- `DELETE /comment/:id` - 删除评论

## 贡献

欢迎贡献代码！请 fork 本仓库并提交 pull request。

## 许可证

本项目使用 MIT 许可证。详情请参阅 [LICENSE](./LICENSE) 文件。