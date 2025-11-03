package configs

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var cfg *Config

func Load() *Config {
	// Return cached config if already loaded
	if cfg != nil {
		return cfg
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	BindAllKeys()

	viper.SetDefault("SERVER_PORT", ":8080")

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	cfg = &c // cache it
	return cfg
}
