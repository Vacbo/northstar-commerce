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

	// Bound delivery timeout to 2 seconds so broker unavailability does not block indefinitely
	saramaConfig.Producer.Timeout = 2 * time.Second

	// Set channel buffer sizes to provide backpressure when broker is slow
	saramaConfig.ChannelBufferSize = 256

	// Sarama has an issue in a single broker kafka if the kafka broker is restarted.
	// This setting is to prevent that issue from manifesting itself, but may swallow failed messages.
	saramaConfig.Producer.RequiredAcks = sarama.NoResponse

	// Set metadata refresh frequency to detect broker availability changes
	saramaConfig.Metadata.RefreshFrequency = 10 * time.Second

	saramaConfig.Version = ProtocolVersion

	// So we can know the partition and offset of messages.
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	// We will log to STDOUT if we're not able to produce messages.
	go func() {
		for err := range producer.Errors() {
			logger.Error(fmt.Sprintf("Failed to write message: %+v", err))

		}
	}()
	return producer, nil
}
