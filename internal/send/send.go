package sendImplementation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type TransferSubmissionInformationPackage struct {
	Date            string `json:"date"`
	ContentCategory string `json:"contentCategory"`
	ContentType     string `json:"contentType"`
	Identifier      string `json:"identifier"`
	Urn             string `json:"urn"`
	Path            string `json:"path"`
}

type sender struct {
	writer *kafka.Writer
}

func PrepareAndSendSubmissionInformationPackage(kafkaEndpoints []string, transferTopicName string, rootPath string) error {
	sender := sender{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(kafkaEndpoints...),
			Topic:    transferTopicName,
			Balancer: &kafka.LeastBytes{},
		},
	}

	defer sender.writer.Close()

	items, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("failed to read root path '%s', original error: '%w'", rootPath, err)
	}
	for _, path := range items {
		if path.IsDir() {
			directoryName := path.Name()
			destinationPath := rootPath + directoryName
			fmt.Printf("Processing directory %s\n", destinationPath)
			transferSubmissionInformationPackage := createSubmissionInformationPackage(destinationPath, directoryName)

			kafkaMessage, err := json.Marshal(transferSubmissionInformationPackage)
			if err != nil {
				return fmt.Errorf("failed to marshal json, original error: '%w'", err)
			}

			err = sender.sendMessageToKafkaTopic(kafkaMessage)
			if err != nil {
				return fmt.Errorf("failed to send message to kafka topic '%s', original error: '%w'", transferTopicName, err)
			}
		} else {
			return fmt.Errorf("found file '%s' in root path '%s', but expected only directories", path.Name(), rootPath)
		}
	}

	return nil
}

func createSubmissionInformationPackage(payloadPath string, payloadDirName string) TransferSubmissionInformationPackage {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.000")
	contentCategory := "nettarkiv"
	contentType := "warc"
	identifier := "no-nb_" + contentCategory + "_" + payloadDirName
	urn := "URN:NBN:" + identifier

	return TransferSubmissionInformationPackage{
		Date:            date,
		ContentCategory: contentCategory,
		ContentType:     contentType,
		Identifier:      identifier,
		Urn:             urn,
		Path:            payloadPath,
	}
}

func (sender *sender) sendMessageToKafkaTopic(payload []byte) error {
	kafkaMessageUuid := uuid.New()
	fmt.Printf("Sending message with uuid %s\n", kafkaMessageUuid)
	kafkaMessageUuidBytes, err := kafkaMessageUuid.MarshalText()
	if err != nil {
		return fmt.Errorf("failed to marshal uuid, original error: '%w'", err)
	}

	err = sender.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   kafkaMessageUuidBytes,
			Value: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write messages, original error: '%w'", err)
	}

	return nil
}
