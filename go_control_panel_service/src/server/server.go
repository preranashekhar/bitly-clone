package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func sendMsgToSns(respShortUrl string) (bool, error){
	mySession := session.Must(session.NewSession())
	svc := sns.New(mySession, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))

	params := &sns.PublishInput{
		Message: aws.String(respShortUrl),
		TopicArn: aws.String(os.Getenv("NEW_SHORTLINKS_SNS_ARN")),
	}

	_, err := svc.Publish(params)

	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	return true, nil
}

func handlePostShorten(c *gin.Context) {
	fmt.Println("Client IP: ", c.ClientIP())
	fmt.Println("user agent: ", c.Request.Header.Get("User-Agent"))
	fmt.Println("referer: ", c.Request.Header.Get("Referer"))
	var reqJson shortURLReq
	if c.ShouldBindJSON(&reqJson); reqJson.LongUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "bad user request"})
		return
	}

	data := []byte(reqJson.LongUrl)
	shortUrl := md5.Sum(data)
	shortUrlStr := hex.EncodeToString(shortUrl[:])[:8]

	respShortUrl := &shortUrlResp{
		reqJson.LongUrl,
		shortUrlStr,
		time.Now(),
		c.ClientIP(),
		c.Request.Header.Get("User-Agent"),
	}

	jsonData, err := json.Marshal(respShortUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "bad user request"})
		return
	}

	_, err = sendMsgToSns(string(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message": "error from server. Try again later."})
		return
	}
	c.JSON(http.StatusOK, respShortUrl)
}

func initServer(port string) {
	router := gin.Default()

	router.POST("/shorten", handlePostShorten)

	router.Run(port)
}
