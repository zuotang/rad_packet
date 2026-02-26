package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	MySQL struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
	JWT struct {
		Secret   string `mapstructure:"secret"`
		TTLHours int    `mapstructure:"ttl_hours"`
	} `mapstructure:"jwt"`
	Admin struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"admin"`
}

func Load() (Config, error) {
	var cfg Config
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetDefault("server.port", "8080")
	v.SetDefault("jwt.ttl_hours", 168)
	v.SetDefault("admin.key", "change-admin-key")

	if err := v.ReadInConfig(); err != nil {
		// Allow env-only startup.
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
