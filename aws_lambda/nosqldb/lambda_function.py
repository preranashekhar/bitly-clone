import json
import requests
import os

def lambda_handler(event, context):
    
    nosql_endpoint = os.environ("NOSQL_ENDPOINT) + "{}"
    
    for record in event["Records"]:
        print("record: ", record)
        payload = json.loads(record['body'])
        print("payload is: ",payload)
        short_url_hash = payload["ShortUrlHash"]
        r = requests.post(nosql_endpoint.format(short_url_hash), json=payload)
        print("Response code: ", r.status_code)
    
