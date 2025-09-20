package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`
	// Conn_Str    string      `yaml:"conn_str" env-required:"true"`
	Http_Server Http_Server `yaml:"http_server"`
}

type Http_Server struct {
	Addr string `yaml:"address" env-required:"true"`
}

func Must_Load() *Config {
	var config_path string
	config_path = os.Getenv("config_path")
	if config_path == "" {
		flags := flag.String("config", "", "this is path for config_path")
		flag.Parse()
		config_path = *flags
		if config_path == "" {
			log.Fatal("path is empty for some reason, maybe its because you didnt send it as a flag")
		}
	}
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Fatal("yo checked the file, its not present")
	}
	var cfg Config
	err := cleanenv.ReadConfig(config_path, &cfg)
	if err != nil {
		log.Fatalf("couldnt read the file: %v", err)
	}
	return &cfg
}
