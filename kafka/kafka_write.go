package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type ClientKafkaWrite struct {
	writer *kafka.Writer
}

func NewKafkaWriterClient(brokers []string, topic string, balancer kafka.Balancer, async bool) *ClientKafkaWrite {
	s := &ClientKafkaWrite{}
	s.writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     balancer,
		BatchTimeout: time.Microsecond,
		Async:        async,
	}
	return s
}

func (s *ClientKafkaWrite) Push(ctx context.Context, key []byte, value []byte) (err error) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		e := recover()
		if e != nil {
			err = fmt.Errorf("[Kafka Push]-runtime error: %v", e)
		}
	}()
	err = s.writer.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
	return
}

func (s *ClientKafkaWrite) Store(ctx context.Context, key []byte, data []interface{}) (err error) {
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		e := recover()
		if e != nil {
			err = fmt.Errorf("[Kafka ClientKafkaWrite]-runtime error: %v", e)
		}
	}()
	var messages []kafka.Message
	for _, txLog := range data {
		val, err := json.Marshal(txLog)
		if err != nil {
			return err
		}
		messages = append(messages, kafka.Message{Key: key, Value: val})
	}
	err = s.writer.WriteMessages(ctx, messages...)
	return
}

func (s *ClientKafkaWrite) Close() error {
	return s.writer.Close()
}

func (s *ClientKafkaWrite) Stats() kafka.WriterStats {
	return s.writer.Stats()
}
