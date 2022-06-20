# mini抖音

## 用户模块
- bcrypt加密存储密码
- 基于JWT的认证中间件
## 视频模块
- 视频上传七牛云
- 使用携程将上传进行异步
- 失败重试
- 使用事务保证云端与mysql的数据一致性
## 评论模块
- 使用redis完成
- 敏感词过滤
## 点赞模块
- 使用redis完成
## 关注模块
- 使用redis完成

## 配置文件示例
```yaml
name: "BytesDanceProject"
mode: "dev"
port: 8080
version: "v0.1.3"
start_time: "2022-05-11"
machine_id: 1

auth:
  jwt_expire: 8760

log:
  level: "debug"
  filename: "web_log"
  max_size: 200
  max_age: 30
  max_backups: 7
  
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "123456"
  dbname: "douyin"
  max_conns: 200
  max_idle_conns: 50
  
redis:
  host: "127.0.0.1"
  port: 6379
  db: 0
  password: ""
  poolsize: 1

qiniuyun:
  access_key: "xxx"
  secret_key: "xxx"
  domain: "xxx"
  bucket: "xxx"
```
