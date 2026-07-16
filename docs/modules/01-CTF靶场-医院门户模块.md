# 01-CTF靶场-医院门户模块

## 模块概述

基于 Golang 的医院门户网站 CTF 靶场。模拟一家名为"仁和医院"的三甲医院官网，包含首页、科室介绍、专家团队、医院简介、联系我们共 5 个页面。

## 核心业务流

### 1. 正常浏览流程

用户访问首页 → 查看各页面 → 浏览医院信息

### 2. CTF 解题流程

1. 访问网站首页，查看页面源码 → 发现 hidden 元素中的 hint
2. 查看 robots.txt → 发现 Disallow: /backup/ 条目
3. 打开浏览器控制台 → 看到 dev hint 日志
4. 使用目录扫描工具（gobuster/dirb/ffuf）扫描 → 发现 /backup/ 路径
5. 根据 hint 猜出文件名 wwwroot.zip → 下载 /backup/wwwroot.zip
6. 解压 → 获得 flag

## Hint 埋点

| 位置 | 形式 | 内容 |
|------|------|------|
| HTML footer | hidden span 元素 + data-backup-path 属性 | 运维备注：源码备份文件位于 backup 目录 |
| main.js | 代码注释 | 开发备忘：源码快照存放在 /backup/wwwroot.zip |
| main.js | JS 变量 | backupPath: '/backup/', backupFile: 'wwwroot.zip' |
| main.js | console.log | Dev hint: source snapshot available at /backup/ |
| robots.txt | Disallow | Disallow: /backup/ (三条假路径 + 一条真路径) |

## 安全设计

- 路径遍历防护：使用 filepath.Base 校验，拒绝 `..` 穿越
- 文件不存在时返回 404，不泄露路径信息
- robots.txt 中 /admin/、/config/、/api/ 为干扰项，/backup/ 为真实下载路径

## 文件结构

```
main.go              — 入口 + 路由 + 下载处理
templates/
  layout.html        — 公共布局（含 header/footer/hint）
  index.html         — 首页
  departments.html   — 科室介绍
  doctors.html       — 专家团队
  about.html         — 医院简介
  contact.html       — 联系我们
static/
  css/style.css      — 样式表
  js/main.js         — 前端脚本（含 hint）
  robots.txt         — 爬虫规则（含 Disallow hint）
backup/
  .gitkeep           — 占位文件
  wwwroot.zip        — 用户自行放置（含 flag）
```

## 构建与运行

```bash
go build -o renhe-hospital .
./renhe-hospital
# 默认监听 :8080，可通过 PORT 环境变量修改
```

## 压缩包说明

- 文件名：`wwwroot.zip`
- 放置路径：`backup/wwwroot.zip`（项目根目录下的 backup 子目录）
- URL：`http://<host>:8080/backup/wwwroot.zip`
- 内容：由用户自定义，flag 放在解压后的某个文件中
