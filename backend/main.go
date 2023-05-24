package main

import (
	"MariaInfoRetrieval/data_process"
	. "MariaInfoRetrieval/maria_types"
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
	cache_capacity := 10
	cache := query_process.NewCache(cache_capacity)

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
		log.Info("queryWords:", queryWords)

		// Search the index and calculate scores
		results, err := query_process.SearchIndex(queryWords, page, resultsPerPage)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching documents"})
			return
		}

		// Store search results in cache
		cache.Set(cacheKey, results)

		c.JSON(200, results)
	})

	// Fetch SearchResult content details
	r.GET("/document", func(c *gin.Context) {
		id := c.Query("id")
		// Search for the document with the given id in docs
		doc, found := query_process.GetFullDoc(id)
		if !found {
			c.JSON(404, gin.H{"error": "Document" + id + " not found"})
		}

		c.JSON(200, doc)
	})

	// Handle feedback
	r.POST("/feedback", func(c *gin.Context) {
		var feedback Feedback
		if err := c.BindJSON(&feedback); err != nil {
			c.JSON(400, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Process the feedback here...
		log.Infof("Received feedback: %v", feedback)

		c.JSON(200, gin.H{"message": "Feedback received successfully"})
	})

	r.Run(":9011")
}
