package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	SecretKey        string `mapstructure:"SECRET_KEY"`
}

func LoadConfig() (*Config, error) {
	var path, fileName string

	flag.StringVar(&path, "path", ".", "path to config file")
	flag.StringVar(&fileName, "name", ".env", "config file name")
	flag.Parse()

	//name := fmt.Sprintf("%s/%s", path, fileName)

	viper.SetConfigFile(".env")

	//viper.AddConfigPath(path)
	viper.AutomaticEnv()

	fmt.Println(path, fileName)

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
