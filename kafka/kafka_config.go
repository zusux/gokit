package kafka

type ConfigKafka struct {
	Brokers []string `yaml:"Brokers"`
	Async   bool     `yaml:"Async"`
	Timeout int      `yaml:"Timeout"`
	Name    string   `yaml:"Name"`
	Group   string   `yaml:"Group"`
	Topic   string   `yaml:"Topic"`
}
