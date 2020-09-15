package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func sendMsgToSns(shortUrlHash string, longUrl string, accessedIP string, userAgent string) (bool, error) {
	data := &snsData{shortUrlHash, longUrl, time.Now(), accessedIP, userAgent}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("jsonData: ", string(jsonData))

	mySession := session.Must(session.NewSession())
	svc := sns.New(mySession, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	params := &sns.PublishInput{
		Message: aws.String(string(jsonData)),
		TopicArn: aws.String(os.Getenv("USED_URL_SNS_ARN")),
	}

	resp, err := svc.Publish(params)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	fmt.Println(resp)
	return true, nil
}

// fetching data from NoSql DB
func fetchLongUrlFromCacheDb(shortUrlHash string) (string, error) {
	fmt.Println("Fetching from Cache DB")
	resp, err := http.Get(os.Getenv("NOSQL_DB_ENDPOINT") + shortUrlHash)

	if err != nil {
		fmt.Println("Error from Nosql, fetchLongUrlFromCacheDb")
		return "", err
	}
	fmt.Println("long url from Nosql", resp)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error from Nosql, fetchLongUrlFromCacheDb")
		return "", err
	}

	var respData noSqlResp
	err = json.Unmarshal([]byte(body), &respData)
	defer resp.Body.Close()
	return respData.LongUrl, nil
}

func fetchLongUrlFromDb(shortUrlHash string) (string, error) {
	fmt.Println("Fetching from DB")
	fmt.Println("datasource: ", os.Getenv("MYSQL_CONNECTION_STRING"))
	db, err := sql.Open("mysql", os.Getenv("MYSQL_CONNECTION_STRING"))
	if err != nil {
		return "", err
	}

	fmt.Println("Connection is good")
	fmt.Println("shortUrlHash : ", shortUrlHash)
	q := "SELECT long_url longUrl, short_url shortUrl FROM BITLY_URLS where short_url=?"
	row := db.QueryRow(q, shortUrlHash)


	var longUrl, shortUrl string
	err = row.Scan(&longUrl, &shortUrl)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}

	fmt.Println("longUrl from DB: ", longUrl)
	fmt.Println("shortUrl from DB: ", shortUrl)

	return longUrl, nil
}

func handleGetRedirect(c *gin.Context) {
	shortUrlHash := c.Param("shortUrlHash")
	fmt.Println("shortUrlHash: ", shortUrlHash)

	longUrl, err := fetchLongUrlFromCacheDb(shortUrlHash)

	if err != nil {
		longUrl, err = fetchLongUrlFromDb(shortUrlHash)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message:": "error from server. Try again later."})
		return
	}

	_, err = sendMsgToSns(shortUrlHash, longUrl, c.ClientIP(), c.Request.Header.Get("User-Agent"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_message:": "error from server. Try again later."})
		return
	}

	fmt.Println("longUrl: ", longUrl)
	c.Redirect(http.StatusMovedPermanently, longUrl)
}


func serverInit(port string) {
	router := gin.Default()

	router.GET("/redirect/:shortUrlHash", handleGetRedirect)

	router.Run(port)
}
