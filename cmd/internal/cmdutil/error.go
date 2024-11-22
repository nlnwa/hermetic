package cmdutil

import (
	"context"
	"log/slog"
	"time"

	"github.com/nlnwa/hermetic/cmd/internal/flags"
	"github.com/nlnwa/hermetic/internal/teams"
)

// HandleError sends error message to Teams and returns the error
func HandleError(err error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := teams.SendMessage(ctx, teams.Error(err), flags.GetTeamsWebhookNotificationUrl()); err != nil {
		slog.Error(err.Error())
	}

	return err
}
