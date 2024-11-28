package dps

import (
	"regexp"
	"testing"
	"time"
)

func TestIsWebArchiveOwnedInvalid(t *testing.T) {
	notWebArchiveResponse := &Message{
		ContentCategory: "something else",
		Date:            "date",
		Identifier:      "identifier",
		Urn:             "urn",
		Path:            "path",
		ContentType:     "contentType",
		Checks: []Check{
			{
				Status:  "status",
				Message: "message",
				Reason:  "reason",
				File:    "file",
			},
		},
	}

	if IsWebArchiveOwned(notWebArchiveResponse) {
		t.Errorf("Expected false, got true")
	}
}

func TestIsWebArchiveOwnedValid(t *testing.T) {
	webArchiveResponse := &Message{
		ContentCategory: "nettarkiv",
		Date:            "date",
		Identifier:      "identifier",
		Urn:             "urn",
		Path:            "path",
		ContentType:     "contentType",
		Checks: []Check{
			{
				Status:  "status",
				Message: "message",
				Reason:  "reason",
				File:    "file",
			},
		},
	}
	if !IsWebArchiveOwned(webArchiveResponse) {
		t.Errorf("Expected true, got false")
	}

}

func TestCreatePackage(t *testing.T) {
	dateFormat := "2006-01-02T15:04:05.000"
	expectedSubmissionInformationPackage := Message{
		Date:            time.Now().UTC().Format(dateFormat),
		ContentCategory: "nettarkiv",
		ContentType:     "warc",
		Identifier:      "no-nb_nettarkiv_nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt_[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}", //https://regex101.com/r/0XbiJ6/1
		Urn:             "URN:NBN:no-nb_nettarkiv_nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt",
		Path:            "/path/to/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt.warc.gz",
	}

	submissionInformationPackage := CreateMessage("/path/to/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt.warc.gz", "nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt", "warc")

	if expectedSubmissionInformationPackage.ContentCategory != submissionInformationPackage.ContentCategory {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.ContentCategory, submissionInformationPackage.ContentCategory)
	}

	if expectedSubmissionInformationPackage.ContentType != submissionInformationPackage.ContentType {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.ContentType, submissionInformationPackage.ContentType)
	}

	var compiledIdentifier = regexp.MustCompile(expectedSubmissionInformationPackage.Identifier)

	if !compiledIdentifier.MatchString(submissionInformationPackage.Identifier) {
		t.Errorf("Expected %s to match %s", submissionInformationPackage.Identifier, compiledIdentifier)
	}

	if expectedSubmissionInformationPackage.Urn != submissionInformationPackage.Urn {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.Urn, submissionInformationPackage.Urn)
	}

	if expectedSubmissionInformationPackage.Path != submissionInformationPackage.Path {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.Path, submissionInformationPackage.Path)
	}

	expectedDate, err := time.Parse(dateFormat, expectedSubmissionInformationPackage.Date)
	if err != nil {
		t.Errorf("Expected %s, got %s", expectedDate, expectedSubmissionInformationPackage.Date)
	}
	date, err := time.Parse(dateFormat, submissionInformationPackage.Date)
	if err != nil {
		t.Errorf("Expected %s, got %s", date, submissionInformationPackage.Date)
	}

	isBeforeCalculatedDate := date.Compare(expectedDate)

	if isBeforeCalculatedDate < 0 {
		t.Errorf("Expected %s to be before %s", date, expectedDate)
	}
}
