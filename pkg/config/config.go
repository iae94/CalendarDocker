package config

import (
	"github.com/spf13/viper"
	"log"
)
type Logger struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"outputPaths"`
	ErrorOutputPaths []string `mapstructure:"errorOutputPaths"`
}
type DB struct {
	User string `mapstructure:"user"`
	Password int `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`
	Host string `mapstructure:"host"`
}

type Rabbit struct {
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Vhost string `mapstructure:"vhost"`
}

type APIConfig struct {
	API struct {
		Port   int
		Logger Logger `mapstructure:"Logger"`
		DB DB `mapstructure:"DB"`
		Rabbit Rabbit `mapstructure:"Rabbit"`
	} `mapstructure:"API"`
}
type SchedulerConfig struct {
	Scheduler struct {
		Cleaner struct{
			CleanDelay int `mapstructure:"cleanDelay"`
		} `mapstructure:"Cleaner"`
		Port   int
		Logger Logger `mapstructure:"Logger"`
		DB DB `mapstructure:"DB"`
		Rabbit Rabbit `mapstructure:"Rabbit"`
	} `mapstructure:"Scheduler"`
}

type NotificatorConfig struct {
	Notificator struct {
		Port   int
		Logger Logger `mapstructure:"Logger"`
		Rabbit Rabbit `mapstructure:"Rabbit"`
	} `mapstructure:"Notificator"`
}

func ReadAPIConfig() (conf *APIConfig, err error) {
	viper.SetConfigName("api")
	viper.AddConfigPath("../../../configs")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Reading api config error: %v \n", err)
		return nil, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("Unmarshaling api config error: %v \n", err)
		return nil, err
	}

	return conf, nil
}

func ReadSchedulerConfig() (conf *SchedulerConfig, err error) {
	viper.SetConfigName("scheduler")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Reading scheduler config error: %v \n", err)
		return nil, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("Unmarshaling scheduler config error: %v \n", err)
		return nil, err
	}

	return conf, nil
}

func ReadNotificatorConfig() (conf *NotificatorConfig, err error) {
	viper.SetConfigName("notificator")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Reading notificator config error: %v \n", err)
		return nil, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("Unmarshaling notificator config error: %v \n", err)
		return nil, err
	}

	return conf, nil
}