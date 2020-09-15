package main

import "time"

type shortUrlResp struct {
	LongUrl			string
	ShortUrlHash    string
	CreatedAt		time.Time
	CreatedIP		string
	UserAgent		string
}

type shortURLReq struct {
	LongUrl 		string
}