import json
import pymongo
import sys

def lambda_handler(event, context):
    print("event: ", event)
    # print("event records: ", event['Records'])
    
    ##Create a MongoDB client, open a connection to Amazon DocumentDB as a replica set and specify the read preference as secondary preferred
    client = pymongo.MongoClient('mongodb+srv://mongo_read_only:pHJAlzOeGinWEPv1@bitlym0-kkugt.mongodb.net/test?retryWrites=true&w=majority') 
    
    for record in event['Records']:
        print("record: ", record)
        payload = json.loads(record['body'])
        print("payload: ", payload )
        
        db = client.bitly_db
        col = db.used_shortlinks
        resp = col.insert_one(payload)
        
        print("resp: ", resp )
    
    ##Close the connection
    client.close()
            
        