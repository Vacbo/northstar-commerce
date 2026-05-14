// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package kafka

import (
	"context"
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
	go func() {
		for err := range producer.Errors() {
			logger.Error(fmt.Sprintf("Failed to write message: %+v", err))

		}
	}()
	return producer, nil
}

// SendWithTimeout sends a message to Kafka with a timeout guarantee.
// This function wraps the blocking select on AsyncProducer channels
// in a goroutine to ensure non-blocking behavior within the timeout.
func SendWithTimeout(ctx context.Context, producer sarama.AsyncProducer, msg *sarama.ProducerMessage, logger *slog.Logger) error {
	resultChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		select {
		case producer.Input() <- msg:
			// Message queued, now wait for result
			select {
			case <-producer.Successes():
				resultChan <- nil
			case errMsg := <-producer.Errors():
				if errMsg != nil {
					resultChan <- fmt.Errorf("kafka producer error: %v", errMsg.Err)
				} else {
					resultChan <- fmt.Errorf("kafka producer unknown error")
				}
			}
		case <-ctx.Done():
			resultChan <- fmt.Errorf("failed to send message to Kafka within timeout: %v", ctx.Err())
		}
	}()

	// Use time.After to guarantee non-blocking timeout
	timeout := 5 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	select {
	case err := <-resultChan:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("kafka send operation timed out after %v", timeout)
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %v", ctx.Err())
	}
}
