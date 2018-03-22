package main

import (
	"flag"
)

type MqttConfig struct {
	Host     string
	User     string
	Password string

	BaseTopic string
}

func (config MqttConfig) BrokerUri() string {
	uri := "tcp://"
	if config.User != "" {
		uri += config.User

		if config.Password != "" {
			uri += ":" + config.Password
		}

		uri += "@"
	}

	uri += config.Host

	return uri
}

type Config struct {
	Mqtt     *MqttConfig
	RepoPath string
}

func parseFlags() *Config {
	repo := flag.String("path", "", "Path to files")
	host := flag.String("host", "localhost:1883", "MQTT broker host")
	user := flag.String("user", "", "MQTT broker host")
	password := flag.String("password", "", "MQTT broker host")
	baseTopic := flag.String("topic", "klang3", "MQTT base topic")

	flag.Parse()

	mqttConfig := &MqttConfig{
		Host:     *host,
		User:     *user,
		Password: *password,

		BaseTopic: *baseTopic,
	}

	config := &Config{
		Mqtt:     mqttConfig,
		RepoPath: *repo,
	}

	return config
}
