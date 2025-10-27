package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
	Env         string `mapstructure:"env"`
	Port        int    `mapstructure:"port"`
	ShowLog     bool   `mapstructure:"showLog"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSL      bool   `mapstructure:"ssl"`
}

type FiberConfig struct {
	BodyLimit             string        `mapstructure:"bodyLimit"`
	CaseSensitive         bool          `mapstructure:"caseSensitive"`
	StrictRouting         bool          `mapstructure:"strictRouting"`
	Prefork               bool          `mapstructure:"prefork"`
	Concurrency           int           `mapstructure:"concurrency"`
	ReadTimeout           time.Duration `mapstructure:"readTimeout"`
	WriteTimeout          time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout           time.Duration `mapstructure:"idleTimeout"`
	EnablePrintRoutes     bool          `mapstructure:"enablePrintRoutes"`
	EnableIPValidation    bool          `mapstructure:"enableIPValidation"`
	Immutable             bool          `mapstructure:"immutable"`
	TrustedProxies        []string      `mapstructure:"trustedProxies"`
	ProxyHeader           string        `mapstructure:"proxyHeader"`
	DisableStartupMessage bool          `mapstructure:"disableStartupMessage"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Fiber    FiberConfig    `mapstructure:"fiber"`
}

func LoadAppConfig() *Config {
	env := os.Getenv("APP_ENV")
	fileName := fmt.Sprintf("env.yaml", env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigFile("./config/env.yaml")

	// Читаем конфигурацию
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ Не удалось загрузить конфигурацию (%s): %v", fileName, err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("❌ Ошибка при парсинге YAML: %v", err)
	}

	log.Printf("✅ Загружена конфигурация: %s (%s)\n", cfg.App.Name, cfg.App.Env)
	return &cfg
}

func (cfg *Config) GetDatabaseURL() string {
	sslMode := "disable"
	if cfg.Database.SSL {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		sslMode,
	)
}
