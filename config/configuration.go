package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

type Configuration struct {
	BardApi    BardApiConfiguration
	WeatherApi WeatherApiConfiguration
	SpotifyApi SpotifyApiConfiguration
	Memcache   MemcacheConfiguration
}

type BardApiConfiguration struct {
	Port string
}

type WeatherApiConfiguration struct {
	Appid string
	Url   string
}

type SpotifyApiConfiguration struct {
	ClientId  string
	SecretKey string
}

type MemcacheConfiguration struct {
	Host string
}

var Config Configuration

func init() {
	var err error
	configPaths := "."

	pwd, _ := os.Getwd()
	if basePath := path.Base(pwd); basePath != "bard" {
		configPaths = "../"
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(configPaths)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	env := os.Getenv("ENV")

	err = viper.UnmarshalKey(env, &Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
