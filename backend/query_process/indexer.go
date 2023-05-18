package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"errors"
	"sort"
	"strings"
)

type DocumentVector struct {
	Doc Document
	Vector map[string]float64
}

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

	// Rank the documents based on all query words.
	results := []SearchResult{}
	for _, word := range queryWords {
		// Check if the word is in the index
		if docs, ok := index[word]; ok {
			for _, doc := range docs {
				score := tfIdf(word, doc, docs)
				results = append(results, SearchResult{Doc: doc, Score: score})
			}
		}
	}

	// Sort results by score
		// TODO: calculate a score based on the relevance of the document to the
		// query.

		// This could involve many factors, such as the frequency of
		// query words in the document, the document length, the positions of
		// query words in the document, and so on.
	sort.Slice(results, func(i, j int) bool {i, j int} bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func tfIdf(word string, doc Document, docs []Document) float64 {
	tf := termFrequency(word, doc)
	idf := inverseDocumentFrequency(word, docs)
	return tf * idf
}

func termFrequency(word string, doc Document) float64 {
	// Term frequency is the count of the word in the document divided by the
	// total number of words in the document
	words := strings.Fields(doc.Content)
	wordCount := float64(len(words))
	termCount := float64(0)
	for _, w := range words {
		if w == word {
			termCount++
		}
	}
	return termCount / wordCount
}

func inverseDocumentFrequency(word string, docs []Document) float64 {
	// Inverse document frequency is the log of the total number of documents
	// divided by the numebr of documents containing the word
	totalDocs := float64(len(docs))
	docsWithWord := float64(0)
	for _, doc := range docs {
		if strings.Contains(doc.Content, word) {
			docsWithWord++
		}
	}
	
	return math.Log(totalDocs / docsWithWord)
}