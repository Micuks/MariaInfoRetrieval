package main

import (
	"MariaInfoRetrieval/data_process"
	"MariaInfoRetrieval/query_process"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	// log.SetReportCaller(true)
}

// Global variable to hold our "index"
// var index map[string][]Document

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// Load and index documents
	docs, err := data_process.LoadDocuments("./data")
	if err != nil {
		log.Error("Failed to load documents: %v", err)
	}
	query_process.BuildIndex(docs)

	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")

		// Process the query
		queryWords := query_process.WordSplit(q)
		log.Debug("queryWords:", queryWords)

		// Search the index and calculate scores
		results, err := query_process.SearchIndex(queryWords)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching documents"})
			return
		}

		log.Debug(results)
		c.JSON(200, results)
	})

	r.Run(":9011")
}
