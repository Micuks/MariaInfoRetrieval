package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"errors"
	"math"
	"sort"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type DocumentVector struct {
	Doc    Document
	Vector map[string]float64
}

var index = make(map[string][]DocumentVector)
var idDocMap = make(map[string]Document)
var docs []Document
var idfMap = make(map[string]float64) // To hold IDF values of all terms
var epsilon = 1e-10

// Build the index and calculate the IDF for all terms
func BuildIndex(documents []Document) {
	// Store documents
	docs = documents

	totalDocs := float64(len(documents))
	docIndex := make(map[string][]Document)

	// First, build the index, and doc's summary at the same time
	for _, doc := range documents {
		// Build the idDocMap
		idDocMap[doc.Id] = doc

		// Build the index
		words := WordSplit(doc.Content)
		for _, word := range words {
			docIndex[word] = append(docIndex[word], doc)
		}
	}

	log.Info("totalDocs:", totalDocs)
	// Second, calculate the IDF values for all words
	for word := range docIndex {
		if _, ok := idfMap[word]; !ok {
			idfMap[word] = math.Log(totalDocs / float64(len(docIndex[word])))
		}
	}

	// Third, build index of []DocumentVector
	for _, doc := range documents {
		docVector := buildDocumentVector(doc)
		words := WordSplit(doc.Content)

		for _, word := range words {
			index[word] = append(index[word], docVector)
		}
	}

}

func buildDocumentVector(doc Document) DocumentVector {
	vector := make(map[string]float64)
	words := WordSplit(doc.Content)
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

	if magnitude > 0.0 {

		sqrtMagnitude := math.Sqrt(magnitude + epsilon)

		// Divide each term's TF-IDF score with the magnitude to get the unit vector
		for word := range vector {
			vector[word] /= sqrtMagnitude
		}
	}

	return DocumentVector{Doc: doc, Vector: vector}
}

func buildSummaryDocument(doc Document) SummaryDocument {
	summaryDoc := SummaryDocument{
		Id:      doc.Id,
		Title:   doc.Title,
		URL:     doc.URL,
		Date:    doc.Date,
		Content: calculateSummary(doc.Content),
	}
	return summaryDoc
}

func SearchIndex(queryWords []string, page, resultsPerPage int) ([]SearchResult, error) {
	if len(queryWords) == 0 {
		return nil, errors.New("empty query")
	}

	queryVector := buildQueryVector(queryWords)
	log.Info("queryVector:", queryVector)

	// Handle the situation when the query words exist in all or none documents
	magnigude := 0.0
	for _, tfidf := range queryVector {
		magnigude += tfidf * tfidf
	}
	if magnigude == 0 {
		log.Info("Query made up of words in every or no documents. Returning all documents.")
		results := make([]SearchResult, 0, len(docs))
		for _, doc := range docs {
			results = append(results, SearchResult{Doc: buildSummaryDocument(doc), Score: 1.0})
		}

		return results, nil
	}

	vectorCounts := make(map[string]int)
	// Count the total number of vectors for each document ID across all query words
	for _, word := range queryWords {
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {
				vectorCounts[vector.Doc.Id]++
			}
		}
	}

	// Parallel result computation
	scoresChansMap := make(map[string]chan float64)
	for id, count := range vectorCounts {
		scoresChansMap[id] = make(chan float64, count)
	}

	var wg sync.WaitGroup

	for _, word := range queryWords {
		if vectors, ok := index[word]; ok {
			for _, vector := range vectors {

				wg.Add(1)
				go func(v DocumentVector, scoresChan chan float64) {
					defer wg.Done()

					// Calculate score
					score := cosineSimilarity(queryVector, v.Vector)

					// Adjust the score based on:
					// -  frequency of query words
					// - document length
					// - position of first query word
					frequency := float64(len(WordSplit(v.Doc.Content)))
					position := float64(strings.Index(v.Doc.Content, word))
					length := float64(len(v.Doc.Content))

					adjustment := (1 + math.Log(frequency+1)) * (1 / (1 + math.Log(length+1)) * (1 / (1 + math.Log(position+1))))
					score *= adjustment

					scoresChan <- score
				}(vector, scoresChansMap[vector.Doc.Id])
			}
		}
	}

	// Wait for all goroutines to finish, then close the results channel
	go func() {
		wg.Wait()
		for _, scoresChan := range scoresChansMap {
			close(scoresChan)
		}
	}()

	// Collect the results
	scoreMap := make(map[string]*SearchResult)
	for id, scoresChan := range scoresChansMap {
		totalScore := 0.0
		for score := range scoresChan {
			totalScore += score
		}
		summaryDoc := buildSummaryDocument(idDocMap[id])
		scoreMap[id] = &SearchResult{Doc: summaryDoc, Score: totalScore}
	}

	log.Info(">>> scoreMap")
	for k, v := range scoreMap {
		log.Info(k, ":", "Doc:", v.Doc, "Score:", v.Score)
	}
	log.Info("<<< scoreMap")

	// Convert the scoreMap to a slice
	results := make([]SearchResult, 0, len(scoreMap))
	for _, result := range scoreMap {
		results = append(results, *result)
	}

	// Sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Apply pagination
	start := (page - 1) * resultsPerPage
	end := start + resultsPerPage
	if start > len(results) {
		start = len(results)
	}
	if end > len(results) {
		end = len(results)
	}

	results = results[start:end]

	return results, nil
}

func GetFullDoc(id string) (Document, bool) {
	doc, ok := idDocMap[id]
	return doc, ok
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
		idf, ok := idfMap[word]
		if !ok {
			// Skip non-indexed words
			continue
		}
		tfIdf := idf * tf
		vector[word] = tfIdf
		magnitude += tfIdf * tfIdf
	}

	if magnitude > 0.0 {
		sqrtMagnitude := math.Sqrt(magnitude + epsilon)

		// Divide each term's TF-IDF score with the magnitude to get the unit vector
		// Only if magnitude is non-zero
		for word := range vector {
			vector[word] /= sqrtMagnitude
		}
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

	sqrtEpsMag1 := math.Sqrt(magnitude1 + epsilon)
	sqrtEpsMag2 := math.Sqrt(magnitude2 + epsilon)
	return dotProduct / (sqrtEpsMag1 * sqrtEpsMag2)
}

func calculateSummary(content string) string {
	if len(content) > 100 {
		return content[:100] + "..."
	}
	return content
}
