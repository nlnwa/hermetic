package dps

import (
	"testing"
)

func TestIsWebArchiveOwnedInvalid(t *testing.T) {
	notWebArchive := KafkaResponse{
		Offset: int64(0),
		Key:    "key",

		DPSResponse: DigitalPreservationSystemResponse{
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
		},
	}
	if IsWebArchiveOwned(&notWebArchive.DPSResponse) {
		t.Errorf("Expected false, got true")
	}
}
func TestIsWebArchiveOwnedValid(t *testing.T) {
	webArchiveResponse := KafkaResponse{
		Offset: int64(0),
		Key:    "key",

		DPSResponse: DigitalPreservationSystemResponse{
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
		},
	}
	if !IsWebArchiveOwned(&webArchiveResponse.DPSResponse) {
		t.Errorf("Expected true, got false")
	}

}
