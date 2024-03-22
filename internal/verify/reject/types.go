package rejectImplementation

type check struct {
	Status  string
	Message string
	Reason  string
	File    string
}

type digitalPreservationSystemResponse struct {
	Date            string
	Identifier      string
	Urn             string
	Path            string
	ContentType     string
	ContentCategory string
	Checks          []check
}

type kafkaResponse struct {
	Offset      int64
	Key         string
	DPSResponse digitalPreservationSystemResponse
}
