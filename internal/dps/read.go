package dps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	ContentTypeWarc        = "warc"
	ContentTypeAcquisition = "acquisition"
)

func getKafkaPartitionLeader(kafkaEndpoints []string, kafkaTopic string) (*kafka.Conn, error) {
	if len(kafkaEndpoints) == 0 {
		return nil, errors.New("no kafka endpoints provided")
	}

	conn, err := net.Dial("tcp", kafkaEndpoints[0])
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	kafkaConn := kafka.NewConn(conn, kafkaTopic, 0)
	defer kafkaConn.Close()

	partitions, err := kafkaConn.ReadPartitions(kafkaTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	if len(partitions) != 1 {
		return nil, fmt.Errorf("expected exactly 1 partition, got '%d'", len(partitions))
	}

	connLeader, err := net.Dial("tcp", fmt.Sprintf("%s:%d", partitions[0].Leader.Host, partitions[0].Leader.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to dial tcp: %w", err)
	}

	return kafka.NewConn(connLeader, kafkaTopic, 0), nil
}

func ReadLatestMessages(ctx context.Context, kafkaEndpoints []string, kafkaTopic string, fn func(*Message) error) error {
	conn, err := getKafkaPartitionLeader(kafkaEndpoints, kafkaTopic)
	if err != nil {
		return fmt.Errorf("failed to get kafka partition leader: %w", err)
	}
	firstOffset, lastOffset, err := conn.ReadOffsets()
	if err != nil {
		return fmt.Errorf("failed to get first and last offset: %w", err)
	}
	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection to leader: %w", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaEndpoints,
		Topic:   kafkaTopic,
	})
	defer reader.Close()

	readTimeout := 5 * time.Minute
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	err = reader.SetOffset(firstOffset)
	if err != nil {
		return fmt.Errorf("failed to set kafka reader offset '%d': %w", firstOffset, err)
	}
	for reader.Offset() < lastOffset {
		kafkaMsg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if kafkaMsg.Value == nil {
			continue
		}
		var msg Message
		err = json.Unmarshal(kafkaMsg.Value, &msg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal kafka message: %w", err)
		}
		err = fn(&msg)
		if err != nil {
			return fmt.Errorf("failed to process message: %w", err)
		}
	}
	return nil
}

func NextMessage(ctx context.Context, reader *kafka.Reader, filter func(*Message) bool) (*KafkaMessage, error) {
	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		var response Message

		err = json.Unmarshal(message.Value, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal kafka message: %w", err)
		}

		if !filter(&response) {
			continue
		}

		return &KafkaMessage{
			Offset: message.Offset,
			Key:    string(message.Key),
			Value:  response,
		}, nil
	}
}
