package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path"`
	GRPC GRPCConf `yaml:"grpc"`
	InMemoryStorage bool
}

type GRPCConf struct {
	Port    int `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil{
		panic("Error loading .env file")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == ""{
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err){
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err !=nil{
		panic("failed to read config: " + err.Error())
	}

	var isInMemoryStorage bool

	flag.BoolVar(&isInMemoryStorage, "in_memory", false, "use in-memory storage")
	flag.Parse()

	cfg.InMemoryStorage = isInMemoryStorage

	return &cfg
}