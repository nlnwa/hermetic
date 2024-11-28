package send

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
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

	excludeFlagName    string = "exclude"
	excludeHelpMessage string = `comma separated list of regular expressions to match directories that should be excluded from preloading to cache`
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().String(dirFlagName, "", dirFlagHelpMessage)
	if err := cmd.MarkFlagRequired(dirFlagName); err != nil {
		panic(err)
	}
	cmd.Flags().StringSlice(excludeFlagName, nil, excludeHelpMessage)
}

func toOptions() (SendOptions, error) {
	excludes := viper.GetStringSlice(excludeFlagName)

	var exclude []*regexp.Regexp
	for _, e := range excludes {
		r, err := regexp.Compile(e)
		if err != nil {
			return SendOptions{}, fmt.Errorf("failed to compile regexp '%s': %w", e, err)
		}
		exclude = append(exclude, r)
	}

	return SendOptions{
		KafkaEndpoints: flags.GetKafkaEndpoints(),
		KafkaTopic:     flags.GetKafkaTopic(),
		Dir:            viper.GetString(dirFlagName),
		Exclude:        exclude,
	}, nil
}

type SendOptions struct {
	KafkaTopic      string
	KafkaEndpoints  []string
	TeamsWebhookUrl string
	Dir             string
	Exclude         []*regexp.Regexp
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Continuously sends data to digital storage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := toOptions()
			if err != nil {
				return err
			}
			return cmdutil.HandleError(opts.Run())
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

	loadCache := func(msg *dps.Message) error {
		// Skip messages that are not web archive messages
		if !dps.IsWebArchiveMessage(msg) {
			return nil
		}
		// Skip messages that are not from the root directory
		if !strings.HasPrefix(msg.Path, o.Dir) {
			return nil
		}
		// Skip messages that are excluded explicitly
		for _, re := range o.Exclude {
			if re.MatchString(msg.Path) {
				return nil
			}
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

			msg := dps.CreateMessage(path, entry.Name(), dps.ContentTypeWarc)

			if err := dps.Send(ctx, writer, msg); err != nil {
				return fmt.Errorf("failed to send message to kafka topic '%s': %w", o.KafkaTopic, err)
			}
			if err := cache.Set(path, []byte("Sent")); err != nil {
				return fmt.Errorf("failed to set '%s' in cache: %w", path, err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
