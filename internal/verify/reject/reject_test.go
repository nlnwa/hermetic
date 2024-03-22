package rejectImplementation

import (
	"testing"
)

func TestIsWebArchiveOwnedInvalid(t *testing.T) {
	notWebArchive := kafkaResponse{
		Offset: int64(0),
		Key:    "key",

		DPSResponse: digitalPreservationSystemResponse{
			ContentCategory: "something else",
			Date:            "date",
			Identifier:      "identifier",
			Urn:             "urn",
			Path:            "path",
			ContentType:     "contentType",
			Checks: []check{
				{
					Status:  "status",
					Message: "message",
					Reason:  "reason",
					File:    "file",
				},
			},
		},
	}
	if isWebArchiveOwned(notWebArchive) {
		t.Errorf("Expected false, got true")
	}
}
func TestIsWebArchiveOwnedValid(t *testing.T) {
	webArchiveResponse := kafkaResponse{
		Offset: int64(0),
		Key:    "key",

		DPSResponse: digitalPreservationSystemResponse{
			ContentCategory: "nettarkiv",
			Date:            "date",
			Identifier:      "identifier",
			Urn:             "urn",
			Path:            "path",
			ContentType:     "contentType",
			Checks: []check{
				{
					Status:  "status",
					Message: "message",
					Reason:  "reason",
					File:    "file",
				},
			},
		},
	}
	if !isWebArchiveOwned(webArchiveResponse) {
		t.Errorf("Expected true, got false")
	}

}
