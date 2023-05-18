package maria_types

type Document struct {
	Id      string
	Title   string
	Content string
	URL     string
	Date    string
}

type SearchResult struct {
	Score float64
	Doc   Document
}
