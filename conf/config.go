package conf

// Config 整个项目的配置
type Config struct {
	Kafka Kafka `ini:"kafka"`
	Log   Log   `ini:"log"`
}

// Kafka kafka的配置
type Kafka struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

// Log 日志的配置
type Log struct {
	FileName string `ini:"filename"`
}