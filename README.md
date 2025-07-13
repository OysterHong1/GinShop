# GinShop 电商系统

基于 Go 语言 Gin 框架开发的现代化电商系统，提供完整的前台购物体验和后台管理功能。

## 📋 项目简介

GinShop 是一个功能完善的电商系统，包含用户购物前台、管理员后台以及 RESTful API 接口。系统采用前后端分离架构，支持商品管理、订单处理、用户管理、支付集成等核心电商功能。

## 🚀 技术栈

### 后端技术
- **框架**: Gin (Go Web 框架)
- **数据库**: MySQL + GORM ORM
- **缓存**: Redis
- **搜索引擎**: Elasticsearch 7.9.0
- **会话管理**: Redis Session
- **配置管理**: INI 配置文件

### 集成服务
- **支付系统**: 支付宝支付、微信支付
- **图像处理**: 图片上传和处理
- **验证码**: Base64 验证码
- **二维码**: QR Code 生成
- **云存储**: 阿里云 OSS

## ✨ 核心功能

### 用户前台功能
- 🛍️ **商品浏览**: 商品分类、详情展示、图片画廊
- 🛒 **购物车**: 添加/删除商品、数量调整、价格计算
- 👤 **用户系统**: 注册登录、个人中心、地址管理
- 📦 **订单管理**: 下单结算、订单跟踪、支付状态
- 💳 **支付集成**: 支付宝、微信支付
- 🔍 **搜索功能**: 基于 Elasticsearch 的商品搜索
- 📱 **短信验证**: 注册/登录短信验证码

### 管理后台功能
- 🏪 **商品管理**: 商品增删改查、分类管理、属性配置
- 👥 **用户管理**: 用户列表、状态管理
- 📋 **订单管理**: 订单查看、状态更新、发货管理
- 🎨 **内容管理**: 轮播图、导航菜单、页面配置
- 🔐 **权限系统**: 角色权限、管理员管理、访问控制
- 📊 **系统设置**: 网站配置、缓存管理

### API 接口
- 🔌 **RESTful API**: v1/v2 版本接口
- 📱 **移动端支持**: JSON 数据格式
- 🔒 **接口鉴权**: Token 验证机制

## 📁 项目结构

```
GinShop/
├── conf/                   # 配置文件
│   └── app.ini            # 应用配置
├── controllers/           # 控制器
│   ├── admin/            # 后台管理控制器
│   ├── api/              # API 接口控制器
│   └── oystershop/       # 前台购物控制器
├── middlewares/          # 中间件
│   ├── adminAuth.go      # 管理员权限验证
│   ├── userAuth.go       # 用户权限验证
│   └── init.go           # 中间件初始化
├── models/               # 数据模型
│   ├── mysqlCore.go      # MySQL 连接
│   ├── esClientCore.go   # Elasticsearch 客户端
│   ├── goods.go          # 商品模型
│   ├── user.go           # 用户模型
│   ├── order.go          # 订单模型
│   └── ...               # 其他模型文件
├── routers/              # 路由配置
│   ├── adminRouters.go   # 后台路由
│   ├── defaultRouters.go # 前台路由
│   └── apiRouters.go     # API 路由
├── static/               # 静态资源
│   ├── admin/            # 后台静态文件
│   ├── oystershop/       # 前台静态文件
│   └── upload/           # 上传文件目录
├── templates/            # 模板文件
│   ├── admin/            # 后台模板
│   └── oystershop/       # 前台模板
├── main.go              # 程序入口
├── go.mod              # Go 模块文件
└── go.sum              # 依赖校验文件
```

## 🛠️ 安装配置

### 环境要求
- Go 1.22.0+
- MySQL 5.7+
- Redis 6.0+
- Elasticsearch 7.9.0

### 1. 克隆项目
```bash
git clone <repository-url>
cd GinShop
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置数据库
编辑 `conf/app.ini` 文件：
```ini
[mysql]
ip       = 127.0.0.1
port     = 3306
user     = root
password = 123456
database = ginxiaomi

[redis]
ip   = 127.0.0.1
port = 6379
redisEnable = true
```

### 4. 启动服务

#### 启动 Redis
```bash
redis-server
```

#### 启动 Elasticsearch
```bash
cd elasticsearch-7.9.0
./bin/elasticsearch
```

#### 启动应用
```bash
go run main.go
```

### 5. 访问地址
- **前台地址**: http://localhost:8080
- **后台地址**: http://localhost:8080/admin
- **API 接口**: http://localhost:8080/api/v1 或 http://localhost:8080/api/v2

## 📦 核心依赖

```go
require (
    github.com/gin-gonic/gin v1.10.0           // Web 框架
    github.com/gin-contrib/sessions v1.0.1     // Session 管理
    github.com/gin-contrib/cors v1.7.2         // 跨域支持
    gorm.io/gorm v1.25.8                       // ORM 框架
    gorm.io/driver/mysql v1.5.7                // MySQL 驱动
    github.com/go-redis/redis/v8 v8.11.5       // Redis 客户端
    github.com/olivere/elastic/v7 v7.0.32      // Elasticsearch 客户端
    github.com/smartwalle/alipay/v3 v3.2.24    // 支付宝支付
    github.com/objcoding/wxpay v1.0.6          // 微信支付
    github.com/mojocn/base64Captcha v1.3.6     // 验证码
    gopkg.in/ini.v1 v1.67.0                   // INI 配置
)
```

## 🔧 功能模块详解

### 权限管理系统
- **多级权限**: 支持菜单权限、操作权限
- **角色管理**: 灵活的角色权限分配
- **管理员管理**: 管理员账号管理和权限控制

### 商品管理系统
- **分类管理**: 无限级商品分类
- **属性管理**: 商品规格、颜色、属性配置
- **库存管理**: 库存数量、预警设置
- **图片管理**: 多图上传、图片处理

### 订单支付系统
- **订单流程**: 下单 -> 支付 -> 发货 -> 完成
- **支付集成**: 支付宝网页支付、微信支付
- **物流管理**: 发货单号、物流跟踪

### 缓存策略
- **Redis 缓存**: 热点数据缓存，提升访问速度
- **缓存更新**: 自动缓存失效和更新机制
- **缓存预热**: 系统启动时预加载热点数据

## 🚦 API 接口说明

### 用户接口 (v1)
- `GET /api/v1/navList` - 获取导航列表
- `POST /api/v1/doLogin` - 用户登录
- `PUT /api/v1/editArticle` - 编辑文章
- `DELETE /api/v1/deleteNav` - 删除导航

### 商品接口 (v2)
- `GET /api/v2/userlist` - 获取用户列表
- `GET /api/v2/plist` - 获取商品列表

## 🔍 搜索功能

基于 Elasticsearch 7.9.0 实现：
- **全文搜索**: 商品标题、描述全文检索
- **分词搜索**: 支持中文分词 (IK 分词器)
- **搜索建议**: 自动完成和搜索推荐
- **搜索过滤**: 价格、分类、品牌等条件过滤

## 💡 开发说明

### 项目特色
1. **模块化设计**: 清晰的分层架构，易于维护和扩展
2. **缓存优化**: Redis 缓存机制，提升系统性能
3. **搜索集成**: Elasticsearch 强大的搜索能力
4. **支付集成**: 完整的支付流程和回调处理
5. **权限控制**: 细粒度的权限管理系统

### 代码规范
- 遵循 Go 语言标准代码规范
- 使用 GORM 进行数据库操作
- 统一的错误处理和日志记录
- RESTful API 设计原则

