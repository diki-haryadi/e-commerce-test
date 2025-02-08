package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"path/filepath"
	"runtime"
)

type Config struct {
	App AppConfig
}

var BaseConfig *Config

type AppConfig struct {
	AppEnv    string `json:"app_env" envconfig:"APP_ENV"`
	AppName   string `json:"app_name" envconfig:"APP_NAME"`
	JWTSecret string `json:"jwt_secret" envconfig:"JWT_SECRET"`
}

func LoadConfig() *Config {
	_, callerDir, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Error generating env dir")
	}

	// Define the possible paths to the .env file
	envPaths := []string{
		filepath.Join(filepath.Dir(callerDir), "..", "envs/.env"),
	}
	_ = godotenv.Overload(envPaths[0])
	var configLoader Config

	if err := envconfig.Process("BaseConfig", &configLoader); err != nil {
		log.Printf("error load config: %v", err)
	}

	BaseConfig = &configLoader
	spew.Dump(configLoader)
	return &configLoader
}
