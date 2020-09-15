### MongoDB

This DB is used to perform analytics query on new and used short links. MongoDB performs better and is well suited for aggregation operations compared to a KV store. The DB is hosted on MongoDB Atlas secured with IP whitelisting.

### Cluster

![Screen Shot 2020-05-02 at 10 25 53 PM](https://user-images.githubusercontent.com/55044852/80899367-26c2ca00-8cc4-11ea-9e96-8c741b667ff9.png)

### Primary and Replicas
![Screen Shot 2020-05-02 at 10 26 04 PM](https://user-images.githubusercontent.com/55044852/80899365-26c2ca00-8cc4-11ea-81d5-cf7665f11cb5.png)


### DB and collections
![Screen Shot 2020-05-02 at 10 26 27 PM](https://user-images.githubusercontent.com/55044852/80899364-262a3380-8cc4-11ea-9c4d-8634bf432731.png)

### NAT IP Whitelisting
![Screen Shot 2020-05-02 at 10 26 56 PM](https://user-images.githubusercontent.com/55044852/80899360-232f4300-8cc4-11ea-8c77-9b5c13d8e6d1.png)

### DB Access users
![Screen Shot 2020-05-02 at 10 26 34 PM](https://user-images.githubusercontent.com/55044852/80899362-24f90680-8cc4-11ea-8a62-030005dee960.png)


### Security

- The DB is secured through whitelisting AWS NAT Gateway IP
