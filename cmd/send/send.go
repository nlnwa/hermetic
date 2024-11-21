package send

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/nlnwa/hermetic/cmd/internal/cmdutil"
	"github.com/nlnwa/hermetic/cmd/internal/flags"
	"github.com/nlnwa/hermetic/internal/dps"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	dirFlagName        string = "dir"
	dirFlagHelpMessage string = `path to the root directory with the following content structure:
/root-directory
├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm
│   ├── /kommuner_2023-20230611002729-0035-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
│   └── /checksum_transferred.md5
└── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm
    ├── /kommuner_2023-20230611002730-0036-veidemann-contentwriter-568c6f8545-frvcm.warc.gz
    └── /checksum_transferred.md5
... etc`
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().String(dirFlagName, "", dirFlagHelpMessage)
	if err := cmd.MarkFlagRequired(dirFlagName); err != nil {
		panic(err)
	}
}

func toOptions() SendOptions {
	return SendOptions{
		KafkaEndpoints: flags.GetKafkaEndpoints(),
		KafkaTopic:     flags.GetKafkaTopic(),
		Dir:            viper.GetString(dirFlagName),
	}
}

type SendOptions struct {
	KafkaTopic      string
	KafkaEndpoints  []string
	TeamsWebhookUrl string
	Dir             string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Continuously sends data to digital storage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdutil.HandleError(toOptions().Run())
		},
	}
	addFlags(cmd)

	return cmd
}

func (o SendOptions) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(o.KafkaEndpoints...),
		Topic:    o.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	config := bigcache.Config{
		Shards:      1024,
		LifeWindow:  24 * 7 * time.Hour,
		CleanWindow: 1 * time.Hour,
	}
	cache, err := bigcache.New(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	loadCache := func(msg *dps.Package) error {
		// Skip messages that are not web archive messages
		if !dps.IsWebArchiveMessage(msg) {
			return nil
		}
		// Skip messages that are not from the root directory
		if !strings.HasPrefix(msg.Path, o.Dir) {
			return nil
		}
		if err := cache.Set(msg.Path, []byte("Sent")); err != nil {
			return fmt.Errorf("failed to set '%s' in cache: %w", msg.Path, err)
		}
		return nil
	}

	err = dps.ReadLatestMessages(ctx, o.KafkaEndpoints, o.KafkaTopic, loadCache)
	if err != nil {
		return fmt.Errorf("failed to read latest messages: %w", err)
	}

	for {
		items, err := os.ReadDir(o.Dir)
		if err != nil {
			return fmt.Errorf("failed to read root path '%s': %w", o.Dir, err)
		}
		for _, entry := range items {
			if !entry.IsDir() {
				return fmt.Errorf("found file '%s', but expected only directories", filepath.Join(o.Dir, entry.Name()))
			}

			path := filepath.Join(o.Dir, entry.Name())
			_, err := cache.Get(path)
			if err == nil {
				continue
			}
			if !errors.Is(err, bigcache.ErrEntryNotFound) {
				return fmt.Errorf("failed to get '%s' from cache: %w", path, err)
			}

			slog.Info("Processing directory", "path", path)

			pkg := dps.CreatePackage(path, entry.Name(), dps.ContentTypeWarc)

			if err := dps.Send(ctx, writer, pkg); err != nil {
				return fmt.Errorf("failed to send message to kafka topic '%s': %w", o.KafkaTopic, err)
			}
			if err := cache.Set(path, []byte("Sent")); err != nil {
				return fmt.Errorf("failed to set '%s' in cache: %w", path, err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
