package app

type Config struct {
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
	HTTPPort   string `yaml:"httpport"`
	GRPCGWPort string `yaml:"grpcgwport"`
	GRPCPort   string `yaml:"grpcport"`
}

type DB struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}
