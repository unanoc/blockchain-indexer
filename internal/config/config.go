package config

import (
	"time"
)

// Default is a config instance.
var Default Config //nolint:gochecknoglobals // config must be global

type Config struct {
	LogLevel string `mapstructure:"log_level"`

	Port string `mapstructure:"port"`

	Gin struct {
		Mode string `mapstructure:"mode"`
	} `mapstructure:"gin"`

	Swagger struct {
		Hostname string `mapstructure:"hostname"`
	} `mapstructure:"swagger"`

	Sentry struct {
		DSN        string  `mapstructure:"dsn"`
		SampleRate float32 `mapstructure:"sample_rate"`
	} `mapstructure:"sentry"`

	Database struct {
		URL string `mapstructure:"url"`
		Log bool   `mapstructure:"log"`
	} `mapstructure:"database"`

	RabbitMQ struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"rabbitmq"`

	Kafka struct {
		Brokers           string        `mapstructure:"brokers"`
		BlocksTopicPrefix string        `mapstructure:"blocks_topic_prefix"`
		MaxAttempts       int           `mapstructure:"max_attempts"`
		MessageMaxBytes   int           `mapstructure:"message_max_bytes"`
		RetentionTime     time.Duration `mapstructure:"retention_time"`
	} `mapstructure:"kafka"`

	Prometheus struct {
		NameSpace string `mapstructure:"namespace"`
		SubSystem string `mapstructure:"subsystem"`

		PushGateway struct {
			URL          string        `mapstructure:"url"`
			Key          string        `mapstructure:"key"`
			PushInterval time.Duration `mapstructure:"push_interval"`
		} `mapstructure:"pushgateway"`
	} `mapstructure:"prometheus"`

	BlockProducer struct {
		Interval           time.Duration `mapstructure:"interval"`
		BackoffInterval    time.Duration `mapstructure:"backoff_interval"`
		FetchBlocksMax     int64         `mapstructure:"fetch_blocks_max"`
		StartFromLastBlock bool          `mapstructure:"start_from_last_block"`
		BlockRetryNum      int           `mapstructure:"block_retry"`
		BlockRetryInterval time.Duration `mapstructure:"block_retry_interval"`
	} `mapstructure:"block_producer"`

	BlockConsumer struct {
		Interval time.Duration `mapstructure:"interval"`
	} `mapstructure:"block_consumer"`

	Nodes struct {
		Interval  time.Duration `mapstructure:"interval"`
		InitNodes bool          `mapstructure:"init_nodes"`
	} `mapstructure:"nodes"`

	TransactionConsumer struct {
		Workers int `mapstructure:"workers"`
	} `mapstructure:"transaction_consumer"`

	Platforms struct {
		Binance struct {
			Node string `mapstructure:"node"`
			Dex  string `mapstructure:"dex"`
		} `mapstructure:"binance"`
		Bitcoin struct {
			Node string `mapstructure:"node"`
		} `mapstructure:"bitcoin"`
		Cosmos struct {
			Node string `mapstructure:"node"`
		} `mapstructure:"cosmos"`
		Ethereum struct {
			Node string `mapstructure:"node"`
		} `mapstructure:"ethereum"`
		Smartchain struct {
			Node string `mapstructure:"node"`
		} `mapstructure:"smartchain"`
	} `mapstructure:"platforms"`
}
