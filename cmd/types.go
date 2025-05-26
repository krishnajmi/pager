package cmd

import "github.com/kp/pager/databases/sql"

type AWSConfig struct {
	AccessKey string `json:"AWS_ACCESS_KEY_ID"`
	SecretKey string `json:"AWS_SECRET_ACCESS_KEY"`
	Region    string `json:"AWS_REGION"`
}

type RedisConfig struct {
	Host     string `json:"REDIS_HOST"`
	Port     string `json:"REDIS_PORT"`
	Password string `json:"REDIS_PASSWORD"`
	DB       string `json:"REDIS_DB"`
}

type KafkaConfig struct {
	Brokers  string `json:"KAFKA_BROKERS"`
	Topic    string `json:"KAFKA_TOPIC"`
	Username string `json:"KAFKA_USERNAME"`
	Password string `json:"KAFKA_PASSWORD"`
}

type AppConfig struct {
	AWSConfig
	sql.DatabaseConfigType
	RedisConfig
	KafkaConfig
}
