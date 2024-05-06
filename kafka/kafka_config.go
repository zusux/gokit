package kafka

type ConfigKafka struct {
	Brokers []string
	Async   bool
	Timeout int
	Name    string
	Group   string
	Topic   string
}
