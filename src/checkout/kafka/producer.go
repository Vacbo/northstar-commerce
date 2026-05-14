// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package kafka

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
)

var (
	Topic           = "orders"
	ProtocolVersion = sarama.V3_0_0_0
)

type saramaLogger struct {
	logger *slog.Logger
}

func (l *saramaLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}
func (l *saramaLogger) Println(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}
func (l *saramaLogger) Print(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}

func CreateKafkaProducer(brokers []string, logger *slog.Logger) (sarama.AsyncProducer, error) {
	// Set the logger for sarama to use.
	sarama.Logger = &saramaLogger{logger: logger}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	// Bound async publish path with explicit timeout and backpressure
	saramaConfig.Producer.Timeout = 2 * time.Second
	saramaConfig.Producer.MaxMessageBytes = 1000000
	saramaConfig.ChannelBufferSize = 256

	// Balance throughput and latency with sensible producer properties
	saramaConfig.Producer.Flush.Frequency = 100 * time.Millisecond
	saramaConfig.Producer.Flush.Messages = 100
	saramaConfig.Producer.Flush.MaxMessages = 1000

	// Retry configuration for transient failures
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond

	// Sarama has an issue in a single broker kafka if the kafka broker is restarted.
	// This setting is to prevent that issue from manifesting itself, but may swallow failed messages.
	saramaConfig.Producer.RequiredAcks = sarama.NoResponse

	saramaConfig.Version = ProtocolVersion

	// So we can know the partition and offset of messages.
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	// We will log to STDOUT if we're not able to produce messages.
	// Handle delivery errors with bounded error processing
	go func() {
		for err := range producer.Errors() {
			logger.Error(fmt.Sprintf("Failed to write message: %+v", err))
		}
	}()

	// Drain successes to prevent channel blocking
	go func() {
		for range producer.Successes() {
			// Success messages drained to prevent backpressure
		}
	}()

	return producer, nil
}
