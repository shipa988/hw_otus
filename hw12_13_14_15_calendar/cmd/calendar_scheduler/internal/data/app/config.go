package app

type Config struct {
	Log      Log       `yaml:"log"`
	RepoType RepoType  `yaml:"repotype"`
	DB       DB        `yaml:"db"`
	Queue    QueueType `yaml:"queue"`
	Rabbit   Rabbit    `yaml:"rabbit"`
}
type Log struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

type Rabbit struct {
	Addr         string `yaml:"addr"`
	ExchangeName string `yaml:"exchangeName"`
	QueueName    string `yaml:"queueName"`
}

type DB struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}
