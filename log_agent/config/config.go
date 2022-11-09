package config

type kafkaConfig struct {
	Addr  []string `ini:"addr"`
	ChannelMaxSize int `ini:"channel_max_size"`
	ProcessNum int `ini:"process_num"`
}

type etcdConfig struct {
	Addr []string `ini:"addr"`
	Timeout int `ini:"timeout"`
	Key string `ini:"key"`
}

type Config struct {
	Etcd  etcdConfig  `ini:"etcd"`
	KafKa kafkaConfig `ini:"kafka"`
}

type LogConf struct {
	Filename  string `json:"path"`  // tail文件路径
	Topic string `json:"topic"` // kafka主题
}
