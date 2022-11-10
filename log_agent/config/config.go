package config

type taskConfig struct {
	Conf string `ini:"conf"`
}

type esConfig struct {
	Addr    string `ini:"addr"`
	Process int    `ini:"process"`
}

type Config struct {
	Task     taskConfig `ini:"task"`
	ESConfig esConfig   `ini:"es"`
}

type TaskConf struct {
	Path  string `json:"path"`
	Index string `json:"index"`
}
