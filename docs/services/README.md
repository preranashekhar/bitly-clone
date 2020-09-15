### Control Panel service

This is the service to create short url links given long urls

#### Endpoints

 - https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/shorten
 - Method: POST
 - Request: {"longUrl": "<long_url>"}

#### Request

```
curl --location --request POST 'https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
	"longUrl": "https://github.com"
}'
```

#### Response
```
{
    "LongUrl": "https://github.com",
    "ShortUrlHash": "3097fca9",
    "CreatedAt": "2020-05-03T05:33:48.536418628Z",
    "CreatedIP": "98.210.57.147",
    "UserAgent": "PostmanRuntime/7.22.0"
}
```

#### How short url hash is calculated

I am using the MD5 hash algorithm (https://en.wikipedia.org/wiki/MD5) to get the short url hash which produces a 128 bit and I am taking the first 8 hex characters.

#### Screenshots



### Redirect Service
The service to redirect short url links to their corresponding long urls. This service is behind a different api gateway than control panel and trend service since this is a service that is most likely used in a browser. The service returns a http code 301 with the long url.

#### Endpoints
 - https://5c0mgyoj1b.execute-api.us-west-2.amazonaws.com/prod/redirect/<short_url_hash>
 - method: GET
 
#### Request

```
curl --location --request GET 'https://5c0mgyoj1b.execute-api.us-west-2.amazonaws.com/prod/redirect/90072175' \
--header 'Content-Type: application/json' \
--data-raw '{
	"longUrl": "https://github.com"
}'
```

#### Response
- Status code: 301
- Redirected to long url

#### Screenshots


### Trend Service

Service that displays analytics data about urls

#### Endpoints
 - https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/recent_clicks
   - method: GET
   - provides recent 10 links
 
 - https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/links/<short_url_hash>/clicks
   -  Method: GET
   - provides recent clicks for a short hash
 
 - https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/links/<short_url_hash>/summary
   - Method: GET
   - provides the total number of clicks for a short url hash
 
 - https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/links/<short_url_hash>/clicks_by_date
   - Method: GET
   - provides the clicks for a short url hash by date


#### Requests
 - Request: https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/recent_clicks
 - Response:
 ```
  [
    {
        "AccessedAt": "2020-05-03T05:43:03.003Z",
        "AccessedIP": "98.210.57.147",
        "LongUrl": "https://web.whatsapp.com/",
        "ShortUrlHash": "d4021167",
        "UserAgent": "PostmanRuntime/7.24.1"
    },
    {
        "AccessedAt": "2020-05-03T05:41:32.263Z",
        "AccessedIP": "98.210.57.147",
        "LongUrl": "https://www.coursera.org/",
        "ShortUrlHash": "90072175",
        "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.122 Safari/537.36"
    },
    {
        "AccessedAt": "2020-05-03T05:40:25.023Z",
        "AccessedIP": "98.210.57.147",
        "LongUrl": "https://www.coursera.org/",
        "ShortUrlHash": "90072175",
        "UserAgent": "PostmanRuntime/7.22.0"
    },
 ]
 ```
 - Request: https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/links/d4021167/clicks
 - Response
   ```
   [
    {
        "AccessedAt": "2020-05-03T05:43:03.003Z",
        "AccessedIP": "98.210.57.147",
        "LongUrl": "https://web.whatsapp.com/",
        "ShortUrlHash": "d4021167",
        "UserAgent": "PostmanRuntime/7.24.1"
    }
    ]
   ```
- Request: https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/links/d4021167/summary
- Response
```
{
    "TotalClicks": 2,
    "Description": "All clicks for the short url rolled up into a single field of clicks"
}
```
- Request: https://n1ccua350b.execute-api.us-west-2.amazonaws.com/prod/data/links/d4021167/clicks_by_date
- Reponse
```
[
    {
        "Clicks": 1,
        "Date": "2020-05-03"
    },
    {
        "Clicks": 1,
        "Date": "2020-05-02"
    }
]
```

