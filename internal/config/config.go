package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	DB         `yaml:"db"`
	SMTP       `yaml:"smtp"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"password123"`
	DBName   string `yaml:"db_name" env-default:"debBubble"`
	Port     string `yaml:"port" env-default:"5632"`
}

type SMTP struct {
	Host     string `yaml:"host" env-default:"smtp.gmail.com"`
	Port     string `yaml:"port" env-default:"587"`
	From     string `yaml:"from" env-default:"devbubbleinc@gmail.com"`
	Password string `yaml:"password" env-default:"myxieiamsfplpdxo"`
}

type JWT struct {
	Secret              string `yaml:"secret" env-default:"ffgdgfhhgfgfhghfhgdhffhdghfgfhgfhghdgf3t43tt33te"`
	ExpirationTimeHours int    `yaml:"expirationTimeHours" env-default:72`
}

func MustLoad() *Config {
	configPath := "../../config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
