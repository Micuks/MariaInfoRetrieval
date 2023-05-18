package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"errors"
)

var index = make(map[string][]Document)

func BuildIndex(documents []Document) {
	for _, doc := range documents {
		words := ProcessQuery(doc.Content)
		for _, word := range words {
			index[word] = append(index[word], doc)
		}
	}
}

func SearchIndex(queryWords []string) ([]SearchResult, error) {
	if len(queryWords) == 0 {
		return nil, errors.New("empty query")
	}

	// Just a demo

	// TODO: You would need to implement a more sophisticated algorithm to rank
	// the documents based on all query words.

	documents := index[queryWords[0]]
	results := make([]SearchResult, len(documents))

	for i, doc := range documents {
		// TODO: calculate a score based on the relevance of the document to the
		// query.

		// This could involve many factors, such as the frequency of
		// query words in the document, the document length, the positions of
		// query words in the document, and so on.
		results[i] = SearchResult{Doc: doc, Score: 1}
	}

	return results, nil
}
