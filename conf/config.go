package conf

// Config 整个项目的配置
type Config struct {
	Kafka Kafka `ini:"kafka"`
	Etcd  Etcd  `ini:"etcd"`
}

// Kafka kafka的配置
type Kafka struct {
	Address     string `ini:"address"`
	ChanMaxSize int    `ini:"chan_max_size"`
	// Topic   string `ini:"topic"`
}

// Etcd etcd的配置
type Etcd struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeout"`
}
