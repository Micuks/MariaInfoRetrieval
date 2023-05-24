package maria_types

type Document struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Keywords string `json:"keywords"`
	URL      string `json:"url"`
	Date     string `json:"date"`
}

type SummaryDocument struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
	Date    string `json:"date"`
}

type SearchResult struct {
	Score float64
	Doc   SummaryDocument
}

// Struct of feedback
type Feedback struct {
	ResultId string `json:"resultId"`
	Score    int    `json:"Score"`
}

// Struct of PerformSearch's response
type Response struct {
	Code    int            `json:"code"`
	Results []SearchResult `json:"results"`
}
