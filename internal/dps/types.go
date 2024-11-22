package dps

import (
	"time"

	"github.com/google/uuid"
)

type Check struct {
	Status  string
	Message string
	Reason  string
	File    string
}

type Response struct {
	Date            string
	Identifier      string
	Urn             string
	Path            string
	ContentType     string
	ContentCategory string
	Checks          []Check
}

type KafkaResponse struct {
	Offset   int64
	Key      string
	Response Response
}

type Package struct {
	Date            string `json:"date"`
	ContentCategory string `json:"contentCategory"`
	ContentType     string `json:"contentType"`
	Identifier      string `json:"identifier"`
	Urn             string `json:"urn"`
	Path            string `json:"path"`
}

func CreatePackage(path string, payloadDirName string, contentType string) Package {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.000")
	contentCategory := "nettarkiv"
	commonPart := "no-nb_" + contentCategory + "_" + payloadDirName
	identifier := commonPart + "_" + uuid.New().String()
	urn := "URN:NBN:" + commonPart

	return Package{
		Date:            date,
		ContentCategory: contentCategory,
		ContentType:     contentType,
		Identifier:      identifier,
		Urn:             urn,
		Path:            path,
	}
}

func IsWebArchiveOwned(message *Response) bool {
	return message.ContentCategory == "nettarkiv"
}
