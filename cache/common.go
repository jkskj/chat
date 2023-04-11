package cache

import (
	"chat/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
	"sync"
	"time"
)

// RedisClient Redis缓存客户端单例
var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

// Init Redis 初始化redis链接
func Init() {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		fmt.Println("Redis 配置文件读取错误，请检查文件路径:", err)
	}
	LoadRedisData(file)
	Redis()
	SingleStream()
	GroupStream()
}

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",     //网络类型，tcp or unix，默认tcp
		Addr:     RedisAddr, //主机名+冒号+端口，默认localhost:6379
		Password: "",        //密码
		DB:       int(db),   // redis数据库index
		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		//Dialer: func() (net.Conn, error) {
		//	netDialer := &net.Dialer{
		//		Timeout:   5 * time.Second,
		//		KeepAlive: 5 * time.Minute,
		//	}
		//	return netDialer.Dial("tcp", "127.0.0.1:6379")
		//},

		//钩子函数
		//OnConnect: func(conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		//	fmt.Printf("conn=%v\n", conn)
		//	return nil
		//},
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
func SingleStream() {
	SingleStreamMQ = NewStreamMQ(RedisClient, 100, true)
	topic := "single"
	count := 10
	var wg sync.WaitGroup
	wg.Add(count * 4)
	go SingleStreamMQ.Consume(context.Background(), topic, "group1", "consumer1", "0", 5, func(msg *Msg) error {
		fmt.Printf("consume group1 consumer1: %+v\n", msg)
		var msg1 model.SingleMessage
		err := json.Unmarshal(msg.Body, &msg1)
		if err != nil {
			fmt.Println("json转换失败", err)
		} else {
			singleMysqlSave(msg1)
		}
		wg.Done()
		return nil
	})
}
func GroupStream() {
	GroupStreamMQ = NewStreamMQ(RedisClient, 100, true)
	topic := "group"
	count := 10
	var wg sync.WaitGroup
	wg.Add(count * 4)
	go GroupStreamMQ.Consume(context.Background(), topic, "group2", "consumer2", "0", 5, func(msg *Msg) error {
		fmt.Printf("consume group1 consumer1: %+v\n", msg)
		var msg1 model.GroupMessage
		err := json.Unmarshal(msg.Body, &msg1)
		if err != nil {
			fmt.Println("json转换失败", err)
		} else {
			groupMysqlSave(msg1)
		}
		wg.Done()
		return nil
	})
}
func singleMysqlSave(replyMsg model.SingleMessage) {
	err := model.DB.Save(&replyMsg).Error
	if err != nil {
		fmt.Println("消息mysql存储失败！！！")
	}
}
func groupMysqlSave(replyMsg model.GroupMessage) {
	err := model.DB.Save(&replyMsg).Error
	if err != nil {
		fmt.Println("消息mysql存储失败！！！")
	}
}
