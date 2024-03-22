package confirmImplmentation

type digitalPreservationSystemResponse struct {
	Date            string
	Identifier      string
	Urn             string
	Path            string
	ContentType     string
	ContentCategory string
}

type kafkaResponse struct {
	Offset      int64
	Key         string
	DPSResponse digitalPreservationSystemResponse
}
