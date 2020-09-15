package main

import "go.mongodb.org/mongo-driver/bson"

type summaryResp struct {
	TotalClicks int64
	Description string
}

var mongoProjection = bson.M{
	"AccessedAt": 1,
	"ShortUrlHash": 1,
	"LongUrl": 1,
	"AccessedIP": 1,
	"UserAgent": 1,
	"_id": 0,
}