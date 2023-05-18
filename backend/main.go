package main

import (
	data_process "MariaInfoRetrieval/data_process"
	"MariaInfoRetrieval/query_process"
	"log"

	"github.com/gin-gonic/gin"
)

// Global variable to hold our "index"
// var index map[string][]Document

func main() {
	r := gin.Default()

	// Load and index documents
	docs, err := data_process.LoadDocuments("./data")
	if err != nil {
		log.Fatalf("Failed to load documents: %v", err)
	}
	query_process.BuildIndex(docs)

	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")

		// Process the query
		queryWords := query_process.ProcessQuery(q)

		// Search the index and calculate scores
		results, err := query_process.SearchIndex(queryWords)
		if err != nil {
			c.JSON(500, gin.H{"error": "Eror fetching documents"})
			return
		}

		c.JSON(200, results)
	})

	r.Run(":9011")
}
