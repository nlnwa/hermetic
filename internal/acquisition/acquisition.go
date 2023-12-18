package acquisitionImplementation

import (
	"encoding/json"
	"fmt"
	kafkaHelpers "hermetic/internal/kafka"
	"hermetic/internal/path"
	"hermetic/internal/submission_information_package"
	"os"
	"path/filepath"
	"time"

	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

const (
	contentType                 = "acquisition"
	supportedAcquisitionVersion = "0.1.0"
)

type DataModel struct {
	AcquisitionVersion string `yaml:"__acquisition_version__"`
	Acquisition        struct {
		Name                 string `yaml:"name"`
		Date                 string `yaml:"date"`
		OriginalPurpose      string `yaml:"original-purpose"`
		AcquisitionPurpose   string `yaml:"acquisition-purpose"`
		AccessConsiderations string `yaml:"access-considerations"`
	} `yaml:"acquisition"`

	Files []struct {
		Name        string `yaml:"name"`
		Format      string `yaml:"format"`
		Path        string `yaml:"path"`
		Description string `yaml:"description"`
	} `yaml:"files"`
	AcquisitionHandling struct {
		Responsible string `yaml:"responsible"`
		Author      string `yaml:"author"`
	} `yaml:"acquisition-handling"`
}

func PrepareAndSendSubmissionInformationPackage(kafkaEndpoints []string, transferTopicName string, acquisitionRoot string) error {
	sender := kafkaHelpers.Sender{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(kafkaEndpoints...),
			Topic:    transferTopicName,
			Balancer: &kafka.LeastBytes{},
		},
	}
	defer sender.Writer.Close()
	isDir, err := path.IsDirectory(acquisitionRoot)
	if err != nil {
		return fmt.Errorf("failed to check if '%s' is a directory, original error: '%w'", acquisitionRoot, err)
	}
	if !isDir {
		return fmt.Errorf("acquisitionRoot ('%s') is not a directory", acquisitionRoot)
	}

	dataModel, err := deserializeYamlFile(filepath.Join(acquisitionRoot, "acquisition.yaml"))
	if err != nil {
		return fmt.Errorf("failed to deserialize yaml file, original error: '%w'", err)
	}

	if err := validate(acquisitionRoot, dataModel); err != nil {
		return fmt.Errorf("failed to process yaml file, original error: '%w'", err)
	}

	identifier := dataModel.Acquisition.Name + "-" + dataModel.Acquisition.Date

	submissionInformationPackage := submission_information_package.CreatePackage(acquisitionRoot, identifier, contentType)

	expectedURN := "URN:NBN:no-nb_nettarkiv_" + dataModel.Acquisition.Name + "-" + dataModel.Acquisition.Date

	if submissionInformationPackage.Urn != expectedURN {
		return fmt.Errorf("failed to create URN, expected %s, got %s", expectedURN, submissionInformationPackage.Urn)
	}

	kafkaMessage, err := json.Marshal(submissionInformationPackage)
	if err != nil {
		return fmt.Errorf("failed to marshal json, original error: '%w'", err)
	}

	err = sender.SendMessageToKafkaTopic(kafkaMessage)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka topic '%s', original error: '%w'", transferTopicName, err)
	}

	return nil
}

func deserializeYamlFile(metadataFilePath string) (DataModel, error) {
	var dataModel DataModel

	content, err := os.ReadFile(metadataFilePath)
	if err != nil {
		return dataModel, fmt.Errorf("failed to read file '%s', original error: '%w'", metadataFilePath, err)
	}

	err = yaml.Unmarshal(content, &dataModel)
	if err != nil {
		return dataModel, fmt.Errorf("failed to unmarshal yaml, original error: '%w'", err)
	}
	return dataModel, nil
}

func validate(rootPath string, dataModel DataModel) error {
	// TODO(https://github.com/nlnwa/hermetic/issues/32): This YAML file should
	// be replaced by more sustainable solution based on industry standards,
	// such as METS https://www.loc.gov/standards/mets/mets-schemadocs.html and
	// PREMIS https://www.loc.gov/standards/premis/

	if dataModel.AcquisitionVersion != supportedAcquisitionVersion {
		return fmt.Errorf("acquisition version '%s' is not supported, expected %s", dataModel.AcquisitionVersion, supportedAcquisitionVersion)
	}

	if err := directoryValidation(rootPath); err != nil {
		return fmt.Errorf("failed to validate directories, original error: '%w'", err)
	}

	if err := fileValidation(rootPath, dataModel); err != nil {
		return fmt.Errorf("failed to validate files, original error: '%w'", err)
	}

	if err := otherMetadataValidation(dataModel); err != nil {
		return fmt.Errorf("failed to validate other metadata, original error: '%w'", err)
	}

	return nil
}

func directoryValidation(rootPath string) error {
	items, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("failed to read root path '%s', original error: '%w'", rootPath, err)
	}

	for _, path := range items {
		if path.IsDir() {
			return fmt.Errorf("found directory '%s' in root path '%s', but expected only files", path.Name(), rootPath)
		}
	}

	return nil
}

func fileValidation(rootPath string, metadata DataModel) error {
	for _, file := range metadata.Files {
		resolvedPath := filepath.Join(rootPath, file.Path)
		isFile, err := path.IsFile(resolvedPath)
		if err != nil {
			return fmt.Errorf("failed to check if '%s' is a file, original error: '%w'", resolvedPath, err)
		}
		if !isFile {
			return fmt.Errorf("file '%s' is not a file", resolvedPath)
		}
	}

	items, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("failed to read root path '%s', original error: '%w'", rootPath, err)
	}

	for _, path := range items {
		found := false
		if !path.IsDir() {
			for _, yamlFile := range metadata.Files {
				if yamlFile.Path == path.Name() {
					found = true
					continue
				}
			}
			if !found {
				return fmt.Errorf("found file '%s' in root path '%s', but expected only files specified in yaml file", path.Name(), rootPath)
			}

		}
	}

	return nil
}

func otherMetadataValidation(metadata DataModel) error {
	test := map[string]string{}
	test["metadata.Acquisition.Name"] = metadata.Acquisition.Name
	test["metadata.Acquisition.Date"] = metadata.Acquisition.Date
	test["metadata.Acquisition.OriginalPurpose"] = metadata.Acquisition.OriginalPurpose
	test["metadata.Acquisition.AcquisitionPurpose"] = metadata.Acquisition.AcquisitionPurpose
	test["metadata.Acquisition.AccessConsiderations"] = metadata.Acquisition.AccessConsiderations
	test["metadata.AcquisitionHandling.Responsible"] = metadata.AcquisitionHandling.Responsible
	test["metadata.AcquisitionHandling.Author"] = metadata.AcquisitionHandling.Author

	for fieldName, value := range test {
		if value == "" {
			return fmt.Errorf("field '%s' is empty", fieldName)
		}
	}
	_, err := time.Parse(time.RFC3339, metadata.Acquisition.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date '%s', original error: '%w'", metadata.Acquisition.Date, err)
	}

	return nil
}
