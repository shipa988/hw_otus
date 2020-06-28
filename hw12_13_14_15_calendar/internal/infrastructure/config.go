package infrastructure

type AppConfig struct {
	Log      Log    `yaml:"log"`
	API      API    `yaml:"api"`
	RepoType string `yaml:"repotype"`
	DB       DB     `yaml:"db"`
}
type Log struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

type API struct {
	Port string `yaml:"port"`
}

type DB struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}
