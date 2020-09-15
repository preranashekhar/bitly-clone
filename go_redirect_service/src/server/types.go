package main

import "time"

type snsData struct {
	ShortUrlHash    string
	LongUrl			string
	AccessedAt		time.Time
	AccessedIP		string
	UserAgent		string
}

type noSqlResp struct {
	LongUrl 		string
}
