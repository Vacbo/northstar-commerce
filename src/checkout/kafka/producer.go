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

	// Bound broker metadata fetch timeout to prevent indefinite blocking
	saramaConfig.Metadata.Timeout = 3 * time.Second
	saramaConfig.Metadata.Retry.Max = 2
	saramaConfig.Metadata.Retry.Backoff = 100 * time.Millisecond

	// Bound network timeouts for producer operations
	saramaConfig.Net.DialTimeout = 3 * time.Second
	saramaConfig.Net.ReadTimeout = 5 * time.Second
	saramaConfig.Net.WriteTimeout = 5 * time.Second

	// Limit in-flight messages and configure backpressure
	saramaConfig.Producer.MaxMessageBytes = 1000000
	saramaConfig.ChannelBufferSize = 256

	// Configure delivery guarantees with bounded retry
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond

	// WaitForLocal provides at-least-once delivery while bounding broker acknowledgment latency
	// This maintains delivery guarantees while preventing indefinite blocking during broker degradation
	saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal

	saramaConfig.Version = ProtocolVersion

	producer, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	// Log delivery errors asynchronously without blocking caller
	go func() {
		for err := range producer.Errors() {
			logger.Error(fmt.Sprintf("Failed to write message: %+v", err))
		}
	}()

	// Consume successes asynchronously to prevent channel blocking
	go func() {
		for range producer.Successes() {
			// Drain success channel to prevent backpressure
		}
	}()

	return producer, nil
}
