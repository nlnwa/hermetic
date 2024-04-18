package dps

type Check struct {
	Status  string
	Message string
	Reason  string
	File    string
}

type DigitalPreservationSystemResponse struct {
	Date            string
	Identifier      string
	Urn             string
	Path            string
	ContentType     string
	ContentCategory string
	Checks          []Check
}

type KafkaResponse struct {
	Offset      int64
	Key         string
	DPSResponse DigitalPreservationSystemResponse
}

func IsWebArchiveOwned(message *DigitalPreservationSystemResponse) bool {
	return message.ContentCategory == "nettarkiv"
}
