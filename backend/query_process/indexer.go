package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"errors"
	"math"
	"sort"
	"strings"
)

type DocumentVector struct {
	Doc    Document
	Vector map[string]float64
}

var index = make(map[string][]DocumentVector)
var idfMap = make(map[string]float64) // To hold IDF values of all terms

// Build the index and calculate the IDF for all terms
func BuildIndex(documents []Document) {
	totalDocs := float64(len(documents))

	for _, doc := range documents {
		vector := buildDocumentVector(doc)
		words := ProcessQuery(doc.Content)
		for _, word := range words {
			if _, ok := idfMap[word]; !ok {
				idfMap[word] = math.Log(totalDocs / float64(len(index[word])))
			}
			index[word] = append(index[word], vector)
		}
	}
}

func buildDocumentVector(doc Document) DocumentVector {
	vector := make(map[string]float64)
	words := ProcessQuery(doc.Content)
	wordCount := float64(len(words))

	// Calculate TF for each term in the document
	for _, word := range words {
		vector[word] += 1.0 / wordCount
	}

	// Multiply TF with IDF for each term to get TF-IDF score and normalize the
	// vector
	magnitude := 0.0
	for word, tf := range vector {
		tfIdf := tf * idfMap[word]
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	// Divide each term's TF-IDF score with the magnitude to get the unit vector
	for word := range vector {
		vector[word] /= math.Sqrt(magnitude)
	}

	return DocumentVector{Doc: doc, Vector: vector}
}

func SearchIndex(queryWords []string) ([]SearchResult, error) {
	if len(queryWords) == 0 {
		return nil, errors.New("empty query")
	}

	queryVector := buildQueryVector(queryWords)

	// Use a map to eliminate duplicates and update scores
	scoreMap := make(map[string]*SearchResult)

	// Rank the documents based on all query words.
	for _, word := range queryWords {
		// Check if the word is in the index
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {
				score := cosineSimilarity(queryVector, vector.Vector)

				// Adjust the score based on frequency of query words, document
				// length, and position of first query word
				frequency := float64(len(ProcessQuery(vector.Doc.Content)))
				position := float64(strings.Index(vector.Doc.Content, word))
				length := float64(len(vector.Doc.Content))

				adjustment := (1 + math.Log(frequency+1)) * (1 / (1 + math.Log(length+1)) * (1 / (1 + math.Log(position+1))))
				score *= adjustment

				if result, ok := scoreMap[vector.Doc.Id]; ok {
					// If the document is already in the scoreMap, update its score
					result.Score += score
				} else {
					// If the document is not in the scoreMap, add it
					scoreMap[vector.Doc.Id] = &SearchResult{Doc: vector.Doc, Score: score}
				}
			}
		}
	}

	// Convert the scoreMap to a slice
	results := make([]SearchResult, 0, len(scoreMap))
	for _, result := range scoreMap {
		results = append(results, *result)
	}

	// Sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func buildQueryVector(queryWords []string) map[string]float64 {
	vector := make(map[string]float64)
	wordCount := float64(len(queryWords))

	// Calculate TF for each term in the query
	for _, word := range queryWords {
		vector[word] += 1.0 / wordCount
	}

	// Multiply TF with IDF for each term to get TF-IDF score and normalize the vector
	magnitude := 0.0
	for word, tf := range vector {
		tfIdf := idfMap[word] * tf
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	// Divide each term's TF-IDF score with the magnitude to get the unit vector
	for word := range vector {
		vector[word] /= magnitude
	}

	return vector
}

func cosineSimilarity(vector1, vector2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0
	for word, value := range vector1 {
		dotProduct += value * vector2[word]
		magnitude1 += value * value
	}
	for _, value := range vector2 {
		magnitude2 += value * value
	}
	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt((magnitude2)))
}
