package send

import (
	"context"
	"encoding/json"
	"fmt"
	"hermetic/internal/common"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

const (
	transferTopicFlagName      string = "transfer-topic"
	stageArtifactsRootFlagName string = "stage-artifacts-root"
)

type TransferSubmissionInformationPackage struct {
	Date            string `json:"date"`
	ContentCategory string `json:"contentCategory"`
	ContentType     string `json:"contentType"`
	Identifier      string `json:"identifier"`
	Urn             string `json:"urn"`
	Path            string `json:"path"`
}

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "send",
		Short: "Continuously sends data to digital storage",
		Args:  cobra.NoArgs,
		RunE:  parseArgumentsAndCallSend,
	}
	command.Flags().String(transferTopicFlagName, "", "name of transfer-topic")
	if markTransferRequiredError := command.MarkFlagRequired(transferTopicFlagName); markTransferRequiredError != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", transferTopicFlagName, markTransferRequiredError))
	}

	command.Flags().String(stageArtifactsRootFlagName, "", `path to the root directory with the following content structure:
/stage-artifacts-root
├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm
│   ├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
│   └── /checksum_transferred.md5
└── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm
    ├── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
    └── /checksum_transferred.md5
... etc`)
	if markRootPathRequiredError := command.MarkFlagRequired(stageArtifactsRootFlagName); markRootPathRequiredError != nil {
		panic(fmt.Sprintf("failed to mark flag %s as required, original error: '%s'", stageArtifactsRootFlagName, markRootPathRequiredError))
	}

	return command
}

func parseArgumentsAndCallSend(cmd *cobra.Command, args []string) error {
	transferTopicName, transferTopicError := cmd.Flags().GetString(transferTopicFlagName)
	if transferTopicError != nil {
		return fmt.Errorf("getting transfer topic name failed, original error: '%w'", transferTopicError)
	}
	kafkaEndpoints, parseGlobalFlagsError := common.ParseKafkaEndpointsFlags(cmd)
	if parseGlobalFlagsError != nil {
		return fmt.Errorf("getting kafka endpoints failed, original error: '%w'", parseGlobalFlagsError)
	}
	stageArtifactsRoot, stageArtifactsRootError := cmd.Flags().GetString("stage-artifacts-root")
	if stageArtifactsRootError != nil {
		return fmt.Errorf("getting stage artifacts root failed, original error: '%w'", stageArtifactsRootError)
	}

	return send(kafkaEndpoints, transferTopicName, stageArtifactsRoot)
}

func createSubmissionInformationPackage(payloadPath string, payloadDirName string) TransferSubmissionInformationPackage {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.000")
	contentCategory := "nettarkiv"
	contentType := "warc"
	identifier := "no-nb_" + contentCategory + "_" + payloadDirName
	urn := "URN:NBN:" + identifier
	destinationPath := payloadPath

	submissionInformationPackage := TransferSubmissionInformationPackage{date, contentCategory, contentType, identifier, urn, destinationPath}
	return submissionInformationPackage
}

func send(kafkaEndpoints []string, transferTopicName string, rootPath string) error {
	items, readRootError := os.ReadDir(rootPath)
	if readRootError != nil {
		return fmt.Errorf("failed to read root path '%s', original error: '%w'", rootPath, readRootError)
	}
	for _, path := range items {
		if path.IsDir() {
			directoryName := path.Name()
			destinationPath := rootPath + directoryName
			fmt.Printf("Processing directory %s\n", destinationPath)
			transferSubmissionInformationPackage := createSubmissionInformationPackage(destinationPath, directoryName)

			kafkaMessage, marshalToJsonError := json.Marshal(transferSubmissionInformationPackage)
			if marshalToJsonError != nil {
				return fmt.Errorf("failed to marshal json, original error: '%w'", marshalToJsonError)
			}

			sendMessageError := sendMessageToKafkaTopic(kafkaEndpoints, kafkaMessage, transferTopicName)
			if sendMessageError != nil {
				return fmt.Errorf("failed to send message to kafka topic '%s', original error: '%w'", transferTopicName, sendMessageError)
			}
		} else {
			return fmt.Errorf("found file '%s' in root path '%s', but expected only directories", path.Name(), rootPath)
		}
	}

	return nil
}

func sendMessageToKafkaTopic(kafkaEndpoints []string, payload []byte, transferTopicName string) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaEndpoints...),
		Topic:    transferTopicName,
		Balancer: &kafka.LeastBytes{},
	}

	kafkaMessageUuid := uuid.New()
	fmt.Printf("Sending message with uuid %s\n", kafkaMessageUuid)

	writeMessageError := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   kafkaMessageUuid[:],
			Value: payload,
		},
	)
	if writeMessageError != nil {
		return fmt.Errorf("failed to write messages, original error: '%w'", writeMessageError)
	}

	if closeWriterError := writer.Close(); closeWriterError != nil {
		return fmt.Errorf("failed to close writer, original error: '%w'", closeWriterError)
	}

	return nil
}
