package config

import "github.com/spf13/viper"

type Config struct {
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	SecretKey        string `mapstructure:"SECRET_KEY"`
}

func LoadConfig(path, fileName string) (*Config, error) {
	viper.SetConfigFile(fileName)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	conf := Config{}

	err := viper.ReadInConfig()
	if err != nil {
		return &conf, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return &conf, err
	}

	return &conf, nil
}
