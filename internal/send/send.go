package sendImplementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	kafkaHelpers "hermetic/internal/kafka"
	"hermetic/internal/submission_information_package"

	"github.com/allegro/bigcache/v3"
	"github.com/segmentio/kafka-go"
)

type offsets struct {
	first int64
	last  int64
}

const (
	warcContentType        = "warc"
	acquisitionContentType = "acquisition"
)

func getFirstAndLastOffsets(kafkaEndpoints []string, transferTopicName string) (offsets, error) {
	conn, err := net.Dial("tcp", kafkaEndpoints[0])
	if err != nil {
		return offsets{}, fmt.Errorf("failed to dial tcp, original error: '%w'", err)
	}
	kafkaConn := kafka.NewConn(conn, transferTopicName, 0)
	partitions, err := kafkaConn.ReadPartitions(transferTopicName)
	if err != nil {
		return offsets{}, fmt.Errorf("failed to read partitions, original error: '%w'", err)
	}
	if len(partitions) != 1 {
		return offsets{}, fmt.Errorf("expected exactly 1 partition, got '%d'", len(partitions))
	}
	connLeader, err := net.Dial("tcp", fmt.Sprintf("%s:%d", partitions[0].Leader.Host, partitions[0].Leader.Port))
	if err != nil {
		return offsets{}, fmt.Errorf("failed to dial tcp, original error: '%w'", err)
	}
	kafkaLeaderConn := kafka.NewConn(connLeader, transferTopicName, 0)

	firstOffset, lastOffset, err := kafkaLeaderConn.ReadOffsets()
	if err != nil {
		return offsets{}, fmt.Errorf("failed to read offsets, original error: '%w'", err)
	}
	return offsets{first: firstOffset, last: lastOffset}, nil
}

func readLatestMessages(kafkaEndpoints []string, transferTopicName string) ([]submission_information_package.Package, error) {
	messageReader := kafkaHelpers.MessageReader{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: kafkaEndpoints,
			Topic:   transferTopicName,
		}),
	}

	defer messageReader.Reader.Close()

	offsets, err := getFirstAndLastOffsets(kafkaEndpoints, transferTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get first and last offset, original error: '%w'", err)
	}
	readTimeout := 10 * time.Second

	var messages []submission_information_package.Package

	for offsetToReadFrom := offsets.first; offsetToReadFrom < offsets.last; offsetToReadFrom++ {
		fmt.Printf("Reading message at offset '%d'\n", offsetToReadFrom)
		err := messageReader.Reader.SetOffset(offsetToReadFrom)
		if err != nil {
			return nil, fmt.Errorf("failed to set offset '%d', original error: '%w'", offsetToReadFrom, err)
		}

		message, err := messageReader.ReadMessageWithTimeout(readTimeout)
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Could not read message at offset '%d', read timeout '%s' exceeded, skipping offset\n", offsetToReadFrom, readTimeout)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read message at offset '%d', original error: '%w'", offsetToReadFrom, err)
		}
		if message.Value == nil {
			fmt.Printf("Message at offset '%d' is nil, skipping offset\n", offsetToReadFrom)
			continue
		}
		var submissionInformationPackage submission_information_package.Package

		err = json.Unmarshal(message.Value, &submissionInformationPackage)
		if err != nil {
			syntaxError := new(json.SyntaxError)
			if errors.As(err, &syntaxError) {
				fmt.Printf("Could not read message at offset '%d', syntax error in message, skipping offset\n", offsetToReadFrom)
				continue
			}
			return nil, fmt.Errorf("failed to unmarshal json, original error: '%w'", err)
		}
		messages = append(messages, submissionInformationPackage)

	}

	return messages, nil
}

func webArchiveRelevantMessages(messages []submission_information_package.Package) ([]submission_information_package.Package, error) {
	var relevantMessages []submission_information_package.Package
	for _, message := range messages {
		if message.ContentCategory == "nettarkiv" {
			switch message.ContentType {
			case warcContentType:
			case acquisitionContentType:
			default:
				return nil, fmt.Errorf("found content type '%s' in message '%+v', expected '%s' or '%s'", message.ContentType, message, warcContentType, acquisitionContentType)
			}
			relevantMessages = append(relevantMessages, message)
		}
	}
	return relevantMessages, nil
}

func PrepareAndSendSubmissionInformationPackage(kafkaEndpoints []string, transferTopicName string, rootPath string) error {
	sender := kafkaHelpers.Sender{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(kafkaEndpoints...),
			Topic:    transferTopicName,
			Balancer: &kafka.LeastBytes{},
		},
	}

	defer sender.Writer.Close()

	latestMessages, err := readLatestMessages(kafkaEndpoints, transferTopicName)
	if err != nil {
		return fmt.Errorf("failed to read latest messages, original error: '%w'", err)
	}

	relevantMessages, err := webArchiveRelevantMessages(latestMessages)
	if err != nil {
		return fmt.Errorf("failed to filter out relevant messages, original error: '%w'", err)
	}
	config := bigcache.Config{
		Shards:      1024,
		LifeWindow:  24 * 7 * time.Hour,
		CleanWindow: 1 * time.Hour,
	}
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create cache, original error: '%w'", err)
	}

	for _, message := range relevantMessages {
		fmt.Printf("Pushing '%s' to cache\n", message.Path)
		err := cache.Set(message.Path, []byte("Sent"))
		if err != nil {
			return fmt.Errorf("failed to set '%s' in cache, original error: '%w'", message.Path, err)
		}
	}

	for {
		items, err := os.ReadDir(rootPath)
		if err != nil {
			return fmt.Errorf("failed to read root path '%s', original error: '%w'", rootPath, err)
		}
		for _, path := range items {
			directoryName := path.Name()
			destinationPath := rootPath + directoryName
			_, err := cache.Get(destinationPath)
			if err == nil {
				fmt.Printf("Skipping directory '%s' as it has already been processed.\n", destinationPath)
				continue
			} else {
				if !errors.Is(err, bigcache.ErrEntryNotFound) {
					return fmt.Errorf("failed to get '%s' from cache, original error: '%w'", destinationPath, err)
				}
			}
			if !path.IsDir() {
				return fmt.Errorf("found file '%s' in root path '%s', but expected only directories", path.Name(), rootPath)
			}
			fmt.Printf("Processing directory %s\n", destinationPath)
			submissionInformationPackage := submission_information_package.CreatePackage(destinationPath, directoryName, warcContentType)

			kafkaMessage, err := json.Marshal(submissionInformationPackage)
			if err != nil {
				return fmt.Errorf("failed to marshal json, original error: '%w'", err)
			}

			err = sender.SendMessageToKafkaTopic(kafkaMessage)
			if err != nil {
				return fmt.Errorf("failed to send message to kafka topic '%s', original error: '%w'", transferTopicName, err)
			}
			err = cache.Set(destinationPath, []byte("Sent"))
			if err != nil {
				return fmt.Errorf("failed to set '%s' in cache, original error: '%w'", destinationPath, err)
			}

		}
		time.Sleep(1 * time.Minute)
	}
}
