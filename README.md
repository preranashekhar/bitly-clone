Bitly clone - URL shortner

## Summary

This project implements a URL shortening service like bitly (bitly.com) which enables users to create short versions of their URLs for easy sharing. Some of the key features of the project are:

- A public facing API to create a short URL given a long URL
- A redirect service to redirect to the long URL given a short URL
- An analytics service API to get analytics and stats data on the URLs

## Demo
https://www.youtube.com/watch?v=AgEydaSlMAI


## Organization

 - bitly - the main folder of the bitly project
 - docs - folder that contains detailed documentation of multiple sections of the project
 - aws_lambda - contains AWS lamda functions
 - go_control_panel_service - control panel api service implemented in go
 - go_redirect_service - go redirect api service implemented in go
 - go_trending_service - go trend api service implemented in go

## UML Diagram

![Screen Shot 2020-05-02 at 9 45 49 PM](https://user-images.githubusercontent.com/55044852/80899411-96d15000-8cc4-11ea-8f90-872147319e8b.png)



## System Diagram (AWS)

![Bitly_System_Design_AWS](https://user-images.githubusercontent.com/55044852/80907083-64295680-8cc8-11ea-86bd-a1b33dacf36e.png)

## System Diagram (AWS + GCP)(Extra Credit)

![AWS+GCP](https://user-images.githubusercontent.com/55044852/80907081-5ffd3900-8cc8-11ea-841f-955952a44e48.png)


## Key Concepts

### CQRS
The project implements CQRS architecture using AWS SNS, SQS and Lambda triggers.

#### SNS topics

##### Why use SNS and not just SQS?
SQS currently does not support multiple topics similar to other systems like Kafka or Kinesis. To achieve an effect of fan-out, I
created multiple SQS queues and added them as subsribers to the SNS topic. This will allow the events to fan out the individual queues. Each queue will have independent downstream responders thereby achieving separation of concerns. 
There are two SNS topics
 - New-shortlinks - topic to publish events when new short links are created
 - Used-shortlinks - topic to publish events when short links are used

#### SQS queues
Below are the subscribers to each of the SNS topics:

SQS Subscribers of new-shortlinks
 - new-shortlinks-mysql - queue that contains new shortlink events to be stored in mysql db
 - new-shortlinks-nosql - queue that contains new shortlink events to be stored in nosql db

SQS Subscribers of new-shortlinks
 - used-shortlinks-documentdb - queue that contains used shortlink events to be stored in mongo db
 - used-shortlinks-mysql - queue that contains used shortlink events to be stored in mysql db

##### Why do we have multiple queues for each database
I wanted to make the components as independent as possible. Having multiple queues will enable that.


#### Lambda triggers

Each of the SQS queues has an attached lambda function trigger that will perform the action of storing data in database
 
 - shortlinksPollerForMysql - lambda trigger for new-shortlinks-mysql and used-shortlinks-mysql SQS queues that will insert the new and used shortlinks to mysql db
 
 - shortlinksPollerForNosqlDB - lambda trigger for new-shortlinks-nosql SQS queue that will insert new shortlinks into nosql db

 - usedShortLinksPollerForDocumentDB - lambda trigger for nused-shortlinks-documentdb SQS queue that will insert used shortlinks into nosql db


### Databases

#### MySQL DB
  - This is the main database that will store new and used short links. This db is mainly accessed by Lambda functions to store data. The DB is setup using AWS RDS service.
  - Docs: https://github.com/preranashekhar/bitly-clone/tree/master/docs/mysql

#### NoSQL DB
 - This DB is used as lookup cache by the redirect server when looking up for a long url given a short url. This is an implementation a 5 node database. This DB is deployed as a 5 pod service on AWS EKS.
 - Docs: https://github.com/preranashekhar/bitly-clone/tree/master/nosql

#### MongoDB
 - This DB is used to perform analytics query on new and used short links. MongoDB performs better and is well suited for aggregation operations compared to a KV store. The DB is hosted on MongoDB Atlas secured with IP whitelisting.
 - Docs: https://github.com/preranashekhar/bitly-clone/tree/master/docs/mongodb


### Service endpoints
The services are implemented in GO using GIN framework https://github.com/gin-gonic/gin
- Docs: https://github.com/preranashekhar/bitly-clone/tree/master/docs/services
- Control panel service
  - Service to create short url links given long urls
- Redirect service
  - Service to redirect short url links to their corresponding long urls
- Trend/Analytics service: 
  - Service to get analytics and stats data of the urls


### Docker Hub Images
- Control Panel Service: https://hub.docker.com/r/preranashekhar/controlpanel
- Redirect Service: https://hub.docker.com/r/preranashekhar/redirectserver
- Trend Service: https://hub.docker.com/r/preranashekhar/trendserver
- NoSQL DB modified for AWS deployment: https://hub.docker.com/r/preranashekhar/nosql-aws
