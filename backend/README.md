# Maria Information Retrieval

- [Maria Information Retrieval](#maria-information-retrieval)
  - [Overall design](#overall-design)
    - [About the search engine](#about-the-search-engine)
    - [Process the query](#process-the-query)
    - [Ranking documents](#ranking-documents)
      - [Cosine similarity for vector space model](#cosine-similarity-for-vector-space-model)

> Information retrieval system in Go

## Overall design

I use Gin framework and sets up a route handler for GET requests at "/search".
Inside the route handler, we use the hand-made search client to create a new
search request, search the documents, and return the response.

### About the search engine

The search engine needs to support splitting Chinese text into words (with or
without a library), treating English text as space-separated words, creating an
inverted index, and supporting a basic vector space model for searching. User
queries should be natural language strings and search results should be sorted
by relevance, with each result containing the relevance score, document title,
matched content, URL, and date.

It can be splitted into the following steps:

1. Indexing: Parse your documents and build an inverted index. An inverted index is a data structure where for each unique word in all documents, you store a list of documents in which that word appears.

2. Query Processing: When a user query comes in, it needs to be processed in the same way as your documents were during indexing (e.g., word splitting for Chinese text, treating English text as space-separated words).

3. Ranking: For each word in the processed query, look up the corresponding list
   of documents in the inverted index, and calculate a score for each document
   based on the vector space model.

### Process the query

In the context of the question, we need to achieve a few things:

1. Process the query text (split Chinese text with a library like Jieba in Python or another equivalent in Go, and English text with whitespace).
2. Build an inverted index.
3. Implement a vector space model for search and ranking.

The Go language does not directly support Chinese text splitting in its standard
library, and the Jieba library is specific to Python. However, there is a Go
version of Jieba available named "gojieba", I use this for Chinese word
segmentation.

### Ranking documents

The vector space model represents documents and queries as vectors in a high-dimensional space, where each unique word in the corpus is a dimension. The relevance of a document to a query is then computed as the cosine of the angle between the document vector and the query vector.

#### Cosine similarity for vector space model

```go
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
```

In this code, DocumentVector is a struct that contains a Document and a Vector, which is a map from words to their term frequencies. BuildIndex now calls buildDocumentVector to construct the document vectors.

The SearchIndex function has been modified to build a query vector and calculate the cosine similarity between the query vector and each document vector, rather than using the TF-IDF model.

buildQueryVector is a helper function to construct the query vector, and
cosineSimilarity calculates the cosine similarity between two vectors.
