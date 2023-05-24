# Maria Information Retrieval

- [Maria Information Retrieval](#maria-information-retrieval)
  - [Overall design](#overall-design)
    - [About the search engine](#about-the-search-engine)
    - [Process the query](#process-the-query)
    - [Ranking documents](#ranking-documents)
      - [Cosine similarity for vector space model](#cosine-similarity-for-vector-space-model)
      - [TF-IDF](#tf-idf)
        - [TF](#tf)
        - [IDF](#idf)
    - [Results sorting](#results-sorting)
    - [Parallel computing](#parallel-computing)

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

Based on the basic vector space model, I improved it by introducing term
weighting and normalization as below.

1. Use the TF-IDF score for term weighting.
2. Normalize document and query vectors to unit length.

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

#### TF-IDF

##### TF

Term frequency works by looking at the frequency of a particular term you are concerned with relative to the document.

$$
tf(t, d) = \frac{count(t)}{\sum_{w\in d}{count(w)}}
$$

##### IDF

Inverse document frequency looks at how common (or uncommon) a word is amongst the corpus. IDF is calculated as follows where t is the term (word) we are looking to measure the commonness of and N is the number of documents (d) in the corpus (D).. The denominator is simply the number of documents in which the term, t, appears in.

$$
idf(t, D) = log(\frac{N}{count(d\in D:t\in d)})
$$

The reason we need IDF is to help correct for words like “of”, “as”, “the”, etc. since they appear frequently in an English corpus. Thus by taking inverse document frequency, we can minimize the weighting of frequent terms while making infrequent terms have a higher impact.

### Results sorting

When processing a search request, after gettings result scores computed by
vector space model introduced as above, the following things are also taken into
account: the frequency of query words, the document length, and the positions of query words.

In `SearchIndex` function, it introduces a score adjustment based on the
frequency of query words, the document length, and the position of the first
query word in the document.

It uses a map to eliminate duplicate documents and to update the score of a
document each time it is encountered.

It adjusts the score of each document based on several factors:

- The frequency of the query words in the document. The more often the query words appear in the document, the higher the score.
- The length of the document. The longer the document, the lower the score. This helps to prevent very long documents from getting artificially high scores simply because they have more words.
- The position of the first query word in the document. The sooner a query word appears in the document, the higher the score.
- These adjustments are added to the cosine similarity score to produce the final score for each document.

### Parallel computing

When processing query, it is time-consuming for computing a TF-IDF score for
each Query word vector. So I parallel compute the match scores for all
documents.

Create a channel once for each document ID, with a size based on the total
number of instances that document ID appears in the vectors of all query words.

The first loop counts the total number of vectors for each document ID across
all query words, and then the channels are created based on these counts. This
ensures that each channel is large enough to accommodate all the scores
corresponding to a specific document ID, and none of the channels will be
replaced and prematurely closed.
