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

type Message struct {
	Date            string  `json:"date"`
	Identifier      string  `json:"identifier"`
	Urn             string  `json:"urn"`
	Path            string  `json:"path"`
	ContentType     string  `json:"contentType"`
	ContentCategory string  `json:"contentCategory"`
	Checks          []Check `json:"checks,omitempty"`
}

type KafkaMessage struct {
	Offset int64
	Key    string
	Value  Message
}

func CreateMessage(path string, payloadDirName string, contentType string) Message {
	date := time.Now().UTC().Format("2006-01-02T15:04:05.000")
	contentCategory := "nettarkiv"
	commonPart := "no-nb_" + contentCategory + "_" + payloadDirName
	identifier := commonPart + "_" + uuid.New().String()
	urn := "URN:NBN:" + commonPart

	return Message{
		Date:            date,
		ContentCategory: contentCategory,
		ContentType:     contentType,
		Identifier:      identifier,
		Urn:             urn,
		Path:            path,
	}
}

func IsWebArchiveOwned(message *Message) bool {
	return message.ContentCategory == "nettarkiv"
}

func IsWebArchiveMessage(message *Message) bool {
	if !IsWebArchiveOwned(message) {
		return false
	}

	switch message.ContentType {
	case ContentTypeWarc, ContentTypeAcquisition:
		return true
	default:
		return false
	}
}
