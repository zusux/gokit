package kafka

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CheckOrCreateTopic(brokers []string, topic string, numPartitions int, replicationFactor int) error {
	if len(brokers) <= 0 {
		return fmt.Errorf("kafka address empty")
	}
	var conn *kafka.Conn
	var err error
	for _, v := range brokers {
		conn, err = kafka.Dial("tcp", v)
		if err != nil {
			continue
		}
	}
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}
	m := map[string]struct{}{}
	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	//已经存在topic
	if _, ok := m[topic]; ok {
		return nil
	}

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}
	return controllerConn.CreateTopics(topicConfigs...)
}
