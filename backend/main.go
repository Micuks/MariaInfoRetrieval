package main

import (
	. "MariaInfoRetrieval/maria_types"
	"MariaInfoRetrieval/query_process"

	"github.com/gin-gonic/gin"
)

// Global variable to hold our "index"
var index map[string][]Document

func main() {
	index = make(map[string][]Document)

	r := gin.Default()

	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")

		// Process the query
		words := query_process.ProcessQuery(q)

		// Search the index and calculate scores
		results, err := query_process.SearchIndex(words)
	})

	r.Run(":9011")
}
