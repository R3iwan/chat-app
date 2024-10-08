package config

import "github.com/spf13/viper"

type Config struct {
	PostgresURL string
	KafkaURL    string
	Port        string
	JWTSecret   string
}

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	config := &Config{
		PostgresURL: viper.GetString("POSTGRES_URL"),
		KafkaURL:    viper.GetString("KAFKA_URL"),
		Port:        viper.GetString("PORT"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
	}

	return config, nil
}
