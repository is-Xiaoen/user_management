🚀 UserHub - 现代化用户管理系统

<div align="center">   <img src="https://img.shields.io/badge/Go-1.22.5-00ADD8?style=for-the-badge&logo=go" alt="Go Version">   <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL">   <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License">   <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge" alt="PRs Welcome"> </div> <div align="center">   <h3>🎯 安全 • 高效 • 易用</h3>   <p>基于 Go 语言构建的企业级用户管理解决方案</p> </div>

---

✨ 特性亮点

🛡️ 安全至上

- BCrypt 密码加密 - 业界标准的密码哈希算法
- CSRF 防护 - 完善的跨站请求伪造防护机制
- 会话管理 - 服务器端会话存储，支持记住登录状态
- SQL 注入防护 - 参数化查询，杜绝注入风险

⚡ 极致性能

- Go 语言驱动 - 编译型语言，毫秒级响应
- 连接池优化 - 智能管理数据库连接
- 并发安全 - 精心设计的锁机制，支持高并发
- 内存优化 - 高效的内存使用和垃圾回收

🎨 现代化 UI

- 响应式设计 - 完美适配各种设备
- 玻璃态效果 - 时尚的毛玻璃视觉效果
- 流畅动画 - 细腻的交互动画体验
- 深色主题 - 护眼的暗色界面设计

🔧 开发友好

- 清晰架构 - MVC + Repository 分层设计
- 模块化 - 高内聚低耦合的代码组织
- 完整日志 - 详细的操作日志和错误追踪
- 易于扩展 - 灵活的接口设计

🛠️ 技术栈

<table> <tr> <td align="center" width="96"> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" width="48" height="48" alt="Go" /> <br>Go </td> <td align="center" width="96"> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mysql/mysql-original.svg" width="48" height="48" alt="MySQL" /> <br>MySQL </td> <td align="center" width="96"> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/html5/html5-original.svg" width="48" height="48" alt="HTML5" /> <br>HTML5 </td> <td align="center" width="96"> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/css3/css3-original.svg" width="48" height="48" alt="CSS3" /> <br>CSS3 </td> <td align="center" width="96"> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/javascript/javascript-original.svg" width="48" height="48" alt="JavaScript" /> <br>JavaScript </td> </tr> </table>

核心依赖

- golang.org/x/crypto - 密码学支持
- go-sql-driver/mysql - MySQL 驱动
- Bootstrap 5 - CSS 框架
- Font Awesome - 图标库

🚀 快速开始

前置要求

- Go 1.22.5+
- MySQL 8.0+
- Git

1. 克隆项目

    git clone https://github.com/yourusername/userhub.git
    cd userhub

2. 配置数据库

    CREATE DATABASE user_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

3. 修改配置

编辑 config/config.go 文件，更新数据库连接信息：

    return &Config{
        DBHost:     "localhost",
        DBPort:     "3306", 
        DBUser:     "your_username",
        DBPassword: "your_password",
        DBName:     "user_management",
        ServerPort: "8080",
    }

4. 安装依赖

    go mod download

5. 运行项目

    go run main.go

访问 http://localhost:8080 即可看到系统界面！

6. 默认管理员账号

首次运行后，使用以下 SQL 创建管理员账号：

    INSERT INTO users (username, password, email, role) VALUES 
    ('admin', '$2a$10$YourHashedPasswordHere', 'admin@example.com', 'admin');

或使用测试账号：

- 用户名：admin
- 密码：admin123

📁 项目结构

    user-management-system/
    ├── 📱 app/                 # 应用核心
    │   └── app.go             # 依赖注入容器
    ├── ⚙️ config/              # 配置管理
    ├── 🎮 controllers/         # 控制器层
    │   ├── auth.go            # 认证控制
    │   └── user.go            # 用户管理
    ├── 🗄️ database/            # 数据库连接
    ├── ⚠️ errors/              # 错误处理
    ├── 📝 logger/              # 日志系统
    ├── 🔒 middleware/          # 中间件
    ├── 📊 models/              # 数据模型
    ├── 💾 repository/          # 数据访问层
    │   ├── interfaces/        # 接口定义
    │   └── mysql/             # MySQL实现
    ├── 🛣️ router/              # 路由配置
    ├── 🔐 session/             # 会话管理
    ├── 🎨 static/              # 静态资源
    │   ├── css/               # 样式文件
    │   └── js/                # JavaScript
    ├── 🖼️ views/               # 视图模板
    └── 🚀 main.go             # 程序入口

