package submission_information_package

import (
	"time"

	"github.com/google/uuid"
)

type Package struct {
	Date            string `json:"date"`
	ContentCategory string `json:"contentCategory"`
	ContentType     string `json:"contentType"`
	Identifier      string `json:"identifier"`
	Urn             string `json:"urn"`
	Path            string `json:"path"`
}

func CreatePackage(payloadPath string, payloadDirName string) Package {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.000")
	contentCategory := "nettarkiv"
	contentType := "warc"
	commonPart := "no-nb_" + contentCategory + "_" + payloadDirName
	identifier := commonPart + "_" + uuid.New().String()
	urn := "URN:NBN:" + commonPart

	return Package{
		Date:            date,
		ContentCategory: contentCategory,
		ContentType:     contentType,
		Identifier:      identifier,
		Urn:             urn,
		Path:            payloadPath,
	}
}
