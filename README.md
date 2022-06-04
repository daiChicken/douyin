# BytesDanceProject

以下 config.yaml 供参考test

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
