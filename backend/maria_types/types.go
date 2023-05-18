package maria_types

type Document struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
	Date    string `json:"date"`
}

type SearchResult struct {
	Score float64
	Doc   Document
}
