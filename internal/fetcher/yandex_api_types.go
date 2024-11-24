package fetcher

type YandexSearchResponse struct {
	Header struct {
		Text string `json:"text"`
	} `json:"header"`
	Blocks []struct {
		Type    string `json:"type"`
		Payload []struct {
			Slug  string `json:"slug"`
			Title string `json:"title"`
			Tags  []struct {
				Title string `json:"title"`
			} `json:"tags"`
		} `json:"payload"`
	} `json:"blocks"`
}
