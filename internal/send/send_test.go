package sendImplementation

import (
	"testing"
	"time"
)

func TestCreateSubmissionInformationPackage(t *testing.T) {
	dateFormat := "2006-01-02T15:04:05.000"
	expectedSubmissionInformationPackage := TransferSubmissionInformationPackage{
		Date:            time.Now().UTC().Format(dateFormat),
		ContentCategory: "nettarkiv",
		ContentType:     "warc",
		Identifier:      "no-nb_nettarkiv_nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt",
		Urn:             "URN:NBN:no-nb_nettarkiv_nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt",
		Path:            "/path/to/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt.warc.gz",
	}

	submissionInformationPackage := createSubmissionInformationPackage("/path/to/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt/nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt.warc.gz", "nettaviser_SCREENSHOT_2023-20230718002403-0216-veidemann-contentwriter-5bb4677d67-qwtmt")

	if expectedSubmissionInformationPackage.ContentCategory != submissionInformationPackage.ContentCategory {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.ContentCategory, submissionInformationPackage.ContentCategory)
	}

	if expectedSubmissionInformationPackage.ContentType != submissionInformationPackage.ContentType {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.ContentType, submissionInformationPackage.ContentType)
	}

	if expectedSubmissionInformationPackage.Identifier != submissionInformationPackage.Identifier {
		t.Errorf("Expected %s, got %s", expectedSubmissionInformationPackage.Identifier, submissionInformationPackage.Identifier)
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
