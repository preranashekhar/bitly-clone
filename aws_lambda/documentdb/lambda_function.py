import json
import pymongo
import sys
import os
from dateutil.parser import parse

def lambda_handler(event, context):
    print("event: ", event)
    # print("event records: ", event['Records'])
    
    ##Create a MongoDB client, open a connection to Amazon DocumentDB as a replica set and specify the read preference as secondary preferred
    client = pymongo.MongoClient(os.environ['MONGODB_URL']) 
    
    for record in event['Records']:
        payload = json.loads(record['body'])
        payload["AccessedAt"] = parse(payload["AccessedAt"])
        
        db = client.bitly_db
        col = db.used_shortlinks
        resp = col.insert_one(payload)
        
        print("resp: ", resp )
    
    ##Close the connection
    client.close()
            
