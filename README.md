# chat
## conf/conf.ini
```
#debug开发模式,release生产模式
[service]
AppMode = debug
HttpPort = :3000 
# 运行端口号 3000端口

[mysql]
Db = mysql
DbHost = "" 
# mysql的ip地址
DbPort = ""
# mysql的端口号,默认3306
DbUser = ""
# mysql user
DbPassWord = ""
# mysql password
DbName = ""
# 数据库名字

[redis]
RedisDb = ""
# redis 名字
RedisAddr = ""
# redis 地址
RedisPw = ""
# redis 密码
RedisDbName = ""
# redis 数据库名
```
