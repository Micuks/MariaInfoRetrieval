package main

import (
	"MariaInfoRetrieval/data_process"
	"MariaInfoRetrieval/query_process"
	"fmt"
	"os"
	"strconv"

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

	// Global search result cache
	cache := query_process.NewCache()

	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		cacheKey := fmt.Sprintf("%s-%s-%s", q, c.Query("page"), c.Query("results_per_page"))

		// Return if hit cache
		if cachedResults, found := cache.Get(cacheKey); found {
			c.JSON(200, cachedResults)
			return
		}

		// Else search
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page number"})
			return
		}

		resultsPerPage, err := strconv.Atoi(c.DefaultQuery("results_per_page", "10"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid number of results per page"})
		}

		// Process the query
		queryWords := query_process.WordSplit(q)
		log.Debug("queryWords:", queryWords)

		// Search the index and calculate scores
		results, err := query_process.SearchIndex(queryWords, page, resultsPerPage)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching documents"})
			return
		}

		// Store search results in cache
		cache.Set(cacheKey, results)

		log.Debug(results)
		c.JSON(200, results)
	})

	r.Run(":9011")
}
