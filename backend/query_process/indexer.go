package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"errors"
	"math"
	"sort"
)

type DocumentVector struct {
	Doc    Document
	Vector map[string]float64
}

var index = make(map[string][]DocumentVector)

func BuildIndex(documents []Document) {
	for _, doc := range documents {
		vector := buildDocumentVector(doc)
		words := ProcessQuery(doc.Content)
		for _, word := range words {
			index[word] = append(index[word], vector)
		}
	}
}

func buildDocumentVector(doc Document) DocumentVector {
	vector := make(map[string]float64)
	words := ProcessQuery(doc.Content)
	for _, word := range words {
		vector[word] += 1.0
	}
	return DocumentVector{Doc: doc, Vector: vector}
}

func SearchIndex(queryWords []string) ([]SearchResult, error) {
	if len(queryWords) == 0 {
		return nil, errors.New("empty query")
	}

	queryVector := buildQueryVector(queryWords)

	// Rank the documents based on all query words.
	results := []SearchResult{}
	for _, word := range queryWords {
		// Check if the word is in the index
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {
				score := cosineSimilarity(queryVector, vector.Vector)
				results = append(results, SearchResult{Doc: vector.Doc, Score: score})
			}
		}
	}

	// Sort results by score
	// TODO: calculate a score based on the relevance of the document to the
	// query.

	// This could involve many factors, such as the frequency of
	// query words in the document, the document length, the positions of
	// query words in the document, and so on.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func buildQueryVector(queryWords []string) map[string]float64 {
	vector := make(map[string]float64)
	for _, word := range queryWords {
		vector[word] += 1.0
	}
	return vector
}

func cosineSimilarity(vector1, vector2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0
	for _, value := range vector1 {
		dotProduct += value * vector2[value]
		magnitude1 += value * value
	}
	for _, value := range vector2 {
		magnitude2 += value * value
	}
	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt((magnitude2)))
}
