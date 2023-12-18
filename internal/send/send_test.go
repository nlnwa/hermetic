package sendImplementation

import (
	"hermetic/internal/submission_information_package"
	"testing"
	"time"
)

func TestWebArchiveRelevantMessages(t *testing.T) {
	nettarkivetMessage := submission_information_package.Package{
		Date:            time.Now().UTC().Format("2006-01-02T15:04:05.000"),
		ContentCategory: "nettarkiv",
		ContentType:     "warc",
		Identifier:      "not-important",
		Urn:             "not-important",
		Path:            "not-important",
	}
	otherMessage := submission_information_package.Package{
		Date:            time.Now().UTC().Format("2006-01-02T15:04:05.000"),
		ContentCategory: "other",
		ContentType:     "other",
		Identifier:      "other",
		Urn:             "other",
		Path:            "other",
	}
	messages := []submission_information_package.Package{nettarkivetMessage, otherMessage}
	filteredResults, err := webArchiveRelevantMessages(messages)
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	if len(filteredResults) != 1 {
		t.Errorf("Expected 1 message, got %d", len(filteredResults))
	}
	if filteredResults[0] != nettarkivetMessage {
		t.Errorf("Expected '%+v', got '%s'", nettarkivetMessage, filteredResults[0])
	}
}
