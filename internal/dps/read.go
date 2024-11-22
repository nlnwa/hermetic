package dps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

func ReadLatestMessages(ctx context.Context, kafkaEndpoints []string, kafkaTopic string, fn func(*Package) error) error {
	conn, err := getKafkaPartitionLeader(kafkaEndpoints, kafkaTopic)
	if err != nil {
		return fmt.Errorf("failed to get kafka partition leader: %w", err)
	}
	first, last, err := conn.ReadOffsets()
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

	for offset := first; offset < last; offset++ {
		err := reader.SetOffset(offset)
		if err != nil {
			return fmt.Errorf("failed to set kafka reader offset '%d': %w", offset, err)
		}

		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if message.Value == nil {
			continue
		}

		var pkg Package

		err = json.Unmarshal(message.Value, &pkg)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				slog.Warn("Could not read message at offset, syntax error in message, skipping offset", "offset", offset)
				continue
			}
			return fmt.Errorf("failed to unmarshal: %w", err)
		}
		err = fn(&pkg)
		if err != nil {
			return fmt.Errorf("failed to process message: %w", err)
		}
	}
	return nil
}

func IsWebArchiveMessage(message *Package) bool {
	if message.ContentCategory == "nettarkiv" {
		switch message.ContentType {
		case ContentTypeWarc, ContentTypeAcquisition:
			return true
		default:
			return false
		}
	}
	return false
}

func NextMessage(ctx context.Context, reader *kafka.Reader, filter func(*Response) bool) (*KafkaResponse, error) {
	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		var response Response

		err = json.Unmarshal(message.Value, &response)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				slog.Warn("Could not read message, skipping...", "offset", message.Offset, "value", string(message.Value), "error", err)
				continue
			}
			return nil, fmt.Errorf("failed to unmarshal json: %w", err)
		}

		if !filter(&response) {
			continue
		}

		return &KafkaResponse{
			Offset:   message.Offset,
			Key:      string(message.Key),
			Response: response,
		}, nil
	}
}
