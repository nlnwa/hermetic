package acquisition

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nlnwa/hermetic/cmd/internal/flags"
	"github.com/nlnwa/hermetic/internal/dps"
	"github.com/nlnwa/hermetic/internal/path"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	dirFlagName     string = "dir"
	dirFlagNameHelp string = `path to the root directory with the following content structure:
/acquisition-root
├── checksums.md5
├── checksum_transferred.md5
├── /acquisition.yaml
├── <other-small-and-few-files>
└── /<other-files-and-directories>
`
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().String(dirFlagName, "", dirFlagNameHelp)
	if err := cmd.MarkFlagRequired(dirFlagName); err != nil {
		panic(err)
	}
}

func toOptions() AcquisitionOptions {
	return AcquisitionOptions{
		KafkaEndpoints: flags.GetKafkaEndpoints(),
		KafkaTopic:     flags.GetKafkaTopic(),
		Dir:            viper.GetString(dirFlagName),
	}
}

type AcquisitionOptions struct {
	KafkaEndpoints []string
	KafkaTopic     string
	Dir            string
}

func (o AcquisitionOptions) Run() error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(o.KafkaEndpoints...),
		Topic:    o.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	isDir, err := path.IsDirectory(o.Dir)
	if err != nil {
		return fmt.Errorf("failed to check if '%s' is a directory: %w", o.Dir, err)
	}
	if !isDir {
		return fmt.Errorf("'%s' is not a directory", o.Dir)
	}

	dataModel, err := deserializeYamlFile(filepath.Join(o.Dir, "acquisition.yaml"))
	if err != nil {
		return fmt.Errorf("failed to deserialize yaml file: %w", err)
	}

	if err := validate(o.Dir, dataModel); err != nil {
		return fmt.Errorf("failed to process yaml file: %w", err)
	}

	identifier := dataModel.ArchiveUnit.Name + "-" + dataModel.ArchiveUnit.Deposit.Date

	parcel := dps.CreatePackage(o.Dir, identifier, contentType)

	expectedURN := "URN:NBN:no-nb_nettarkiv_" + dataModel.ArchiveUnit.Name + "-" + dataModel.ArchiveUnit.Deposit.Date

	if parcel.Urn != expectedURN {
		return fmt.Errorf("failed to create URN, expected %s, got %s", expectedURN, parcel.Urn)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = dps.Send(ctx, writer, parcel)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka topic '%s': %w", o.KafkaTopic, err)
	}

	return nil
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acquisition",
		Short: "Uploads data to digital storage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return toOptions().Run()
		},
	}

	addFlags(cmd)

	return cmd
}

func deserializeYamlFile(metadataFilePath string) (DataModel, error) {
	var dataModel DataModel

	content, err := os.ReadFile(metadataFilePath)
	if err != nil {
		return dataModel, fmt.Errorf("failed to read file '%s': %w", metadataFilePath, err)
	}

	err = yaml.Unmarshal(content, &dataModel)
	if err != nil {
		return dataModel, fmt.Errorf("failed to unmarshal yaml: %w", err)
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
		return fmt.Errorf("failed to validate directories: %w", err)
	}

	if err := fileValidation(rootPath, dataModel); err != nil {
		return fmt.Errorf("failed to validate files: %w", err)
	}

	if err := otherMetadataValidation(dataModel); err != nil {
		return fmt.Errorf("failed to validate other metadata: %w", err)
	}

	return nil
}

func directoryValidation(rootPath string) error {
	items, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("failed to read root path '%s': %w", rootPath, err)
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
			return fmt.Errorf("failed to check if '%s' is a file: %w", resolvedPath, err)
		}
		if !isFile {
			return fmt.Errorf("file '%s' is not a file", resolvedPath)
		}
	}

	items, err := os.ReadDir(rootPath)
	if err != nil {
		return fmt.Errorf("failed to read root path '%s': %w", rootPath, err)
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
	test["metadata.ArchiveUnit.Name"] = metadata.ArchiveUnit.Name
	test["metadata.ArchiveUnit.Type"] = metadata.ArchiveUnit.Type
	test["metadata.ArchiveUnit.Creator"] = metadata.ArchiveUnit.Creator
	test["metadata.ArchiveUnit.Description"] = metadata.ArchiveUnit.Description
	test["metadata.ArchiveUnit.CopyrightClearance"] = metadata.ArchiveUnit.CopyrightClearance
	test["metadata.ArchiveUnit.AccessConsiderations"] = metadata.ArchiveUnit.AccessConsiderations
	test["metadata.ArchiveUnit.Deposit.Depositor"] = metadata.ArchiveUnit.Deposit.Depositor
	test["metadata.ArchiveUnit.Deposit.Date"] = metadata.ArchiveUnit.Deposit.Date
	test["metadata.ArchiveUnit.Deposit.AcquisitionPurpose"] = metadata.ArchiveUnit.Deposit.AcquisitionPurpose
	test["metadata.ArchiveUnit.Handling.Author"] = metadata.ArchiveUnit.Handling.Author

	for fieldName, value := range test {
		if value == "" {
			return fmt.Errorf("field '%s' is empty", fieldName)
		}
	}
	_, err := time.Parse(time.RFC3339, metadata.ArchiveUnit.Deposit.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date '%s': %w", metadata.ArchiveUnit.Deposit.Date, err)
	}

	if metadata.ArchiveUnit.Type != "acquisition" {
		return fmt.Errorf("field 'metadata.ArchiveUnit.Type' is '%s', but expected 'acquisition'", metadata.ArchiveUnit.Type)
	}

	return nil
}
