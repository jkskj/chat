package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var errBusyGroup = "BUSYGROUP Consumer Group name already exists"

// Msg 消息
type Msg struct {
	ID        string // 消息的编号
	Topic     string // 消息的主题
	Body      []byte // 消息的Body
	Partition int    // 分区号
	Group     string // 消费者组
	Consumer  string // 消费者组里的消费者
}

var SingleStreamMQ *StreamMQ
var GroupStreamMQ *StreamMQ

// Handler 返回值代表消息是否消费成功
type Handler func(msg *Msg) error

type StreamMQ struct {
	// Redis客户端
	client *redis.Client
	// 最大消息数量，如果大于这个数量，旧消息会被删除，0表示不管
	maxLen int64
	// Approx是配合MaxLen使用的，表示几乎精确的删除消息，也就是不完全精确，由于stream内部是流，所以设置此参数xadd会更加高效
	approx bool
}

func NewStreamMQ(client *redis.Client, maxLen int, approx bool) *StreamMQ {
	return &StreamMQ{
		client: client,
		maxLen: int64(maxLen),
		approx: approx,
	}
}

// Stream：表示流的名字
// ID：表示消息的唯一编号，它由两部分组成<unix_time_milliseconds>-<sequence_number_in_same_millisecond>，
//也就是当前时间戳毫秒数加上当前毫秒数的序列号，类似于雪花算法。

// 可以像这样1526919030474-55自己指定时间戳和序列号；
// 也可以像这样1526919030474-*自己指定时间戳和让Redis生成序列号；
// 还可以像这样*让Redis生成时间戳和序列号。我们这里就是使用这种形式。

// Values：表示消息内容键值对。我们的消息只有body，所以只设置了一个键值对。
// MaxLen和Approx：指定了MaxLen可以让Redis在消息数量大于此值时删除旧消息，避免内存溢出；
//而Approx是配合MaxLen使用的，表示几乎精确的删除消息，也就是不完全精确，由于stream内部是流，所以设置此参数xadd会更加高效。

func (q *StreamMQ) SendMsg(ctx context.Context, msg *Msg) error {
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: msg.Topic,
		MaxLen: q.maxLen,
		Approx: q.approx,
		ID:     "*",
		Values: []interface{}{"body", msg.Body},
	}).Err()
}

//首先，我们使用xgroup create命令创建消费者组，当然这里可能会重复创建，重复创建会报错，
//我们忽略这个错误。其中的start参数表示该消费者组从哪个位置开始消费消息，可以指定为ID或$，其中$表示从最后一条消息开始消费。
//然后我们使用xreadgroup命令阻塞的获取消息，其中参数的含义是：
//
//Group：消费者组
//Consumer：消费者组里的消费者
//Streams：消费的流。
//
//这里后面还有一个">"其实是属于ID参数，表示只接收未投递给其他消费者的消息;
//如果指定ID为数值，则表示只接收大于这个ID的已经被拉取却没有被ACK的消息。
//所以我们这里先使用>拉取一次最新消息，再使用0拉取已经投递却没有ACK的消息，保证消息都能够成功消费。
//
//
//Count：一次性读取消息的条数，减少网络传输时间。
//
//如果成功消费，我们再使用xack指令提交消费位点，这样这条消息就不会再次被投递了。

// Consume 返回值代表消费过程中遇到的无法处理的错误
// group 消费者组
// consumer 消费者组里的消费者
// batchSize 每次批量获取一批的大小
// start 用于创建消费者组的时候指定起始消费ID，0表示从头开始消费，$表示从最后一条消息开始消费
func (q *StreamMQ) Consume(ctx context.Context, topic, group, consumer, start string, batchSize int, h Handler) error {
	err := q.client.XGroupCreateMkStream(ctx, topic, group, start).Err()
	if err != nil && err.Error() != errBusyGroup {
		return err
	}
	for {
		// 拉取新消息
		if err1 := q.consume(ctx, topic, group, consumer, ">", batchSize, h); err1 != nil {
			return err
		}
		// 拉取已经投递却未被ACK的消息，保证消息至少被成功消费1次
		if err2 := q.consume(ctx, topic, group, consumer, "0", batchSize, h); err2 != nil {
			return err
		}
	}
}

func (q *StreamMQ) consume(ctx context.Context, topic, group, consumer, id string, batchSize int, h Handler) error {
	// 阻塞的获取消息
	result, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{topic, id},
		Count:    int64(batchSize),
	}).Result()
	if err != nil {
		return err
	}
	// 处理消息
	for _, msg := range result[0].Messages {
		err1 := h(&Msg{
			ID:       msg.ID,
			Topic:    topic,
			Body:     []byte(msg.Values["body"].(string)),
			Group:    group,
			Consumer: consumer,
		})
		if err1 == nil {
			err2 := q.client.XAck(ctx, topic, group, msg.ID).Err()
			if err2 != nil {
				return err
			}
		}
	}
	return nil
}
