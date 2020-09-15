package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func handleGetSummary(c *gin.Context) {
	shortUrlHash := c.Param("shortUrlHash")

	// make request to DocumentDB to get stats
	client, err := getDocumentDBClient()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}
	collection := client.Database(os.Getenv("BITLY_DB")).Collection(os.Getenv("BITLY_COLLECTION"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"ShortUrlHash": shortUrlHash}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}
	defer client.Disconnect(ctx)

	summaryRespData := &summaryResp{count, "All clicks for the short url rolled up into a single field of clicks"}
	c.JSON(http.StatusOK, summaryRespData)
}

func handleGetRecentClicks(c *gin.Context) {
	// make request to DocumentDB to get stats
	client, err := getDocumentDBClient()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(os.Getenv("BITLY_DB")).Collection(os.Getenv("BITLY_COLLECTION"))

	options := options.Find()
	options.SetProjection(mongoProjection)
	options.SetSort(bson.D{{"AccessedAt", -1}})
	options.SetLimit(10)

	cur, err := collection.Find(ctx, bson.M{}, options)
	var clicksResp []bson.M
	if err = cur.All(ctx, &clicksResp); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}

	c.JSON(http.StatusOK, clicksResp)
}

func handleGetClicks(c *gin.Context) {
	shortUrlHash := c.Param("shortUrlHash")

	// make request to DocumentDB to get stats
	client, err := getDocumentDBClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(os.Getenv("BITLY_DB")).Collection(os.Getenv("BITLY_COLLECTION"))

	filter := bson.M{"ShortUrlHash": shortUrlHash}
	cur, err := collection.Find(ctx, filter, options.Find().SetProjection(mongoProjection))
	var clicksResp []bson.M
	if err = cur.All(ctx, &clicksResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}

	c.JSON(http.StatusOK, clicksResp)
}

func handleGetClicksByDate(c *gin.Context) {
	shortUrlHash := c.Param("shortUrlHash")

	// make request to DocumentDB to get stats
	client, err := getDocumentDBClient()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}
	collection := client.Database(os.Getenv("BITLY_DB")).Collection(os.Getenv("BITLY_COLLECTION"))

	matchStage := bson.D{{"$match", bson.D{{"ShortUrlHash", shortUrlHash}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$AccessedAt"},
		{"Clicks", bson.D{{"$sum", 1 }}}}}}

	projectStage := bson.D{{"$project", bson.D{{"Date", bson.D{{"$dateToString", bson.D{{"format", "%Y-%m-%d"},
		{"date", "$_id"}}}}},
		{"Clicks", "$Clicks"}, {"_id", 0}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}
	var clicksResp []bson.M
	if err = cur.All(ctx, &clicksResp); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}

	defer client.Disconnect(ctx)
	c.JSON(http.StatusOK, clicksResp)
}

func getDocumentDBClient() (*mongo.Client, error) {
	connectionURI := os.Getenv("MONGODB_CONNECTION_URL")
	fmt.Println("mongo conn: ", connectionURI)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("Connected to DocumentDB!")

	return client, nil
}

func initServer(port string) {
	router := gin.Default()

	router.GET("/links/:shortUrlHash/summary", handleGetSummary)
	router.GET("/links/:shortUrlHash/clicks", handleGetClicks)
	router.GET("/links/:shortUrlHash/clicks_by_date", handleGetClicksByDate)
	router.GET("/recent_clicks", handleGetRecentClicks)

	router.Run(port)
}
