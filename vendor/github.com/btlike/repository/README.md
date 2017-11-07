## repository
[![Build Status](https://drone.io/github.com/btlike/repository/status.png)](https://drone.io/github.com/btlike/repository/latest)


接口化的数据仓库，默认实现了Mysql，你可以实现自己的数据仓库并注册进来（Postgresql等）
- 定义顶层数据操作接口
- 实现默认数据仓库（Mysql）

### Mysql表结构

- torrent0 ~ torrentf 分表存储metadata
- history 存储搜索历史
- infohash 存储爬虫抓取到的活跃infohash
- recommend 存储网站首页显示的推荐关键词


## 初始化

- 创建名为 torrent的库
- 执行[mysql.sql](https://github.com/btlike/repository/blob/master/mysql.sql)中的sql语句

### 安装
`go get github.com/btlike/repository`