🌟 核心功能

用户认证

- ✅ 用户注册（邮箱验证）
- ✅ 用户登录（记住我功能）
- ✅ 安全登出
- ✅ 密码强度检测

用户管理

- ✅ 用户列表展示
- ✅ 用户搜索过滤
- ✅ 用户信息编辑
- ✅ 用户删除（权限控制）

权限系统

- ✅ 角色管理（管理员/普通用户）
- ✅ 权限中间件
- ✅ 操作日志记录

安全特性

- ✅ CSRF Token 保护
- ✅ SQL 注入防护
- ✅ XSS 防护
- ✅ 会话劫持防护

📸 界面预览

🏠 首页

- 现代化的 Hero 区域
- 渐变文字效果
- 浮动卡片动画
- 响应式布局

🔐 登录/注册

- 玻璃态卡片设计
- 实时密码强度提示
- 优雅的错误提示
- 记住登录状态

👥 用户管理

- 实时搜索过滤
- 角色标签显示
- 快速编辑功能
- 批量操作支持

🔧 高级配置

数据库连接池

    db.SetMaxOpenConns(25)          // 最大连接数
    db.SetMaxIdleConns(5)           // 最大空闲连接
    db.SetConnMaxLifetime(5 * time.Minute)  // 连接生命周期

会话配置

    sessionManager := session.NewManager(
        "session_id",      // Cookie 名称
        2*time.Hour,       // 默认过期时间
    )

日志配置

- 自动按日期轮转
- 分级日志记录
- 用户操作审计

🚦 API 文档

认证接口

  方法  	路径       	描述  	权限  
  GET 	/login   	登录页面	无   
  POST	/login   	用户登录	无   
  GET 	/register	注册页面	无   
  POST	/register	用户注册	无   
  POST	/logout  	用户登出	登录用户

用户管理接口

  方法  	路径           	描述  	权限  
  GET 	/users       	用户列表	登录用户
  POST	/users/update	更新用户	管理员 
  POST	/users/delete	删除用户	管理员 

🤝 贡献指南

我们欢迎所有形式的贡献！无论是新功能、bug 修复还是文档改进。

开发流程

1. Fork 本项目
2. 创建特性分支 (git checkout -b feature/AmazingFeature)
3. 提交更改 (git commit -m 'Add some AmazingFeature')
4. 推送到分支 (git push origin feature/AmazingFeature)
5. 提交 Pull Request

代码规范

- 遵循 Go 官方编码规范
- 保持代码简洁清晰
- 添加必要的注释
- 编写单元测试

📈 性能优化

- 数据库索引优化 - 为常用查询字段添加索引
- 查询优化 - 使用预编译语句，减少数据库压力
- 静态资源缓存 - 浏览器缓存策略
- GZIP 压缩 - 减少传输数据量

🔐 安全最佳实践

1. 定期更新依赖 - 保持所有依赖包为最新版本
2. 强密码策略 - 实施密码复杂度要求
3. 限制登录尝试 - 防止暴力破解
4. HTTPS 部署 - 生产环境使用 SSL/TLS
5. 定期备份 - 数据库定期备份策略

📜 更新日志

v1.0.0 (2024-01-15)

- 🎉 首次发布
- ✨ 完整的用户管理功能
- 🔐 CSRF 防护实现
- 🎨 现代化 UI 设计

🏆 致谢

感谢所有为这个项目做出贡献的开发者！

特别感谢以下开源项目：

- Go
- MySQL
- Bootstrap
- Font Awesome

📄 许可证

本项目采用 MIT 许可证 - 查看 LICENSE 文件了解详情

---

<div align="center">   <p>如果这个项目对您有帮助，请给我一个 ⭐️</p>   <p>Made with ❤️ by UserHub Team</p> </div>


