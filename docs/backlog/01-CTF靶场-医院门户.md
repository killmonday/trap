# 01-CTF靶场-医院门户 Backlog

## BL-CTF-001 医院门户网站靶场搭建

- 模块：CTF 靶场-医院门户
- 优先级：P0
- 状态：已完成
- 当前阶段：阶段一
- 价值：为 CTF 参赛者提供一个真实的医院门户网站靶场环境，通过目录扫描发现源码备份文件获取 flag。
- 关联模块文档：docs/modules/01-CTF靶场-医院门户模块.md
- 关联计划：plans/2026-07-15-BL-CTF-001-医院门户靶场搭建.md
- 关联代码：main.go, templates/, static/
- 验收状态：待验收

### 范围

- 搭建基于 Golang 的 Web 应用，模拟医院门户网站（首页、科室、医生、简介、联系）
- 在 HTML/JS/robots.txt 中埋入 hint，引导做题人发现源码备份文件
- 提供 /backup/ 路径下的文件下载功能
- 做好路径遍历安全防护

### 不做

- 不制作真实的后端数据库
- 不实现用户登录/注册
- 不制作压缩包（由用户自行创建）
- 不做后端管理面板

### 验收条件

- [x] 访问 /, /index, /departments, /doctors, /about, /contact 均返回 200
- [x] robots.txt 包含 Disallow: /backup/ 条目
- [x] 首页 HTML 源码中包含 backup 路径的 hint（hidden 元素）
- [x] main.js 中包含 backup 路径的 hint（注释 + console.log）
- [x] /backup/wwwroot.zip 在文件不存在时返回 404
- [x] 路径遍历攻击被阻止（/backup/../../../etc/passwd 返回 404）
- [x] 所有测试通过（7/7）

### 依赖

- Go 1.22+
- 无外部依赖

### 状态记录

- 2026-07-15：创建并完成。包含 5 个页面、CSS/JS 静态文件、robots.txt、7 个测试用例。
